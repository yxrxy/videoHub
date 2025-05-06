package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/yxrxy/videoHub/app/video/domain/repository"
)

type OpenAILLM struct {
	client *openai.Client
	model  string
}

func NewOpenAILLM(apiKey string) repository.LLMService {
	client := openai.NewClient(apiKey)
	return &OpenAILLM{
		client: client,
		model:  openai.GPT3Dot5Turbo,
	}
}

func (o *OpenAILLM) GenerateResponse(ctx context.Context, query string, documents []string) (string, error) {
	// 构建提示
	systemPrompt := "你是一个视频搜索助手，基于提供的视频信息，帮助用户找到最相关的内容。请提供简洁有用的摘要。"
	userPrompt := fmt.Sprintf("用户搜索：%s\n\n相关视频信息：\n%s",
		query, strings.Join(documents, "\n\n"))

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: userPrompt,
		},
	}

	req := openai.ChatCompletionRequest{
		Model:     o.model,
		Messages:  messages,
		MaxTokens: 500,
	}

	// 调用API
	resp, err := o.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("生成响应失败：没有返回结果")
	}

	return resp.Choices[0].Message.Content, nil
}

func (o *OpenAILLM) GenerateRelatedQueries(ctx context.Context, query string) ([]string, error) {
	systemPrompt := "你是一个视频搜索助手，请基于用户的搜索词，生成5个相关的搜索建议，每行一个。"
	userPrompt := fmt.Sprintf("用户搜索：%s\n\n生成相关搜索建议：", query)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: userPrompt,
		},
	}

	req := openai.ChatCompletionRequest{
		Model:     o.model,
		Messages:  messages,
		MaxTokens: 200,
	}

	resp, err := o.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("生成相关查询失败：没有返回结果")
	}

	// 解析结果
	content := resp.Choices[0].Message.Content
	lines := strings.Split(strings.TrimSpace(content), "\n")

	// 清理结果
	var queries []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			// 移除可能的数字前缀 (如 "1. ")
			if idx := strings.Index(line, ". "); idx > 0 && idx < 4 {
				line = strings.TrimSpace(line[idx+2:])
			}
			queries = append(queries, line)
		}
	}

	return queries, nil
}
