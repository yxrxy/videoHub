package usecase

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/yxrxy/videoHub/app/social/domain/model"
	"github.com/yxrxy/videoHub/app/social/infrastructure/ws"
)

// RegisterWebSocketClient 注册WebSocket客户端
func (s *useCase) RegisterWebSocketClient(userID int64) error {
	return s.wsService.RegisterClient(userID)
}

// JoinChatRoom 加入聊天室
func (s *useCase) JoinChatRoom(userID int64, roomID int64) error {
	return s.wsService.JoinChatRoom(userID, roomID)
}

// LeaveChatRoom 离开聊天室
func (s *useCase) LeaveChatRoom(userID int64, roomID int64) error {
	return s.wsService.LeaveChatRoom(userID, roomID)
}

// GetOnlineUsers 获取在线用户列表
func (s *useCase) GetOnlineUsers() []int64 {
	return s.wsService.GetOnlineUsers()
}

// IsUserOnline 检查用户是否在线
func (s *useCase) IsUserOnline(userID int64) bool {
	return s.wsService.IsUserOnline(userID)
}

// BroadcastSystemMessage 广播系统消息
func (s *useCase) BroadcastSystemMessage(content string) {
	s.wsService.Broadcast(content)
}

// SendPrivateMessage 私信相关
func (s *useCase) SendPrivateMessage(ctx context.Context, senderID, receiverID int64, content string) error {
	if err := s.svc.SavePrivateMessage(ctx, senderID, receiverID, content); err != nil {
		return err
	}
	return nil
}

// GetPrivateMessages 获取私信列表
func (s *useCase) GetPrivateMessages(ctx context.Context, senderID, receiverID int64, page, size int32) ([]*model.PrivateMessage, error) {
	messages, _, err := s.db.GetPrivateMessages(ctx, senderID, receiverID, int(page), int(size))
	if err != nil {
		return nil, err
	}

	// 转换为domain类型
	domainMessages := make([]*model.PrivateMessage, 0, len(messages))
	for _, msg := range messages {
		domainMsg := &model.PrivateMessage{
			ID:         msg.ID,
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Content:    msg.Content,
			IsRead:     msg.IsRead,
		}
		domainMessages = append(domainMessages, domainMsg)
	}

	return domainMessages, nil
}

// CreateChatRoom 聊天室相关
func (s *useCase) CreateChatRoom(ctx context.Context, name string, creatorID int64, roomType int8, memberIDs []int64) (*model.ChatRoom, error) {
	room := &model.ChatRoom{
		Name:      name,
		CreatorID: creatorID,
		Type:      roomType,
	}

	members, err := s.svc.CreateChatRoom(ctx, name, creatorID, memberIDs)
	if err != nil {
		return nil, err
	}

	room.Members = members
	id, err := s.db.CreateChatRoom(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = id

	// 转换为domain类型
	domainRoom := &model.ChatRoom{
		ID:        room.ID,
		Name:      room.Name,
		CreatorID: room.CreatorID,
		Type:      room.Type,
	}

	return domainRoom, nil
}

func (s *useCase) GetChatRoom(ctx context.Context, roomID int64) (*model.ChatRoom, error) {
	room, err := s.db.GetChatRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}

	// 转换为domain类型
	domainRoom := &model.ChatRoom{
		ID:        room.ID,
		Name:      room.Name,
		CreatorID: room.CreatorID,
		Type:      room.Type,
	}

	return domainRoom, nil
}

