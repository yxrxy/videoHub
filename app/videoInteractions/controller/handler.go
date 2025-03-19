package controller

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/yxrrxy/videoHub/kitex_gen/videoInteractions/interactionservice"
	pkgcontext "github.com/yxrrxy/videoHub/pkg/context"
	"github.com/yxrrxy/videoHub/pkg/errno"
	"github.com/yxrrxy/videoHub/pkg/response"
)

type InteractionController struct {
	client interactionservice.Client
}

func NewInteractionController(client interactionservice.Client) *InteractionController {
	return &InteractionController{
		client: client,
	}
}

func (s *InteractionController) Like(ctx context.Context, c *app.RequestContext) {
	userID, exists := pkgcontext.GetUserID(ctx)
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	var likeReq struct {
		VideoID int64 `json:"video_id"`
	}

	if err := c.BindAndValidate(&likeReq); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	success, err := s.client.Like(ctx, userID, likeReq.VideoID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InteractionError.ErrCode, err.Error()))
		return
	}
	if !success {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InteractionError.ErrCode, errno.InteractionError.ErrMsg))
		return
	}

	c.JSON(consts.StatusOK, response.Success(nil))
}

func (s *InteractionController) GetLikes(ctx context.Context, c *app.RequestContext) {
	var req struct {
		VideoID int64 `json:"video_id"`
		Page    int32 `json:"page" default:"1"`
		Size    int32 `json:"size" default:"20"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	likes, err := s.client.GetLikes(ctx, req.VideoID, req.Page, req.Size)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InteractionError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(likes))
}

// Comment 发表评论
func (s *InteractionController) Comment(ctx context.Context, c *app.RequestContext) {
	userID, exists := pkgcontext.GetUserID(ctx)
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	var req struct {
		VideoID  int64  `json:"video_id"`
		Content  string `json:"content"`
		ParentID *int64 `json:"parent_id,omitempty"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	var parentID int64
	if req.ParentID != nil {
		parentID = *req.ParentID
	}

	success, err := s.client.Comment(ctx, userID, req.VideoID, req.Content, parentID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InteractionError.ErrCode, err.Error()))
		return
	}
	if !success {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InteractionError.ErrCode, errno.InteractionError.ErrMsg))
		return
	}

	c.JSON(consts.StatusOK, response.Success(nil))
}

// GetComments 获取评论列表
func (s *InteractionController) GetComments(ctx context.Context, c *app.RequestContext) {
	var req struct {
		VideoID int64 `json:"video_id"`
		Page    int32 `json:"page" default:"1"`
		Size    int32 `json:"size" default:"20"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	comments, err := s.client.GetComments(ctx, req.VideoID, req.Page, req.Size)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InteractionError.ErrCode, err.Error()))
		return
	}

	c.JSON(consts.StatusOK, response.Success(comments))
}

// DeleteComment 删除评论
func (s *InteractionController) DeleteComment(ctx context.Context, c *app.RequestContext) {
	userID, exists := pkgcontext.GetUserID(ctx)
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	var req struct {
		CommentID int64 `json:"comment_id"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	success, err := s.client.DeleteComment(ctx, userID, req.CommentID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InteractionError.ErrCode, err.Error()))
		return
	}
	if !success {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InteractionError.ErrCode, errno.InteractionError.ErrMsg))
		return
	}

	c.JSON(consts.StatusOK, response.Success(nil))
}

// LikeComment 点赞评论
func (s *InteractionController) LikeComment(ctx context.Context, c *app.RequestContext) {
	userID, exists := pkgcontext.GetUserID(ctx)
	if !exists {
		c.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	var req struct {
		CommentID int64 `json:"comment_id"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	success, err := s.client.LikeComment(ctx, userID, req.CommentID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InteractionError.ErrCode, err.Error()))
		return
	}
	if !success {
		c.JSON(consts.StatusInternalServerError, response.Error(errno.InteractionError.ErrCode, errno.InteractionError.ErrMsg))
		return
	}

	c.JSON(consts.StatusOK, response.Success(nil))
}
