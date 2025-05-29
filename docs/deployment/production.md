# 生产环境部署指南

本指南详细介绍如何在生产环境中部署 Go Admin 系统，包括安全配置、性能优化和监控设置。

## 前置要求

- Linux 服务器 (Ubuntu 20.04+ / CentOS 8+ / RHEL 8+)
- MySQL 5.7+ 或 MySQL 8.0+
- Redis 6.0+
- Nginx (可选，用于反向代理)
- SSL 证书 (生产环境推荐)

## 系统准备

### 1. 创建专用用户

```bash
# 创建系统用户
sudo useradd --system --shell /bin/false --home-dir /opt/go-admin --create-home goadmin

# 创建目录结构
sudo mkdir -p /opt/go-admin/{bin,configs,logs,backup,tmp}
sudo chown -R goadmin:goadmin /opt/go-admin
```

### 2. 系统优化

```bash
# 文件句柄限制
echo "* soft nofile 65536" | sudo tee -a /etc/security/limits.conf
echo "* hard nofile 65536" | sudo tee -a /etc/security/limits.conf
echo "goadmin soft nproc 4096" | sudo tee -a /etc/security/limits.conf
echo "goadmin hard nproc 4096" | sudo tee -a /etc/security/limits.conf

# TCP 优化
cat << 'EOF' | sudo tee -a /etc/sysctl.conf
net.core.somaxconn = 1024
net.ipv4.tcp_max_syn_backlog = 1024
net.core.netdev_max_backlog = 5000
net.ipv4.tcp_keepalive_time = 600
net.ipv4.tcp_keepalive_intvl = 60
net.ipv4.tcp_keepalive_probes = 3
EOF

# 应用配置
sudo sysctl -p
```

## 应用部署

### 1. 构建和传输

```bash
# 本地构建
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/go-admin-linux-amd64 cmd/server/main.go

# 传输到服务器
scp build/go-admin-linux-amd64 user@server:/opt/go-admin/bin/go-admin
scp configs/config.example.yaml user@server:/opt/go-admin/configs/config.yaml

# 设置权限
sudo chmod +x /opt/go-admin/bin/go-admin
sudo chmod 600 /opt/go-admin/configs/config.yaml
sudo chown goadmin:goadmin /opt/go-admin/bin/go-admin /opt/go-admin/configs/config.yaml
```

### 2. 生产配置

```yaml
# /opt/go-admin/configs/config.yaml
app:
  name: go-admin
  mode: release
  port: 8080
  max_connections: 1000
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s

database:
  driver: mysql
  host: localhost
  port: 3306
  database: go_admin
  username: go_admin_user
  password: your_secure_password
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600

redis:
  host: localhost
  port: 6379
  password: your_redis_password
  db: 0
  pool_size: 10
  min_idle_conns: 5
  dial_timeout: 5s
  read_timeout: 3s
  write_timeout: 3s

log:
  level: info
  format: json
  output: /opt/go-admin/logs/app.log
  max_size: 100
  max_backups: 10
  max_age: 30
  compress: true

jwt:
  secret: your_very_secure_jwt_secret_key_here
  expire: 24h

security:
  cors_origins: ["https://your-frontend-domain.com"]
  rate_limit: 1000
  max_upload_size: 10485760  # 10MB
```

## Systemd 服务配置

### 1. 创建服务文件

```bash
sudo tee /etc/systemd/system/go-admin.service > /dev/null << 'EOF'
[Unit]
Description=Go Admin Backend Service
Documentation=https://github.com/yourusername/go-admin-scaffold
After=network.target mysql.service redis.service
Wants=network.target
Requires=mysql.service redis.service

[Service]
Type=simple
User=goadmin
Group=goadmin

# 服务配置
ExecStart=/opt/go-admin/bin/go-admin
ExecReload=/bin/kill -HUP $MAINPID
ExecStop=/bin/kill -TERM $MAINPID
Restart=always
RestartSec=5
StartLimitBurst=3
StartLimitInterval=60

# 工作目录和环境变量
WorkingDirectory=/opt/go-admin
Environment=GIN_MODE=release
Environment=CONFIG_PATH=/opt/go-admin/configs/config.yaml

# 安全设置
NoNewPrivileges=true
PrivateTmp=true
PrivateDevices=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/go-admin/logs /opt/go-admin/tmp
CapabilityBoundingSet=CAP_NET_BIND_SERVICE
SystemCallFilter=@system-service
MemoryDenyWriteExecute=true
LockPersonality=true

# 资源限制
LimitNOFILE=65536
LimitNPROC=4096
MemoryMax=1G
CPUQuota=50%

# 日志配置
StandardOutput=journal
StandardError=journal
SyslogIdentifier=go-admin

[Install]
WantedBy=multi-user.target
EOF
```

