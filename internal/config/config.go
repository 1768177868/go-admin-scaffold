package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"app/pkg/i18n"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App        AppConfig        `mapstructure:"app"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Log        LogConfig        `mapstructure:"log"`
	Cache      CacheConfig      `mapstructure:"cache"`
	Queue      QueueConfig      `mapstructure:"queue"`
	I18n       i18n.Config      `mapstructure:"i18n"`
	CORS       CORSConfig       `mapstructure:"cors"`
	Server     ServerConfig     `mapstructure:"server"`
	Storage    StorageConfig    `mapstructure:"storage"`
	SuperAdmin SuperAdminConfig `mapstructure:"super_admin"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Address string `mapstructure:"address"`
	Mode    string `mapstructure:"mode"`
}

// AppConfig holds application configuration
type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Port      int    `mapstructure:"port"`
	APIPrefix string `mapstructure:"api_prefix"`
	Env       string `mapstructure:"env"`
	Debug     bool   `mapstructure:"debug"`
	BaseURL   string `mapstructure:"baseUrl"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Database        string `mapstructure:"database"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Charset         string `mapstructure:"charset"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// LogConfig represents logging configuration
type LogConfig struct {
	Level      string `yaml:"level"`       // 日志级别
	Filename   string `yaml:"filename"`    // 日志文件路径
	MaxSize    int    `yaml:"max_size"`    // 每个日志文件最大尺寸，单位MB
	MaxBackups int    `yaml:"max_backups"` // 保留的旧日志文件最大数量
	MaxAge     int    `yaml:"max_age"`     // 保留的旧日志文件最大天数
	Compress   bool   `yaml:"compress"`    // 是否压缩旧日志文件
	Daily      bool   `yaml:"daily"`       // 是否按天切割日志
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Driver  string                 `mapstructure:"driver"`
	Prefix  string                 `mapstructure:"prefix"`
	Options map[string]interface{} `mapstructure:"options"`
}

// QueueConfig holds queue configuration
type QueueConfig struct {
	Driver     string `mapstructure:"driver"`
	Queue      string `mapstructure:"queue"`
	Connection struct {
		Redis    string `mapstructure:"redis"`
		Database string `mapstructure:"database"`
	} `mapstructure:"connection"`
	Worker struct {
		Sleep   int `mapstructure:"sleep"`
		MaxJobs int `mapstructure:"max_jobs"`
		MaxTime int `mapstructure:"max_time"`
		Rest    int `mapstructure:"rest"`
		Memory  int `mapstructure:"memory"`
		Tries   int `mapstructure:"tries"`
		Timeout int `mapstructure:"timeout"`
	} `mapstructure:"worker"`
	Queues map[string]QueueDetail `mapstructure:"queues"`
}

// QueueDetail holds configuration for individual queues
type QueueDetail struct {
	Priority   int   `mapstructure:"priority"`
	Processes  int   `mapstructure:"processes"`
	Timeout    int   `mapstructure:"timeout"`
	Tries      int   `mapstructure:"tries"`
	RetryAfter int   `mapstructure:"retry_after"`
	Backoff    []int `mapstructure:"backoff"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins     []string      `mapstructure:"allow_origins"`
	AllowMethods     []string      `mapstructure:"allow_methods"`
	AllowHeaders     []string      `mapstructure:"allow_headers"`
	ExposeHeaders    []string      `mapstructure:"expose_headers"`
	AllowCredentials bool          `mapstructure:"allow_credentials"`
	MaxAge           time.Duration `mapstructure:"max_age"`
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	Driver string      `mapstructure:"driver"`
	Local  LocalConfig `mapstructure:"local"`
	S3     S3Config    `mapstructure:"s3"`
}

// LocalConfig holds local storage configuration
type LocalConfig struct {
	Path string `mapstructure:"path"`
}

// S3Config holds S3 storage configuration
type S3Config struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Bucket          string `mapstructure:"bucket"`
	Region          string `mapstructure:"region"`
	UseSSL          bool   `mapstructure:"use_ssl"`
}

// SuperAdminConfig holds SuperAdmin configuration
type SuperAdminConfig struct {
	UserIDs []string `mapstructure:"user_ids"`
}

