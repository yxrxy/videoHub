package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/smartystreets/goconvey/convey"
	"github.com/yxrxy/videoHub/app/video/domain/model"
	"github.com/yxrxy/videoHub/app/video/domain/repository"
	"github.com/yxrxy/videoHub/app/video/domain/service"
)

// 测试视频发布功能
//
// export MOCKEY_CHECK_GCFLAGS=false
func TestPublish(t *testing.T) {
	type TestCase struct {
		Name              string
		VideoData         []byte
		Title             string
		Description       string
		Category          string
		Tags              []string
		MockCheckVideo    bool
		MockSaveVideoPath string
		MockSaveVideoErr  error
		ExpectedPath      string
		ExpectedError     error
	}

	testCases := []TestCase{
		{
			Name:              "发布视频成功",
			VideoData:         []byte("test video data"),
			Title:             "测试视频",
			Description:       "测试描述",
			Category:          "测试分类",
			Tags:              []string{"标签1", "标签2"},
			MockCheckVideo:    true,
			MockSaveVideoPath: "video/path",
			MockSaveVideoErr:  nil,
			ExpectedPath:      "video/path",
			ExpectedError:     nil,
		},
		{
			Name:           "视频格式检查失败",
			VideoData:      []byte("invalid video data"),
			Title:          "测试视频",
			Description:    "测试描述",
			Category:       "测试分类",
			Tags:           []string{"标签1", "标签2"},
			MockCheckVideo: false,
			ExpectedPath:   "",
			ExpectedError:  fmt.Errorf("[1001] video format error"),
		},
		{
			Name:              "保存视频失败",
			VideoData:         []byte("test video data"),
			Title:             "测试视频",
			Description:       "测试描述",
			Category:          "测试分类",
			Tags:              []string{"标签1", "标签2"},
			MockCheckVideo:    true,
			MockSaveVideoPath: "",
			MockSaveVideoErr:  fmt.Errorf("failed to save video"),
			ExpectedPath:      "",
			ExpectedError:     fmt.Errorf("failed to save video"),
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// 创建 useCase 实例
			uc := &useCase{
				svc:   new(service.VideoService),
				db:    *new(repository.VideoDB),
				es:    *new(repository.VideoElastic),
				cache: *new(repository.VideoCache),
			}

			ctx := context.Background()

			// Mock CheckVideo 方法
			mockey.Mock((*service.VideoService).CheckVideo).
				Return(tc.MockCheckVideo).
				Build()

			// 只有当 CheckVideo 返回 true 时才会调用 SaveVideo
			if tc.MockCheckVideo {
				mockey.Mock((*service.VideoService).SaveVideo).
					Return(tc.MockSaveVideoPath, tc.MockSaveVideoErr).
					Build()
			}

			path, err := uc.Publish(ctx, 1, tc.VideoData, "video/mp4", tc.Title, &tc.Description, &tc.Category, tc.Tags, false)

			if tc.ExpectedError != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldBeNil)
			}
			convey.So(path, convey.ShouldEqual, tc.ExpectedPath)
		})
	}
}

// 测试语义搜索功能
func TestSemanticSearch(t *testing.T) {
	type TestCase struct {
		Name             string
		Query            string
		Limit            int32
		Score            float64
		UserID           int32
		MockSearchResult *model.SemanticSearchResultItem
		MockSearchError  error
		ExpectedResult   []*model.SemanticSearchResultItem
		ExpectedError    error
	}

	// 准备测试数据
	expectedVideo := &model.Video{
		ID:          1,
		Title:       "测试视频",
		UserID:      1,
		Category:    "测试分类",
		Tags:        "标签1,标签2",
		Description: "测试描述",
	}

	testCases := []TestCase{
		{
			Name:   "搜索成功",
			Query:  "测试查询",
			Limit:  10,
			Score:  0.5,
			UserID: 1,
			MockSearchResult: &model.SemanticSearchResultItem{
				Videos:         []*model.Video{expectedVideo},
				Summary:        "测试摘要",
				RelatedQueries: []string{"相关查询1"},
				FromCache:      false,
			},
			MockSearchError: nil,
			ExpectedResult: []*model.SemanticSearchResultItem{
				{
					Videos:         []*model.Video{expectedVideo},
					Summary:        "测试摘要",
					RelatedQueries: []string{"相关查询1"},
					FromCache:      false,
				},
			},
			ExpectedError: nil,
		},
		{
			Name:             "搜索失败",
			Query:            "测试查询",
			Limit:            10,
			Score:            0.5,
			UserID:           1,
			MockSearchResult: nil,
			MockSearchError:  fmt.Errorf("search failed"),
			ExpectedResult:   nil,
			ExpectedError:    fmt.Errorf("search failed"),
		},
		{
			Name:             "空查询",
			Query:            "",
			Limit:            10,
			Score:            0.5,
			UserID:           1,
			MockSearchResult: nil,
			MockSearchError:  fmt.Errorf("empty query"),
			ExpectedResult:   nil,
			ExpectedError:    fmt.Errorf("empty query"),
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// 创建 useCase 实例
			uc := &useCase{
				svc:   new(service.VideoService),
				db:    *new(repository.VideoDB),
				es:    *new(repository.VideoElastic),
				cache: *new(repository.VideoCache),
			}

			ctx := context.Background()

			// Mock Search 方法
			mockey.Mock((*service.VideoService).Search).
				Return(tc.MockSearchResult, tc.MockSearchError).
				Build()

			if tc.MockSearchResult != nil {
				mockey.Mock((*service.VideoService).GetVideoDetail).
					Return(expectedVideo, nil).
					Build()
			}

			// 执行测试
			result, err := uc.SemanticSearch(ctx, tc.Query, tc.Limit, tc.UserID, tc.Score)

			// 验证结果
			if tc.ExpectedError != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldBeNil)
			}

			if tc.ExpectedResult != nil {
				convey.So(result, convey.ShouldNotBeNil)
				convey.So(len(result), convey.ShouldEqual, len(tc.ExpectedResult))
				convey.So(result[0].Videos[0].ID, convey.ShouldEqual, tc.ExpectedResult[0].Videos[0].ID)
				convey.So(result[0].Summary, convey.ShouldEqual, tc.ExpectedResult[0].Summary)
				convey.So(result[0].RelatedQueries, convey.ShouldResemble, tc.ExpectedResult[0].RelatedQueries)
			} else {
				convey.So(result, convey.ShouldBeNil)
			}
		})
	}
}

