package rpc

import (
	"context"
	"log"

	"github.com/yxrxy/videoHub/kitex_gen/model"
	"github.com/yxrxy/videoHub/kitex_gen/social"
	"github.com/yxrxy/videoHub/pkg/base/client"
	"github.com/yxrxy/videoHub/pkg/errno"
)

// InitSocialRPC 初始化社交服务客户端
func InitSocialRPC() {
	c, err := client.InitSocialRPC()
	if err != nil {
		log.Fatalf("初始化社交服务客户端失败: %v", err)
	}
	socialClient = *c
}

// SendPrivateMessageRPC 发送私信
func SendPrivateMessageRPC(ctx context.Context, req *social.SendPrivateMessageRequest) (*model.PrivateMessage, error) {
	resp, err := socialClient.SendPrivateMessage(ctx, req)
	if err != nil {
		log.Printf("发送私信RPC调用失败: %v", err)
		return nil, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Message, nil
}

// GetPrivateMessagesRPC 获取私信列表
func GetPrivateMessagesRPC(ctx context.Context, req *social.GetPrivateMessagesRequest) ([]*model.PrivateMessage, int64, error) {
	resp, err := socialClient.GetPrivateMessages(ctx, req)
	if err != nil {
		log.Printf("获取私信列表RPC调用失败: %v", err)
		return nil, 0, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, 0, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.MessageList, resp.Total, nil
}

// CreateChatRoomRPC 创建聊天室
func CreateChatRoomRPC(ctx context.Context, req *social.CreateChatRoomRequest) (*model.ChatRoom, error) {
	resp, err := socialClient.CreateChatRoom(ctx, req)
	if err != nil {
		log.Printf("创建聊天室RPC调用失败: %v", err)
		return nil, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Room, nil
}

// GetUserChatRoomsRPC 获取用户聊天室列表
func GetUserChatRoomsRPC(ctx context.Context, req *social.GetUserChatRoomsRequest) ([]*model.ChatRoom, int64, error) {
	resp, err := socialClient.GetUserChatRooms(ctx, req)
	if err != nil {
		log.Printf("获取用户聊天室列表RPC调用失败: %v", err)
		return nil, 0, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, 0, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.RoomList, resp.Total, nil
}

// GetChatRoomRPC 获取聊天室信息
func GetChatRoomRPC(ctx context.Context, req *social.GetChatRoomRequest) (*model.ChatRoom, error) {
	resp, err := socialClient.GetChatRoom(ctx, req)
	if err != nil {
		log.Printf("获取聊天室信息RPC调用失败: %v", err)
		return nil, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Room, nil
}

// SendChatMessageRPC 发送聊天室消息
func SendChatMessageRPC(ctx context.Context, req *social.SendChatMessageRequest) (*model.ChatMessage, error) {
	resp, err := socialClient.SendChatMessage(ctx, req)
	if err != nil {
		log.Printf("发送聊天室消息RPC调用失败: %v", err)
		return nil, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Message, nil
}

// GetChatMessagesRPC 获取聊天室消息列表
func GetChatMessagesRPC(ctx context.Context, req *social.GetChatMessagesRequest) ([]*model.ChatMessage, int64, error) {
	resp, err := socialClient.GetChatMessages(ctx, req)
	if err != nil {
		log.Printf("获取聊天室消息列表RPC调用失败: %v", err)
		return nil, 0, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, 0, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.MessageList, resp.Total, nil
}

// AddFriendRPC 添加好友
func AddFriendRPC(ctx context.Context, req *social.AddFriendRequest) error {
	resp, err := socialClient.AddFriend(ctx, req)
	if err != nil {
		log.Printf("添加好友RPC调用失败: %v", err)
		return errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return nil
}

// GetFriendshipRPC 获取好友关系
func GetFriendshipRPC(ctx context.Context, req *social.GetFriendshipRequest) (*model.Friendship, error) {
	resp, err := socialClient.GetFriendship(ctx, req)
	if err != nil {
		log.Printf("获取好友关系RPC调用失败: %v", err)
		return nil, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Friend, nil
}

// GetUserFriendsRPC 获取用户好友列表
func GetUserFriendsRPC(ctx context.Context, req *social.GetUserFriendsRequest) ([]*model.Friendship, int64, error) {
	resp, err := socialClient.GetUserFriends(ctx, req)
	if err != nil {
		log.Printf("获取用户好友列表RPC调用失败: %v", err)
		return nil, 0, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, 0, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.FriendList, resp.Total, nil
}

// CreateFriendRequestRPC 创建好友请求
func CreateFriendRequestRPC(ctx context.Context, req *social.CreateFriendRequestRequest) (*model.FriendRequest, error) {
	resp, err := socialClient.CreateFriendRequest(ctx, req)
	if err != nil {
		log.Printf("创建好友请求RPC调用失败: %v", err)
		return nil, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Request, nil
}

// GetFriendRequestsRPC 获取好友请求列表
func GetFriendRequestsRPC(ctx context.Context, req *social.GetFriendRequestsRequest) ([]*model.FriendRequest, int64, error) {
	resp, err := socialClient.GetFriendRequests(ctx, req)
	if err != nil {
		log.Printf("获取好友请求列表RPC调用失败: %v", err)
		return nil, 0, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return nil, 0, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.RequestList, resp.Total, nil
}

// HandleFriendRequestRPC 处理好友请求
func HandleFriendRequestRPC(ctx context.Context, req *social.HandleFriendRequestRequest) error {
	resp, err := socialClient.HandleFriendRequest(ctx, req)
	if err != nil {
		log.Printf("处理好友请求RPC调用失败: %v", err)
		return errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return nil
}

// MarkMessageReadRPC 标记消息已读
func MarkMessageReadRPC(ctx context.Context, req *social.MarkMessageReadRequest) error {
	resp, err := socialClient.MarkMessageRead(ctx, req)
	if err != nil {
		log.Printf("标记消息已读RPC调用失败: %v", err)
		return errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return nil
}

// GetUnreadMessageCountRPC 获取未读消息数量
func GetUnreadMessageCountRPC(ctx context.Context, req *social.GetUnreadMessageCountRequest) (int64, error) {
	resp, err := socialClient.GetUnreadMessageCount(ctx, req)
	if err != nil {
		log.Printf("获取未读消息数量RPC调用失败: %v", err)
		return 0, errno.InternalServiceError.WithError(err)
	}
	if resp.Base.Code != errno.SuccessCode {
		return 0, errno.InternalServiceError.WithMessage(resp.Base.Msg)
	}
	return resp.Count, nil
}
