package controller

import (
	"context"
	"io"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/yxrrxy/videoHub/kitex_gen/video"
	"github.com/yxrrxy/videoHub/kitex_gen/video/videoservice"
	pkgcontext "github.com/yxrrxy/videoHub/pkg/context"
	"github.com/yxrrxy/videoHub/pkg/errno"
	"github.com/yxrrxy/videoHub/pkg/response"
)

type VideoController struct {
	client videoservice.Client
}

func NewVideoController(client videoservice.Client) *VideoController {
	return &VideoController{client: client}
}

func (c *VideoController) Publish(ctx context.Context, req *app.RequestContext) {
	userID, exists := pkgcontext.GetUserID(ctx)
	if !exists {
		req.JSON(consts.StatusUnauthorized, response.Error(errno.ErrUnauthorized.ErrCode, errno.ErrUnauthorized.ErrMsg))
		return
	}

	file, err := req.FormFile("video")
	if err != nil {
		req.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, errno.ErrInvalidParam.ErrMsg))
		return
	}

	title := string(req.FormValue("title"))
	if title == "" {
		req.JSON(consts.StatusBadRequest, response.Error(errno.ErrInvalidParam.ErrCode, "视频标题不能为空"))
		return
	}

	f, err := file.Open()
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

	publishReq := &video.PublishRequest{
		UserId:      userID,
		VideoData:   fileData,
		ContentType: file.Header.Get("Content-Type"),
		Title:       title,
		Description: string(req.FormValue("description")),
	}

	resp, err := c.client.Publish(ctx, publishReq)
	if err != nil {
		req.JSON(consts.StatusInternalServerError, response.Error(errno.InternalServerError.ErrCode, err.Error()))
		return
	}

	data := map[string]interface{}{
		"video_url": resp.VideoUrl,
	}

	req.JSON(consts.StatusOK, response.Success(data))
}
