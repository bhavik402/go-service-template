package locking

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// releaseScript only deletes the key if it still holds the token this
// process set, so a lock released after its TTL expired (and was
// re-acquired by someone else) never gets deleted out from under them.
const releaseScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end
`

type RedisLocker struct {
	client *redis.Client
}

func NewRedisLocker(client *redis.Client) *RedisLocker {
	return &RedisLocker{client: client}
}

func (l *RedisLocker) Lock(ctx context.Context, key string, ttl time.Duration) (func(context.Context) error, error) {
	token := uuid.NewString()

	ok, err := l.client.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNotAcquired
	}

	release := func(releaseCtx context.Context) error {
		return l.client.Eval(releaseCtx, releaseScript, []string{key}, token).Err()
	}
	return release, nil
}
