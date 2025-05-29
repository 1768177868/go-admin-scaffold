# Go Admin Backend Template

A modern, production-ready admin backend template built with Go.

## Features

- 🔐 JWT Authentication
- 👥 Role-Based Access Control (RBAC)
- 📝 Operation Logging
- 🌐 I18n Support
- 🔄 Redis Queue System
- 📦 AWS S3 Integration
- 🗄️ MySQL Database
- 📊 API Documentation

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
├── configs/             # Configuration files
├── internal/            # Private application code
│   ├── api/            # API handlers
│   ├── config/         # Configuration structures
│   ├── core/           # Core business logic
│   └── routes/         # Route definitions
├── pkg/                # Public libraries
│   ├── database/       # Database utilities
│   ├── logger/         # Logging utilities
│   ├── queue/          # Queue implementation
│   └── response/       # API response helpers
└── scripts/            # Build and deployment scripts
```

## API Documentation

The API documentation is available at `/swagger/index.html` when running in development mode.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.