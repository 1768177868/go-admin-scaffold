package redis

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

// Config represents Redis configuration
type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// Setup initializes the Redis client
func Setup(cfg *Config) error {
	var err error
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Host + ":" + cfg.Port,
			Password: cfg.Password,
			DB:       cfg.DB,
		})

		// Test connection
		err = client.Ping(context.Background()).Err()
	})
	return err
}

// GetClient returns the Redis client instance
func GetClient() *redis.Client {
	if client == nil {
		// If client is not initialized, create a default one
		client = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}
	return client
}

// Close closes the Redis client connection
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
