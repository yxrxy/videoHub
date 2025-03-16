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
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, "视频文件读取失败"))
			return
		}
	}(f)

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

// 添加请求结构体
type VideoListRequest struct {
	Page     int64  `form:"page" query:"page"`         // 页码
	Size     int32  `form:"size" query:"size"`         // 每页数量
	Category string `form:"category" query:"category"` // 分类
}

func (c *VideoController) GetUserVideoList(ctx context.Context, req *app.RequestContext) {
	userID, exists := pkgcontext.GetUserID(ctx)
	if !exists {
		req.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	var formReq VideoListRequest
	if err := req.BindAndValidate(&formReq); err != nil {
		req.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	// 构建RPC请求
	rpcReq := &video.VideoListRequest{
		UserId:   userID,
		Page:     &formReq.Page,
		Size:     &formReq.Size,
		Category: &formReq.Category,
	}

	resp, err := c.client.GetVideoList(ctx, rpcReq)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	items := make([]map[string]interface{}, 0)
	for _, v := range resp.Videos {
		video2 := map[string]interface{}{
			"id":            strconv.FormatInt(v.Id, 10),
			"user_id":       strconv.FormatInt(v.UserId, 10),
			"video_url":     v.VideoUrl,
			"cover_url":     v.CoverUrl,
			"title":         v.Title,
			"description":   v.Description,
			"duration":      v.Duration,
			"category":      v.Category,
			"tags":          v.Tags,
			"visit_count":   v.VisitCount,
			"like_count":    v.LikeCount,
			"comment_count": v.CommentCount,
			"is_private":    v.IsPrivate,
			"created_at":    utils.FormatTimestamp(v.CreatedAt),
			"updated_at":    utils.FormatTimestamp(v.UpdatedAt),
		}
		if v.DeletedAt != nil {
			video2["deleted_at"] = utils.FormatTimestamp(*v.DeletedAt)
		}
		items = append(items, video2)
	}

	data := map[string]interface{}{
		"items": items,
		"total": resp.Total,
	}

	req.JSON(consts.StatusOK, response.Success(data))
}

func (c *VideoController) GetHotVideoList(ctx context.Context, req *app.RequestContext) {
	var formReq video.HotVideoRequest
	if err := req.BindAndValidate(&formReq); err != nil {
		req.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	rpcReq := &video.HotVideoRequest{
		Limit:     formReq.Limit,
		Category:  formReq.Category,
		LastVisit: formReq.LastVisit,
		LastLike:  formReq.LastLike,
		LastId:    formReq.LastId,
	}

	resp, err := c.client.GetHotVideos(ctx, rpcReq)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	items := make([]map[string]interface{}, 0)
	for _, v := range resp.Videos {
		video2 := map[string]interface{}{
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
		if v.DeletedAt != nil {
			video2["deleted_at"] = utils.FormatTimestamp(*v.DeletedAt)
		}
		items = append(items, video2)
	}

	data := map[string]interface{}{
		"items":      items,
		"total":      resp.Total,
		"next_visit": resp.NextVisit,
		"next_like":  resp.NextLike,
		"next_id":    resp.NextId,
	}

	req.JSON(consts.StatusOK, response.Success(data))
}

func (c *VideoController) IncrementLikeCount(ctx context.Context, req *app.RequestContext) {
	var formReq struct {
		VideoID int64 `json:"video_id" binding:"required"`
	}

	if err := req.BindAndValidate(&formReq); err != nil {
		req.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	rpcReq := &video.IncrementLikeCountRequest{
		VideoId: formReq.VideoID,
	}

	resp, err := c.client.IncrementLikeCount(ctx, rpcReq)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	req.JSON(consts.StatusOK, response.Success(map[string]interface{}{
		"success": resp.Success,
	}))
}

func (c *VideoController) IncrementVisitCount(ctx context.Context, req *app.RequestContext) {
	var formReq struct {
		VideoID int64 `json:"video_id" binding:"required"`
	}

	if err := req.BindAndValidate(&formReq); err != nil {
		req.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, err.Error()))
		return
	}

	rpcReq := &video.IncrementVisitCountRequest{
		VideoId: formReq.VideoID,
	}

	resp, err := c.client.IncrementVisitCount(ctx, rpcReq)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	req.JSON(consts.StatusOK, response.Success(map[string]interface{}{
		"success": resp.Success,
	}))
}