// 测试获取热门视频功能
func TestGetHotVideos(t *testing.T) {
	type TestCase struct {
		Name           string
		Limit          int32
		Category       string
		LastVisit      int64
		LastLike       int64
		LastID         int64
		MockVideos     []*model.Video
		MockLastVisit  int64
		MockLastLike   int64
		MockLastID     int64
		MockTotal      int64
		MockError      error
		ExpectedVideos []*model.Video
		ExpectedTotal  int64
		ExpectedError  error
	}

	expectedVideos := []*model.Video{
		{ID: 1, Title: "热门视频1"},
		{ID: 2, Title: "热门视频2"},
	}

	testCases := []TestCase{
		{
			Name:           "获取热门视频成功",
			Limit:          10,
			Category:       "测试分类",
			LastVisit:      0,
			LastLike:       0,
			LastID:         0,
			MockVideos:     expectedVideos,
			MockLastVisit:  100,
			MockLastLike:   200,
			MockLastID:     300,
			MockTotal:      2,
			MockError:      nil,
			ExpectedVideos: expectedVideos,
			ExpectedTotal:  2,
			ExpectedError:  nil,
		},
		{
			Name:           "无热门视频",
			Limit:          10,
			Category:       "测试分类",
			LastVisit:      0,
			LastLike:       0,
			LastID:         0,
			MockVideos:     []*model.Video{},
			MockLastVisit:  0,
			MockLastLike:   0,
			MockLastID:     0,
			MockTotal:      0,
			MockError:      nil,
			ExpectedVideos: []*model.Video{},
			ExpectedTotal:  0,
			ExpectedError:  nil,
		},
		{
			Name:           "获取失败",
			Limit:          10,
			Category:       "测试分类",
			LastVisit:      0,
			LastLike:       0,
			LastID:         0,
			MockVideos:     nil,
			MockLastVisit:  0,
			MockLastLike:   0,
			MockLastID:     0,
			MockTotal:      0,
			MockError:      fmt.Errorf("failed to get hot videos"),
			ExpectedVideos: nil,
			ExpectedTotal:  0,
			ExpectedError:  fmt.Errorf("failed to get hot videos"),
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// 创建 useCase 实例
			uc := &useCase{
				svc:   new(service.VideoService),
				db:    *new(repository.VideoDB),
				es:    *new(repository.VideoElastic),
				cache: *new(repository.VideoCache),
			}

			ctx := context.Background()

			// Mock GetHotVideos 方法
			mockey.Mock((*service.VideoService).GetHotVideos).
				Return(tc.MockVideos, tc.MockLastVisit, tc.MockLastLike, tc.MockLastID, tc.MockTotal, tc.MockError).
				Build()

			// 执行测试
			videos, lastVisit, lastLike, lastID, total, err := uc.GetHotVideos(
				ctx, tc.Limit, tc.Category, tc.LastVisit, tc.LastLike, tc.LastID)

			// 验证结果
			if tc.ExpectedError != nil {
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
			} else {
				convey.So(err, convey.ShouldBeNil)
			}

			convey.So(videos, convey.ShouldResemble, tc.ExpectedVideos)
			convey.So(lastVisit, convey.ShouldEqual, tc.MockLastVisit)
			convey.So(lastLike, convey.ShouldEqual, tc.MockLastLike)
			convey.So(lastID, convey.ShouldEqual, tc.MockLastID)
			convey.So(total, convey.ShouldEqual, tc.ExpectedTotal)
		})
	}
}
