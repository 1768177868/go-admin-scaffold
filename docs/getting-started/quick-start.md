# 快速开始

## 环境要求

- Go 1.23+
- MySQL 5.7+
- Redis 6.0+
- 操作系统：Linux, macOS 或 Windows

## 安装步骤

1. 克隆项目：

```bash
git clone https://github.com/yourusername/go-admin-scaffold.git
cd go-admin-scaffold
```

2. 安装依赖：

```bash
go mod download
```

3. 配置环境：

```bash
# 复制配置文件
cp configs/config.example.yaml configs/config.yaml

# 编辑配置文件
vim configs/config.yaml
```

配置文件示例：

```yaml
app:
  name: go-admin
  mode: development
  port: 8080

database:
  driver: mysql
  host: localhost
  port: 3306
  database: go_admin
  username: root
  password: your_password

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-secret-key
  expire: 24h
```

4. 初始化数据库：

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE go_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 运行数据库迁移和初始化数据
go run cmd/migrate/main.go -seed
```

初始化后会创建默认管理员账户：
- 用户名：admin
- 密码：admin123
- 角色：管理员（具有所有权限）

5. 启动服务：

```bash
go run cmd/server/main.go
```

现在可以访问 http://localhost:8080 了。

## Docker 部署

1. 构建镜像：

```bash
docker build -t go-admin .
```

2. 运行容器：

```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/configs:/app/configs \
  --name go-admin \
  go-admin
```

## 开发模式

### 热重载

使用 air 实现热重载：

```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 运行开发服务器
air
```

### 调试模式

在 VS Code 中调试：

1. 安装 Go 扩展
2. 创建 launch.json 配置
3. 设置断点
4. 按 F5 启动调试

## 项目结构

```
.
├── cmd/                  # 应用入口
├── configs/             # 配置文件
├── internal/            # 内部代码
├── pkg/                # 公共库
└── scripts/            # 脚本文件
```

## 下一步

- [项目结构详解](structure.md)
- [配置说明](configuration.md)
- [API 文档](../api/README.md)
- [开发指南](../development/README.md) 