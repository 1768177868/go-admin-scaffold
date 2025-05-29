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

4. Run the application:
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

# Copy binary and configuration
scp build/go-admin-linux-amd64 user@your-server:/opt/go-admin/bin/go-admin
scp configs/config.prod.yaml user@your-server:/opt/go-admin/configs/config.yaml
```

3. Set up systemd service (Linux):
```bash
# Create systemd service file
sudo cat > /etc/systemd/system/go-admin.service << EOF
[Unit]
Description=Go Admin Backend Service
After=network.target mysql.service redis.service

[Service]
Type=simple
User=go-admin
WorkingDirectory=/opt/go-admin
ExecStart=/opt/go-admin/bin/go-admin
Restart=always
RestartSec=3
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and start service
sudo systemctl daemon-reload
sudo systemctl enable go-admin
sudo systemctl start go-admin
```

4. Monitor the service:
```bash
# Check service status
sudo systemctl status go-admin

# View logs
sudo journalctl -u go-admin -f
```

For detailed deployment instructions and best practices, please refer to our [Deployment Guide](docs/deployment.md).

## Project Structure

```
.
├── cmd/                  # Application entry points
│   └── server/          # Main server application
├── configs/             # Configuration files
│   ├── config.yaml      # Main configuration file
│   └── config.example.yaml  # Example configuration
├── docs/               # Documentation files
│   ├── api/           # API documentation
│   ├── features/      # Feature documentation
│   └── getting-started/ # Getting started guides
├── internal/            # Private application code
│   ├── api/            # API handlers
│   ├── config/         # Configuration structures
│   ├── core/           # Core business logic
│   ├── middleware/     # HTTP middleware
│   ├── models/         # Database models
│   └── routes/         # Route definitions
├── pkg/                # Public libraries
│   ├── database/       # Database utilities
│   ├── logger/         # Logging utilities
│   ├── queue/          # Queue implementation
│   ├── auth/           # Authentication utilities
│   └── response/       # API response helpers
├── scripts/            # Build and deployment scripts
├── static/             # Static assets
├── locales/            # I18n translation files
└── deploy/             # Deployment configurations
```

## Development

### Testing

For running tests and test coverage, please refer to our [Testing Guide](docs/testing.md).

### API Documentation

The API documentation is available at `/swagger/index.html` when running in development mode. For detailed API documentation, check [API Reference](docs/api/README.md).

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.