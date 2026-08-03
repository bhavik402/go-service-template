// Package cache is a read-through cache abstraction backed by Redis.
// It is never a source of truth — callers always fall back to a
// repository on a miss. See internal/locking for Redis used as a
// distributed lock instead.
package cache

import (
	"context"
	"errors"
	"time"
)

var ErrMiss = errors.New("cache: key not found")

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
