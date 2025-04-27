package rpc

import (
	"context"

	"github.com/yxrxy/videoHub/app/social/controllers/rpc/pack"
	"github.com/yxrxy/videoHub/app/social/usecase"
	"github.com/yxrxy/videoHub/kitex_gen/social"
	"github.com/yxrxy/videoHub/pkg/base"
)

type SocialHandler struct {
	useCase usecase.SocialUseCase
}

func NewSocialHandler(useCase usecase.SocialUseCase) *SocialHandler {
	return &SocialHandler{useCase: useCase}
}

// 发送私信
func (h *SocialHandler) SendPrivateMessage(ctx context.Context, req *social.SendPrivateMessageRequest) (r *social.SendPrivateMessageResponse, err error) {
	r = new(social.SendPrivateMessageResponse)

	if err = h.useCase.SendPrivateMessage(ctx, req.SenderId, req.ReceiverId, req.Content); err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	return
}

// 获取私信列表
func (h *SocialHandler) GetPrivateMessages(ctx context.Context, req *social.GetPrivateMessagesRequest) (r *social.GetPrivateMessagesResponse, err error) {
	r = new(social.GetPrivateMessagesResponse)

	messages, err := h.useCase.GetPrivateMessages(ctx, req.UserId, req.OtherUserId, *req.Page, *req.Size)
	if err != nil {
		return
	}
	r.MessageList = pack.PackPrivateMessageList(messages)
	r.Base = base.BuildBaseResp(err)
	return
}

// 创建聊天室
func (h *SocialHandler) CreateChatRoom(ctx context.Context, req *social.CreateChatRoomRequest) (r *social.CreateChatRoomResponse, err error) {
	r = new(social.CreateChatRoomResponse)

	room, err := h.useCase.CreateChatRoom(ctx, req.Name, req.CreatorId, req.Type, req.MemberIds)
	if err != nil {
		return
	}
	r.Room = pack.PackChatRoom(room)
	r.Base = base.BuildBaseResp(err)
	return
}

// 获取聊天室信息
func (h *SocialHandler) GetChatRoom(ctx context.Context, req *social.GetChatRoomRequest) (r *social.GetChatRoomResponse, err error) {
	r = new(social.GetChatRoomResponse)

	room, err := h.useCase.GetChatRoom(ctx, req.RoomId)
	if err != nil {
		return
	}
	r.Room = pack.PackChatRoom(room)
	r.Base = base.BuildBaseResp(err)
	return
}

// 获取用户的聊天室列表
func (h *SocialHandler) GetUserChatRooms(ctx context.Context, req *social.GetUserChatRoomsRequest) (r *social.GetUserChatRoomsResponse, err error) {
	r = new(social.GetUserChatRoomsResponse)

	rooms, err := h.useCase.GetUserChatRooms(ctx, req.UserId)
	if err != nil {
		return
	}
	r.RoomList = pack.PackChatRoomList(rooms)
	r.Base = base.BuildBaseResp(err)
	return
}

// SendChatMessage 发送聊天消息
func (h *SocialHandler) SendChatMessage(ctx context.Context, req *social.SendChatMessageRequest) (r *social.SendChatMessageResponse, err error) {
	r = new(social.SendChatMessageResponse)

	if err = h.useCase.SendChatMessage(ctx, req.RoomId, req.SenderId, req.Content, *req.Type); err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	return
}

// 获取聊天消息列表
func (h *SocialHandler) GetChatMessages(ctx context.Context, req *social.GetChatMessagesRequest) (r *social.GetChatMessagesResponse, err error) {
	r = new(social.GetChatMessagesResponse)

	messages, err := h.useCase.GetChatMessages(ctx, req.RoomId, int32(*req.Page), *req.Size)
	if err != nil {
		return
	}
	r.MessageList = pack.PackChatMessageList(messages)
	r.Base = base.BuildBaseResp(err)
	return
}

// 添加好友
func (h *SocialHandler) AddFriend(ctx context.Context, req *social.AddFriendRequest) (r *social.AddFriendResponse, err error) {
	r = new(social.AddFriendResponse)

	if err = h.useCase.AddFriend(ctx, req.UserId, req.FriendId); err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	return
}

// 获取好友列表
func (h *SocialHandler) GetUserFriends(ctx context.Context, req *social.GetUserFriendsRequest) (r *social.GetUserFriendsResponse, err error) {
	r = new(social.GetUserFriendsResponse)

	friends, err := h.useCase.GetUserFriends(ctx, req.UserId)
	if err != nil {
		return
	}
	r.FriendList = pack.PackFriendList(friends)
	r.Base = base.BuildBaseResp(err)
	return
}

