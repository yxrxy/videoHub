package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yxrxy/videoHub/app/video/domain/repository"
	"github.com/yxrxy/videoHub/app/video/domain/service"
	"github.com/yxrxy/videoHub/pkg/constants"
)

type VideoCache struct {
	client *redis.Client
}

func NewVideoCache(client *redis.Client) repository.VideoCache {
	return &VideoCache{client: client}
}

const (
	VideoHotKey    = "video:hot"      // 视频热度 zset
	VideoHotExpire = 30 * time.Minute // 热度数据过期时间
	CategoryHotKey = "video:hot:%s"   // 分类下的热门视频
)

// UpdateVideoScore 更新视频分数（同步更新总榜和分类榜）
func (v *VideoCache) UpdateVideoScore(ctx context.Context, videoID int64, visitDelta, likeDelta int64, category string) error {
	score := float64(visitDelta) + float64(likeDelta)*constants.VideoScoreWeight
	videoIDStr := strconv.FormatInt(videoID, 10)

	// 1. 更新总榜
	_, err := v.client.ZIncrBy(ctx, VideoHotKey, score, videoIDStr).Result()
	if err != nil {
		return err
	}
	v.client.Expire(ctx, VideoHotKey, VideoHotExpire)

	// 2. 更新分类榜（如果有分类）
	if category != "" {
		categoryKey := fmt.Sprintf(CategoryHotKey, category)
		if err := v.client.ZIncrBy(ctx, categoryKey, score, videoIDStr).Err(); err != nil {
			return err
		}
		v.client.Expire(ctx, categoryKey, VideoHotExpire)
	}

	return nil
}

// GetHotVideos 获取热门视频ID列表（支持分类和总榜）
func (v *VideoCache) GetHotVideos(ctx context.Context, category string, limit int, lastVisitCount, lastLikeCount, lastID int64) ([]string, error) {
	key := VideoHotKey // 默认查询总榜
	if category != "" {
		key = fmt.Sprintf(CategoryHotKey, category) // 查询分类榜
	}

	lastScore := float64(lastVisitCount) + float64(lastLikeCount)*constants.VideoScoreWeight

	var _min, _max string
	if lastScore > 0 {
		_min = "-inf"
		_max = fmt.Sprintf("%f", lastScore)
	} else {
		_min = "-inf"
		_max = "+inf"
	}

	// 获取热门视频ID
	videoIDs, err := v.client.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:   _min,
		Max:   _max,
		Count: int64(limit + 1),
	}).Result()
	if err != nil {
		return nil, err
	}

	// 如果 lastID 存在，跳过第一个（避免重复）
	if len(videoIDs) > 0 && lastID > 0 {
		firstID, _ := strconv.ParseInt(videoIDs[0], 10, 64)
		if firstID == lastID {
			videoIDs = videoIDs[1:]
		}
	}

	// 限制返回数量
	if len(videoIDs) > limit {
		videoIDs = videoIDs[:limit]
	}

	return videoIDs, nil
}

func (v *VideoCache) Load(key string) (interface{}, bool) {
	value, err := v.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, false
	}
	var entry service.CacheEntry
	if err := json.Unmarshal([]byte(value), &entry); err != nil {
		return nil, false
	}
	return entry, true
}

func (v *VideoCache) Delete(key string) error {
	return v.client.Del(context.Background(), key).Err()
}

func (v *VideoCache) Store(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return v.client.Set(context.Background(), key, data, 0).Err()
}

func (v *VideoCache) Range(f func(key, value interface{}) bool) {
	iter := v.client.Scan(context.Background(), 0, "", 0).Iterator()
	for iter.Next(context.Background()) {
		key := iter.Val()
		value, err := v.client.Get(context.Background(), key).Result()
		if err != nil {
			continue
		}
		if !f(key, value) {
			break
		}
	}
}
