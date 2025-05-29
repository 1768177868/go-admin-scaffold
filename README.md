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

# Copy binary
scp build/go-admin-linux-amd64 user@your-server:/opt/go-admin/bin/go-admin

# Copy example config and customize it
scp configs/config.example.yaml user@your-server:/opt/go-admin/configs/config.yaml
ssh user@your-server "vim /opt/go-admin/configs/config.yaml"  # Edit according to your environment
```

3. Set up systemd service (Linux):
```