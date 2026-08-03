// Package locking provides distributed locks (Redis SET NX EX) for
// coordinating work across multiple instances of the service — e.g.
// preventing two replicas from processing the same background job at
// once.
package locking

import (
	"context"
	"errors"
	"time"
)

var ErrNotAcquired = errors.New("locking: lock not acquired")

// Locker acquires a distributed lock identified by key. It returns a
// release function that must be called to unlock (typically deferred).
type Locker interface {
	Lock(ctx context.Context, key string, ttl time.Duration) (release func(context.Context) error, err error)
}
