package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(c *redis.Client) *Cache {
	return &Cache{client: c}
}

func (c *Cache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if expiration <= 0 {
		return c.client.Set(ctx, key, value, 0).Err()
	}
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	res, err := c.client.Exists(ctx, key).Result()
	return res > 0, err
}

func (c *Cache) Evaluate(ctx context.Context, hashKey, script string, keys []string, args []any) (any, error) {
	hash, err := c.Get(ctx, hashKey)
	if err != nil {
		hash, err = c.load(ctx, hashKey, script)
		if err != nil {
			return nil, err
		}
	}
	return c.client.EvalSha(ctx, hash, keys, args...).Result()
}

func (c *Cache) load(ctx context.Context, key, script string) (string, error) {
	res, err := c.client.ScriptLoad(ctx, script).Result()
	if err != nil {
		return "", err
	}
	if err := c.Set(ctx, key, res, -1); err != nil {
		return "", err
	}
	return res, nil
}
