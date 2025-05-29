# 配置系统

本框架采用简洁的YAML配置系统，便于管理和部署。

## 配置文件

### 配置文件结构

项目只使用两个配置文件：

- `configs/config.yaml` - 实际配置文件（生产使用）
- `configs/config.example.yaml` - 配置示例文件（模板）

### 配置示例

`configs/config.example.yaml`：

```yaml
app:
  name: "Go Admin"
  env: "development"  # development, production, test
  mode: "development"
  debug: true
  baseUrl: "http://localhost:8080"
  api_prefix: "/api/v1"
  port: 8080

server:
  address: "0.0.0.0:8080"
  mode: "debug"  # debug, release, test

mysql:
  host: "localhost"
  port: 3306
  username: "root"
  password: ""
  database: "go_admin"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600  # seconds

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-secret-key-here"  # Change this in production
  expire_time: 86400  # 24 hours

log:
  level: "debug"  # debug, info, warn, error
  filename: "storage/logs/app.log"
  max_size: 100    # megabytes
  max_backups: 3
  max_age: 28      # days

cors:
  allow_origins: ["*"]  # Use specific domains in production
  allow_methods: ["GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"]
  allow_headers: ["Origin", "Content-Type", "Accept", "Authorization"]
  allow_credentials: true
  max_age: 86400  # seconds

i18n:
  default_locale: "en"
  load_path: "./locales"
  available_locales: ["en", "zh"]

storage:
  driver: "local"  # local, s3
  options:
    local:
      path: "storage/uploads"
    s3:
      region: "us-west-2"
      bucket: "your-bucket"
      access_key: ""
      secret_key: ""
```

## 配置使用

### 初始化配置

```bash
# 复制配置示例文件
cp configs/config.example.yaml configs/config.yaml

# 根据环境修改配置
vim configs/config.yaml
```

### 在代码中访问配置

```go
package main

import "app/internal/config"

func main() {
    // 加载配置
    cfg, err := config.LoadConfig()
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

## 环境特定配置

针对不同环境，直接修改 `configs/config.yaml` 中的相应字段：

### 开发环境

```yaml
app:
  env: "development"
  mode: "development"
  debug: true

server:
  mode: "debug"

log:
  level: "debug"
```

### 生产环境

```yaml
app:
  env: "production"
  mode: "production"
  debug: false
  baseUrl: "https://your-domain.com"

server:
  mode: "release"

mysql:
  host: "prod-db-host"
  username: "prod-user"
  password: "secure-password"
  database: "prod_database"

jwt:
  secret: "production-secret-key"

log:
  level: "info"
  filename: "/var/log/go-admin/app.log"

cors:
  allow_origins: ["https://your-domain.com"]
```

## 配置验证

配置加载时会进行基本验证：

```go
func (c *Config) Validate() error {
    if c.App.Name == "" {
        return errors.New("app.name is required")
    }
    
    if c.MySQL.Host == "" {
        return errors.New("mysql.host is required")
    }
    
    if c.JWT.Secret == "" || c.JWT.Secret == "your-secret-key-here" {
        return errors.New("jwt.secret must be set to a secure value")
    }
    
    return nil
}
```

## 敏感信息处理

1. **不要提交包含敏感信息的 config.yaml**
2. **生产环境配置建议**：
   - 使用强密码和随机密钥
   - 限制数据库用户权限
   - 使用HTTPS
   - 设置具体的CORS域名

3. **密钥生成**：
```bash
# 生成JWT密钥
openssl rand -base64 32

# 生成App密钥
openssl rand -hex 32
```

## 部署最佳实践

1. **配置文件管理**：
   - 在服务器上直接创建 `config.yaml`
   - 或通过配置管理工具部署
   - 确保文件权限安全（600）

2. **不同环境的配置**：
```bash
# 开发环境
configs/config.yaml  # 开发配置

# 测试环境  
configs/config.yaml  # 测试配置

# 生产环境
configs/config.yaml  # 生产配置
```

3. **配置备份**：
```bash
# 备份当前配置
cp configs/config.yaml configs/config.backup.yaml

# 恢复配置
cp configs/config.backup.yaml configs/config.yaml
```

这种简化的配置系统让部署和管理变得更加直观和可靠。 