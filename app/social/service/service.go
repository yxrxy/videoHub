package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/hertz-contrib/websocket"
	"github.com/yxrrxy/videoHub/app/social/model"
	"github.com/yxrrxy/videoHub/app/social/repository"
	"github.com/yxrrxy/videoHub/app/social/ws"
	"github.com/yxrrxy/videoHub/kitex_gen/social"
	"github.com/yxrrxy/videoHub/kitex_gen/user/userservice"
)

type SocialService struct {
	repo       *repository.Social
	wsManager  *ws.Manager
	userClient userservice.Client
}

func NewSocialService(repo *repository.Social, wsManager *ws.Manager, userClient userservice.Client) *SocialService {
	return &SocialService{
		repo:       repo,
		wsManager:  wsManager,
		userClient: userClient,
	}
}

// GetWSManager 获取 WebSocket 管理器
func (s *SocialService) GetWSManager() *ws.Manager {
	return s.wsManager
}

// RegisterWebSocketClient 注册WebSocket客户端
func (s *SocialService) RegisterWebSocketClient(userID int64, conn *websocket.Conn) {
	s.wsManager.RegisterClient(userID, conn)
	// 启动客户端的消息读取循环
	go ws.NewExtendedClient(userID, conn).ReadPump(s.wsManager)
}

// JoinChatRoom 让用户加入聊天室
func (s *SocialService) JoinChatRoom(userID int64, roomID int64) error {
	// 检查用户是否在线
	client, exists := s.wsManager.GetClient(userID)
	if exists {
		s.wsManager.JoinRoom(client, roomID)
		return nil
	}
	return errors.New("user not online")
}

// LeaveChatRoom 让用户离开聊天室
func (s *SocialService) LeaveChatRoom(userID int64, roomID int64) error {
	// 检查用户是否在线
	client, exists := s.wsManager.GetClient(userID)
	if exists {
		s.wsManager.LeaveRoom(client, roomID)
		return nil
	}
	return errors.New("user not online")
}

// GetOnlineUsers 获取在线用户列表
func (s *SocialService) GetOnlineUsers() []int64 {
	return s.wsManager.GetOnlineUsers()
}

// IsUserOnline 检查用户是否在线
func (s *SocialService) IsUserOnline(userID int64) bool {
	return s.wsManager.IsUserOnline(userID)
}

// BroadcastSystemMessage 广播系统消息
func (s *SocialService) BroadcastSystemMessage(content string) {
	msg := &ws.Message{
		Type:      ws.MessageTypeSystem,
		Content:   content,
		Timestamp: time.Now().Unix(),
	}
	s.wsManager.Broadcast(msg)
}

// 私信相关
func (s *SocialService) SendPrivateMessage(ctx context.Context, senderID, receiverID int64, content string) error {
	msg := &model.PrivateMessage{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		IsRead:     false,
	}
	if err := s.repo.SendPrivateMessage(ctx, msg); err != nil {
		return err
	}

	// 序列化消息
	wsMsg := &ws.Message{
		Type:    ws.MessageTypePrivate,
		From:    senderID,
		To:      receiverID,
		Content: content,
	}
	data, err := json.Marshal(wsMsg)
	if err != nil {
		return err
	}

	// 发送序列化后的消息
	s.wsManager.SendToUser(receiverID, data)
	return nil
}

// GetPrivateMessages 获取私信列表
func (s *SocialService) GetPrivateMessages(ctx context.Context, senderID, receiverID int64, page, size int32) (*social.PrivateMessagesResponse, error) {
	messages, total, err := s.repo.GetPrivateMessages(ctx, senderID, receiverID, int(page), int(size))
	if err != nil {
		return nil, err
	}

	// 转换为Thrift类型
	thriftMessages := make([]*social.PrivateMessage, 0, len(messages))
	for _, msg := range messages {
		thriftMsg := &social.PrivateMessage{
			Id:         msg.ID,
			SenderId:   msg.SenderID,
			ReceiverId: msg.ReceiverID,
			Content:    msg.Content,
			IsRead:     msg.IsRead,
			CreatedAt:  msg.CreatedAt.Unix(),
			UpdatedAt:  msg.UpdatedAt.Unix(),
		}
		if msg.DeletedAt != nil {
			deletedAt := msg.DeletedAt.Unix()
			thriftMsg.DeletedAt = &deletedAt
		}
		thriftMessages = append(thriftMessages, thriftMsg)
	}

	return &social.PrivateMessagesResponse{
		Messages: thriftMessages,
		Total:    total,
	}, nil
}

