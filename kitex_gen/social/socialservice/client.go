// Code generated by Kitex v0.13.1. DO NOT EDIT.

package socialservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	social "github.com/yxrxy/videoHub/kitex_gen/social"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	SendPrivateMessage(ctx context.Context, req *social.SendPrivateMessageRequest, callOptions ...callopt.Option) (r *social.SendPrivateMessageResponse, err error)
	GetPrivateMessages(ctx context.Context, req *social.GetPrivateMessagesRequest, callOptions ...callopt.Option) (r *social.GetPrivateMessagesResponse, err error)
	CreateChatRoom(ctx context.Context, req *social.CreateChatRoomRequest, callOptions ...callopt.Option) (r *social.CreateChatRoomResponse, err error)
	GetChatRoom(ctx context.Context, req *social.GetChatRoomRequest, callOptions ...callopt.Option) (r *social.GetChatRoomResponse, err error)
	GetUserChatRooms(ctx context.Context, req *social.GetUserChatRoomsRequest, callOptions ...callopt.Option) (r *social.GetUserChatRoomsResponse, err error)
	SendChatMessage(ctx context.Context, req *social.SendChatMessageRequest, callOptions ...callopt.Option) (r *social.SendChatMessageResponse, err error)
	GetChatMessages(ctx context.Context, req *social.GetChatMessagesRequest, callOptions ...callopt.Option) (r *social.GetChatMessagesResponse, err error)
	AddFriend(ctx context.Context, req *social.AddFriendRequest, callOptions ...callopt.Option) (r *social.AddFriendResponse, err error)
	GetFriendship(ctx context.Context, req *social.GetFriendshipRequest, callOptions ...callopt.Option) (r *social.GetFriendshipResponse, err error)
	GetUserFriends(ctx context.Context, req *social.GetUserFriendsRequest, callOptions ...callopt.Option) (r *social.GetUserFriendsResponse, err error)
	CreateFriendRequest(ctx context.Context, req *social.CreateFriendRequestRequest, callOptions ...callopt.Option) (r *social.CreateFriendRequestResponse, err error)
	GetFriendRequests(ctx context.Context, req *social.GetFriendRequestsRequest, callOptions ...callopt.Option) (r *social.GetFriendRequestsResponse, err error)
	HandleFriendRequest(ctx context.Context, req *social.HandleFriendRequestRequest, callOptions ...callopt.Option) (r *social.HandleFriendRequestResponse, err error)
	MarkMessageRead(ctx context.Context, req *social.MarkMessageReadRequest, callOptions ...callopt.Option) (r *social.MarkMessageReadResponse, err error)
	GetUnreadMessageCount(ctx context.Context, req *social.GetUnreadMessageCountRequest, callOptions ...callopt.Option) (r *social.GetUnreadMessageCountResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfoForClient(), options...)
	if err != nil {
		return nil, err
	}
	return &kSocialServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kSocialServiceClient struct {
	*kClient
}

func (p *kSocialServiceClient) SendPrivateMessage(ctx context.Context, req *social.SendPrivateMessageRequest, callOptions ...callopt.Option) (r *social.SendPrivateMessageResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SendPrivateMessage(ctx, req)
}

func (p *kSocialServiceClient) GetPrivateMessages(ctx context.Context, req *social.GetPrivateMessagesRequest, callOptions ...callopt.Option) (r *social.GetPrivateMessagesResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetPrivateMessages(ctx, req)
}

func (p *kSocialServiceClient) CreateChatRoom(ctx context.Context, req *social.CreateChatRoomRequest, callOptions ...callopt.Option) (r *social.CreateChatRoomResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CreateChatRoom(ctx, req)
}

func (p *kSocialServiceClient) GetChatRoom(ctx context.Context, req *social.GetChatRoomRequest, callOptions ...callopt.Option) (r *social.GetChatRoomResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetChatRoom(ctx, req)
}

func (p *kSocialServiceClient) GetUserChatRooms(ctx context.Context, req *social.GetUserChatRoomsRequest, callOptions ...callopt.Option) (r *social.GetUserChatRoomsResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetUserChatRooms(ctx, req)
}

func (p *kSocialServiceClient) SendChatMessage(ctx context.Context, req *social.SendChatMessageRequest, callOptions ...callopt.Option) (r *social.SendChatMessageResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SendChatMessage(ctx, req)
}

func (p *kSocialServiceClient) GetChatMessages(ctx context.Context, req *social.GetChatMessagesRequest, callOptions ...callopt.Option) (r *social.GetChatMessagesResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetChatMessages(ctx, req)
}

func (p *kSocialServiceClient) AddFriend(ctx context.Context, req *social.AddFriendRequest, callOptions ...callopt.Option) (r *social.AddFriendResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.AddFriend(ctx, req)
}

func (p *kSocialServiceClient) GetFriendship(ctx context.Context, req *social.GetFriendshipRequest, callOptions ...callopt.Option) (r *social.GetFriendshipResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetFriendship(ctx, req)
}

func (p *kSocialServiceClient) GetUserFriends(ctx context.Context, req *social.GetUserFriendsRequest, callOptions ...callopt.Option) (r *social.GetUserFriendsResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetUserFriends(ctx, req)
}

func (p *kSocialServiceClient) CreateFriendRequest(ctx context.Context, req *social.CreateFriendRequestRequest, callOptions ...callopt.Option) (r *social.CreateFriendRequestResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CreateFriendRequest(ctx, req)
}

func (p *kSocialServiceClient) GetFriendRequests(ctx context.Context, req *social.GetFriendRequestsRequest, callOptions ...callopt.Option) (r *social.GetFriendRequestsResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetFriendRequests(ctx, req)
}

func (p *kSocialServiceClient) HandleFriendRequest(ctx context.Context, req *social.HandleFriendRequestRequest, callOptions ...callopt.Option) (r *social.HandleFriendRequestResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.HandleFriendRequest(ctx, req)
}

func (p *kSocialServiceClient) MarkMessageRead(ctx context.Context, req *social.MarkMessageReadRequest, callOptions ...callopt.Option) (r *social.MarkMessageReadResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.MarkMessageRead(ctx, req)
}

func (p *kSocialServiceClient) GetUnreadMessageCount(ctx context.Context, req *social.GetUnreadMessageCountRequest, callOptions ...callopt.Option) (r *social.GetUnreadMessageCountResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetUnreadMessageCount(ctx, req)
}
