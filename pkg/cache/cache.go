package cache

import (
	"context"
	"errors"
	"time"
)

var (
	ErrKeyNotFound = errors.New("key not found in cache")
	ErrKeyExpired  = errors.New("key has expired")
)

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

// Config represents cache configuration
type Config struct {
	Driver  string                 `yaml:"driver"`  // Cache driver (file/redis)
	Prefix  string                 `yaml:"prefix"`  // Key prefix
	Options map[string]interface{} `yaml:"options"` // Driver-specific options
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
