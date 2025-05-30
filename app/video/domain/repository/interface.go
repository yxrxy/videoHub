package repository

import (
	"context"

	"github.com/olivere/elastic/v7"
	"github.com/yxrxy/videoHub/app/video/domain/model"
	"github.com/yxrxy/videoHub/pkg/kafka"
)

type VideoDB interface {
	CreateVideo(ctx context.Context, video *model.Video) error
	UpdateVideo(ctx context.Context, video *model.Video) error
	GetVideoList(ctx context.Context, userID, page int64, size int32, category *string) ([]*model.Video, int64, error)
	GetHotVideos(
		ctx context.Context,
		limit int32,
		category string,
		lastVisitCount, lastLikeCount, lastID int64,
	) ([]*model.Video, int64, int64, int64, int64, error)
	IncrementVisitCount(ctx context.Context, videoID int64) error
	IncrementLikeCount(ctx context.Context, videoID int64) error
	GetVideoByID(ctx context.Context, videoID int64) (*model.Video, error)
	DeleteVideo(ctx context.Context, videoID int64) error
}

type VideoCache interface {
	UpdateVideoScore(ctx context.Context, videoID int64, visitDelta, likeDelta int64, category string) error
	GetHotVideos(ctx context.Context, category string, limit int, lastVisitCount, lastLikeCount, lastID int64) ([]string, error)
	Load(key string) (interface{}, bool)
	Delete(key string) error
	Range(f func(key, value interface{}) bool)
	Store(key string, value interface{}) error
}

type VideoMQ interface {
	SendProcessVideo(ctx context.Context, videoID int64, videoPath string) error
	ConsumeProcessVideo(ctx context.Context) <-chan *kafka.Message
}

type VideoElastic interface {
	IsExist(ctx context.Context, indexName string) bool
	CreateIndex(ctx context.Context, indexName string) error
	AddItem(ctx context.Context, indexName string, video *model.Video, name string) error
	RemoveItem(ctx context.Context, indexName string, id int64) error
	UpdateItem(ctx context.Context, indexName string, video *model.VideoES, name string) error
	SearchItems(ctx context.Context, indexName string, query *model.VideoES) ([]int64, int64, error)
	BuildQuery(req *model.VideoES) *elastic.BoolQuery
}

type VectorDB interface {
	StoreVector(ctx context.Context, videoID int64, vector []float32, metadata *model.VideoMetadata) error
	SearchSimilar(ctx context.Context, queryVector []float32, limit int32, filter *model.VectorSearchFilter) ([]int64, []float32, error)
	DeleteEmbedding(ctx context.Context, videoID int64) error
}

type EmbeddingService interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
}

type LLMService interface {
	GenerateResponse(ctx context.Context, query string, documents []string) (string, error)
	GenerateRelatedQueries(ctx context.Context, query string) ([]string, error)
}