func (s *useCase) GetUserChatRooms(ctx context.Context, userID int64) ([]*model.ChatRoom, error) {
	rooms, err := s.db.GetUserChatRooms(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 转换为domain类型
	domainRooms := make([]*model.ChatRoom, 0, len(rooms))
	for _, room := range rooms {
		domainRoom := &model.ChatRoom{
			ID:        room.ID,
			Name:      room.Name,
			CreatorID: room.CreatorID,
			Type:      room.Type,
		}
		domainRooms = append(domainRooms, domainRoom)
	}

	return domainRooms, nil
}

// 聊天消息相关
func (s *useCase) SendChatMessage(ctx context.Context, roomID, senderID int64, content string, msgType int8) error {
	msg := &model.ChatMessage{
		RoomID:   roomID,
		SenderID: senderID,
		Content:  content,
		Type:     msgType,
	}
	if err := s.db.SendChatMessage(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (s *useCase) GetChatMessages(ctx context.Context, roomID int64, page, size int32) ([]*model.ChatMessage, error) {
	messages, _, err := s.db.GetChatMessages(ctx, roomID, int(page), int(size))
	if err != nil {
		return nil, err
	}

	// 转换为domain类型
	domainMessages := make([]*model.ChatMessage, 0, len(messages))
	for _, msg := range messages {
		domainMsg := &model.ChatMessage{
			ID:       msg.ID,
			RoomID:   msg.RoomID,
			SenderID: msg.SenderID,
			Content:  msg.Content,
			Type:     msg.Type,
		}
		domainMessages = append(domainMessages, domainMsg)
	}

	return domainMessages, nil
}

// 好友关系相关
func (s *useCase) AddFriend(ctx context.Context, userID, friendID int64) error {
	// 检查是否已经是好友
	if friendship, err := s.db.GetFriendship(ctx, userID, friendID); err == nil && friendship != nil {
		return errors.New("already friends")
	}

	friendship := &model.Friendship{
		UserID:   userID,
		FriendID: friendID,
		Status:   model.FriendStatusPending,
	}
	return s.db.AddFriend(ctx, friendship)
}

func (s *useCase) GetFriendship(ctx context.Context, userID, friendID int64) (*model.Friendship, error) {
	friendship, err := s.db.GetFriendship(ctx, userID, friendID)
	if err != nil {
		return nil, err
	}

	// 转换为domain类型
	domainFriendship := &model.Friendship{
		ID:       friendship.ID,
		UserID:   friendship.UserID,
		FriendID: friendship.FriendID,
		Status:   friendship.Status,
	}

	return domainFriendship, nil
}

func (s *useCase) GetUserFriends(ctx context.Context, userID int64) ([]*model.Friendship, error) {
	friendships, err := s.db.GetUserFriends(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 转换为domain类型
	domainFriendships := make([]*model.Friendship, 0, len(friendships))
	for _, friendship := range friendships {
		domainFriendship := &model.Friendship{
			ID:       friendship.ID,
			UserID:   friendship.UserID,
			FriendID: friendship.FriendID,
			Status:   friendship.Status,
		}
		domainFriendships = append(domainFriendships, domainFriendship)
	}

	return domainFriendships, nil
}

// 好友申请相关
func (s *useCase) CreateFriendRequest(ctx context.Context, senderID, receiverID int64, message string) error {
	request := &model.FriendRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Message:    message,
		Status:     model.FriendRequestStatusPending,
	}
	if err := s.db.CreateFriendRequest(ctx, request); err != nil {
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
	err = s.wsService.SendPrivateMessage(senderID, receiverID, string(data))
	if err != nil {
		return err
	}
	return nil
}

func (s *useCase) GetFriendRequests(ctx context.Context, userID int64, status int8) ([]*model.FriendRequest, error) {
	requests, err := s.db.GetFriendRequests(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	// 转换为domain类型
	domainRequests := make([]*model.FriendRequest, 0, len(requests))
	for _, request := range requests {
		domainRequest := &model.FriendRequest{
			ID:         request.ID,
			SenderID:   request.SenderID,
			ReceiverID: request.ReceiverID,
			Message:    request.Message,
			Status:     request.Status,
		}
		domainRequests = append(domainRequests, domainRequest)
	}

	return domainRequests, nil
}

func (s *useCase) HandleFriendRequest(ctx context.Context, requestID int64, status int8) error {
	return s.db.UpdateFriendRequest(ctx, requestID, status)
}

// MarkMessageRead 消息已读状态相关
func (s *useCase) MarkMessageRead(ctx context.Context, messageID, userID int64) error {
	return s.db.MarkMessageRead(ctx, messageID, userID)
}

func (s *useCase) GetUnreadMessageCount(ctx context.Context, userID int64) (int64, error) {
	return s.db.GetUnreadMessageCount(ctx, userID)
}
