package usecase

import (
	"context"

	"github.com/yxrxy/videoHub/app/interaction/domain/model"
	"github.com/yxrxy/videoHub/app/interaction/domain/repository"
	"github.com/yxrxy/videoHub/app/interaction/domain/service"
)

type InteractionUseCase interface {
	Like(ctx context.Context, userID int64, videoID int64) (bool, error)
	GetLikes(ctx context.Context, videoID int64, page int32, size int32) ([]*model.Like, error)
	Comment(ctx context.Context, userID int64, videoID int64, content string, parentID int64) (bool, error)
	GetComments(ctx context.Context, videoID int64, page int32, size int32) ([]*model.Comment, error)
	DeleteComment(ctx context.Context, userID int64, commentID int64) (bool, error)
	LikeComment(ctx context.Context, userID int64, commentID int64) (bool, error)
}

type useCase struct {
	db  repository.InteractionRepository
	svc *service.InteractionService
}

func NewInteractionCase(db repository.InteractionRepository, svc *service.InteractionService) *useCase {
	return &useCase{
		db:  db,
		svc: svc,
	}
}
