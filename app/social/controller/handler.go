package handler

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/yxrrxy/videoHub/kitex_gen/social/socialservice"
	pkgcontext "github.com/yxrrxy/videoHub/pkg/context"
	"github.com/yxrrxy/videoHub/pkg/errno"
	"github.com/yxrrxy/videoHub/pkg/response"
)

type SocialHandler struct {
	client socialservice.Client
}

func NewSocialHandler(client socialservice.Client) *SocialHandler {
	return &SocialHandler{
		client: client,
	}
}

// SendPrivateMessage 发送私信
func (s *SocialHandler) SendPrivateMessage(ctx context.Context, c *app.RequestContext) {
	userID, exists := pkgcontext.GetUserID(ctx)
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	var req struct {
		ReceiverID int64  `json:"receiver_id"`
		Content    string `json:"content"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	if err := s.client.SendPrivateMessage(ctx, userID, req.ReceiverID, req.Content); err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(nil))
}

// GetPrivateMessages 获取私信列表
func (s *SocialHandler) GetPrivateMessages(ctx context.Context, c *app.RequestContext) {
	userID, exists := pkgcontext.GetUserID(ctx)
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	receiverID, err := strconv.ParseInt(c.Query("receiver_id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	page := int32(1)
	size := int32(20)

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = int32(p)
		}
	}
	if sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			size = int32(s)
		}
	}

	resp, err := s.client.GetPrivateMessages(ctx, userID, receiverID, page, size)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(map[string]interface{}{
		"messages": resp.Messages,
		"total":    resp.Total,
	}))
}

// CreateChatRoom 创建聊天室
func (s *SocialHandler) CreateChatRoom(ctx context.Context, c *app.RequestContext) {
	userID, exists := pkgcontext.GetUserID(ctx)
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	var req struct {
		Name      string  `json:"name"`
		Type      int8    `json:"type"`
		MemberIDs []int64 `json:"member_ids"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	room, err := s.client.CreateChatRoom(ctx, req.Name, userID, req.Type, req.MemberIDs)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(room))
}

// GetChatRoom 获取聊天室信息
func (s *SocialHandler) GetChatRoom(ctx context.Context, c *app.RequestContext) {
	roomID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	room, err := s.client.GetChatRoom(ctx, roomID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(room))
}

// GetUserChatRooms 获取用户的聊天室列表
func (s *SocialHandler) GetUserChatRooms(ctx context.Context, c *app.RequestContext) {
	userID, exists := pkgcontext.GetUserID(ctx)
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	rooms, err := s.client.GetUserChatRooms(ctx, userID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(rooms))
}

// SendChatMessage 发送聊天室消息
func (s *SocialHandler) SendChatMessage(ctx context.Context, c *app.RequestContext) {
	userID, exists := pkgcontext.GetUserID(ctx)
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	roomID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	var req struct {
		Content string `json:"content"`
		Type    int8   `json:"type"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	if err := s.client.SendChatMessage(ctx, roomID, userID, req.Content, req.Type); err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(nil))
}

// GetChatMessages 获取聊天室消息列表
func (s *SocialHandler) GetChatMessages(ctx context.Context, c *app.RequestContext) {
	roomID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	page := int32(1)
	size := int32(20)

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = int32(p)
		}
	}
	if sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			size = int32(s)
		}
	}

	resp, err := s.client.GetChatMessages(ctx, roomID, page, size)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(map[string]interface{}{
		"messages": resp.Messages,
		"total":    resp.Total,
	}))
}

// AddFriend 添加好友
func (s *SocialHandler) AddFriend(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	var req struct {
		FriendID int64 `json:"friend_id"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	if err := s.client.AddFriend(ctx, userID.(int64), req.FriendID); err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(nil))
}

// GetFriendship 获取好友关系信息
func (s *SocialHandler) GetFriendship(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	friendID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	friendship, err := s.client.GetFriendship(ctx, userID.(int64), friendID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(friendship))
}

// GetUserFriends 获取用户的好友列表
func (s *SocialHandler) GetUserFriends(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	friends, err := s.client.GetUserFriends(ctx, userID.(int64))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(friends))
}

// CreateFriendRequest 创建好友申请
func (s *SocialHandler) CreateFriendRequest(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	var req struct {
		ReceiverID int64  `json:"receiver_id"`
		Message    string `json:"message"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	if err := s.client.CreateFriendRequest(ctx, userID.(int64), req.ReceiverID, req.Message); err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(nil))
}

// GetFriendRequests 获取好友申请列表
func (s *SocialHandler) GetFriendRequests(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	statusStr := c.Query("status")
	status := int8(0)
	if statusStr != "" {
		if s, err := strconv.ParseInt(statusStr, 10, 8); err == nil {
			status = int8(s)
		}
	}

	requests, err := s.client.GetFriendRequests(ctx, userID.(int64), status)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(requests))
}

// HandleFriendRequest 处理好友申请
func (s *SocialHandler) HandleFriendRequest(ctx context.Context, c *app.RequestContext) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	requestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	var req struct {
		Status int8 `json:"status"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	if err := s.client.HandleFriendRequest(ctx, requestID, req.Status); err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(nil))
}

// MarkMessageRead 标记消息为已读
func (s *SocialHandler) MarkMessageRead(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	messageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	if err := s.client.MarkMessageRead(ctx, messageID, userID.(int64)); err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(nil))
}

// GetUnreadMessageCount 获取未读消息数量
func (s *SocialHandler) GetUnreadMessageCount(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	count, err := s.client.GetUnreadMessageCount(ctx, userID.(int64))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(map[string]interface{}{
		"count": count,
	}))
}
