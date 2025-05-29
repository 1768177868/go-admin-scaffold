# 项目结构

本框架采用清晰的目录结构，遵循 Go 项目最佳实践和领域驱动设计（DDD）原则。

## 目录结构

```
.
├── cmd/                    # 命令行入口
│   ├── server/            # HTTP 服务器
│   ├── tools/             # 数据库工具
│   └── artisan/           # 命令行工具
├── config/                # 配置文件
│   ├── config.yaml        # 主配置文件
│   └── config.example.yaml # 配置文件示例
├── internal/              # 内部包
│   ├── commands/         # 命令行命令
│   ├── config/          # 配置加载
│   ├── controllers/     # HTTP 控制器
│   ├── middleware/      # HTTP 中间件
│   ├── models/          # 数据模型
│   ├── repositories/    # 数据仓储
│   ├── schedule/        # 任务调度
│   └── services/        # 业务服务
├── pkg/                  # 公共包
│   ├── cache/           # 缓存组件
│   ├── console/         # 命令行框架
│   ├── database/        # 数据库组件
│   ├── logger/          # 日志组件
│   └── utils/           # 工具函数
├── web/                 # 前端资源
│   ├── static/         # 静态文件
│   └── templates/      # 模板文件
├── docs/                # 文档
├── tests/               # 测试文件
├── go.mod              # Go 模块文件
├── go.sum              # Go 依赖版本
└── README.md           # 项目说明
```

## 核心目录说明

### cmd/

包含项目的主要入口点：

- `server/`: HTTP 服务器入口
- `tools/`: 数据库迁移和填充工具
- `artisan/`: 命令行工具入口

### config/

配置文件目录：

```yaml
# config/config.yaml
app:
  name: "Go Admin"
  env: "development"
  debug: true
  port: 8080

mysql:
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "go_admin"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

### internal/

内部应用代码，不对外暴露：

#### commands/

命令行命令实现：

```go
// internal/commands/make_command.go
type MakeCommand struct {
    *console.BaseCommand
}

func (c *MakeCommand) Handle(ctx context.Context) error {
    // 命令实现
    return nil
}
```

#### controllers/

HTTP 控制器：

```go
// internal/controllers/user_controller.go
type UserController struct {
    userService *services.UserService
}

func (c *UserController) List(ctx *gin.Context) {
    // 处理请求
}
```

#### middleware/

HTTP 中间件：

```go
// internal/middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 中间件逻辑
    }
}
```

#### models/

数据模型定义：

```go
// internal/models/user.go
type User struct {
    ID        uint      `gorm:"primarykey"`
    Name      string    `gorm:"size:255;not null"`
    Email     string    `gorm:"size:255;not null;unique"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### services/

业务逻辑实现：

```go
// internal/services/user_service.go
type UserService struct {
    repo *repositories.UserRepository
}

func (s *UserService) Create(user *models.User) error {
    // 业务逻辑
    return nil
}
```

### pkg/

可重用的公共包：

#### cache/

缓存实现：

```go
// pkg/cache/redis.go
type RedisCache struct {
    client *redis.Client
}

func (c *RedisCache) Get(key string) (interface{}, error) {
    // 缓存实现
    return nil, nil
}
```

#### console/

命令行框架：

```go
// pkg/console/command.go
type Command interface {
    GetName() string
    GetDescription() string
    Handle(context.Context) error
}
```

## 开发规范

1. 包的组织：
   - 相关的代码放在同一个包中
   - 包名使用小写
   - 避免循环依赖

2. 命名约定：
   - 使用有意义的名称
   - 遵循 Go 命名规范
   - 避免缩写（除非非常常见）

3. 接口定义：
   - 在使用者包中定义接口
   - 保持接口小而精确
   - 使用接口进行依赖注入

4. 错误处理：
   - 总是检查错误
   - 使用有意义的错误消息
   - 在适当的层级处理错误

5. 注释规范：
   - 为导出的类型和函数添加注释
   - 解释复杂的算法和业务逻辑
   - 使用示例说明用法

## 依赖注入

项目使用依赖注入管理组件依赖：

```go
// 服务注册
container := dig.New()
container.Provide(NewUserRepository)
container.Provide(NewUserService)
container.Provide(NewUserController)

// 服务解析
err := container.Invoke(func(controller *UserController) {
    // 使用控制器
})
```

## 测试组织

```
tests/
├── unit/              # 单元测试
├── integration/       # 集成测试
└── fixtures/          # 测试数据
```

## 部署结构

```
/var/www/go-admin/     # 应用根目录
├── bin/              # 编译后的二进制文件
├── config/           # 配置文件
├── storage/          # 存储目录
│   ├── logs/        # 日志文件
│   └── uploads/     # 上传文件
└── .env             # 环境变量
``` 