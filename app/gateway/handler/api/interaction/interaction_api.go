// Code generated by hertz generator.

package interaction

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	api "github.com/yxrxy/videoHub/app/gateway/model/interaction"
	"github.com/yxrxy/videoHub/app/gateway/pack"
	"github.com/yxrxy/videoHub/app/gateway/rpc"
	"github.com/yxrxy/videoHub/kitex_gen/interaction"
	"github.com/yxrxy/videoHub/pkg/errno"
)

// Like .
// @router /api/v1/video/:video_id/like [POST]
func Like(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.LikeRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.RespError(c, errno.ParamVerifyError.WithError(err))
		return
	}
	err = rpc.LikeVideoRPC(ctx, &interaction.LikeRequest{
		VideoId: req.VideoID,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	pack.RespSuccess(c)
}

// GetLikes .
// @router /api/v1/video/:video_id/likes [GET]
func GetLikes(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetLikesRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.RespError(c, errno.ParamVerifyError.WithError(err))
		return
	}

	resp, _, err := rpc.GetLikesRPC(ctx, &interaction.GetLikesRequest{
		VideoId: req.VideoID,
		Page:    req.Page,
		Size:    req.Size,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	pack.RespData(c, resp)
}

// Comment .
// @router /api/v1/video/:video_id/comment [POST]
func Comment(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.CommentRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.RespError(c, errno.ParamVerifyError.WithError(err))
		return
	}
	if req.ParentID == nil {
		defaultValue := int64(-1)
		req.ParentID = &defaultValue
	}
	err = rpc.CommentVideoRPC(ctx, &interaction.CommentRequest{
		VideoId:  req.VideoID,
		Content:  req.Content,
		ParentId: req.ParentID,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	pack.RespSuccess(c)
}

// GetComments .
// @router /api/v1/video/:video_id/comments [GET]
func GetComments(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.GetCommentsRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.RespError(c, errno.ParamVerifyError.WithError(err))
		return
	}

	resp, _, err := rpc.GetCommentsRPC(ctx, &interaction.GetCommentsRequest{
		VideoId: req.VideoID,
		Page:    req.Page,
		Size:    req.Size,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	pack.RespData(c, resp)
}

// DeleteComment .
// @router /api/v1/comment/:comment_id [DELETE]
func DeleteComment(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DeleteCommentRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.RespError(c, errno.ParamVerifyError.WithError(err))
		return
	}

	err = rpc.DeleteCommentRPC(ctx, &interaction.DeleteCommentRequest{
		CommentId: req.CommentID,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	pack.RespSuccess(c)
}

// LikeComment .
// @router /api/v1/comment/:comment_id/like [POST]
func LikeComment(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.LikeCommentRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		pack.RespError(c, errno.ParamVerifyError.WithError(err))
		return
	}

	err = rpc.LikeCommentRPC(ctx, &interaction.LikeCommentRequest{
		CommentId: req.CommentID,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}
	pack.RespSuccess(c)
}
