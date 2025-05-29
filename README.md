# Go Admin Backend Template

A modern, production-ready admin backend template built with Go.

## Features

- 🔐 JWT Authentication - [Authentication Guide](docs/features/authentication.md)
- 👥 Role-Based Access Control (RBAC) - [RBAC Documentation](docs/features/rbac.md)
- 📝 Operation Logging - [Logging Guide](docs/features/logging.md)
- 🌐 I18n Support - [Internationalization](docs/features/i18n.md)
- 🔄 Redis Queue System - [Queue System](docs/features/queue.md)
- 📦 AWS S3 Integration - [S3 Integration](docs/features/s3.md)
- 🗄️ MySQL Database - [Database Guide](docs/database/README.md)
- 📊 API Documentation - [API Reference](docs/api/README.md)

## Documentation

- [Getting Started Guide](docs/getting-started/README.md)
- [Feature Documentation](docs/features/README.md)
- [API Reference](docs/api/README.md)
- [Database Schema](docs/database/README.md)
- [Testing Guide](docs/testing.md)
- [Advanced Topics](docs/advanced/README.md)
- [Example Usage](docs/examples/README.md)

## Quick Start

### Prerequisites

- Go 1.23+
- MySQL 5.7+
- Redis 6.0+

### Local Development

1. Clone the repository:
```bash
git clone <your-repo-url>
cd app
```

2. Install dependencies:
```bash
go mod download
```

3. Configure your environment:
```bash
cp configs/config.example.yaml configs/config.yaml
# Edit configs/config.yaml with your settings
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
│   ├── database/      # Database documentation
│   ├── features/      # Feature documentation
│   └── examples/      # Usage examples
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