package embedding

import (
	"context"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/sashabaranov/go-openai"
	"github.com/smartystreets/goconvey/convey"
)

func TestOpenAIEmbedding_GenerateEmbedding(t *testing.T) {
	type TestCase struct {
		Name          string
		InputText     string
		MockResponse  openai.EmbeddingResponse
		MockError     error
		ExpectedError error
		ExpectedEmbed []float32
	}

	testCases := []TestCase{
		{
			Name:      "成功场景",
			InputText: "测试文本",
			MockResponse: openai.EmbeddingResponse{
				Data: []openai.Embedding{
					{
						Embedding: []float32{0.1, 0.2, 0.3},
					},
				},
			},
			ExpectedEmbed: []float32{0.1, 0.2, 0.3},
		},
		{
			Name:          "空文本错误",
			InputText:     "",
			ExpectedError: fmt.Errorf("嵌入文本不能为空"),
		},
		{
			Name:          "OpenAI API 错误",
			InputText:     "测试文本",
			MockError:     fmt.Errorf("API 调用失败"),
			ExpectedError: fmt.Errorf("API 调用失败"),
		},
		{
			Name:      "空响应错误",
			InputText: "测试文本",
			MockResponse: openai.EmbeddingResponse{
				Data: []openai.Embedding{},
			},
			ExpectedError: fmt.Errorf("嵌入生成失败：没有返回结果"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// 创建 OpenAI 客户端实例
			svc := NewOpenAIEmbedding("test-key", "", "")

			// Mock CreateEmbeddings 方法
			mockey.Mock((*openai.Client).CreateEmbeddings).Return(tc.MockResponse, tc.MockError).Build()

			// 执行测试
			result, err := svc.GenerateEmbedding(context.Background(), tc.InputText)

			// 验证结果
			if tc.ExpectedError != nil {
				convey.So(err, convey.ShouldNotBeNil)
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
				return
			}

			convey.So(err, convey.ShouldBeNil)
			convey.So(result, convey.ShouldResemble, tc.ExpectedEmbed)
		})
	}
}

func TestNewOpenAIEmbedding(t *testing.T) {
	type TestCase struct {
		Name      string
		APIKey    string
		BaseURL   string
		ProxyURL  string
		CheckFunc func(*OpenAIEmbedding)
	}

	testCases := []TestCase{
		{
			Name:     "基本配置",
			APIKey:   "test-key",
			BaseURL:  "",
			ProxyURL: "",
			CheckFunc: func(e *OpenAIEmbedding) {
				convey.So(e.model, convey.ShouldEqual, "text-embedding-ada-002")
				convey.So(e.client, convey.ShouldNotBeNil)
			},
		},
		{
			Name:     "自定义基础URL",
			APIKey:   "test-key",
			BaseURL:  "https://custom.openai.api",
			ProxyURL: "",
			CheckFunc: func(e *OpenAIEmbedding) {
				convey.So(e.client, convey.ShouldNotBeNil)
			},
		},
		{
			Name:     "配置代理",
			APIKey:   "test-key",
			BaseURL:  "",
			ProxyURL: "http://proxy.example.com",
			CheckFunc: func(e *OpenAIEmbedding) {
				convey.So(e.client, convey.ShouldNotBeNil)
			},
		},
	}

	for _, tc := range testCases {
		convey.Convey(tc.Name, t, func() {
			svc := NewOpenAIEmbedding(tc.APIKey, tc.BaseURL, tc.ProxyURL)
			embedSvc, ok := svc.(*OpenAIEmbedding)
			convey.So(ok, convey.ShouldBeTrue)
			tc.CheckFunc(embedSvc)
		})
	}
}
