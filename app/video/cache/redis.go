package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
)

func Init(addr string) {
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // 如果有密码，在这里设置
		DB:       0,  // 使用默认 DB
	})
}

func GetClient() *redis.Client {
	return client
}

// 视频热度相关的 key
const (
	VideoHotKey    = "video:hot"      // 视频热度 zset
	VideoHotExpire = 30 * time.Minute // 热度数据过期时间
	VideoScoreKey  = "video:score:%d" // 单个视频分数
	CategoryHotKey = "video:hot:%s"   // 分类下的热门视频
)

// UpdateVideoScore 更新视频分数
func UpdateVideoScore(ctx context.Context, videoID int64, visitDelta, likeDelta int64) error {
	score := float64(visitDelta) + float64(likeDelta)*1.5 // 点赞权重1.5倍

	// 更新总榜单
	_, err := client.ZIncrBy(ctx, VideoHotKey, score, string(videoID)).Result()
	if err != nil {
		return err
	}

	// 更新过期时间
	client.Expire(ctx, VideoHotKey, VideoHotExpire)
	return nil
}

// GetHotVideos 获取热门视频ID列表
func GetHotVideos(ctx context.Context, category string, limit int, cursor int64) ([]string, error) {
	key := VideoHotKey
	if category != "" {
		key = CategoryHotKey + category
	}

	opt := &redis.ZRangeBy{
		Min:    "-inf",
		Max:    string(cursor),
		Offset: 0,
		Count:  int64(limit),
	}

	// 按分数从高到低获取
	return client.ZRevRangeByScore(ctx, key, opt).Result()
}
