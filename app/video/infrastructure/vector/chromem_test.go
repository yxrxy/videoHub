package vector

import (
	"context"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/philippgille/chromem-go"
	"github.com/smartystreets/goconvey/convey"
	"github.com/yxrxy/videoHub/app/video/domain/model"
)

func TestChromemDB_StoreVector(t *testing.T) {
	type TestCase struct {
		Name          string
		VideoID       int64
		Vector        []float32
		Metadata      *model.VideoMetadata
		MockError     error
		ExpectedError error
	}

	testCases := []TestCase{
		{
			Name:    "成功场景",
			VideoID: 1,
			Vector:  []float32{0.1, 0.2, 0.3},
			Metadata: &model.VideoMetadata{
				Title:       "测试视频",
				Description: "这是一个测试视频",
				Tags:        []string{"测试", "视频"},
				Category:    "教育",
				UserID:      100,
			},
		},
		{
			Name:    "添加文档失败",
			VideoID: 2,
			Vector:  []float32{0.1, 0.2, 0.3},
			Metadata: &model.VideoMetadata{
				Title:    "测试视频",
				Category: "教育",
				UserID:   100,
			},
			MockError:     fmt.Errorf("添加文档失败"),
			ExpectedError: fmt.Errorf("添加文档失败"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// Mock AddDocument 方法
			mockey.Mock((*chromem.Collection).AddDocument).Return(tc.MockError).Build()

			// 创建 ChromemDB 实例
			db := &ChromemDB{
				db:         chromem.NewDB(),
				collection: &chromem.Collection{},
			}

			// 执行测试
			err := db.StoreVector(context.Background(), tc.VideoID, tc.Vector, tc.Metadata)

			// 验证结果
			if tc.ExpectedError != nil {
				convey.So(err, convey.ShouldNotBeNil)
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
				return
			}
			convey.So(err, convey.ShouldBeNil)
		})
	}
}

func TestChromemDB_SearchSimilar(t *testing.T) {
	type TestCase struct {
		Name             string
		QueryVector      []float32
		Limit            int32
		Filter           *model.VectorSearchFilter
		MockResults      []chromem.Result
		MockError        error
		ExpectedIDs      []int64
		ExpectedScores   []float32
		ExpectedError    error
		InvalidResultIDs bool // 用于测试无效的ID转换场景
	}

	testCases := []TestCase{
		{
			Name:        "成功场景",
			QueryVector: []float32{0.1, 0.2, 0.3},
			Limit:       2,
			Filter: &model.VectorSearchFilter{
				Category: strPtr("教育"),
			},
			MockResults: []chromem.Result{
				{ID: "1", Similarity: 0.9},
				{ID: "2", Similarity: 0.8},
			},
			ExpectedIDs:    []int64{1, 2},
			ExpectedScores: []float32{0.9, 0.8},
		},
		{
			Name:          "查询失败",
			QueryVector:   []float32{0.1, 0.2, 0.3},
			Limit:         2,
			MockError:     fmt.Errorf("查询失败"),
			ExpectedError: fmt.Errorf("查询失败"),
		},
		{
			Name:        "无效的ID格式",
			QueryVector: []float32{0.1, 0.2, 0.3},
			Limit:       2,
			MockResults: []chromem.Result{
				{ID: "invalid", Similarity: 0.9},
			},
			InvalidResultIDs: true,
			ExpectedError:    fmt.Errorf("strconv.ParseInt: parsing \"invalid\": invalid syntax"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// Mock QueryWithOptions 方法
			mockey.Mock((*chromem.Collection).QueryWithOptions).Return(tc.MockResults, tc.MockError).Build()

			// 创建 ChromemDB 实例
			db := &ChromemDB{
				db:         chromem.NewDB(),
				collection: &chromem.Collection{},
			}

			// 执行测试
			ids, scores, err := db.SearchSimilar(context.Background(), tc.QueryVector, tc.Limit, tc.Filter)

			// 验证结果
			if tc.ExpectedError != nil {
				convey.So(err, convey.ShouldNotBeNil)
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
				return
			}
			convey.So(err, convey.ShouldBeNil)
			convey.So(ids, convey.ShouldResemble, tc.ExpectedIDs)
			convey.So(scores, convey.ShouldResemble, tc.ExpectedScores)
		})
	}
}

func TestChromemDB_DeleteEmbedding(t *testing.T) {
	type TestCase struct {
		Name          string
		VideoID       int64
		MockError     error
		ExpectedError error
	}

	testCases := []TestCase{
		{
			Name:    "成功场景",
			VideoID: 1,
		},
		{
			Name:          "删除失败",
			VideoID:       2,
			MockError:     fmt.Errorf("删除失败"),
			ExpectedError: fmt.Errorf("删除失败"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// Mock Delete 方法
			mockey.Mock((*chromem.Collection).Delete).Return(tc.MockError).Build()

			// 创建 ChromemDB 实例
			db := &ChromemDB{
				db:         chromem.NewDB(),
				collection: &chromem.Collection{},
			}

			// 执行测试
			err := db.DeleteEmbedding(context.Background(), tc.VideoID)

			// 验证结果
			if tc.ExpectedError != nil {
				convey.So(err, convey.ShouldNotBeNil)
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
				return
			}
			convey.So(err, convey.ShouldBeNil)
		})
	}
}

func TestNewChromemDB(t *testing.T) {
	type TestCase struct {
		Name           string
		CollectionName string
		MockError      error
		ExpectedError  error
	}

	testCases := []TestCase{
		{
			Name:           "成功场景",
			CollectionName: "test_collection",
		},
		{
			Name:           "创建集合失败",
			CollectionName: "test_collection",
			MockError:      fmt.Errorf("创建集合失败"),
			ExpectedError:  fmt.Errorf("创建集合失败"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// Mock CreateCollection 方法
			mockey.Mock((*chromem.DB).CreateCollection).Return(nil, tc.MockError).Build()

			// 执行测试
			db, err := NewChromemDB(tc.CollectionName)

			// 验证结果
			if tc.ExpectedError != nil {
				convey.So(err, convey.ShouldNotBeNil)
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
				convey.So(db, convey.ShouldBeNil)
				return
			}
			convey.So(err, convey.ShouldBeNil)
			convey.So(db, convey.ShouldNotBeNil)
		})
	}
}

// 辅助函数：创建字符串指针
func strPtr(s string) *string {
	return &s
}
