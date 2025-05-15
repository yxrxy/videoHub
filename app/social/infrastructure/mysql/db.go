package mysql

import (
	"context"

	"github.com/yxrxy/videoHub/app/social/domain/model"
	"github.com/yxrxy/videoHub/app/social/domain/repository"
	"gorm.io/gorm"
)

type SocialDB struct {
	db *gorm.DB
}

func NewSocialDB(db *gorm.DB) repository.SocialDB {
	return &SocialDB{db: db}
}

// 私信相关
func (s *SocialDB) SendPrivateMessage(ctx context.Context, msg *model.PrivateMessage) error {
	dbMsg := &PrivateMessage{
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		IsRead:     msg.IsRead,
	}
	return s.db.WithContext(ctx).Create(dbMsg).Error
}

func (s *SocialDB) GetPrivateMessages(ctx context.Context, senderID, receiverID int64, page, size int) ([]model.PrivateMessage, int64, error) {
	var dbMessages []PrivateMessage
	var total int64

	query := s.db.WithContext(ctx).Model(&PrivateMessage{}).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			senderID, receiverID, receiverID, senderID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(size).
		Find(&dbMessages).Error; err != nil {
		return nil, 0, err
	}

	messages := make([]model.PrivateMessage, len(dbMessages))
	for i, dbMsg := range dbMessages {
		messages[i] = model.PrivateMessage{
			ID:         dbMsg.ID,
			SenderID:   dbMsg.SenderID,
			ReceiverID: dbMsg.ReceiverID,
			Content:    dbMsg.Content,
			IsRead:     dbMsg.IsRead,
		}
	}

	return messages, total, nil
}

// 聊天室相关
func (s *SocialDB) CreateChatRoom(ctx context.Context, room *model.ChatRoom) (int64, error) {
	dbRoom := &ChatRoom{
		Name:      room.Name,
		CreatorID: room.CreatorID,
		Type:      room.Type,
	}
	err := s.db.WithContext(ctx).Create(dbRoom).Error
	if err != nil {
		return 0, err
	}
	return dbRoom.ID, nil
}

func (s *SocialDB) GetChatRoom(ctx context.Context, roomID int64) (*model.ChatRoom, error) {
	var dbRoom ChatRoom
	if err := s.db.WithContext(ctx).
		Preload("Members").
		First(&dbRoom, roomID).Error; err != nil {
		return nil, err
	}

	members := make([]model.ChatRoomMember, len(dbRoom.Members))
	for i, dbMember := range dbRoom.Members {
		members[i] = model.ChatRoomMember{
			ID:       dbMember.ID,
			RoomID:   dbMember.RoomID,
			UserID:   dbMember.UserID,
			Nickname: dbMember.Nickname,
			Role:     dbMember.Role,
		}
	}

	return &model.ChatRoom{
		ID:        dbRoom.ID,
		Name:      dbRoom.Name,
		CreatorID: dbRoom.CreatorID,
		Type:      dbRoom.Type,
		Members:   members,
	}, nil
}

func (s *SocialDB) GetUserChatRooms(ctx context.Context, userID int64) ([]model.ChatRoom, error) {
	var dbRooms []ChatRoom
	if err := s.db.WithContext(ctx).
		Where("NOT (chat_rooms.id = -1 AND chat_rooms.creator_id != ?)", userID).
		Preload("Members").
		Find(&dbRooms).Error; err != nil {
		return nil, err
	}

	rooms := make([]model.ChatRoom, len(dbRooms))
	for i, dbRoom := range dbRooms {
		members := make([]model.ChatRoomMember, len(dbRoom.Members))
		for j, dbMember := range dbRoom.Members {
			members[j] = model.ChatRoomMember{
				ID:       dbMember.ID,
				RoomID:   dbMember.RoomID,
				UserID:   dbMember.UserID,
				Nickname: dbMember.Nickname,
				Role:     dbMember.Role,
			}
		}

		rooms[i] = model.ChatRoom{
			ID:        dbRoom.ID,
			Name:      dbRoom.Name,
			CreatorID: dbRoom.CreatorID,
			Type:      dbRoom.Type,
			Members:   members,
		}
	}

	return rooms, nil
}

// 聊天消息相关
func (s *SocialDB) SendChatMessage(ctx context.Context, msg *model.ChatMessage) error {
	dbMsg := &ChatMessage{
		RoomID:   msg.RoomID,
		SenderID: msg.SenderID,
		Content:  msg.Content,
		Type:     msg.Type,
	}
	return s.db.WithContext(ctx).Create(dbMsg).Error
}

