package repository

import (
	"context"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/yxrrxy/videoHub/app/social/model"
	"github.com/yxrrxy/videoHub/config"
)

type Social struct {
	db *gorm.DB
}

func NewSocial(db *gorm.DB) *Social {
	return &Social{db: db}
}

// 私信相关
func (s *Social) SendPrivateMessage(ctx context.Context, msg *model.PrivateMessage) error {
	return s.db.WithContext(ctx).Create(msg).Error
}

func (s *Social) GetPrivateMessages(ctx context.Context, senderID, receiverID int64, page, size int) ([]model.PrivateMessage, int64, error) {
	var messages []model.PrivateMessage
	var total int64

	query := s.db.WithContext(ctx).Model(&model.PrivateMessage{}).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			senderID, receiverID, receiverID, senderID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(size).
		Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// 聊天室相关
func (s *Social) CreateChatRoom(ctx context.Context, room *model.ChatRoom) error {
	return s.db.WithContext(ctx).Create(room).Error
}

func (s *Social) GetChatRoom(ctx context.Context, roomID int64) (*model.ChatRoom, error) {
	var room model.ChatRoom
	if err := s.db.WithContext(ctx).
		Preload("Members").
		First(&room, roomID).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *Social) GetUserChatRooms(ctx context.Context, userID int64) ([]model.ChatRoom, error) {
	var rooms []model.ChatRoom
	if err := s.db.WithContext(ctx).
		Joins("JOIN chat_room_members ON chat_rooms.id = chat_room_members.room_id").
		Where("chat_room_members.user_id = ?", userID).
		Preload("Members").
		Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

// 聊天消息相关
func (s *Social) SendChatMessage(ctx context.Context, msg *model.ChatMessage) error {
	return s.db.WithContext(ctx).Create(msg).Error
}

func (s *Social) GetChatMessages(ctx context.Context, roomID int64, page, size int) ([]model.ChatMessage, int64, error) {
	var messages []model.ChatMessage
	var total int64

	query := s.db.WithContext(ctx).Model(&model.ChatMessage{}).
		Where("room_id = ?", roomID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	if err := query.Order("created_at DESC").
		Offset(offset).Limit(size).
		Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// 好友关系相关
func (s *Social) AddFriend(ctx context.Context, friendship *model.Friendship) error {
	return s.db.WithContext(ctx).Create(friendship).Error
}

func (s *Social) GetFriendship(ctx context.Context, userID, friendID int64) (*model.Friendship, error) {
	var friendship model.Friendship
	if err := s.db.WithContext(ctx).
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
			userID, friendID, friendID, userID).
		First(&friendship).Error; err != nil {
		return nil, err
	}
	return &friendship, nil
}

func (s *Social) GetUserFriends(ctx context.Context, userID int64) ([]model.Friendship, error) {
	var friendships []model.Friendship
	if err := s.db.WithContext(ctx).
		Where("(user_id = ? OR friend_id = ?) AND status = ?",
			userID, userID, model.FriendStatusAccepted).
		Find(&friendships).Error; err != nil {
		return nil, err
	}
	return friendships, nil
}

// 好友申请相关
func (s *Social) CreateFriendRequest(ctx context.Context, request *model.FriendRequest) error {
	return s.db.WithContext(ctx).Create(request).Error
}

func (s *Social) GetFriendRequests(ctx context.Context, userID int64, status int8) ([]model.FriendRequest, error) {
	var requests []model.FriendRequest
	query := s.db.WithContext(ctx).Where("receiver_id = ?", userID)
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	if err := query.Order("created_at DESC").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

func (s *Social) UpdateFriendRequest(ctx context.Context, requestID int64, status int8) error {
	return s.db.WithContext(ctx).Model(&model.FriendRequest{}).
		Where("id = ?", requestID).
		Update("status", status).Error
}

// 消息已读状态相关
func (s *Social) MarkMessageRead(ctx context.Context, messageID, userID int64) error {
	now := time.Now()
	return s.db.WithContext(ctx).Create(&model.MessageReadStatus{
		MessageID: messageID,
		UserID:    userID,
		IsRead:    true,
		ReadAt:    &now,
	}).Error
}

func (s *Social) GetUnreadMessageCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.PrivateMessage{}).
		Where("receiver_id = ? AND is_read = ?", userID, false).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// 初始化数据库
func InitDB() *gorm.DB {
	dsn := config.GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 自动迁移
	if err := db.AutoMigrate(
		&model.PrivateMessage{},
		&model.ChatRoom{},
		&model.ChatRoomMember{},
		&model.ChatMessage{},
		&model.Friendship{},
		&model.FriendRequest{},
		&model.MessageReadStatus{},
	); err != nil {
		panic(err)
	}

	return db
}
