# Go Admin Backend Template

A modern, production-ready admin backend template built with Go.

## Features

- 🔐 JWT Authentication - [Authentication Guide](docs/features/authentication.md)
- 👥 Role-Based Access Control (RBAC) - [RBAC Documentation](docs/features/rbac.md)
- 📝 Operation Logging
- 🌐 I18n Support
- 💫 Redis Queue System - [Cache & Queue](docs/features/cache.md)
- 📦 AWS S3 Integration
- 🗄️ MySQL Database
- 📊 API Documentation - [API Reference](docs/api/README.md)
- ⏱️ Task Scheduling - [Scheduling Guide](docs/features/scheduling.md)

## Documentation

- [Getting Started Guide](docs/getting-started/quick-start.md)
- [Project Structure](docs/getting-started/structure.md)
- [Configuration Guide](docs/getting-started/configuration.md)
- [API Documentation](docs/api/README.md)
- [Testing Guide](docs/testing.md)

## Quick Start

### Prerequisites

- Go 1.23+
- MySQL 5.7+
- Redis 6.0+

### Local Development

1. Clone the repository:
```bash
git clone <your-repo-url>
cd go-admin-scaffold
```

2. Install dependencies:
```bash
go mod download
```

3. Configure your environment:
```bash
# Copy configuration file
cp configs/config.example.yaml configs/config.yaml

# Edit configuration file
vim configs/config.yaml
```

4. Set up the database:
```bash
# Create database (if not exists)
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS go_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# Run all pending migrations
go run cmd/tools/main.go migrate run

# Run database seeding
go run cmd/tools/main.go seed run

# Other available commands:
go run cmd/tools/main.go migrate status    # Check migration status
go run cmd/tools/main.go migrate rollback  # Rollback last batch of migrations
go run cmd/tools/main.go migrate reset     # Reset all migrations
go run cmd/tools/main.go migrate refresh   # Reset and re-run all migrations
go run cmd/tools/main.go seed status       # Check seeding status
go run cmd/tools/main.go seed reset        # Reset seeding data
```

After successful setup, the following default accounts are created:

**Admin Account:**
- Username: `admin`
- Password: `admin123`
- Role: Administrator (all permissions)

**Test Accounts:**
- Username: `manager` / `user`
- Password: `admin123`
- Roles: Manager / Regular User (limited permissions)

5. Run the application:
```bash
go run cmd/server/main.go
```

### Docker Deployment

1. Build the Docker image:
```bash
docker build -t go-admin .
```

2. Run the container:
```bash
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/configs:/app/configs \
  --name go-admin \
  go-admin
```

### Binary Deployment

For production environments, you can deploy the application as a binary:

1. Build the binary for your target platform:
```bash
# Build for Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o build/go-admin-linux-amd64 cmd/server/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o build/go-admin-windows-amd64.exe cmd/server/main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o build/go-admin-darwin-amd64 cmd/server/main.go
```

2. Transfer the binary and configuration to your server:
```bash
# Create directory structure
ssh user@your-server "mkdir -p /opt/go-admin/{bin,configs,logs}"

# Copy binary
scp build/go-admin-linux-amd64 user@your-server:/opt/go-admin/bin/go-admin

# Copy example config and customize it
scp configs/config.example.yaml user@your-server:/opt/go-admin/configs/config.yaml
ssh user@your-server "vim /opt/go-admin/configs/config.yaml"  # Edit according to your environment
```

3. Set up systemd service (Linux):

```bash
# Create a dedicated user for the service
sudo useradd --system --shell /bin/false --home-dir /opt/go-admin --create-home goadmin
sudo chown -R goadmin:goadmin /opt/go-admin

# Make binary executable
sudo chmod +x /opt/go-admin/bin/go-admin

# Create systemd service file
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

# Service configuration
ExecStart=/opt/go-admin/bin/go-admin
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=5

# Working directory and environment
WorkingDirectory=/opt/go-admin
Environment=GIN_MODE=release
Environment=CONFIG_PATH=/opt/go-admin/configs/config.yaml

# Security settings
NoNewPrivileges=true
PrivateTmp=true
PrivateDevices=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/go-admin/logs

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=go-admin

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and enable the service
sudo systemctl daemon-reload
sudo systemctl enable go-admin.service

# Start the service
sudo systemctl start go-admin.service

# Check service status
sudo systemctl status go-admin.service
```

**Service Management Commands:**

```bash
# Start service
sudo systemctl start go-admin

# Stop service
sudo systemctl stop go-admin

# Restart service
sudo systemctl restart go-admin

# Check status
sudo systemctl status go-admin

# Enable auto-start on boot
sudo systemctl enable go-admin

# Disable auto-start
sudo systemctl disable go-admin

# View logs (real-time)
sudo journalctl -u go-admin -f

# View recent logs
sudo journalctl -u go-admin --since "1 hour ago"

# View all logs
sudo journalctl -u go-admin --no-pager
```

**Database Migration in Production:**

