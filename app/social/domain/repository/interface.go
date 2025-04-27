package repository

import (
	"context"

	"github.com/yxrxy/videoHub/app/social/domain/model"
)

type SocialDB interface {
	// 私信相关
	SendPrivateMessage(ctx context.Context, msg *model.PrivateMessage) error
	GetPrivateMessages(ctx context.Context, senderID, receiverID int64, page, size int) ([]model.PrivateMessage, int64, error)

	// 聊天室相关
	CreateChatRoom(ctx context.Context, room *model.ChatRoom) (int64, error)
	GetChatRoom(ctx context.Context, roomID int64) (*model.ChatRoom, error)
	GetUserChatRooms(ctx context.Context, userID int64) ([]model.ChatRoom, error)

	// 聊天消息相关
	SendChatMessage(ctx context.Context, msg *model.ChatMessage) error
	GetChatMessages(ctx context.Context, roomID int64, page, size int) ([]model.ChatMessage, int64, error)

	// 好友关系相关
	AddFriend(ctx context.Context, friendship *model.Friendship) error
	GetFriendship(ctx context.Context, userID, friendID int64) (*model.Friendship, error)
	GetUserFriends(ctx context.Context, userID int64) ([]model.Friendship, error)

	// 好友申请相关
	CreateFriendRequest(ctx context.Context, request *model.FriendRequest) error
	GetFriendRequests(ctx context.Context, userID int64, status int8) ([]model.FriendRequest, error)
	UpdateFriendRequest(ctx context.Context, requestID int64, status int8) error

	// 消息已读状态相关
	MarkMessageRead(ctx context.Context, messageID, userID int64) error
	GetUnreadMessageCount(ctx context.Context, userID int64) (int64, error)
}

type SocialCache interface {
	// 在线状态相关
	SetUserOnline(ctx context.Context, userID int64) error
	SetUserOffline(ctx context.Context, userID int64) error
	IsUserOnline(ctx context.Context, userID int64) (bool, error)
}

type SocialWebSocket interface {
	Start(ctx context.Context)
	RegisterClient(userID int64) error
	JoinChatRoom(userID int64, roomID int64) error
	LeaveChatRoom(userID int64, roomID int64) error
	SendChatMessage(userID int64, roomID int64, content string, msgType int8) error
	SendPrivateMessage(fromID int64, toID int64, content string) error
	Broadcast(content string)
	GetOnlineUsers() []int64
	IsUserOnline(userID int64) bool
}