func (s *SocialDB) GetChatMessages(ctx context.Context, roomID int64, page, size int) ([]model.ChatMessage, int64, error) {
	var dbMessages []ChatMessage
	var total int64

	query := s.db.WithContext(ctx).Model(&ChatMessage{}).
		Where("room_id = ?", roomID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(size).
		Find(&dbMessages).Error; err != nil {
		return nil, 0, err
	}

	messages := make([]model.ChatMessage, len(dbMessages))
	for i, dbMsg := range dbMessages {
		messages[i] = model.ChatMessage{
			ID:       dbMsg.ID,
			RoomID:   dbMsg.RoomID,
			SenderID: dbMsg.SenderID,
			Content:  dbMsg.Content,
			Type:     dbMsg.Type,
		}
	}

	return messages, total, nil
}

// 好友关系相关
func (s *SocialDB) AddFriend(ctx context.Context, friendship *model.Friendship) error {
	dbFriendship := &Friendship{
		UserID:   friendship.UserID,
		FriendID: friendship.FriendID,
		Status:   friendship.Status,
		Remark:   friendship.Remark,
	}
	return s.db.WithContext(ctx).Create(dbFriendship).Error
}

func (s *SocialDB) GetFriendship(ctx context.Context, userID, friendID int64) (*model.Friendship, error) {
	var dbFriendship Friendship
	if err := s.db.WithContext(ctx).
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
			userID, friendID, friendID, userID).
		First(&dbFriendship).Error; err != nil {
		return nil, err
	}

	return &model.Friendship{
		ID:       dbFriendship.ID,
		UserID:   dbFriendship.UserID,
		FriendID: dbFriendship.FriendID,
		Status:   dbFriendship.Status,
		Remark:   dbFriendship.Remark,
	}, nil
}

func (s *SocialDB) GetUserFriends(ctx context.Context, userID int64) ([]model.Friendship, error) {
	var dbFriendships []Friendship
	if err := s.db.WithContext(ctx).
		Where("(user_id = ? OR friend_id = ?)",
			userID, userID).
		Find(&dbFriendships).Error; err != nil {
		return nil, err
	}

	friendships := make([]model.Friendship, len(dbFriendships))
	for i, dbFriendship := range dbFriendships {
		friendships[i] = model.Friendship{
			ID:       dbFriendship.ID,
			UserID:   dbFriendship.UserID,
			FriendID: dbFriendship.FriendID,
			Status:   dbFriendship.Status,
			Remark:   dbFriendship.Remark,
		}
	}

	return friendships, nil
}

// 好友申请相关
func (s *SocialDB) CreateFriendRequest(ctx context.Context, request *model.FriendRequest) error {
	dbRequest := &FriendRequest{
		SenderID:   request.SenderID,
		ReceiverID: request.ReceiverID,
		Message:    request.Message,
		Status:     request.Status,
	}
	return s.db.WithContext(ctx).Create(dbRequest).Error
}

func (s *SocialDB) GetFriendRequests(ctx context.Context, userID int64, status int8) ([]model.FriendRequest, error) {
	var dbRequests []FriendRequest
	query := s.db.WithContext(ctx).Where("receiver_id = ?", userID)
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	if err := query.Order("created_at DESC").Find(&dbRequests).Error; err != nil {
		return nil, err
	}

	requests := make([]model.FriendRequest, len(dbRequests))
	for i, dbRequest := range dbRequests {
		requests[i] = model.FriendRequest{
			ID:         dbRequest.ID,
			SenderID:   dbRequest.SenderID,
			ReceiverID: dbRequest.ReceiverID,
			Message:    dbRequest.Message,
			Status:     dbRequest.Status,
		}
	}

	return requests, nil
}

func (s *SocialDB) UpdateFriendRequest(ctx context.Context, requestID int64, status int8) error {
	return s.db.WithContext(ctx).Model(&FriendRequest{}).
		Where("id = ?", requestID).
		Update("status", status).Error
}

// 消息已读状态相关
func (s *SocialDB) MarkMessageRead(ctx context.Context, messageID, userID int64) error {
	dbStatus := &MessageReadStatus{
		MessageID: messageID,
		UserID:    userID,
		IsRead:    true,
	}
	return s.db.WithContext(ctx).Create(dbStatus).Error
}

func (s *SocialDB) GetUnreadMessageCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	if err := s.db.WithContext(ctx).Model(&PrivateMessage{}).
		Where("receiver_id = ? AND is_read = ?", userID, false).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
