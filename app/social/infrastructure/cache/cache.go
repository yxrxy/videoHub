package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/yxrxy/videoHub/app/social/domain/repository"
)

type socialCache struct {
	client *redis.Client
}

func NewSocialCache(client *redis.Client) repository.SocialCache {
	return &socialCache{client: client}
}

func (c *socialCache) IsUserOnline(ctx context.Context, userID int64) (bool, error) {
	return c.client.SIsMember(ctx, "online_users", userID).Result()
}

func (c *socialCache) SetUserOnline(ctx context.Context, userID int64) error {
	return c.client.SAdd(ctx, "online_users", userID).Err()
}

func (c *socialCache) SetUserOffline(ctx context.Context, userID int64) error {
	return c.client.SRem(ctx, "online_users", userID).Err()
}