// 聊天室相关
func (s *SocialService) CreateChatRoom(ctx context.Context, name string, creatorID int64, roomType int8, memberIDs []int64) (*social.ChatRoom, error) {
	room := &model.ChatRoom{
		Name:      name,
		CreatorID: creatorID,
		Type:      roomType,
	}

	// 创建聊天室成员
	members := make([]model.ChatRoomMember, 0, len(memberIDs)+1)
	// 添加创建者
	members = append(members, model.ChatRoomMember{
		UserID:   creatorID,
		Role:     model.ChatRoomRoleCreator,
		Nickname: name, // 可以从用户服务获取昵称
	})

	// 添加其他成员
	for _, memberID := range memberIDs {
		if memberID != creatorID {
			members = append(members, model.ChatRoomMember{
				UserID:   memberID,
				Role:     model.ChatRoomRoleMember,
				Nickname: name, // 可以从用户服务获取昵称
			})
		}
	}

	room.Members = members
	if err := s.repo.CreateChatRoom(ctx, room); err != nil {
		return nil, err
	}

	// 转换为Thrift类型
	return convertChatRoomToThrift(room), nil
}

func (s *SocialService) GetChatRoom(ctx context.Context, roomID int64) (*social.ChatRoom, error) {
	room, err := s.repo.GetChatRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}

	// 转换为Thrift类型
	return convertChatRoomToThrift(room), nil
}