```bash
# Create a migration script
sudo tee /opt/go-admin/bin/migrate.sh > /dev/null << 'EOF'
#!/bin/bash
cd /opt/go-admin
./bin/go-admin -migrate
EOF

sudo chmod +x /opt/go-admin/bin/migrate.sh
sudo chown goadmin:goadmin /opt/go-admin/bin/migrate.sh

# Run migrations
sudo -u goadmin /opt/go-admin/bin/migrate.sh
```

**Health Check and Monitoring:**

```bash
# Check if service is running
curl -f http://localhost:8080/api/open/v1/public/health

# Monitor service with systemd
sudo systemctl is-active go-admin
sudo systemctl is-enabled go-admin

# Check resource usage
sudo systemctl show go-admin --property=ActiveState,SubState,MainPID
ps aux | grep go-admin
```

**Troubleshooting:**

```bash
# Check service logs for errors
sudo journalctl -u go-admin --since "10 minutes ago" --priority=err

# Check configuration
sudo -u goadmin /opt/go-admin/bin/go-admin --check-config

# Test configuration syntax
sudo -u goadmin cat /opt/go-admin/configs/config.yaml

# Check port binding
sudo netstat -tlnp | grep :8080
sudo ss -tlnp | grep :8080

# Check file permissions
ls -la /opt/go-admin/
sudo -u goadmin test -r /opt/go-admin/configs/config.yaml && echo "Config readable" || echo "Config not readable"
```

**Log Rotation Setup:**

```bash
# Create logrotate configuration
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
}
EOF
```

**Firewall Configuration (if using UFW):**

```bash
# Allow HTTP traffic
sudo ufw allow 8080/tcp

# Allow HTTPS if using reverse proxy
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
```

**Nginx Reverse Proxy (Optional):**

```bash
# Install nginx
sudo apt update && sudo apt install nginx

# Create nginx configuration
sudo tee /etc/nginx/sites-available/go-admin > /dev/null << 'EOF'
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
EOF

# Enable site
sudo ln -s /etc/nginx/sites-available/go-admin /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

**Security Best Practices:**

```bash
# Set proper file permissions
sudo chmod 600 /opt/go-admin/configs/config.yaml
sudo chmod 755 /opt/go-admin/bin/go-admin
sudo chmod -R 755 /opt/go-admin/logs

# Create secure directories
sudo mkdir -p /opt/go-admin/{backup,tmp}
sudo chown goadmin:goadmin /opt/go-admin/{backup,tmp}
sudo chmod 750 /opt/go-admin/{backup,tmp}

# Limit systemd service capabilities
sudo systemctl edit go-admin --full
# Add these security settings to [Service] section:
# CapabilityBoundingSet=CAP_NET_BIND_SERVICE
# SystemCallFilter=@system-service
# MemoryDenyWriteExecute=true
# LockPersonality=true
```

**Production Environment Configuration:**

```bash
# Update your production config.yaml
sudo tee -a /opt/go-admin/configs/config.yaml > /dev/null << 'EOF'

# Production optimizations
app:
  mode: release
  max_connections: 1000
  read_timeout: 30s
  write_timeout: 30s

# Logging configuration
log:
  level: info
  format: json
  output: /opt/go-admin/logs/app.log
  max_size: 100
  max_backups: 10
  max_age: 30

# Database connection pool
database:
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600

# Redis connection pool
redis:
  pool_size: 10
  min_idle_conns: 5
  dial_timeout: 5s
  read_timeout: 3s
  write_timeout: 3s
EOF
```

**Monitoring and Alerting:**

```bash
# Create monitoring script
sudo tee /opt/go-admin/bin/monitor.sh > /dev/null << 'EOF'
#!/bin/bash

# Health check
if ! curl -f -s http://localhost:8080/api/open/v1/public/health > /dev/null; then
    echo "$(date): Health check failed" >> /opt/go-admin/logs/monitor.log
    # Send alert (email, slack, etc.)
fi

# Check disk space
DISK_USAGE=$(df /opt/go-admin | tail -1 | awk '{print $5}' | sed 's/%//')
if [ $DISK_USAGE -gt 80 ]; then
    echo "$(date): Disk usage high: ${DISK_USAGE}%" >> /opt/go-admin/logs/monitor.log
fi

# Check memory usage
MEMORY_USAGE=$(ps -o pid,ppid,cmd,%mem,%cpu --sort=-%mem -C go-admin | tail -1 | awk '{print $4}')
if (( $(echo "$MEMORY_USAGE > 80" | bc -l) )); then
    echo "$(date): Memory usage high: ${MEMORY_USAGE}%" >> /opt/go-admin/logs/monitor.log
fi
EOF

sudo chmod +x /opt/go-admin/bin/monitor.sh
sudo chown goadmin:goadmin /opt/go-admin/bin/monitor.sh

# Add to crontab for regular monitoring
sudo crontab -u goadmin -e
# Add this line:
# */5 * * * * /opt/go-admin/bin/monitor.sh
```

**Backup and Recovery:**

```bash
# Create backup script
sudo tee /opt/go-admin/bin/backup.sh > /dev/null << 'EOF'
#!/bin/bash

BACKUP_DIR="/opt/go-admin/backup"
DATE=$(date +%Y%m%d_%H%M%S)