### 2. 启动服务

```bash
# 重新加载 systemd
sudo systemctl daemon-reload

# 启用开机自启
sudo systemctl enable go-admin.service

# 启动服务
sudo systemctl start go-admin.service

# 检查状态
sudo systemctl status go-admin.service
```

## 数据库初始化

### 1. 创建数据库用户

```sql
-- 连接到 MySQL
mysql -u root -p

-- 创建数据库
CREATE DATABASE go_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建专用用户
CREATE USER 'go_admin_user'@'localhost' IDENTIFIED BY 'your_secure_password';
GRANT ALL PRIVILEGES ON go_admin.* TO 'go_admin_user'@'localhost';
FLUSH PRIVILEGES;
```

### 2. 运行迁移

```bash
# 创建迁移工具
sudo tee /opt/go-admin/bin/migrate.sh > /dev/null << 'EOF'
#!/bin/bash
cd /opt/go-admin
exec ./bin/go-admin migrate "$@"
EOF

sudo chmod +x /opt/go-admin/bin/migrate.sh
sudo chown goadmin:goadmin /opt/go-admin/bin/migrate.sh

# 执行迁移
sudo -u goadmin /opt/go-admin/bin/migrate.sh run
sudo -u goadmin /opt/go-admin/bin/migrate.sh seed
```

## 反向代理配置

### 1. Nginx 配置

```bash
# 安装 Nginx
sudo apt update && sudo apt install nginx

# 创建配置文件
sudo tee /etc/nginx/sites-available/go-admin > /dev/null << 'EOF'
upstream go_admin {
    server 127.0.0.1:8080;
    keepalive 32;
}

server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL 配置
    ssl_certificate /path/to/your/certificate.pem;
    ssl_certificate_key /path/to/your/private.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    # 安全头
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload";

    # 上传大小限制
    client_max_body_size 10M;

    # 代理配置
    location / {
        proxy_pass http://go_admin;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_connect_timeout 30;
        proxy_send_timeout 30;
        proxy_read_timeout 30;
    }

    # 静态文件
    location /static/ {
        alias /opt/go-admin/static/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # 健康检查
    location /health {
        access_log off;
        proxy_pass http://go_admin/api/open/v1/public/health;
    }
}
EOF

# 启用站点
sudo ln -s /etc/nginx/sites-available/go-admin /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 监控和日志

### 1. 日志轮转

```bash
sudo tee /etc/logrotate.d/go-admin > /dev/null << 'EOF'
/opt/go-admin/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    copytruncate
    su goadmin goadmin
    postrotate
        systemctl reload go-admin || true
    endscript
}
EOF
```

### 2. 监控脚本

```bash
sudo tee /opt/go-admin/bin/monitor.sh > /dev/null << 'EOF'
#!/bin/bash

LOG_FILE="/opt/go-admin/logs/monitor.log"

# 健康检查
health_check() {
    if ! curl -f -s --max-time 10 http://localhost:8080/api/open/v1/public/health > /dev/null; then
        echo "$(date): Health check failed" >> $LOG_FILE
        return 1
    fi
    return 0
}

# 检查磁盘空间
disk_check() {
    DISK_USAGE=$(df /opt/go-admin | tail -1 | awk '{print $5}' | sed 's/%//')
    if [ $DISK_USAGE -gt 80 ]; then
        echo "$(date): Disk usage high: ${DISK_USAGE}%" >> $LOG_FILE
    fi
}

# 检查内存使用
memory_check() {
    MEMORY_USAGE=$(ps -o pid,ppid,cmd,%mem --sort=-%mem -C go-admin | tail -1 | awk '{print $4}')
    if (( $(echo "$MEMORY_USAGE > 80" | bc -l) )); then
        echo "$(date): Memory usage high: ${MEMORY_USAGE}%" >> $LOG_FILE
    fi
}

# 执行检查
health_check
disk_check
memory_check
EOF

sudo chmod +x /opt/go-admin/bin/monitor.sh
sudo chown goadmin:goadmin /opt/go-admin/bin/monitor.sh
```

### 3. Cron 任务

```bash
# 编辑 crontab
sudo crontab -u goadmin -e

# 添加以下任务
*/5 * * * * /opt/go-admin/bin/monitor.sh
0 2 * * * /opt/go-admin/bin/backup.sh
0 3 * * 0 /usr/sbin/logrotate -f /etc/logrotate.d/go-admin
```

## 备份策略

### 1. 备份脚本

```bash
sudo tee /opt/go-admin/bin/backup.sh > /dev/null << 'EOF'
#!/bin/bash

BACKUP_DIR="/opt/go-admin/backup"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=7

# 创建备份目录
mkdir -p $BACKUP_DIR

