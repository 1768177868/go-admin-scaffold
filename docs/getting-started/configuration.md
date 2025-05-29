# 配置系统

本框架采用分层配置系统，支持多环境配置、环境变量覆盖和默认值。

## 配置文件

### 基本配置文件

主配置文件位于 `configs/config.yaml`：

```yaml
app:
  name: "Go Admin"
  env: "development"  # development, production, test
  debug: true
  base_url: "http://localhost:8080"
  api_prefix: "/api/v1"
  port: 8080
  key: "your-app-key"
  timezone: "Asia/Shanghai"

http:
  read_timeout: 60
  write_timeout: 60
  idle_timeout: 60
  shutdown_timeout: 30

mysql:
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "go_admin"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 100
  min_idle_conns: 10

log:
  level: "debug"  # debug, info, warn, error
  format: "json"  # json, text
  output: "storage/logs/app.log"
  max_size: 100    # MB
  max_backups: 3
  max_age: 28      # days
  compress: true

jwt:
  secret: "your-jwt-secret"
  ttl: 86400       # seconds
  refresh_ttl: 604800

mail:
  host: "smtp.mailtrap.io"
  port: 2525
  username: "your-username"
  password: "your-password"
  from_address: "admin@example.com"
  from_name: "Go Admin"
```

### 环境配置

1. 开发环境：`config/config.development.yaml`
2. 测试环境：`config/config.testing.yaml`
3. 生产环境：`config/config.production.yaml`

## 使用配置

### 在代码中访问配置

```go
package main

import "app/internal/config"

func main() {
    // 加载配置
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    // 访问配置值
    appName := cfg.App.Name
    dbHost := cfg.MySQL.Host
    redisPort := cfg.Redis.Port
}
```

### 配置注入

使用依赖注入在服务中使用配置：

```go
type UserService struct {
    config *config.Config
    db     *gorm.DB
}

func NewUserService(cfg *config.Config, db *gorm.DB) *UserService {
    return &UserService{
        config: cfg,
        db:     db,
    }
}
```

## 环境变量

配置值可以通过环境变量覆盖。环境变量的优先级高于配置文件：

```bash
# 应用配置
export APP_NAME=MyAdmin
export APP_ENV=production
export APP_DEBUG=false
export APP_URL=http://admin.example.com

# 数据库配置
export DB_HOST=db.example.com
export DB_PORT=3306
export DB_DATABASE=admin
export DB_USERNAME=root
export DB_PASSWORD=secret

# Redis 配置
export REDIS_HOST=redis.example.com
export REDIS_PORT=6379
export REDIS_PASSWORD=secret
export REDIS_DB=0

# JWT 配置
export JWT_SECRET=your-secret-key
export JWT_EXPIRE=24

# 日志配置
export LOG_LEVEL=info
export LOG_FILENAME=storage/logs/app.log
```

## 配置加载顺序

配置系统按以下顺序加载和合并配置：

1. 默认配置
2. 配置文件（configs/config.yaml）
3. 环境变量

后加载的配置会覆盖先加载的配置。

## 配置缓存

在生产环境中，配置会被缓存以提高性能：

```bash
# 生成配置缓存
go run cmd/tools/main.go config:cache

# 清除配置缓存
go run cmd/tools/main.go config:clear
```

## 敏感信息处理

1. 不要在配置文件中存储敏感信息
2. 使用环境变量存储敏感数据
3. 使用加密配置存储敏感数据：

```go
// 加密配置值
func (c *Config) EncryptValue(key string, value string) error {
    encrypted, err := crypto.Encrypt(value, c.App.Key)
    if err != nil {
        return err
    }
    return c.Set(key, encrypted)
}

// 解密配置值
func (c *Config) DecryptValue(key string) (string, error) {
    value := c.Get(key)
    return crypto.Decrypt(value, c.App.Key)
}
```

## 配置验证

配置加载时会进行验证：

```go
func (c *Config) Validate() error {
    if c.App.Name == "" {
        return errors.New("app name is required")
    }
    if c.App.Key == "" {
        return errors.New("app key is required")
    }
    if c.MySQL.Host == "" {
        return errors.New("database host is required")
    }
    // ... 其他验证
    return nil
}
```

## 最佳实践

1. 配置分组：
   - 按功能模块组织配置
   - 使用有意义的配置键名
   - 保持配置结构清晰

2. 环境管理：
   - 每个环境使用独立的配置文件
   - 使用环境变量覆盖敏感配置
   - 不同环境使用不同的密钥

3. 安全性：
   - 不要提交敏感配置到版本控制
   - 使用环境变量或加密存储敏感信息
   - 定期轮换密钥和凭证

4. 维护性：
   - 记录配置变更
   - 保持配置文件的版本控制
   - 定期审查和清理未使用的配置 