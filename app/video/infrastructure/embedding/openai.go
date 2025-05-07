package embedding

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/yxrxy/videoHub/app/video/domain/repository"
)

type OpenAIEmbedding struct {
	client *openai.Client
	model  string
}

func NewOpenAIEmbedding(apiKey string, baseURL string, proxyURL string) repository.EmbeddingService {
	config := openai.DefaultConfig(apiKey)

	// 设置基础 URL（如果在配置中指定）
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	// 设置代理（如果在配置中指定）
	if proxyURL != "" {
		proxyUrl, err := url.Parse(proxyURL)
		if err == nil {
			httpClient := &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(proxyUrl),
				},
			}
			config.HTTPClient = httpClient
		}
	}

	client := openai.NewClientWithConfig(config)
	return &OpenAIEmbedding{
		client: client,
		model:  "text-embedding-ada-002",
	}
}

func (o *OpenAIEmbedding) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// 清理和准备文本
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("嵌入文本不能为空")
	}

	// 创建嵌入请求
	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.AdaEmbeddingV2,
	}

	resp, err := o.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("嵌入生成失败：没有返回结果")
	}

	return resp.Data[0].Embedding, nil
}
