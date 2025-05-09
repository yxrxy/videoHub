package llm

import (
	"context"
	"fmt"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/sashabaranov/go-openai"
	"github.com/smartystreets/goconvey/convey"
)

func TestOpenAILLM_GenerateResponse(t *testing.T) {
	type TestCase struct {
		Name           string
		Query          string
		Documents      []string
		MockResponse   openai.ChatCompletionResponse
		MockError      error
		ExpectedResult string
		ExpectedError  error
	}

	testCases := []TestCase{
		{
			Name:      "成功场景",
			Query:     "编程教程",
			Documents: []string{"Go语言教程", "Python基础入门"},
			MockResponse: openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: "找到两个相关的编程教程视频：Go语言教程和Python基础入门",
						},
					},
				},
			},
			ExpectedResult: "找到两个相关的编程教程视频：Go语言教程和Python基础入门",
		},
		{
			Name:          "API调用错误",
			Query:         "编程教程",
			Documents:     []string{"Go语言教程"},
			MockError:     fmt.Errorf("API调用失败"),
			ExpectedError: fmt.Errorf("API调用失败"),
		},
		{
			Name:      "空响应错误",
			Query:     "编程教程",
			Documents: []string{"Go语言教程"},
			MockResponse: openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{},
			},
			ExpectedError: fmt.Errorf("生成响应失败：没有返回结果"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// 创建 OpenAI 客户端实例
			svc := NewOpenAILLM("test-key", "", "")

			// Mock CreateChatCompletion 方法
			mockey.Mock((*openai.Client).CreateChatCompletion).Return(tc.MockResponse, tc.MockError).Build()

			// 执行测试
			result, err := svc.GenerateResponse(context.Background(), tc.Query, tc.Documents)

			// 验证结果
			if tc.ExpectedError != nil {
				convey.So(err, convey.ShouldNotBeNil)
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
				return
			}

			convey.So(err, convey.ShouldBeNil)
			convey.So(result, convey.ShouldEqual, tc.ExpectedResult)
		})
	}
}

func TestOpenAILLM_GenerateRelatedQueries(t *testing.T) {
	type TestCase struct {
		Name           string
		Query          string
		MockResponse   openai.ChatCompletionResponse
		MockError      error
		ExpectedResult []string
		ExpectedError  error
	}

	testCases := []TestCase{
		{
			Name:  "成功场景",
			Query: "编程教程",
			MockResponse: openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: "1. Go语言入门教程\n2. Python基础课程\n3. Java编程指南\n4. Web开发教程\n5. 数据结构与算法",
						},
					},
				},
			},
			ExpectedResult: []string{
				"Go语言入门教程",
				"Python基础课程",
				"Java编程指南",
				"Web开发教程",
				"数据结构与算法",
			},
		},
		{
			Name:          "API调用错误",
			Query:         "编程教程",
			MockError:     fmt.Errorf("API调用失败"),
			ExpectedError: fmt.Errorf("API调用失败"),
		},
		{
			Name:  "空响应错误",
			Query: "编程教程",
			MockResponse: openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{},
			},
			ExpectedError: fmt.Errorf("生成相关查询失败：没有返回结果"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.Name, t, func() {
			// 创建 OpenAI 客户端实例
			svc := NewOpenAILLM("test-key", "", "")

			// Mock CreateChatCompletion 方法
			mockey.Mock((*openai.Client).CreateChatCompletion).Return(tc.MockResponse, tc.MockError).Build()

			// 执行测试
			result, err := svc.GenerateRelatedQueries(context.Background(), tc.Query)

			// 验证结果
			if tc.ExpectedError != nil {
				convey.So(err, convey.ShouldNotBeNil)
				convey.So(err.Error(), convey.ShouldEqual, tc.ExpectedError.Error())
				return
			}

			convey.So(err, convey.ShouldBeNil)
			convey.So(result, convey.ShouldResemble, tc.ExpectedResult)
		})
	}
}

func TestNewOpenAILLM(t *testing.T) {
	type TestCase struct {
		Name      string
		APIKey    string
		BaseURL   string
		ProxyURL  string
		CheckFunc func(*OpenAILLM)
	}

	testCases := []TestCase{
		{
			Name:     "基本配置",
			APIKey:   "test-key",
			BaseURL:  "",
			ProxyURL: "",
			CheckFunc: func(e *OpenAILLM) {
				convey.So(e.model, convey.ShouldEqual, openai.GPT4oMini)
				convey.So(e.client, convey.ShouldNotBeNil)
			},
		},
		{
			Name:     "自定义基础URL",
			APIKey:   "test-key",
			BaseURL:  "https://custom.openai.api",
			ProxyURL: "",
			CheckFunc: func(e *OpenAILLM) {
				convey.So(e.client, convey.ShouldNotBeNil)
			},
		},
		{
			Name:     "配置代理",
			APIKey:   "test-key",
			BaseURL:  "",
			ProxyURL: "http://proxy.example.com",
			CheckFunc: func(e *OpenAILLM) {
				convey.So(e.client, convey.ShouldNotBeNil)
			},
		},
	}

	for _, tc := range testCases {
		convey.Convey(tc.Name, t, func() {
			svc := NewOpenAILLM(tc.APIKey, tc.BaseURL, tc.ProxyURL)
			llmSvc, ok := svc.(*OpenAILLM)
			convey.So(ok, convey.ShouldBeTrue)
			tc.CheckFunc(llmSvc)
		})
	}
}