// Load loads configuration from environment variables and config files
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Set default configuration file paths
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 配置文件搜索路径
	viper.AddConfigPath("./configs") // 首选路径
	viper.AddConfigPath(".")         // 当前目录
	viper.AddConfigPath("/etc/app/") // 系统配置目录

	// Load configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %v", err)
		}
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Environment variables take precedence over config file
	// App
	config.App.Name = getEnvOrDefault("APP_NAME", viper.GetString("app.name"))
	config.App.Env = getEnvOrDefault("APP_ENV", viper.GetString("app.env"))
	config.App.Debug = getEnvBoolOrDefault("APP_DEBUG", viper.GetBool("app.debug"))
	config.App.BaseURL = getEnvOrDefault("APP_URL", viper.GetString("app.baseUrl"))
	config.App.Port = getEnvIntOrDefault("APP_PORT", viper.GetInt("app.port"))
	config.App.APIPrefix = getEnvOrDefault("APP_API_PREFIX", viper.GetString("app.api_prefix"))

	// JWT
	config.JWT.Secret = getEnvOrDefault("JWT_SECRET", viper.GetString("jwt.secret"))
	config.JWT.ExpireTime = getEnvIntOrDefault("JWT_EXPIRE", viper.GetInt("jwt.expire_time"))

	// Database
	config.Database.Driver = getEnvOrDefault("DB_DRIVER", viper.GetString("database.driver"))
	config.Database.Host = getEnvOrDefault("DB_HOST", viper.GetString("database.host"))
	config.Database.Port = getEnvOrDefault("DB_PORT", viper.GetString("database.port"))
	config.Database.Username = getEnvOrDefault("DB_USERNAME", viper.GetString("database.username"))
	config.Database.Password = getEnvOrDefault("DB_PASSWORD", viper.GetString("database.password"))
	config.Database.Database = getEnvOrDefault("DB_DATABASE", viper.GetString("database.database"))
	config.Database.Charset = getEnvOrDefault("DB_CHARSET", viper.GetString("database.charset"))
	config.Database.MaxIdleConns = getEnvIntOrDefault("DB_MAX_IDLE_CONNS", viper.GetInt("database.max_idle_conns"))
	config.Database.MaxOpenConns = getEnvIntOrDefault("DB_MAX_OPEN_CONNS", viper.GetInt("database.max_open_conns"))
	config.Database.ConnMaxLifetime = getEnvIntOrDefault("DB_CONN_MAX_LIFETIME", viper.GetInt("database.conn_max_lifetime"))

	// Redis
	config.Redis.Host = getEnvOrDefault("REDIS_HOST", viper.GetString("redis.host"))
	config.Redis.Port = getEnvOrDefault("REDIS_PORT", viper.GetString("redis.port"))
	config.Redis.Password = getEnvOrDefault("REDIS_PASSWORD", viper.GetString("redis.password"))
	config.Redis.DB = getEnvIntOrDefault("REDIS_DB", viper.GetInt("redis.db"))

	// Cache
	config.Cache.Driver = getEnvOrDefault("CACHE_DRIVER", viper.GetString("cache.driver"))
	config.Cache.Prefix = getEnvOrDefault("CACHE_PREFIX", viper.GetString("cache.prefix"))
	config.Cache.Options = viper.GetStringMap("cache.options")

	// Queue
	config.Queue.Driver = getEnvOrDefault("QUEUE_DRIVER", viper.GetString("queue.driver"))
	config.Queue.Queue = getEnvOrDefault("QUEUE_NAME", viper.GetString("queue.queue"))

	// Queue connection
	config.Queue.Connection.Redis = viper.GetString("queue.connection.redis")
	config.Queue.Connection.Database = viper.GetString("queue.connection.database")

	// Queue worker
	config.Queue.Worker.Sleep = viper.GetInt("queue.worker.sleep")
	config.Queue.Worker.MaxJobs = viper.GetInt("queue.worker.max_jobs")
	config.Queue.Worker.MaxTime = viper.GetInt("queue.worker.max_time")
	config.Queue.Worker.Rest = viper.GetInt("queue.worker.rest")
	config.Queue.Worker.Memory = viper.GetInt("queue.worker.memory")
	config.Queue.Worker.Tries = viper.GetInt("queue.worker.tries")
	config.Queue.Worker.Timeout = viper.GetInt("queue.worker.timeout")

	// Queue details
	config.Queue.Queues = make(map[string]QueueDetail)
	queues := viper.GetStringMap("queue.queues")
	for name := range queues {
		var detail QueueDetail
		if err := viper.UnmarshalKey("queue.queues."+name, &detail); err != nil {
			return nil, fmt.Errorf("error unmarshaling queue %s: %v", name, err)
		}
		config.Queue.Queues[name] = detail
	}

	// Server
	config.Server.Address = getEnvOrDefault("SERVER_ADDRESS", viper.GetString("server.address"))
	config.Server.Mode = getEnvOrDefault("SERVER_MODE", viper.GetString("server.mode"))

	// Log
	config.Log.Level = getEnvOrDefault("LOG_LEVEL", viper.GetString("log.level"))
	config.Log.Filename = getEnvOrDefault("LOG_FILENAME", viper.GetString("log.filename"))
	config.Log.MaxSize = getEnvIntOrDefault("LOG_MAX_SIZE", viper.GetInt("log.maxSize"))
	config.Log.MaxBackups = getEnvIntOrDefault("LOG_MAX_BACKUPS", viper.GetInt("log.maxBackups"))
	config.Log.MaxAge = getEnvIntOrDefault("LOG_MAX_AGE", viper.GetInt("log.maxAge"))
	config.Log.Compress = getEnvBoolOrDefault("LOG_COMPRESS", viper.GetBool("log.compress"))
	config.Log.Daily = getEnvBoolOrDefault("LOG_DAILY", viper.GetBool("log.daily"))

	// CORS
	config.CORS.AllowOrigins = viper.GetStringSlice("cors.allow_origins")
	config.CORS.AllowMethods = viper.GetStringSlice("cors.allow_methods")
	config.CORS.AllowHeaders = viper.GetStringSlice("cors.allow_headers")
	config.CORS.ExposeHeaders = viper.GetStringSlice("cors.expose_headers")
	config.CORS.AllowCredentials = viper.GetBool("cors.allow_credentials")
	config.CORS.MaxAge = viper.GetDuration("cors.max_age")

	// I18n
	config.I18n.DefaultLocale = getEnvOrDefault("I18N_DEFAULT_LOCALE", viper.GetString("i18n.default_locale"))
	config.I18n.LoadPath = getEnvOrDefault("I18N_LOAD_PATH", viper.GetString("i18n.load_path"))
	config.I18n.AvailableLocales = viper.GetStringSlice("i18n.available_locales")

	// Storage
	config.Storage.Driver = getEnvOrDefault("STORAGE_DRIVER", viper.GetString("storage.driver"))

	// Local storage
	config.Storage.Local.Path = getEnvOrDefault("STORAGE_LOCAL_PATH", viper.GetString("storage.local.path"))

	// S3 storage
	config.Storage.S3.Endpoint = getEnvOrDefault("STORAGE_S3_ENDPOINT", viper.GetString("storage.s3.endpoint"))
	config.Storage.S3.AccessKeyID = getEnvOrDefault("STORAGE_S3_ACCESS_KEY_ID", viper.GetString("storage.s3.access_key_id"))
	config.Storage.S3.SecretAccessKey = getEnvOrDefault("STORAGE_S3_SECRET_ACCESS_KEY", viper.GetString("storage.s3.secret_access_key"))
	config.Storage.S3.Bucket = getEnvOrDefault("STORAGE_S3_BUCKET", viper.GetString("storage.s3.bucket"))
	config.Storage.S3.Region = getEnvOrDefault("STORAGE_S3_REGION", viper.GetString("storage.s3.region"))
	config.Storage.S3.UseSSL = getEnvBoolOrDefault("STORAGE_S3_USE_SSL", viper.GetBool("storage.s3.use_ssl"))

	// SuperAdmin
	superAdminIDs := viper.GetStringSlice("super_admin.user_ids")
	for _, idStr := range superAdminIDs {
		config.SuperAdmin.UserIDs = append(config.SuperAdmin.UserIDs, idStr)
	}

	return config, nil
}

// getEnvOrDefault gets environment variable value or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault gets environment variable as int or returns default value
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBoolOrDefault gets environment variable as bool or returns default value
func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// ParseSuperAdminIDs converts string IDs to uint
func (c *Config) ParseSuperAdminIDs() []uint {
	var ids []uint
	for _, idStr := range c.SuperAdmin.UserIDs {
		if id, err := strconv.ParseUint(idStr, 10, 32); err == nil {
			ids = append(ids, uint(id))
		}
	}
	return ids
}