# 数据库备份
mysqldump -u go_admin_user -p'your_secure_password' \
    --single-transaction \
    --routines \
    --triggers \
    go_admin > ${BACKUP_DIR}/db_backup_${DATE}.sql

# 配置文件备份
cp /opt/go-admin/configs/config.yaml ${BACKUP_DIR}/config_backup_${DATE}.yaml

# 压缩备份
gzip ${BACKUP_DIR}/db_backup_${DATE}.sql

# 清理旧备份
find ${BACKUP_DIR} -name "*.gz" -mtime +${RETENTION_DAYS} -delete
find ${BACKUP_DIR} -name "config_backup_*.yaml" -mtime +${RETENTION_DAYS} -delete

echo "$(date): Backup completed" >> /opt/go-admin/logs/backup.log
EOF

sudo chmod +x /opt/go-admin/bin/backup.sh
sudo chown goadmin:goadmin /opt/go-admin/bin/backup.sh
```

## 安全加固

### 1. 防火墙配置

```bash
# 使用 UFW
sudo ufw --force reset
sudo ufw default deny incoming
sudo ufw default allow outgoing

# 允许 SSH (根据实际端口修改)
sudo ufw allow 22/tcp

# 允许 HTTP/HTTPS (如果使用 Nginx)
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# 如果直接暴露应用端口
# sudo ufw allow 8080/tcp

sudo ufw --force enable
```

### 2. 系统安全

```bash
# 安装 fail2ban
sudo apt install fail2ban

# 配置 fail2ban
sudo tee /etc/fail2ban/jail.local > /dev/null << 'EOF'
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 3

[sshd]
enabled = true
port = ssh
filter = sshd
logpath = /var/log/auth.log
maxretry = 3

[nginx-http-auth]
enabled = true
filter = nginx-http-auth
port = http,https
logpath = /var/log/nginx/error.log
EOF

sudo systemctl enable fail2ban
sudo systemctl start fail2ban
```

## 性能优化

### 1. Go 应用优化

```bash
# 编译时优化
go build -ldflags="-s -w" -gcflags="-B" -o go-admin cmd/server/main.go

# 运行时环境变量
echo 'GOGC=100' | sudo tee -a /etc/environment
echo 'GOMAXPROCS=4' | sudo tee -a /etc/environment
```

### 2. 数据库优化

```sql
-- MySQL 配置优化 (/etc/mysql/mysql.conf.d/mysqld.cnf)
[mysqld]
innodb_buffer_pool_size = 1G
innodb_log_file_size = 256M
innodb_flush_log_at_trx_commit = 2
innodb_file_per_table = 1
max_connections = 200
query_cache_size = 64M
query_cache_type = 1
```

## 故障排除

### 1. 常见问题

```bash
# 检查服务状态
sudo systemctl status go-admin

# 查看日志
sudo journalctl -u go-admin -f
sudo tail -f /opt/go-admin/logs/app.log

# 检查配置
sudo -u goadmin /opt/go-admin/bin/go-admin --check-config

# 检查端口占用
sudo netstat -tlnp | grep :8080
sudo ss -tlnp | grep :8080

# 检查数据库连接
mysql -u go_admin_user -p -e "SELECT 1"

# 检查 Redis 连接
redis-cli ping
```

### 2. 应急处理

```bash
# 重启服务
sudo systemctl restart go-admin

# 强制杀死进程
sudo pkill -f go-admin

# 回滚数据库
sudo -u goadmin /opt/go-admin/bin/migrate.sh rollback

# 恢复备份
gunzip -c /opt/go-admin/backup/db_backup_YYYYMMDD_HHMMSS.sql.gz | mysql -u go_admin_user -p go_admin
```

## 更新部署

### 1. 蓝绿部署

```bash
# 构建新版本
GOOS=linux GOARCH=amd64 go build -o build/go-admin-new cmd/server/main.go

# 上传到服务器
scp build/go-admin-new user@server:/opt/go-admin/bin/go-admin-new

# 备份当前版本
sudo cp /opt/go-admin/bin/go-admin /opt/go-admin/bin/go-admin-backup

# 原子替换
sudo mv /opt/go-admin/bin/go-admin-new /opt/go-admin/bin/go-admin
sudo chmod +x /opt/go-admin/bin/go-admin
sudo chown goadmin:goadmin /opt/go-admin/bin/go-admin

# 重启服务
sudo systemctl restart go-admin

# 验证部署
curl -f http://localhost:8080/api/open/v1/public/health

# 如果出现问题，回滚
# sudo mv /opt/go-admin/bin/go-admin-backup /opt/go-admin/bin/go-admin
# sudo systemctl restart go-admin
```

这个部署指南涵盖了生产环境的各个方面，确保系统的安全性、可靠性和性能。 