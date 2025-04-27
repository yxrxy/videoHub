package rpc

import (
	"context"

	"github.com/yxrxy/videoHub/app/interaction/usecase"
	"github.com/yxrxy/videoHub/kitex_gen/interaction"
	rpcmodel "github.com/yxrxy/videoHub/kitex_gen/model"
	"github.com/yxrxy/videoHub/pkg/base"
	pkgContext "github.com/yxrxy/videoHub/pkg/base/context"
)

type InteractionHandler struct {
	useCase usecase.InteractionUseCase
}

func NewInteractionHandler(useCase usecase.InteractionUseCase) *InteractionHandler {
	return &InteractionHandler{useCase: useCase}
}

func (h *InteractionHandler) Like(ctx context.Context, req *interaction.LikeRequest) (r *interaction.LikeResponse, err error) {
	r = new(interaction.LikeResponse)
	userID, err := pkgContext.GetUserID(ctx)
	if err != nil {
		return
	}
	if _, err = h.useCase.Like(ctx, userID, req.VideoId); err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	return
}

func (h *InteractionHandler) GetLikes(ctx context.Context, req *interaction.GetLikesRequest) (r *interaction.GetLikesResponse, err error) {
	r = new(interaction.GetLikesResponse)

	likes, err := h.useCase.GetLikes(ctx, req.VideoId, req.Page, req.Size)
	if err != nil {
		return
	}
	r.LikeList = make([]*rpcmodel.LikeInfo, len(likes))
	for i, like := range likes {
		r.LikeList[i] = &rpcmodel.LikeInfo{
			UserId:  like.UserID,
			VideoId: like.VideoID,
		}
	}
	return
}

func (h *InteractionHandler) Comment(ctx context.Context, req *interaction.CommentRequest) (r *interaction.CommentResponse, err error) {
	r = new(interaction.CommentResponse)
	userID, err := pkgContext.GetUserID(ctx)
	if err != nil {
		return
	}
	if req.ParentId == nil {
		*req.ParentId = -1
	}
	if _, err = h.useCase.Comment(ctx, userID, req.VideoId, req.Content, *req.ParentId); err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	return
}

func (h *InteractionHandler) GetComments(ctx context.Context, req *interaction.GetCommentsRequest) (r *interaction.GetCommentsResponse, err error) {
	r = new(interaction.GetCommentsResponse)

	comments, err := h.useCase.GetComments(ctx, req.VideoId, req.Page, req.Size)
	if err != nil {
		return
	}
	r.CommentList = make([]*rpcmodel.Comment, len(comments))
	for i, comment := range comments {
		r.CommentList[i] = &rpcmodel.Comment{
			Id:      comment.ID,
			UserId:  comment.UserID,
			VideoId: comment.VideoID,
			Content: comment.Content,
		}
	}
	return
}

func (h *InteractionHandler) DeleteComment(ctx context.Context, req *interaction.DeleteCommentRequest) (r *interaction.DeleteCommentResponse, err error) {
	r = new(interaction.DeleteCommentResponse)
	userID, err := pkgContext.GetUserID(ctx)
	if err != nil {
		return
	}
	if _, err = h.useCase.DeleteComment(ctx, userID, req.CommentId); err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	return
}

func (h *InteractionHandler) LikeComment(ctx context.Context, req *interaction.LikeCommentRequest) (r *interaction.LikeCommentResponse, err error) {
	r = new(interaction.LikeCommentResponse)
	userID, err := pkgContext.GetUserID(ctx)
	if err != nil {
		return
	}
	if _, err = h.useCase.LikeComment(ctx, userID, req.CommentId); err != nil {
		return
	}
	r.Base = base.BuildBaseResp(err)
	return
}
