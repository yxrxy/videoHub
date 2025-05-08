package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"

	"github.com/yxrxy/videoHub/app/video/domain/model"
)

// TestVideoService_Search 测试视频搜索功能
func TestVideoService_Search(t *testing.T) {
	type TestCase struct {
		Name string
		// 搜索参数
		Query string
		Limit int32
		// Mock 缓存相关
		MockCacheKey  string
		MockCacheHit  bool
		MockCacheData *model.SemanticSearchResultItem
		MockCacheErr  error
		// Mock 向量搜索相关
		MockEmbedding           []float32
		MockEmbeddingErr        error
		MockVectorSearchResults []int64
		MockVectorSearchScores  []float32
		MockVectorSearchErr     error
		// Mock 数据库查询相关
		MockVideos []*model.Video
		MockDBErr  error
		// Mock LLM 相关
		MockSummary    string
		MockSummaryErr error
		MockQueries    []string
		MockQueriesErr error
		// 预期结果
		ExpectedError  error
		ExpectedResult *model.SemanticSearchResultItem
	}

	testCases := []TestCase{
		{
			Name:         "缓存命中成功",
			Query:        "test query",
			Limit:        10,
			MockCacheKey: "test query:10",
			MockCacheHit: true,
			MockCacheData: &model.SemanticSearchResultItem{
				Videos:         []*model.Video{{ID: 1, Title: "Test Video"}},
				Summary:        "Test Summary",
				RelatedQueries: []string{"related query"},
				FromCache:      true,
			},
			ExpectedError: nil,
			ExpectedResult: &model.SemanticSearchResultItem{
				Videos:         []*model.Video{{ID: 1, Title: "Test Video"}},
				Summary:        "Test Summary",
				RelatedQueries: []string{"related query"},
				FromCache:      true,
			},
		},
		{
			Name:                    "缓存未命中但搜索成功",
			Query:                   "test query",
			Limit:                   10,
			MockCacheHit:            false,
			MockEmbedding:           []float32{0.1, 0.2},
			MockVectorSearchResults: []int64{1},
			MockVectorSearchScores:  []float32{0.9},
			MockVideos:              []*model.Video{{ID: 1, Title: "Test Video"}},
			MockSummary:             "Test Summary",
			MockQueries:             []string{"related query"},
			ExpectedError:           nil,
			ExpectedResult: &model.SemanticSearchResultItem{
				Videos:         []*model.Video{{ID: 1, Title: "Test Video"}},
				Summary:        "Test Summary",
				RelatedQueries: []string{"related query"},
				FromCache:      false,
			},
		},
		{
			Name:             "生成向量嵌入失败",
			Query:            "test query",
			Limit:            10,
			MockCacheHit:     false,
			MockEmbeddingErr: errors.New("embedding generation failed"),
			ExpectedError:    fmt.Errorf("failed to generate embedding: embedding generation failed"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// Mock 缓存服务
			mockCache := new(MockCache)
			if tc.MockCacheHit {
				mockCache.On("Load", tc.MockCacheKey).Return(cacheEntry{
					result:    tc.MockCacheData,
					timestamp: time.Now(),
				}, true)
			} else {
				mockCache.On("Load", tc.MockCacheKey).Return(nil, false)
				mockCache.On("Store", tc.MockCacheKey, mock.Anything).Return(nil)
			}

			// Mock 向量数据库服务
			mockVectorDB := new(MockVectorDB)
			mockVectorDB.On("SearchSimilar",
				mock.Anything,
				tc.MockEmbedding,
				tc.Limit,
				mock.Anything,
			).Return(tc.MockVectorSearchResults, tc.MockVectorSearchScores, tc.MockVectorSearchErr)

			// Mock 嵌入服务
			mockEmbedding := new(MockEmbedding)
			mockEmbedding.On("GenerateEmbedding",
				mock.Anything,
				tc.Query,
			).Return(tc.MockEmbedding, tc.MockEmbeddingErr)

			// Mock 数据库服务
			mockDB := new(MockDB)
			for _, id := range tc.MockVectorSearchResults {
				mockDB.On("GetVideoByID", mock.Anything, id).Return(
					tc.MockVideos[0], tc.MockDBErr,
				)
			}

			// Mock LLM 服务
			mockLLM := new(MockLLM)
			mockLLM.On("GenerateResponse",
				mock.Anything,
				tc.Query,
				mock.Anything,
			).Return(tc.MockSummary, tc.MockSummaryErr)
			mockLLM.On("GenerateRelatedQueries",
				mock.Anything,
				tc.Query,
			).Return(tc.MockQueries, tc.MockQueriesErr)

			// 创建服务实例
			svc := &VideoService{
				cache:     mockCache,
				vectorDB:  mockVectorDB,
				embedding: mockEmbedding,
				db:        mockDB,
				llm:       mockLLM,
			}

			// 执行测试
			result, err := svc.Search(context.Background(), tc.Query, tc.Limit)

			// 验证结果
			if err != nil || tc.ExpectedError != nil {
				if err != nil && tc.ExpectedError != nil {
					convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
				} else {
					convey.So(err, convey.ShouldEqual, tc.ExpectedError)
				}
				return
			}

			// 验证成功场景的返回值
			convey.So(result, convey.ShouldResemble, tc.ExpectedResult)

			// 验证 mock 调用
			mockCache.AssertExpectations(t)
			mockVectorDB.AssertExpectations(t)
			mockEmbedding.AssertExpectations(t)
			mockDB.AssertExpectations(t)
			mockLLM.AssertExpectations(t)
		})
	}
}
