package controller

import (
	"context"
	"io"
	"mime/multipart"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/yxrrxy/videoHub/kitex_gen/video"
	"github.com/yxrrxy/videoHub/kitex_gen/video/videoservice"
	pkgcontext "github.com/yxrrxy/videoHub/pkg/context"
	"github.com/yxrrxy/videoHub/pkg/errno"
	"github.com/yxrrxy/videoHub/pkg/response"
	"github.com/yxrrxy/videoHub/pkg/utils"
)

type VideoController struct {
	client videoservice.Client
}

func NewVideoController(client videoservice.Client) *VideoController {
	return &VideoController{client: client}
}

type PublishFormRequest struct {
	Title       string                `form:"title" binding:"required"`
	Description string                `form:"description,omitempty"`
	Category    string                `form:"category,omitempty"`
	Tags        []string              `form:"tags,omitempty"`
	IsPrivate   bool                  `form:"is_private,omitempty"`
	Video       *multipart.FileHeader `form:"video" binding:"required"`
}

func (c *VideoController) Publish(ctx context.Context, req *app.RequestContext) {
	userID, exists := pkgcontext.GetUserID(ctx)
	if !exists {
		req.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	var formReq PublishFormRequest
	if err := req.BindAndValidate(&formReq); err != nil {
		req.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	f, err := formReq.Video.Open()
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}
	defer f.Close()

	fileData, err := io.ReadAll(f)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	rpcReq := &video.PublishRequest{
		UserId:      userID,
		VideoData:   fileData,
		ContentType: formReq.Video.Header.Get("Content-Type"),
		Title:       formReq.Title,
		Description: &formReq.Description,
		Category:    &formReq.Category,
		Tags:        formReq.Tags,
		IsPrivate:   &formReq.IsPrivate,
	}

	resp, err := c.client.Publish(ctx, rpcReq)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	data := map[string]interface{}{
		"video_url": resp.VideoUrl,
	}

	req.JSON(consts.StatusOK, response.Success(data))
}

func (c *VideoController) GetUserVideoList(ctx context.Context, req *app.RequestContext) {
	userID, exists := pkgcontext.GetUserID(ctx)
	if !exists {
		req.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	var formReq video.VideoListRequest
	if err := req.BindAndValidate(&formReq); err != nil {
		req.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	rpcReq := &video.VideoListRequest{
		UserId:   userID,
		Cursor:   formReq.Cursor,
		PageSize: formReq.PageSize,
		Category: formReq.Category,
	}

	resp, err := c.client.GetVideoList(ctx, rpcReq)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	items := make([]map[string]interface{}, 0)
	for _, v := range resp.Videos {
		video := map[string]interface{}{
			"id":            strconv.FormatInt(v.Id, 10),
			"user_id":       strconv.FormatInt(v.UserId, 10),
			"video_url":     v.VideoUrl,
			"cover_url":     v.CoverUrl,
			"title":         v.Title,
			"description":   v.Description,
			"visit_count":   v.VisitCount,
			"like_count":    v.LikeCount,
			"comment_count": v.CommentCount,
			"created_at":    utils.FormatTimestamp(v.CreatedAt),
			"updated_at":    utils.FormatTimestamp(v.UpdatedAt),
			"deleted_at":    utils.FormatTimestamp(*v.DeletedAt),
		}
		items = append(items, video)
	}

	data := map[string]interface{}{
		"items": items,
		"total": resp.Total,
	}

	req.JSON(consts.StatusOK, response.Success(data))
}
