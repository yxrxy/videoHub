package embedding

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/yxrxy/videoHub/app/video/domain/repository"
)

type OpenAIEmbedding struct {
	client *openai.Client
	model  string
}

func NewOpenAIEmbedding(apiKey string) repository.EmbeddingService {
	client := openai.NewClient(apiKey)
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
