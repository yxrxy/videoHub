package usecase

import (
	"context"

	"github.com/yxrxy/videoHub/app/social/domain/model"
	"github.com/yxrxy/videoHub/app/social/domain/repository"
	"github.com/yxrxy/videoHub/app/social/domain/service"
)

type SocialUseCase interface {
	// WebSocket相关
	RegisterWebSocketClient(userID int64) error
	JoinChatRoom(userID int64, roomID int64) error
	LeaveChatRoom(userID int64, roomID int64) error
	GetOnlineUsers() []int64
	IsUserOnline(userID int64) bool
	BroadcastSystemMessage(content string)

	// 私信相关
	SendPrivateMessage(ctx context.Context, senderID, receiverID int64, content string) error
	GetPrivateMessages(ctx context.Context, senderID, receiverID int64, page, size int32) ([]*model.PrivateMessage, error)

	// 聊天室相关
	CreateChatRoom(ctx context.Context, name string, creatorID int64, roomType int8, memberIDs []int64) (*model.ChatRoom, error)
	GetChatRoom(ctx context.Context, roomID int64) (*model.ChatRoom, error)
	GetUserChatRooms(ctx context.Context, userID int64) ([]*model.ChatRoom, error)

	// 聊天消息相关
	SendChatMessage(ctx context.Context, roomID, senderID int64, content string, msgType int8) error
	GetChatMessages(ctx context.Context, roomID int64, page, size int32) ([]*model.ChatMessage, error)

	// 好友关系相关
	AddFriend(ctx context.Context, userID, friendID int64) error
	GetFriendship(ctx context.Context, userID, friendID int64) (*model.Friendship, error)
	GetUserFriends(ctx context.Context, userID int64) ([]*model.Friendship, error)

	// 好友申请相关
	CreateFriendRequest(ctx context.Context, senderID, receiverID int64, message string) error
	GetFriendRequests(ctx context.Context, userID int64, status int8) ([]*model.FriendRequest, error)
	HandleFriendRequest(ctx context.Context, requestID int64, status int8) error

	// 消息已读状态相关
	MarkMessageRead(ctx context.Context, messageID, userID int64) error
	GetUnreadMessageCount(ctx context.Context, userID int64) (int64, error)
}

type useCase struct {
	db        repository.SocialDB
	cache     repository.SocialCache
	svc       *service.SocialService
	wsService repository.SocialWebSocket
}

func NewSocialCase(db repository.SocialDB, cache repository.SocialCache, svc *service.SocialService, wsService repository.SocialWebSocket) *useCase {
	return &useCase{
		db:        db,
		cache:     cache,
		svc:       svc,
		wsService: wsService,
	}
}
