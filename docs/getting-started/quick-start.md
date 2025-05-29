# 快速开始

## 环境要求

- Go 1.23 或更高版本
- MySQL 5.7 或更高版本
- Redis 6.0 或更高版本

## 安装

1. 克隆项目：
```bash
git clone https://github.com/yourusername/go-admin.git
cd go-admin
```

2. 安装依赖：
```bash
go mod download
```

3. 复制配置文件：
```bash
cp config/config.example.yaml config/config.yaml
```

4. 修改配置文件：
```yaml
mysql:
  host: "localhost"
  port: 3306
  username: "your-username"
  password: "your-password"
  database: "go_admin"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

5. 运行数据库迁移：
```bash
go run cmd/tools/main.go migrate
```

6. 填充初始数据：
```bash
go run cmd/tools/main.go db:seed
```

7. 启动服务：
```bash
go run cmd/server/main.go
```

## 默认账户

- 用户名：admin
- 密码：admin123
- 邮箱：admin@example.com

## 基本使用

### 1. 访问管理面板

打开浏览器访问：`http://localhost:8080/admin`

### 2. API 文档

访问 Swagger 文档：`http://localhost:8080/swagger/index.html`

### 3. 命令行工具

查看可用命令：
```bash
go run cmd/tools/main.go
```

常用命令：
```bash
# 创建新的迁移文件
go run cmd/tools/main.go make:migration create_users_table

# 创建新的模型
go run cmd/tools/main.go make:model User

# 创建新的控制器
go run cmd/tools/main.go make:controller UserController
```

### 4. 定时任务

编辑 `internal/schedule/kernel.go` 添加定时任务：
```go
func (k *Kernel) Schedule() {
    // 每天凌晨执行
    k.scheduler.Command("backup:run").Daily().At("00:00").Unique().Register()
    
    // 每30分钟执行
    k.scheduler.Command("cache:clear").EveryThirtyMinutes().Register()
}
```

## 下一步

- 阅读 [项目结构](structure.md) 了解代码组织
- 查看 [配置说明](configuration.md) 了解更多配置选项
- 浏览 [功能特性](../features) 了解所有功能
- 参考 [开发指南](../development/standards.md) 开始开发 