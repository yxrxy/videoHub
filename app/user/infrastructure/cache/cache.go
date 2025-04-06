package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/yxrxy/videoHub/app/user/domain/repository"
	"github.com/yxrxy/videoHub/pkg/constants"
)

type userCache struct {
	client *redis.Client
}

func NewUserCache(client *redis.Client) repository.UserCache {
	return &userCache{client: client}
}

func (c *userCache) IsExist(ctx context.Context, key string) bool {
	return c.client.Exists(ctx, key).Val() == 1
}

func (c *userCache) SetUserAccessToken(ctx context.Context, key string, token string) error {
	return c.client.Set(ctx, key, token, constants.TokenExpiry).Err()
}

func (c *userCache) SetUserRefreshToken(ctx context.Context, key string, token string) error {
	return c.client.Set(ctx, key, token, constants.TokenExpiry2).Err()
}

func (c *userCache) DeleteUserToken(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *userCache) GetToken(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}
