package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements Cache on top of Redis. It's exported (rather than
// returned as the Cache interface) so the composition root can also call
// Ping and Close, which aren't part of the Cache contract services use.
type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) *RedisCache {
	return &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrMiss
	}
	return val, err
}

func (c *RedisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Client exposes the underlying Redis client for other packages built on
// Redis (e.g. locking.RedisLocker) that need more than the Cache
// interface offers.
func (c *RedisCache) Client() *redis.Client {
	return c.client
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}