// 获取好友申请列表
func (h *SocialHandler) GetFriendRequests(ctx context.Context, req *social.GetFriendRequestsRequest) (r *social.GetFriendRequestsResponse, err error) {
	r = new(social.GetFriendRequestsResponse)

	requests, err := h.useCase.GetFriendRequests(ctx, req.UserId, *req.Type)
	if err != nil {
		return
	}
	r.RequestList = pack.PackFriendRequestList(requests)
	r.Base = base.BuildBaseResp(err)
	return
}

// 处理好友申请
func (h *SocialHandler) HandleFriendRequest(ctx context.Context, req *social.HandleFriendRequestRequest) (r *social.HandleFriendRequestResponse, err error) {
	r = new(social.HandleFriendRequestResponse)

	if err = h.useCase.HandleFriendRequest(ctx, req.RequestId, req.Action); err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	return
}

// 获取未读消息数量
func (h *SocialHandler) GetUnreadMessageCount(
	ctx context.Context,
	req *social.GetUnreadMessageCountRequest,
) (r *social.GetUnreadMessageCountResponse, err error) {
	r = new(social.GetUnreadMessageCountResponse)

	count, err := h.useCase.GetUnreadMessageCount(ctx, req.UserId)
	if err != nil {
		return
	}
	r.Count = count
	r.Base = base.BuildBaseResp(err)
	return
}

// 设置消息已读
func (h *SocialHandler) MarkMessageRead(ctx context.Context, req *social.MarkMessageReadRequest) (r *social.MarkMessageReadResponse, err error) {
	r = new(social.MarkMessageReadResponse)

	if err = h.useCase.MarkMessageRead(ctx, req.MessageId, req.UserId); err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	return
}

// 创建好友申请
func (h *SocialHandler) CreateFriendRequest(ctx context.Context, req *social.CreateFriendRequestRequest) (r *social.CreateFriendRequestResponse, err error) {
	r = new(social.CreateFriendRequestResponse)

	message := ""
	if req.Message != nil {
		message = *req.Message
	}

	if err = h.useCase.CreateFriendRequest(ctx, req.SenderId, req.ReceiverId, message); err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	return
}

// 获取好友申请
func (h *SocialHandler) GetFriendRequest(ctx context.Context, req *social.GetFriendRequestsRequest) (r *social.GetFriendRequestsResponse, err error) {
	r = new(social.GetFriendRequestsResponse)

	requests, err := h.useCase.GetFriendRequests(ctx, req.UserId, *req.Type)
	if err != nil {
		return
	}
	r.RequestList = pack.PackFriendRequestList(requests)
	r.Base = base.BuildBaseResp(err)
	return
}

func (h *SocialHandler) GetFriendship(ctx context.Context, req *social.GetFriendshipRequest) (r *social.GetFriendshipResponse, err error) {
	r = new(social.GetFriendshipResponse)

	friendship, err := h.useCase.GetFriendship(ctx, req.UserId, req.FriendId)
	if err != nil {
		return
	}
	r.Friend = pack.PackFriendship(friendship)
	r.Base = base.BuildBaseResp(err)
	return
}

func (h *SocialHandler) JoinChatRoom(ctx context.Context, req *social.JoinChatRoomRequest) (r *social.JoinChatRoomResponse, err error) {
	r = new(social.JoinChatRoomResponse)

	if err := h.useCase.JoinChatRoom(req.UserId, req.RoomId); err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, nil
	}

	r.Base = base.BuildBaseResp(nil)
	return r, nil
}

func (h *SocialHandler) LeaveChatRoom(ctx context.Context, req *social.LeaveChatRoomRequest) (r *social.LeaveChatRoomResponse, err error) {
	r = new(social.LeaveChatRoomResponse)

	if err := h.useCase.LeaveChatRoom(req.UserId, req.RoomId); err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, nil
	}

	r.Base = base.BuildBaseResp(nil)
	return r, nil
}

func (h *SocialHandler) RegisterWebSocketClient(
	ctx context.Context,
	req *social.RegisterWebSocketClientRequest,
) (r *social.RegisterWebSocketClientResponse, err error) {
	r = new(social.RegisterWebSocketClientResponse)
	if err := h.useCase.RegisterWebSocketClient(req.UserId); err != nil {
		r.Base = base.BuildBaseResp(err)
		return r, nil
	}
	r.Base = base.BuildBaseResp(nil)
	return r, nil
}
