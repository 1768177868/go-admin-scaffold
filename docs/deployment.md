# 部署指南

## 1. 环境要求

- Go 1.18 或更高版本
- MySQL 5.7 或更高版本
- Redis (可选，用于缓存和会话管理)
- 支持 systemd 的 Linux 系统（推荐 Ubuntu 20.04 或 CentOS 8）

## 2. 编译

### 2.1 Linux/Mac 环境编译

```bash
# 设置编译环境
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

# 编译
make build
# 或者直接使用 go build
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-admin cmd/server/main.go
```

### 2.2 Windows 环境编译

```powershell
# 设置编译环境
$env:GO111MODULE="on"
$env:GOPROXY="https://goproxy.cn,direct"

# 编译
go build -o go-admin.exe cmd/server/main.go
```

## 3. 部署步骤

### 3.1 准备工作

1. 创建部署目录：
```bash
mkdir -p /opt/go-admin
cd /opt/go-admin
```

2. 上传编译好的二进制文件和配置文件：
```bash
scp go-admin configs/config.prod.yaml user@your-server:/opt/go-admin/
```

3. 创建必要的目录：
```bash
mkdir -p storage/logs
mkdir -p storage/uploads
```

### 3.2 配置文件

1. 重命名并修改配置文件：
```bash
mv config.prod.yaml config.yaml
vim config.yaml
```

2. 修改关键配置：
   - 数据库连接信息
   - Redis 连接信息（如果使用）
   - JWT 密钥
   - 服务器监听地址和端口
   - 日志配置

### 3.3 systemd 服务配置

1. 创建 systemd 服务文件：
```bash
vim /etc/systemd/system/go-admin.service
```

2. 添加以下内容：
```ini
[Unit]
Description=Go Admin Service
After=network.target mysql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/go-admin
ExecStart=/opt/go-admin/go-admin
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
```

3. 启动服务：
```bash
# 重载 systemd
systemctl daemon-reload

# 启动服务
systemctl start go-admin

# 设置开机自启
systemctl enable go-admin

# 查看服务状态
systemctl status go-admin
```

### 3.4 Nginx 配置（可选）

如果需要使用 Nginx 作为反向代理，创建以下配置：

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 3.5 数据库迁移

```bash
# 运行数据库迁移
./go-admin migrate

# 初始化基础数据
./go-admin migrate -seed
```

## 4. 监控和维护

### 4.1 日志查看

```bash
# 查看服务日志
journalctl -u go-admin -f

# 查看应用日志
tail -f /opt/go-admin/storage/logs/app.log
```

### 4.2 备份

1. 数据库备份：
```bash
mysqldump -u root -p go_admin > backup_$(date +%Y%m%d).sql
```

2. 上传文件备份：
```bash
tar -czf uploads_$(date +%Y%m%d).tar.gz storage/uploads/
```

### 4.3 性能监控

推荐使用以下工具进行监控：
- Prometheus + Grafana
- 阿里云监控
- 腾讯云监控

## 5. 安全建议

1. 确保生产环境配置文件权限正确：
```bash
chmod 640 config.yaml
```

2. 使用非 root 用户运行服务

3. 配置防火墙只开放必要端口

4. 定期更新系统和依赖包

5. 启用 HTTPS

## 6. 故障排除

### 6.1 常见问题

1. 服务无法启动
   - 检查配置文件是否正确
   - 检查日志文件权限
   - 检查端口是否被占用

2. 数据库连接失败
   - 检查数据库配置
   - 确认数据库服务是否运行
   - 检查防火墙设置

3. 文件上传失败
   - 检查存储目录权限
   - 确认磁盘空间是否充足

### 6.2 性能优化

1. 数据库优化
   - 配置合适的连接池大小
   - 添加必要的索引
   - 定期维护数据库

2. 应用优化
   - 适当配置日志级别
   - 启用 Redis 缓存
   - 配置合理的 CORS 策略 