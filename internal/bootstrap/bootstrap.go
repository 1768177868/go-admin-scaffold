package bootstrap

import (
	"app/internal/config"
	"app/pkg/cache"
	"app/pkg/redis"
)

// SetupRedis initializes the Redis connection
func SetupRedis(cfg *config.Config) error {
	return redis.Setup(&redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
}

// SetupCache initializes the cache system
func SetupCache(cfg *config.Config) error {
	return cache.Setup(&cache.Config{
		Driver:  cfg.Cache.Driver,
		Prefix:  cfg.Cache.Prefix,
		Options: cfg.Cache.Options,
	})
}
