# 开发环境配置

本文档详细说明了如何搭建 Go Admin Scaffold 的开发环境。

## 环境要求

### 1. 基础环境

- Go 1.21 或更高版本
- MySQL 5.7 或更高版本
- Redis 6.0 或更高版本
- Git
- Make (可选，用于构建)

### 2. 开发工具

- VSCode 或 GoLand
- Postman 或类似 API 测试工具
- MySQL 客户端工具
- Redis 客户端工具

## 环境搭建

### 1. 安装 Go

#### Windows
1. 下载 Go 安装包：https://golang.org/dl/
2. 运行安装程序
3. 设置环境变量：
```bash
# 设置 GOPATH
setx GOPATH "%USERPROFILE%\go"
# 将 Go 工具添加到 PATH
setx PATH "%PATH%;%GOROOT%\bin;%GOPATH%\bin"
```

#### Mac
```bash
# 使用 Homebrew 安装
brew install go

# 设置环境变量
echo 'export GOPATH=$HOME/go' >> ~/.zshrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.zshrc
source ~/.zshrc
```

#### Linux
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# CentOS/RHEL
sudo yum install golang

# 设置环境变量
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc
```

### 2. 安装 MySQL

#### Windows
1. 下载 MySQL 安装包：https://dev.mysql.com/downloads/installer/
2. 运行安装程序
3. 记住设置的 root 密码

#### Mac
```bash
# 使用 Homebrew 安装
brew install mysql

# 启动服务
brew services start mysql

# 设置 root 密码
mysql_secure_installation
```

#### Linux
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server
sudo mysql_secure_installation

# CentOS/RHEL
sudo yum install mysql-server
sudo systemctl start mysqld
sudo mysql_secure_installation
```

### 3. 安装 Redis

#### Windows
1. 下载 Redis for Windows：https://github.com/microsoftarchive/redis/releases
2. 解压并运行 redis-server.exe

#### Mac
```bash
# 使用 Homebrew 安装
brew install redis

# 启动服务
brew services start redis
```

#### Linux
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install redis-server
sudo systemctl start redis-server

# CentOS/RHEL
sudo yum install redis
sudo systemctl start redis
```

## 项目设置

### 1. 获取代码

```bash
# 克隆项目
git clone https://github.com/yourusername/go-admin-scaffold.git
cd go-admin-scaffold

# 安装依赖
go mod download
```

### 2. 配置开发环境

1. 复制开发环境配置：
```bash
cp configs/config.example.yaml configs/config.dev.yaml
```

2. 修改配置文件：
```yaml
# configs/config.dev.yaml
app:
  env: "development"
  debug: true
  port: 8080

mysql:
  host: "localhost"
  port: 3306
  username: "root"
  password: "your_password"
  database: "go_admin_dev"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

### 3. 初始化数据库

```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE go_admin_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 运行迁移
go run cmd/migrate/main.go

# 填充测试数据
go run cmd/tools/main.go seed run
```

### 4. 启动开发服务

```bash
# 安装开发工具
go install github.com/cosmtrek/air@latest

# 启动主服务（带热重载）
air

# 启动队列服务（可选）
go run cmd/worker/main.go
```

## 开发工具配置

### 1. VSCode 配置

安装推荐的扩展：
- Go
- Go Test Explorer
- Go Doc
- Go Outliner
- Go Test
- Go Coverage

settings.json 配置：
```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.formatTool": "goimports",
    "editor.formatOnSave": true,
    "go.testFlags": ["-v"],
    "go.coverOnSave": true
}
```

### 2. GoLand 配置

1. 启用 Go Modules
2. 配置 Go Linter
3. 设置代码风格
4. 配置调试器
5. 设置测试运行器

### 3. 开发工具安装

```bash
# 安装代码生成工具
go install github.com/golang/mock/mockgen@latest
go install github.com/vektra/mockery/v2@latest

# 安装代码检查工具
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 安装测试工具
go install github.com/cweill/gotests/...@latest
go install github.com/ramya-rao-a/go-outline@latest
```

## 开发工作流

### 1. 代码规范

- 遵循 [Go 代码规范](https://golang.org/doc/effective_go)
- 使用 `golangci-lint` 进行代码检查
- 提交前运行测试和检查

### 2. 分支管理

```bash
# 创建功能分支
git checkout -b feature/your-feature

# 提交更改
git add .
git commit -m "feat: add new feature"

# 推送到远程
git push origin feature/your-feature
```

### 3. 测试

```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test ./internal/models -v

# 运行测试覆盖率
go test ./... -cover

# 生成测试覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 4. 调试

1. 使用 VSCode 调试：
   - 设置断点
   - 启动调试会话
   - 查看变量和调用栈

2. 使用 GoLand 调试：
   - 配置运行配置
   - 设置断点
   - 启动调试会话

3. 使用日志调试：
```go
import "github.com/yourusername/go-admin-scaffold/pkg/logger"

logger.Debug("debug message")
logger.Info("info message")
logger.Error("error message")
```

## 常用命令

### 1. 开发命令

```bash
# 启动开发服务器
air

# 运行测试
go test ./...

# 代码检查
golangci-lint run

# 生成 mock
mockgen -source=internal/services/user.go -destination=internal/mocks/user_mock.go

# 数据库迁移
go run cmd/migrate/main.go

# 填充测试数据
go run cmd/tools/main.go seed run
```

### 2. 构建命令

```bash
# 构建所有二进制文件
make build

# 构建特定目标
make build-server
make build-worker
make build-tools

# 清理构建文件
make clean
```

### 3. 工具命令

```bash
# 生成新控制器
go run cmd/tools/main.go make controller user

# 生成新模型
go run cmd/tools/main.go make model user

# 生成新迁移
go run cmd/tools/main.go make migration create_users_table

# 生成新数据填充
go run cmd/tools/main.go make seeder user
```

## 常见问题

### 1. 依赖问题

检查：
- Go 版本是否正确
- 依赖是否完整
- 代理设置是否正确

解决：
```bash
# 清理依赖缓存
go clean -modcache

# 更新依赖
go mod tidy

# 验证依赖
go mod verify
```

### 2. 数据库问题

检查：
- 数据库服务是否运行
- 连接配置是否正确
- 用户权限是否正确

解决：
```bash
# 检查数据库状态
mysql -u root -p -e "SHOW DATABASES;"

# 重置数据库
go run cmd/tools/main.go migrate reset
go run cmd/tools/main.go seed run
```

### 3. Redis 问题

检查：
- Redis 服务是否运行
- 连接配置是否正确
- 内存使用情况

解决：
```bash
# 检查 Redis 状态
redis-cli ping

# 清空 Redis 数据
redis-cli FLUSHALL
```

## 相关文档

- [快速开始指南](../getting-started/quick-start.md)
- [项目结构说明](../getting-started/structure.md)
- [配置说明](../getting-started/configuration.md)
- [测试指南](testing.md)
- [API 文档](../api/README.md) 