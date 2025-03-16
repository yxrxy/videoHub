package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/yxrrxy/videoHub/config"
)

var (
	client *redis.Client
)

func RedisInit() {
	client = config.GetClient()
}

func GetClient() *redis.Client {
	return client
}

const (
	VideoHotKey    = "video:hot"      // 视频热度 zset
	VideoHotExpire = 30 * time.Minute // 热度数据过期时间
	CategoryHotKey = "video:hot:%s"   // 分类下的热门视频
)

// UpdateVideoScore 更新视频分数（同步更新总榜和分类榜）
func UpdateVideoScore(ctx context.Context, videoID int64, visitDelta, likeDelta int64, category string) error {
	score := float64(visitDelta) + float64(likeDelta)*1.5
	videoIDStr := strconv.FormatInt(videoID, 10)

	// 1. 更新总榜
	_, err := client.ZIncrBy(ctx, VideoHotKey, score, videoIDStr).Result()
	if err != nil {
		return err
	}
	client.Expire(ctx, VideoHotKey, VideoHotExpire)

	// 2. 更新分类榜（如果有分类）
	if category != "" {
		categoryKey := fmt.Sprintf(CategoryHotKey, category)
		if err := client.ZIncrBy(ctx, categoryKey, score, videoIDStr).Err(); err != nil {
			return err
		}
		client.Expire(ctx, categoryKey, VideoHotExpire)
	}

	return nil
}

// GetHotVideos 获取热门视频ID列表（支持分类和总榜）
func GetHotVideos(ctx context.Context, category string, limit int, lastVisitCount, lastLikeCount, lastID int64) ([]string, error) {
	key := VideoHotKey // 默认查询总榜
	if category != "" {
		key = fmt.Sprintf(CategoryHotKey, category) // 查询分类榜
	}

	lastScore := float64(lastVisitCount) + float64(lastLikeCount)*1.5

	var _min, _max string
	if lastScore > 0 {
		_min = "-inf"
		_max = fmt.Sprintf("%f", lastScore)
	} else {
		_min = "-inf"
		_max = "+inf"
	}

	// 获取热门视频ID
	videoIDs, err := client.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
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
	if len(videoIDs) > int(limit) {
		videoIDs = videoIDs[:limit]
	}

	return videoIDs, nil
}
