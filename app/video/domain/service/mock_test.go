package service

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/yxrxy/videoHub/app/video/domain/model"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Load(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *MockCache) Store(key string, value interface{}) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockCache) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockCache) Range(f func(key, value interface{}) bool) {
	m.Called(f)
}

func (m *MockCache) GetHotVideos(ctx context.Context, category string, limit int, lastVisitCount, lastLikeCount, lastID int64) ([]string, error) {
	args := m.Called(ctx, category, limit, lastVisitCount, lastLikeCount, lastID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCache) UpdateVideoScore(ctx context.Context, videoID int64, visitDelta, likeDelta int64, category string) error {
	args := m.Called(ctx, videoID, visitDelta, likeDelta, category)
	return args.Error(0)
}

type MockVectorDB struct {
	mock.Mock
}

func (m *MockVectorDB) SearchSimilar(ctx context.Context, embedding []float32, limit int32, filter *model.VectorSearchFilter) ([]int64, []float32, error) {
	args := m.Called(ctx, embedding, limit, filter)
	return args.Get(0).([]int64), args.Get(1).([]float32), args.Error(2)
}

func (m *MockVectorDB) DeleteEmbedding(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVectorDB) StoreVector(ctx context.Context, id int64, vector []float32, metadata *model.VideoMetadata) error {
	args := m.Called(ctx, id, vector, metadata)
	return args.Error(0)
}

type MockEmbedding struct {
	mock.Mock
}

func (m *MockEmbedding) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	args := m.Called(ctx, text)
	return args.Get(0).([]float32), args.Error(1)
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetVideoByID(ctx context.Context, id int64) (*model.Video, error) {
	args := m.Called(ctx, id)
	if v := args.Get(0); v != nil {
		return v.(*model.Video), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDB) CreateVideo(ctx context.Context, video *model.Video) error {
	args := m.Called(ctx, video)
	return args.Error(0)
}

func (m *MockDB) UpdateVideo(ctx context.Context, video *model.Video) error {
	args := m.Called(ctx, video)
	return args.Error(0)
}

func (m *MockDB) GetVideoList(ctx context.Context, userID, page int64, size int32, category *string) ([]*model.Video, int64, error) {
	args := m.Called(ctx, userID, page, size, category)
	return args.Get(0).([]*model.Video), args.Get(1).(int64), args.Error(2)
}

func (m *MockDB) GetHotVideos(ctx context.Context, limit int32, category string, lastVisitCount, lastLikeCount, lastID int64) ([]*model.Video, int64, int64, int64, int64, error) {
	args := m.Called(ctx, limit, category, lastVisitCount, lastLikeCount, lastID)
	return args.Get(0).([]*model.Video), args.Get(1).(int64), args.Get(2).(int64), args.Get(3).(int64), args.Get(4).(int64), args.Error(5)
}

func (m *MockDB) IncrementVisitCount(ctx context.Context, videoID int64) error {
	args := m.Called(ctx, videoID)
	return args.Error(0)
}

func (m *MockDB) IncrementLikeCount(ctx context.Context, videoID int64) error {
	args := m.Called(ctx, videoID)
	return args.Error(0)
}

func (m *MockDB) DeleteVideo(ctx context.Context, videoID int64) error {
	args := m.Called(ctx, videoID)
	return args.Error(0)
}

type MockLLM struct {
	mock.Mock
}

func (m *MockLLM) GenerateResponse(ctx context.Context, query string, texts []string) (string, error) {
	args := m.Called(ctx, query, texts)
	return args.String(0), args.Error(1)
}

func (m *MockLLM) GenerateRelatedQueries(ctx context.Context, query string) ([]string, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]string), args.Error(1)
}
