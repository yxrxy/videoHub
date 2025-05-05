package usecase

import (
	"context"

	"github.com/yxrxy/videoHub/app/video/domain/model"
	"github.com/yxrxy/videoHub/app/video/domain/repository"
	"github.com/yxrxy/videoHub/app/video/domain/service"
)

type VideoUseCase interface {
	Publish(
		ctx context.Context,
		userID int64,
		videoData []byte,
		contentType, title string,
		description, category *string,
		tags []string,
		isPrivate bool,
	) (string, error)
	GetVideoList(ctx context.Context, userID, page int64, size int32, category *string) ([]*model.Video, int64, error)
	GetVideoDetail(ctx context.Context, videoID, userID int64) (*model.Video, error)
	GetHotVideos(
		ctx context.Context,
		limit int32,
		category string,
		lastVisit, lastLike, lastID int64,
	) ([]*model.Video, int64, int64, int64, int64, error)
	DeleteVideo(ctx context.Context, videoID, userID int64) error
	IncrementVisitCount(ctx context.Context, videoID int64) error
	IncrementLikeCount(ctx context.Context, videoID int64) error
	SearchVideo(
		ctx context.Context,
		keywords string,
		pageSize, pageNum int32,
		fromDate, toDate *int64,
		username *string,
	) ([]*model.Video, int64, error)
}

type useCase struct {
	db    repository.VideoDB
	cache repository.VideoCache
	es    repository.VideoElastic
	svc   *service.VideoService
}

func NewVideoCase(db repository.VideoDB, cache repository.VideoCache, es repository.VideoElastic, svc *service.VideoService) *useCase {
	return &useCase{
		db:    db,
		cache: cache,
		es:    es,
		svc:   svc,
	}
}
