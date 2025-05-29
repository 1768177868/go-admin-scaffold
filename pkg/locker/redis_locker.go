package locker

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLocker struct {
	client *redis.Client
}

func NewRedisLocker(client *redis.Client) *RedisLocker {
	return &RedisLocker{
		client: client,
	}
}

// TryLock attempts to acquire a lock with a given key and TTL
// Returns true if lock is acquired, false otherwise
func (l *RedisLocker) TryLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	// Use SET NX to ensure atomic lock acquisition
	ok, err := l.client.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

// Unlock releases a lock
func (l *RedisLocker) Unlock(ctx context.Context, key string) error {
	return l.client.Del(ctx, key).Err()
}

// RefreshLock extends the lock TTL
func (l *RedisLocker) RefreshLock(ctx context.Context, key string, ttl time.Duration) error {
	return l.client.Expire(ctx, key, ttl).Err()
}
