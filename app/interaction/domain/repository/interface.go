package repository

import (
	"context"

	"github.com/yxrxy/videoHub/app/interaction/domain/model"
)

type InteractionRepository interface {
	GetLikeList(ctx context.Context, videoID int64, offset, limit int) ([]*model.Like, error)
	CreateLike(ctx context.Context, userID, videoID int64) error
	GetLike(ctx context.Context, userID, videoID int64) (*model.Like, error)
	DeleteLike(ctx context.Context, like *model.Like) error
	IsVideoExist(ctx context.Context, videoID int64) (bool, error)
	CreateComment(ctx context.Context, userID, videoID int64, content string, parentID *int64) error
	GetComment(ctx context.Context, commentID int64) (*model.Comment, error)
	GetCommentList(ctx context.Context, videoID int64, offset, limit int) ([]*model.Comment, int64, error)
	DeleteComment(ctx context.Context, comment *model.Comment) error
	LikeComment(ctx context.Context, userID, commentID int64) error
	UnlikeComment(ctx context.Context, userID, commentID int64) error
	IsCommentLiked(ctx context.Context, userID, commentID int64) (bool, error)
	GetCommentLikeCount(ctx context.Context, commentID int64) (int64, error)
}