func (s *SocialService) GetUserChatRooms(ctx context.Context, userID int64) ([]*social.ChatRoom, error) {
	rooms, err := s.repo.GetUserChatRooms(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 转换为Thrift类型
	thriftRooms := make([]*social.ChatRoom, 0, len(rooms))
	for i := range rooms {
		thriftRooms = append(thriftRooms, convertChatRoomToThrift(&rooms[i]))
	}

	return thriftRooms, nil
}

// 聊天消息相关
func (s *SocialService) SendChatMessage(ctx context.Context, roomID, senderID int64, content string, msgType int8) error {
	msg := &model.ChatMessage{
		RoomID:   roomID,
		SenderID: senderID,
		Content:  content,
		Type:     msgType,
	}
	if err := s.repo.SendChatMessage(ctx, msg); err != nil {
		return err
	}

	// 通过WebSocket推送消息
	s.wsManager.SendToRoom(roomID, &ws.Message{
		Type:    ws.MessageTypeGroup,
		From:    senderID,
		RoomID:  roomID,
		Content: content,
	})

	return nil
}

func (s *SocialService) GetChatMessages(ctx context.Context, roomID int64, page, size int32) (*social.ChatMessagesResponse, error) {
	messages, total, err := s.repo.GetChatMessages(ctx, roomID, int(page), int(size))
	if err != nil {
		return nil, err
	}

	// 转换为Thrift类型
	thriftMessages := make([]*social.ChatMessage, 0, len(messages))
	for _, msg := range messages {
		thriftMsg := &social.ChatMessage{
			Id:        msg.ID,
			RoomId:    msg.RoomID,
			SenderId:  msg.SenderID,
			Content:   msg.Content,
			Type:      msg.Type,
			CreatedAt: msg.CreatedAt.Unix(),
			UpdatedAt: msg.UpdatedAt.Unix(),
		}
		if msg.DeletedAt != nil {
			deletedAt := msg.DeletedAt.Unix()
			thriftMsg.DeletedAt = &deletedAt
		}
		thriftMessages = append(thriftMessages, thriftMsg)
	}

	return &social.ChatMessagesResponse{
		Messages: thriftMessages,
		Total:    total,
	}, nil
}

// 好友关系相关
func (s *SocialService) AddFriend(ctx context.Context, userID, friendID int64) error {
	// 检查是否已经是好友
	if friendship, err := s.repo.GetFriendship(ctx, userID, friendID); err == nil && friendship != nil {
		return errors.New("already friends")
	}

	friendship := &model.Friendship{
		UserID:   userID,
		FriendID: friendID,
		Status:   model.FriendStatusPending,
	}
	return s.repo.AddFriend(ctx, friendship)
}

func (s *SocialService) GetFriendship(ctx context.Context, userID, friendID int64) (*social.Friendship, error) {
	friendship, err := s.repo.GetFriendship(ctx, userID, friendID)
	if err != nil {
		return nil, err
	}

	// 转换为Thrift类型
	return convertFriendshipToThrift(friendship), nil
}

func (s *SocialService) GetUserFriends(ctx context.Context, userID int64) ([]*social.Friendship, error) {
	friendships, err := s.repo.GetUserFriends(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 转换为Thrift类型
	thriftFriendships := make([]*social.Friendship, 0, len(friendships))
	for i := range friendships {
		thriftFriendships = append(thriftFriendships, convertFriendshipToThrift(&friendships[i]))
	}

	return thriftFriendships, nil
}

// 好友申请相关
func (s *SocialService) CreateFriendRequest(ctx context.Context, senderID, receiverID int64, message string) error {
	request := &model.FriendRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Message:    message,
		Status:     model.FriendRequestStatusPending,
	}
	if err := s.repo.CreateFriendRequest(ctx, request); err != nil {
		return err
	}

	// 序列化消息
	wsMsg := &ws.Message{
		Type:    ws.MessageTypeFriendRequest,
		From:    senderID,
		To:      receiverID,
		Content: message,
	}
	data, err := json.Marshal(wsMsg)
	if err != nil {
		return err
	}

	// 发送序列化后的消息
	s.wsManager.SendToUser(receiverID, data)
	return nil
}

func (s *SocialService) GetFriendRequests(ctx context.Context, userID int64, status int8) ([]*social.FriendRequest, error) {
	requests, err := s.repo.GetFriendRequests(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	// 转换为Thrift类型
	thriftRequests := make([]*social.FriendRequest, 0, len(requests))
	for _, req := range requests {
		thriftReq := &social.FriendRequest{
			Id:         req.ID,
			SenderId:   req.SenderID,
			ReceiverId: req.ReceiverID,
			Status:     req.Status,
			CreatedAt:  req.CreatedAt.Unix(),
			UpdatedAt:  req.UpdatedAt.Unix(),
		}
		if req.Message != "" {
			thriftReq.Message = &req.Message
		}
		if req.DeletedAt != nil {
			deletedAt := req.DeletedAt.Unix()
			thriftReq.DeletedAt = &deletedAt
		}
		thriftRequests = append(thriftRequests, thriftReq)
	}

	return thriftRequests, nil
}

func (s *SocialService) HandleFriendRequest(ctx context.Context, requestID int64, status int8) error {
	return s.repo.UpdateFriendRequest(ctx, requestID, status)
}

// 消息已读状态相关
func (s *SocialService) MarkMessageRead(ctx context.Context, messageID, userID int64) error {
	return s.repo.MarkMessageRead(ctx, messageID, userID)
}

func (s *SocialService) GetUnreadMessageCount(ctx context.Context, userID int64) (int64, error) {
	return s.repo.GetUnreadMessageCount(ctx, userID)
}

// 辅助函数：转换聊天室对象为Thrift类型
func convertChatRoomToThrift(room *model.ChatRoom) *social.ChatRoom {
	thriftRoom := &social.ChatRoom{
		Id:        room.ID,
		Name:      room.Name,
		CreatorId: room.CreatorID,
		Type:      room.Type,
		CreatedAt: room.CreatedAt.Unix(),
		UpdatedAt: room.UpdatedAt.Unix(),
	}

	if room.DeletedAt != nil {
		deletedAt := room.DeletedAt.Unix()
		thriftRoom.DeletedAt = &deletedAt
	}

	if len(room.Members) > 0 {
		thriftMembers := make([]*social.ChatRoomMember, 0, len(room.Members))
		for _, member := range room.Members {
			thriftMember := &social.ChatRoomMember{
				Id:        member.ID,
				RoomId:    member.RoomID,
				UserId:    member.UserID,
				Role:      member.Role,
				CreatedAt: member.CreatedAt.Unix(),
				UpdatedAt: member.UpdatedAt.Unix(),
			}
			if member.Nickname != "" {
				thriftMember.Nickname = &member.Nickname
			}
			if member.DeletedAt != nil {
				deletedAt := member.DeletedAt.Unix()
				thriftMember.DeletedAt = &deletedAt
			}
			thriftMembers = append(thriftMembers, thriftMember)
		}
		thriftRoom.Members = thriftMembers
	}

	return thriftRoom
}

// 辅助函数：转换好友关系对象为Thrift类型
func convertFriendshipToThrift(friendship *model.Friendship) *social.Friendship {
	thriftFriendship := &social.Friendship{
		Id:        friendship.ID,
		UserId:    friendship.UserID,
		FriendId:  friendship.FriendID,
		Status:    friendship.Status,
		CreatedAt: friendship.CreatedAt.Unix(),
		UpdatedAt: friendship.UpdatedAt.Unix(),
	}

	if friendship.Remark != "" {
		thriftFriendship.Remark = &friendship.Remark
	}

	if friendship.DeletedAt != nil {
		deletedAt := friendship.DeletedAt.Unix()
		thriftFriendship.DeletedAt = &deletedAt
	}

	return thriftFriendship
}