# Database backup
mysqldump -u dbuser -p'dbpass' go_admin > ${BACKUP_DIR}/db_backup_${DATE}.sql

# Config backup
cp /opt/go-admin/configs/config.yaml ${BACKUP_DIR}/config_backup_${DATE}.yaml

# Compress and clean old backups
gzip ${BACKUP_DIR}/db_backup_${DATE}.sql
find ${BACKUP_DIR} -name "*.gz" -mtime +7 -delete
find ${BACKUP_DIR} -name "config_backup_*.yaml" -mtime +7 -delete

echo "$(date): Backup completed" >> /opt/go-admin/logs/backup.log
EOF

sudo chmod +x /opt/go-admin/bin/backup.sh
sudo chown goadmin:goadmin /opt/go-admin/bin/backup.sh

# Schedule daily backups
sudo crontab -u goadmin -e
# Add this line:
# 0 2 * * * /opt/go-admin/bin/backup.sh
```

**Performance Tuning:**

```bash
# System optimizations for production
echo "* soft nofile 65536" | sudo tee -a /etc/security/limits.conf
echo "* hard nofile 65536" | sudo tee -a /etc/security/limits.conf
echo "goadmin soft nproc 4096" | sudo tee -a /etc/security/limits.conf
echo "goadmin hard nproc 4096" | sudo tee -a /etc/security/limits.conf

# TCP optimizations
echo 'net.core.somaxconn = 1024' | sudo tee -a /etc/sysctl.conf
echo 'net.ipv4.tcp_max_syn_backlog = 1024' | sudo tee -a /etc/sysctl.conf
echo 'net.core.netdev_max_backlog = 5000' | sudo tee -a /etc/sysctl.conf

# Apply changes
sudo sysctl -p
```

**For detailed production deployment guide, including security hardening, monitoring, and performance optimization, see [Production Deployment Guide](docs/deployment/production.md).**

## API Testing

Once the server is running, you can test the API endpoints:

```bash
# Health check
curl http://localhost:8080/api/open/v1/public/health

# Login with admin account
curl -X POST http://localhost:8080/api/admin/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'

# Get user list (requires authentication token from login response)
curl http://localhost:8080/api/admin/v1/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get roles
curl http://localhost:8080/api/admin/v1/roles \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get permissions
curl http://localhost:8080/api/admin/v1/permissions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**API Documentation:** Visit [API Reference](docs/api/README.md) for complete API documentation.

## Features Overview

### 🔐 Authentication
- JWT-based authentication
- Token refresh mechanism
- Session management
- Password hashing with bcrypt

### 👥 RBAC System
- Role-based permissions
- Granular permission control
- Dynamic permission checking
- User-role associations

### 📝 Logging
- Operation logs
- Login logs
- Request/response logging
- Structured logging with JSON format

### 🌐 I18n Support
- Multi-language support
- Dynamic locale switching
- Extensible translation system

### 💫 Redis Integration
- Caching
- Session storage
- Queue system
- Rate limiting

### 📊 Admin Features
- User management
- Role management
- Permission management
- System monitoring
- Operation logs

## Development

### Project Structure
```
├── cmd/                    # Application entry points
│   ├── server/            # Main server application
│   └── tools/             # CLI tools (migration, seeding)
├── configs/               # Configuration files
├── internal/              # Private application code
│   ├── api/              # API handlers
│   ├── core/             # Business logic
│   ├── database/         # Database layer
│   └── middleware/       # HTTP middleware
├── pkg/                  # Public library code
└── docs/                 # Documentation
```

### Adding New Features

1. **Create Migration:**
```bash
go run cmd/tools/main.go make:migration create_new_table
```

2. **Create Seeder:**
```bash
go run cmd/tools/main.go make:seeder new_table_seeder
```

3. **Add API Endpoints:**
- Create handler in `internal/api/admin/v1/`
- Add routes in router configuration
- Implement business logic in `internal/core/services/`

### Environment Variables

Key environment variables for development:

```bash
export GIN_MODE=debug           # Set to 'release' for production
export CONFIG_PATH=configs/config.yaml
export DB_DEBUG=true           # Enable SQL query logging
```

## Contributing

We welcome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) before submitting pull requests.

### Development Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Add tests for your changes
5. Run tests: `go test ./...`
6. Commit your changes: `git commit -m 'Add amazing feature'`
7. Push to your branch: `git push origin feature/amazing-feature`
8. Submit a pull request

### Code Style

- Follow Go conventions and best practices
- Use `gofmt` to format your code
- Add comments for exported functions and types
- Write meaningful commit messages

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed history of changes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](docs/README.md)
- 🐛 [Issue Tracker](https://github.com/yourusername/go-admin-scaffold/issues)
- 💬 [Discussions](https://github.com/yourusername/go-admin-scaffold/discussions)

## Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://github.com/go-gorm/gorm) - ORM library
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation
- [Redis](https://redis.io/) - In-memory data structure store
- [MySQL](https://www.mysql.com/) - Relational database

---

**⭐ If this project helps you, please consider giving it a star!**