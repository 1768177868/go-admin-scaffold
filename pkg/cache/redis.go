package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore implements Redis-based caching
type RedisStore struct {
	client *redis.Client
	prefix string
}

// NewRedisStore creates a new Redis cache store
func NewRedisStore(config *Config) (*RedisStore, error) {
	redisConfig := config.GetRedisConfig()

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %v", err)
	}

	return &RedisStore{
		client: client,
		prefix: config.Prefix,
	}, nil
}

// Get retrieves a cached value
func (s *RedisStore) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := s.client.Get(ctx, s.getKey(key)).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key not found: %s", key)
	}
	if err != nil {
		return nil, err
	}

	var value interface{}
	if err := json.Unmarshal([]byte(val), &value); err != nil {
		return nil, err
	}

	return value, nil
}

// Set stores a value in cache
func (s *RedisStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, s.getKey(key), data, expiration).Err()
}

// Delete removes a value from cache
func (s *RedisStore) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, s.getKey(key)).Err()
}

// Clear removes all values from cache with the current prefix
func (s *RedisStore) Clear(ctx context.Context) error {
	pattern := s.getKey("*")
	iter := s.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := s.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Has checks if a key exists
func (s *RedisStore) Has(ctx context.Context, key string) bool {
	exists, err := s.client.Exists(ctx, s.getKey(key)).Result()
	return err == nil && exists > 0
}

// Remember gets from cache or stores the result of getter
func (s *RedisStore) Remember(ctx context.Context, key string, expiration time.Duration, getter func() (interface{}, error)) (interface{}, error) {
	// Try to get from cache first
	if value, err := s.Get(ctx, key); err == nil {
		return value, nil
	}

	// Get fresh value
	value, err := getter()
	if err != nil {
		return nil, err
	}

	// Store in cache
	if err := s.Set(ctx, key, value, expiration); err != nil {
		return nil, err
	}

	return value, nil
}

// Increment increments a number value
func (s *RedisStore) Increment(ctx context.Context, key string) error {
	_, err := s.client.Incr(ctx, s.getKey(key)).Result()
	return err
}

// Decrement decrements a number value
func (s *RedisStore) Decrement(ctx context.Context, key string) error {
	_, err := s.client.Decr(ctx, s.getKey(key)).Result()
	return err
}

// Close closes the Redis connection
func (s *RedisStore) Close() error {
	return s.client.Close()
}

// Internal methods

func (s *RedisStore) getKey(key string) string {
	if s.prefix != "" {
		return fmt.Sprintf("%s:%s", s.prefix, key)
	}
	return key
}
