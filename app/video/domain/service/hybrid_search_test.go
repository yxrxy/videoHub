package service

import (
	"context"
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
			Name:  "成功场景-缓存未命中",
			Query: "test query",
			Limit: 10,
			// 缓存配置
			MockCacheKey:  "test query:10",
			MockCacheHit:  false,
			MockCacheData: nil,
			// 向量搜索配置
			MockEmbedding:           []float32{0.1, 0.2, 0.3},
			MockVectorSearchResults: []int64{1, 2},
			MockVectorSearchScores:  []float32{0.9, 0.8},
			// 数据库配置
			MockVideos: []*model.Video{
				{
					ID:          1,
					Title:       "测试视频1",
					Description: "描述1",
					Tags:        "标签1",
				},
				{
					ID:          2,
					Title:       "测试视频2",
					Description: "描述2",
					Tags:        "标签2",
				},
			},
			// LLM配置
			MockSummary: "测试摘要",
			MockQueries: []string{"相关查询1", "相关查询2"},
			// 预期结果
			ExpectedResult: &model.SemanticSearchResultItem{
				Videos: []*model.Video{
					{
						ID:          1,
						Title:       "测试视频1",
						Description: "描述1",
						Tags:        "标签1",
					},
					{
						ID:          2,
						Title:       "测试视频2",
						Description: "描述2",
						Tags:        "标签2",
					},
				},
				Summary:        "测试摘要",
				RelatedQueries: []string{"相关查询1", "相关查询2"},
				FromCache:      false,
			},
		},
		{
			Name:  "成功场景-缓存命中",
			Query: "cached query",
			Limit: 5,
			// 缓存配置
			MockCacheKey: "cached query:5",
			MockCacheHit: true,
			MockCacheData: &model.SemanticSearchResultItem{
				Videos: []*model.Video{
					{
						ID:          3,
						Title:       "缓存视频",
						Description: "缓存描述",
						Tags:        "缓存标签",
					},
				},
				Summary:        "缓存摘要",
				RelatedQueries: []string{"缓存相关查询"},
				FromCache:      true,
			},
			// 预期结果
			ExpectedResult: &model.SemanticSearchResultItem{
				Videos: []*model.Video{
					{
						ID:          3,
						Title:       "缓存视频",
						Description: "缓存描述",
						Tags:        "缓存标签",
					},
				},
				Summary:        "缓存摘要",
				RelatedQueries: []string{"缓存相关查询"},
				FromCache:      true,
			},
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
				mockCache.On("Load", tc.MockCacheKey).Return(cacheEntry{}, false)
				mockCache.On("Store", tc.MockCacheKey, mock.Anything).Return(nil)
			}

			// Mock 向量数据库服务
			mockVectorDB := new(MockVectorDB)
			if !tc.MockCacheHit {
				mockVectorDB.On("SearchSimilar",
					mock.Anything,
					tc.MockEmbedding,
					tc.Limit,
					mock.Anything,
				).Return(tc.MockVectorSearchResults, tc.MockVectorSearchScores, tc.MockVectorSearchErr)
			}

			// Mock 嵌入服务
			mockEmbedding := new(MockEmbedding)
			if !tc.MockCacheHit {
				mockEmbedding.On("GenerateEmbedding",
					mock.Anything,
					tc.Query,
				).Return(tc.MockEmbedding, tc.MockEmbeddingErr)
			}

			// Mock 数据库服务
			mockDB := new(MockDB)
			if !tc.MockCacheHit {
				for i, id := range tc.MockVectorSearchResults {
					mockDB.On("GetVideoByID", mock.Anything, id).Return(
						tc.MockVideos[i], tc.MockDBErr,
					)
				}
			}

			// Mock LLM 服务
			mockLLM := new(MockLLM)
			if !tc.MockCacheHit {
				mockLLM.On("GenerateResponse",
					mock.Anything,
					tc.Query,
					mock.Anything,
				).Return(tc.MockSummary, tc.MockSummaryErr)
				mockLLM.On("GenerateRelatedQueries",
					mock.Anything,
					tc.Query,
				).Return(tc.MockQueries, tc.MockQueriesErr)
			}

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
			if tc.ExpectedError != nil {
				convey.So(err, convey.ShouldNotBeNil)
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
				return
			}

			convey.So(err, convey.ShouldBeNil)
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
