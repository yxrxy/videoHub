package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/yxrxy/videoHub/app/video/domain/model"
)

// CacheEntry represents a cached search result with timestamp
type CacheEntry struct {
	Result    *model.SemanticSearchResultItem `json:"result"`
	Timestamp time.Time                       `json:"timestamp"`
}

func (s *VideoService) Search(
	ctx context.Context,
	query string,
	limit int32,
) (*model.SemanticSearchResultItem, error) {
	// 检查缓存
	cacheKey := fmt.Sprintf("%s:%d", query, limit)
	if entry, ok := s.cache.Load(cacheKey); ok {
		if ce, ok := entry.(CacheEntry); ok {
			// 10分钟内的有效
			if time.Since(ce.Timestamp) < 10*time.Minute {
				ce.Result.FromCache = true
				return ce.Result, nil
			}
			if s.cache.Delete(cacheKey) != nil {
				return nil, fmt.Errorf("删除缓存失败")
			}
		}
	}

	// 并发执行向量搜索和ES搜索
	var (
		vectorResults = make(chan struct {
			ids    []int64
			scores []float32
			err    error
		}, 1)
		/*esResults = make(chan struct {
			ids   []int64
			total int64
			err   error
		}, 1)*/
	)

	// 向量搜索协程
	go func() {
		queryVector, err := s.embedding.GenerateEmbedding(ctx, query)
		if err != nil {
			vectorResults <- struct {
				ids    []int64
				scores []float32
				err    error
			}{nil, nil, err}
			return
		}

		/*vectorFilter := &model.VectorSearchFilter{
			Category: &filter.Category,
		}*/
		var f = "all"
		vectorFilter := &model.VectorSearchFilter{
			Category: &f,
		}
		ids, scores, err := s.vectorDB.SearchSimilar(ctx, queryVector, limit, vectorFilter)
		vectorResults <- struct {
			ids    []int64
			scores []float32
			err    error
		}{ids, scores, err}
	}()

	// ES搜索协程
	/*go func() {
		ids, total, err := s.es.SearchItems(ctx, "video", filter)
		esResults <- struct {
			ids   []int64
			total int64
			err   error
		}{ids, total, err}
	}()*/

	// 等待结果
	vectorResult := <-vectorResults
	if vectorResult.err != nil {
		return nil, vectorResult.err
	}

	/*esResult := <-esResults
	if esResult.err != nil {
		return nil, esResult.err
	}*/

	// 合并和排序结果
	videoScores := make(map[int64]float32)
	for i, id := range vectorResult.ids {
		videoScores[id] = vectorResult.scores[i]
	}

	/*for _, id := range esResult.ids {
		if _, exists := videoScores[id]; !exists {
			videoScores[id] = 0.3
		}
	}*/

	var rankedVideos []struct {
		ID    int64
		Score float32
	}
	for id, score := range videoScores {
		rankedVideos = append(rankedVideos, struct {
			ID    int64
			Score float32
		}{id, score})
	}
	sort.Slice(rankedVideos, func(i, j int) bool {
		return rankedVideos[i].Score > rankedVideos[j].Score
	})

	// 获取视频详情（批量获取）
	var videos []*model.Video
	var videoTexts []string
	for _, vs := range rankedVideos {
		if video, err := s.db.GetVideoByID(ctx, vs.ID); err == nil {
			videos = append(videos, video)
			videoTexts = append(videoTexts,
				fmt.Sprintf("%s %s %s", video.Title, video.Description, video.Tags))
		}
	}

	// 并发生成摘要和相关查询
	var (
		summaryResult = make(chan struct {
			summary string
			err     error
		}, 1)
		queriesResult = make(chan struct {
			queries []string
			err     error
		}, 1)
	)

	go func() {
		summary, err := s.llm.GenerateResponse(ctx, query, videoTexts)
		summaryResult <- struct {
			summary string
			err     error
		}{summary, err}
	}()

	go func() {
		queries, err := s.llm.GenerateRelatedQueries(ctx, query)
		queriesResult <- struct {
			queries []string
			err     error
		}{queries, err}
	}()

	summaryRes := <-summaryResult
	if summaryRes.err != nil {
		return nil, summaryRes.err
	}

	queriesRes := <-queriesResult
	if queriesRes.err != nil {
		return nil, queriesRes.err
	}

	result := &model.SemanticSearchResultItem{
		Videos:         videos,
		Summary:        summaryRes.summary,
		RelatedQueries: queriesRes.queries,
		FromCache:      false,
	}

	// 缓存结果
	err := s.cache.Store(cacheKey, CacheEntry{
		Result:    result,
		Timestamp: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// IndexVideo 添加视频时同时更新向量索引和(ES索引TODO)
func (s *VideoService) IndexVideo(ctx context.Context, video *model.Video, _ string) error {
	// 生成视频内容的向量表示
	err := s.GenerateVideoEmbedding(ctx, video.ID)
	if err != nil {
		return fmt.Errorf("生成向量失败: %w", err)
	}
	return nil
}

func (s *VideoService) DeleteVideoEmbedding(ctx context.Context, videoID int64) error {
	err := s.vectorDB.DeleteEmbedding(ctx, videoID)
	if err != nil {
		return fmt.Errorf("删除向量失败: %w", err)
	}
	return nil
}

func (s *VideoService) ClearRelatedCache(category string) {
	// 遍历并清除相关缓存
	s.cache.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok && strings.Contains(k, category) {
			err := s.cache.Delete(k)
			if err != nil {
				return false
			}
		}
		return true
	})
}
