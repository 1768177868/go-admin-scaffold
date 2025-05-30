package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"app/pkg/redis"
)

var (
	ErrKeyNotFound = errors.New("key not found in cache")
	ErrKeyExpired  = errors.New("key has expired")
)

// Cache interface defines the methods that any cache implementation must provide
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// RedisCache implements Cache interface using Redis
type RedisCache struct {
	prefix string
}

// Config represents cache configuration
type Config struct {
	Driver  string                 `mapstructure:"driver"`
	Prefix  string                 `mapstructure:"prefix"`
	Options map[string]interface{} `mapstructure:"options"`
}

var (
	defaultCache Cache
)

// Setup initializes the cache system
func Setup(cfg *Config) error {
	switch cfg.Driver {
	case "redis":
		defaultCache = &RedisCache{
			prefix: cfg.Prefix,
		}
		return nil
	default:
		return fmt.Errorf("unsupported cache driver: %s", cfg.Driver)
	}
}

// Default returns the default cache instance
func Default() Cache {
	return defaultCache
}

// Get retrieves a value from Redis cache
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return redis.GetClient().Get(ctx, c.prefix+key).Result()
}

// Set stores a value in Redis cache
func (c *RedisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return redis.GetClient().Set(ctx, c.prefix+key, value, expiration).Err()
}

// Delete removes a value from Redis cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return redis.GetClient().Del(ctx, c.prefix+key).Err()
}

// Exists checks if a key exists in Redis cache
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := redis.GetClient().Exists(ctx, c.prefix+key).Result()
	return result > 0, err
}

// Store defines the interface for cache implementations
type Store interface {
	// Get retrieves a value by key
	Get(ctx context.Context, key string) (interface{}, error)

	// Set stores a value by key
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Delete removes a key
	Delete(ctx context.Context, key string) error

	// Clear removes all keys
	Clear(ctx context.Context) error

	// Remember gets from cache or stores the result of getter
	Remember(ctx context.Context, key string, expiration time.Duration, getter func() (interface{}, error)) (interface{}, error)

	// Has checks if a key exists
	Has(ctx context.Context, key string) bool

	// Increment increments a number value
	Increment(ctx context.Context, key string) error

	// Decrement decrements a number value
	Decrement(ctx context.Context, key string) error

	// Close closes the cache store
	Close() error
}

// GetFilePath returns the file path for file cache
func (c *Config) GetFilePath() string {
	if path, ok := c.Options["file_path"].(string); ok {
		return path
	}
	return "storage/cache" // default path
}

// GetRedisConfig returns Redis configuration from options
func (c *Config) GetRedisConfig() RedisConfig {
	config := RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}

	if host, ok := c.Options["host"].(string); ok {
		config.Host = host
	}
	if port, ok := c.Options["port"].(int); ok {
		config.Port = port
	}
	if password, ok := c.Options["password"].(string); ok {
		config.Password = password
	}
	if db, ok := c.Options["db"].(int); ok {
		config.DB = db
	}

	return config
}

// RedisConfig represents Redis cache configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

var (
	DefaultExpiration = 24 * time.Hour
	NoExpiration      = time.Duration(0)
)

// Manager manages cache stores
type Manager struct {
	config *Config
	stores map[string]Store
}

// NewManager creates a new cache manager
func NewManager(config *Config) *Manager {
	return &Manager{
		config: config,
		stores: make(map[string]Store),
	}
}

// Store gets or creates a cache store
func (m *Manager) Store(driver string) (Store, error) {
	if store, exists := m.stores[driver]; exists {
		return store, nil
	}

	var store Store
	var err error

	switch driver {
	case "file":
		store, err = NewFileStore(m.config)
	case "redis":
		store, err = NewRedisStore(m.config)
	default:
		store, err = NewFileStore(m.config) // Default to file store
	}

	if err != nil {
		return nil, err
	}

	m.stores[driver] = store
	return store, nil
}

// Close closes all cache stores
func (m *Manager) Close() error {
	for _, store := range m.stores {
		if err := store.Close(); err != nil {
			return err
		}
	}
	return nil
}
