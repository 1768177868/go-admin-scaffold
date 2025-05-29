package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"app/pkg/i18n"

	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Redis  RedisConfig  `mapstructure:"redis"`
	Log    LogConfig    `mapstructure:"log"`
	Cache  CacheConfig  `mapstructure:"cache"`
	Queue  QueueConfig  `mapstructure:"queue"`
	I18n   i18n.Config  `mapstructure:"i18n"`
	CORS   CORSConfig   `mapstructure:"cors"`
	Server ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	Address string `mapstructure:"address"`
	Mode    string `mapstructure:"mode"`
}

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Port      int    `mapstructure:"port"`
	APIPrefix string `mapstructure:"api_prefix"`
	Env       string `mapstructure:"env"`
	Debug     bool   `mapstructure:"debug"`
	BaseURL   string `mapstructure:"baseUrl"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"`
}

type MySQLConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Database        string `mapstructure:"database"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type CacheConfig struct {
	Driver  string                 `mapstructure:"driver"`
	Prefix  string                 `mapstructure:"prefix"`
	Options map[string]interface{} `mapstructure:"options"`
}

type QueueConfig struct {
	Driver  string                 `mapstructure:"driver"`
	Options map[string]interface{} `mapstructure:"options"`
}

type CORSConfig struct {
	AllowOrigins     []string      `mapstructure:"allow_origins"`
	AllowMethods     []string      `mapstructure:"allow_methods"`
	AllowHeaders     []string      `mapstructure:"allow_headers"`
	ExposeHeaders    []string      `mapstructure:"expose_headers"`
	AllowCredentials bool          `mapstructure:"allow_credentials"`
	MaxAge           time.Duration `mapstructure:"max_age"`
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

	// MySQL
	config.MySQL.Host = getEnvOrDefault("DB_HOST", viper.GetString("mysql.host"))
	config.MySQL.Port = getEnvIntOrDefault("DB_PORT", viper.GetInt("mysql.port"))
	config.MySQL.Username = getEnvOrDefault("DB_USERNAME", viper.GetString("mysql.username"))
	config.MySQL.Password = getEnvOrDefault("DB_PASSWORD", viper.GetString("mysql.password"))
	config.MySQL.Database = getEnvOrDefault("DB_DATABASE", viper.GetString("mysql.database"))

	// Redis
	config.Redis.Host = getEnvOrDefault("REDIS_HOST", viper.GetString("redis.host"))
	config.Redis.Port = getEnvIntOrDefault("REDIS_PORT", viper.GetInt("redis.port"))
	config.Redis.Password = getEnvOrDefault("REDIS_PASSWORD", viper.GetString("redis.password"))
	config.Redis.DB = getEnvIntOrDefault("REDIS_DB", viper.GetInt("redis.db"))

	// Server
	config.Server.Address = getEnvOrDefault("SERVER_ADDRESS", viper.GetString("server.address"))
	config.Server.Mode = getEnvOrDefault("SERVER_MODE", viper.GetString("server.mode"))

	// Log
	config.Log.Level = getEnvOrDefault("LOG_LEVEL", viper.GetString("log.level"))
	config.Log.Filename = getEnvOrDefault("LOG_FILENAME", viper.GetString("log.filename"))
	config.Log.MaxSize = getEnvIntOrDefault("LOG_MAX_SIZE", viper.GetInt("log.maxSize"))
	config.Log.MaxBackups = getEnvIntOrDefault("LOG_MAX_BACKUPS", viper.GetInt("log.maxBackups"))
	config.Log.MaxAge = getEnvIntOrDefault("LOG_MAX_AGE", viper.GetInt("log.maxAge"))

	// CORS
	config.CORS.AllowOrigins = viper.GetStringSlice("cors.allow_origins")
	config.CORS.AllowMethods = viper.GetStringSlice("cors.allow_methods")
	config.CORS.AllowHeaders = viper.GetStringSlice("cors.allow_headers")
	config.CORS.ExposeHeaders = viper.GetStringSlice("cors.expose_headers")
	config.CORS.AllowCredentials = viper.GetBool("cors.allow_credentials")
	config.CORS.MaxAge = viper.GetDuration("cors.max_age")

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
