// Code generated by Kitex v0.13.1. DO NOT EDIT.

package socialservice

import (
	"context"
	"errors"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	social "github.com/yxrxy/videoHub/kitex_gen/social"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"SendPrivateMessage": kitex.NewMethodInfo(
		sendPrivateMessageHandler,
		newSocialServiceSendPrivateMessageArgs,
		newSocialServiceSendPrivateMessageResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetPrivateMessages": kitex.NewMethodInfo(
		getPrivateMessagesHandler,
		newSocialServiceGetPrivateMessagesArgs,
		newSocialServiceGetPrivateMessagesResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"CreateChatRoom": kitex.NewMethodInfo(
		createChatRoomHandler,
		newSocialServiceCreateChatRoomArgs,
		newSocialServiceCreateChatRoomResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetChatRoom": kitex.NewMethodInfo(
		getChatRoomHandler,
		newSocialServiceGetChatRoomArgs,
		newSocialServiceGetChatRoomResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetUserChatRooms": kitex.NewMethodInfo(
		getUserChatRoomsHandler,
		newSocialServiceGetUserChatRoomsArgs,
		newSocialServiceGetUserChatRoomsResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"SendChatMessage": kitex.NewMethodInfo(
		sendChatMessageHandler,
		newSocialServiceSendChatMessageArgs,
		newSocialServiceSendChatMessageResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetChatMessages": kitex.NewMethodInfo(
		getChatMessagesHandler,
		newSocialServiceGetChatMessagesArgs,
		newSocialServiceGetChatMessagesResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"AddFriend": kitex.NewMethodInfo(
		addFriendHandler,
		newSocialServiceAddFriendArgs,
		newSocialServiceAddFriendResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetFriendship": kitex.NewMethodInfo(
		getFriendshipHandler,
		newSocialServiceGetFriendshipArgs,
		newSocialServiceGetFriendshipResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetUserFriends": kitex.NewMethodInfo(
		getUserFriendsHandler,
		newSocialServiceGetUserFriendsArgs,
		newSocialServiceGetUserFriendsResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"CreateFriendRequest": kitex.NewMethodInfo(
		createFriendRequestHandler,
		newSocialServiceCreateFriendRequestArgs,
		newSocialServiceCreateFriendRequestResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetFriendRequests": kitex.NewMethodInfo(
		getFriendRequestsHandler,
		newSocialServiceGetFriendRequestsArgs,
		newSocialServiceGetFriendRequestsResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"HandleFriendRequest": kitex.NewMethodInfo(
		handleFriendRequestHandler,
		newSocialServiceHandleFriendRequestArgs,
		newSocialServiceHandleFriendRequestResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"MarkMessageRead": kitex.NewMethodInfo(
		markMessageReadHandler,
		newSocialServiceMarkMessageReadArgs,
		newSocialServiceMarkMessageReadResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetUnreadMessageCount": kitex.NewMethodInfo(
		getUnreadMessageCountHandler,
		newSocialServiceGetUnreadMessageCountArgs,
		newSocialServiceGetUnreadMessageCountResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
}

var (
	socialServiceServiceInfo                = NewServiceInfo()
	socialServiceServiceInfoForClient       = NewServiceInfoForClient()
	socialServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return socialServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return socialServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return socialServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "SocialService"
	handlerType := (*social.SocialService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "social",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.13.1",
		Extra:           extra,
	}
	return svcInfo
}

func sendPrivateMessageHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceSendPrivateMessageArgs)
	realResult := result.(*social.SocialServiceSendPrivateMessageResult)
	success, err := handler.(social.SocialService).SendPrivateMessage(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceSendPrivateMessageArgs() interface{} {
	return social.NewSocialServiceSendPrivateMessageArgs()
}

func newSocialServiceSendPrivateMessageResult() interface{} {
	return social.NewSocialServiceSendPrivateMessageResult()
}

func getPrivateMessagesHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceGetPrivateMessagesArgs)
	realResult := result.(*social.SocialServiceGetPrivateMessagesResult)
	success, err := handler.(social.SocialService).GetPrivateMessages(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceGetPrivateMessagesArgs() interface{} {
	return social.NewSocialServiceGetPrivateMessagesArgs()
}

func newSocialServiceGetPrivateMessagesResult() interface{} {
	return social.NewSocialServiceGetPrivateMessagesResult()
}

func createChatRoomHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceCreateChatRoomArgs)
	realResult := result.(*social.SocialServiceCreateChatRoomResult)
	success, err := handler.(social.SocialService).CreateChatRoom(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceCreateChatRoomArgs() interface{} {
	return social.NewSocialServiceCreateChatRoomArgs()
}

func newSocialServiceCreateChatRoomResult() interface{} {
	return social.NewSocialServiceCreateChatRoomResult()
}

func getChatRoomHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceGetChatRoomArgs)
	realResult := result.(*social.SocialServiceGetChatRoomResult)
	success, err := handler.(social.SocialService).GetChatRoom(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceGetChatRoomArgs() interface{} {
	return social.NewSocialServiceGetChatRoomArgs()
}

func newSocialServiceGetChatRoomResult() interface{} {
	return social.NewSocialServiceGetChatRoomResult()
}

func getUserChatRoomsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceGetUserChatRoomsArgs)
	realResult := result.(*social.SocialServiceGetUserChatRoomsResult)
	success, err := handler.(social.SocialService).GetUserChatRooms(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceGetUserChatRoomsArgs() interface{} {
	return social.NewSocialServiceGetUserChatRoomsArgs()
}

func newSocialServiceGetUserChatRoomsResult() interface{} {
	return social.NewSocialServiceGetUserChatRoomsResult()
}

func sendChatMessageHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceSendChatMessageArgs)
	realResult := result.(*social.SocialServiceSendChatMessageResult)
	success, err := handler.(social.SocialService).SendChatMessage(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceSendChatMessageArgs() interface{} {
	return social.NewSocialServiceSendChatMessageArgs()
}

func newSocialServiceSendChatMessageResult() interface{} {
	return social.NewSocialServiceSendChatMessageResult()
}

func getChatMessagesHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceGetChatMessagesArgs)
	realResult := result.(*social.SocialServiceGetChatMessagesResult)
	success, err := handler.(social.SocialService).GetChatMessages(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceGetChatMessagesArgs() interface{} {
	return social.NewSocialServiceGetChatMessagesArgs()
}

func newSocialServiceGetChatMessagesResult() interface{} {
	return social.NewSocialServiceGetChatMessagesResult()
}

func addFriendHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceAddFriendArgs)
	realResult := result.(*social.SocialServiceAddFriendResult)
	success, err := handler.(social.SocialService).AddFriend(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceAddFriendArgs() interface{} {
	return social.NewSocialServiceAddFriendArgs()
}

func newSocialServiceAddFriendResult() interface{} {
	return social.NewSocialServiceAddFriendResult()
}

func getFriendshipHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceGetFriendshipArgs)
	realResult := result.(*social.SocialServiceGetFriendshipResult)
	success, err := handler.(social.SocialService).GetFriendship(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceGetFriendshipArgs() interface{} {
	return social.NewSocialServiceGetFriendshipArgs()
}

func newSocialServiceGetFriendshipResult() interface{} {
	return social.NewSocialServiceGetFriendshipResult()
}

func getUserFriendsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceGetUserFriendsArgs)
	realResult := result.(*social.SocialServiceGetUserFriendsResult)
	success, err := handler.(social.SocialService).GetUserFriends(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceGetUserFriendsArgs() interface{} {
	return social.NewSocialServiceGetUserFriendsArgs()
}

func newSocialServiceGetUserFriendsResult() interface{} {
	return social.NewSocialServiceGetUserFriendsResult()
}

func createFriendRequestHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceCreateFriendRequestArgs)
	realResult := result.(*social.SocialServiceCreateFriendRequestResult)
	success, err := handler.(social.SocialService).CreateFriendRequest(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceCreateFriendRequestArgs() interface{} {
	return social.NewSocialServiceCreateFriendRequestArgs()
}

func newSocialServiceCreateFriendRequestResult() interface{} {
	return social.NewSocialServiceCreateFriendRequestResult()
}

func getFriendRequestsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceGetFriendRequestsArgs)
	realResult := result.(*social.SocialServiceGetFriendRequestsResult)
	success, err := handler.(social.SocialService).GetFriendRequests(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceGetFriendRequestsArgs() interface{} {
	return social.NewSocialServiceGetFriendRequestsArgs()
}

func newSocialServiceGetFriendRequestsResult() interface{} {
	return social.NewSocialServiceGetFriendRequestsResult()
}

func handleFriendRequestHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceHandleFriendRequestArgs)
	realResult := result.(*social.SocialServiceHandleFriendRequestResult)
	success, err := handler.(social.SocialService).HandleFriendRequest(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceHandleFriendRequestArgs() interface{} {
	return social.NewSocialServiceHandleFriendRequestArgs()
}

func newSocialServiceHandleFriendRequestResult() interface{} {
	return social.NewSocialServiceHandleFriendRequestResult()
}

func markMessageReadHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceMarkMessageReadArgs)
	realResult := result.(*social.SocialServiceMarkMessageReadResult)
	success, err := handler.(social.SocialService).MarkMessageRead(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceMarkMessageReadArgs() interface{} {
	return social.NewSocialServiceMarkMessageReadArgs()
}

func newSocialServiceMarkMessageReadResult() interface{} {
	return social.NewSocialServiceMarkMessageReadResult()
}

func getUnreadMessageCountHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*social.SocialServiceGetUnreadMessageCountArgs)
	realResult := result.(*social.SocialServiceGetUnreadMessageCountResult)
	success, err := handler.(social.SocialService).GetUnreadMessageCount(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newSocialServiceGetUnreadMessageCountArgs() interface{} {
	return social.NewSocialServiceGetUnreadMessageCountArgs()
}

func newSocialServiceGetUnreadMessageCountResult() interface{} {
	return social.NewSocialServiceGetUnreadMessageCountResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) SendPrivateMessage(ctx context.Context, req *social.SendPrivateMessageRequest) (r *social.SendPrivateMessageResponse, err error) {
	var _args social.SocialServiceSendPrivateMessageArgs
	_args.Req = req
	var _result social.SocialServiceSendPrivateMessageResult
	if err = p.c.Call(ctx, "SendPrivateMessage", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetPrivateMessages(ctx context.Context, req *social.GetPrivateMessagesRequest) (r *social.GetPrivateMessagesResponse, err error) {
	var _args social.SocialServiceGetPrivateMessagesArgs
	_args.Req = req
	var _result social.SocialServiceGetPrivateMessagesResult
	if err = p.c.Call(ctx, "GetPrivateMessages", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) CreateChatRoom(ctx context.Context, req *social.CreateChatRoomRequest) (r *social.CreateChatRoomResponse, err error) {
	var _args social.SocialServiceCreateChatRoomArgs
	_args.Req = req
	var _result social.SocialServiceCreateChatRoomResult
	if err = p.c.Call(ctx, "CreateChatRoom", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetChatRoom(ctx context.Context, req *social.GetChatRoomRequest) (r *social.GetChatRoomResponse, err error) {
	var _args social.SocialServiceGetChatRoomArgs
	_args.Req = req
	var _result social.SocialServiceGetChatRoomResult
	if err = p.c.Call(ctx, "GetChatRoom", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetUserChatRooms(ctx context.Context, req *social.GetUserChatRoomsRequest) (r *social.GetUserChatRoomsResponse, err error) {
	var _args social.SocialServiceGetUserChatRoomsArgs
	_args.Req = req
	var _result social.SocialServiceGetUserChatRoomsResult
	if err = p.c.Call(ctx, "GetUserChatRooms", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) SendChatMessage(ctx context.Context, req *social.SendChatMessageRequest) (r *social.SendChatMessageResponse, err error) {
	var _args social.SocialServiceSendChatMessageArgs
	_args.Req = req
	var _result social.SocialServiceSendChatMessageResult
	if err = p.c.Call(ctx, "SendChatMessage", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetChatMessages(ctx context.Context, req *social.GetChatMessagesRequest) (r *social.GetChatMessagesResponse, err error) {
	var _args social.SocialServiceGetChatMessagesArgs
	_args.Req = req
	var _result social.SocialServiceGetChatMessagesResult
	if err = p.c.Call(ctx, "GetChatMessages", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) AddFriend(ctx context.Context, req *social.AddFriendRequest) (r *social.AddFriendResponse, err error) {
	var _args social.SocialServiceAddFriendArgs
	_args.Req = req
	var _result social.SocialServiceAddFriendResult
	if err = p.c.Call(ctx, "AddFriend", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetFriendship(ctx context.Context, req *social.GetFriendshipRequest) (r *social.GetFriendshipResponse, err error) {
	var _args social.SocialServiceGetFriendshipArgs
	_args.Req = req
	var _result social.SocialServiceGetFriendshipResult
	if err = p.c.Call(ctx, "GetFriendship", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetUserFriends(ctx context.Context, req *social.GetUserFriendsRequest) (r *social.GetUserFriendsResponse, err error) {
	var _args social.SocialServiceGetUserFriendsArgs
	_args.Req = req
	var _result social.SocialServiceGetUserFriendsResult
	if err = p.c.Call(ctx, "GetUserFriends", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) CreateFriendRequest(ctx context.Context, req *social.CreateFriendRequestRequest) (r *social.CreateFriendRequestResponse, err error) {
	var _args social.SocialServiceCreateFriendRequestArgs
	_args.Req = req
	var _result social.SocialServiceCreateFriendRequestResult
	if err = p.c.Call(ctx, "CreateFriendRequest", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetFriendRequests(ctx context.Context, req *social.GetFriendRequestsRequest) (r *social.GetFriendRequestsResponse, err error) {
	var _args social.SocialServiceGetFriendRequestsArgs
	_args.Req = req
	var _result social.SocialServiceGetFriendRequestsResult
	if err = p.c.Call(ctx, "GetFriendRequests", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) HandleFriendRequest(ctx context.Context, req *social.HandleFriendRequestRequest) (r *social.HandleFriendRequestResponse, err error) {
	var _args social.SocialServiceHandleFriendRequestArgs
	_args.Req = req
	var _result social.SocialServiceHandleFriendRequestResult
	if err = p.c.Call(ctx, "HandleFriendRequest", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) MarkMessageRead(ctx context.Context, req *social.MarkMessageReadRequest) (r *social.MarkMessageReadResponse, err error) {
	var _args social.SocialServiceMarkMessageReadArgs
	_args.Req = req
	var _result social.SocialServiceMarkMessageReadResult
	if err = p.c.Call(ctx, "MarkMessageRead", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetUnreadMessageCount(ctx context.Context, req *social.GetUnreadMessageCountRequest) (r *social.GetUnreadMessageCountResponse, err error) {
	var _args social.SocialServiceGetUnreadMessageCountArgs
	_args.Req = req
	var _result social.SocialServiceGetUnreadMessageCountResult
	if err = p.c.Call(ctx, "GetUnreadMessageCount", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
