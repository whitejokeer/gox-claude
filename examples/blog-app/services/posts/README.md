# Posts Service

A microservice for posts functionality in the GOX application.

## Overview

This service handles posts-related operations and can be deployed independently.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)

### Running Locally

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Copy environment variables:
   ```bash
   cp .env.example .env
   ```

3. Run the service:
   ```bash
   go run cmd/server/main.go
   ```

The service will start on the port specified in the .env file (default: 3001).

### Running with Docker

1. Build the Docker image:
   ```bash
   docker build -t posts-service .
   ```

2. Run the container:
   ```bash
   docker run -p 3001:3001 --env-file .env posts-service
   ```

## API Endpoints

### Health Check
- **GET** /health - Returns service health status

### Service Endpoints
Add your service-specific endpoints documentation here.

## Configuration

See .env.example for all available configuration options.

## Development

### Project Structure

```
posts/
├── cmd/
│   └── server/         # Application entry point
├── internal/           # Private application code
│   ├── config/         # Configuration management
│   ├── handlers/       # HTTP handlers
│   ├── models/         # Data models
│   ├── repository/     # Data access layer
│   └── service/        # Business logic
├── api/                # API definitions (proto files, etc.)
├── tests/              # Test files
├── .env                # Environment variables (not in git)
├── .env.example        # Environment variables template
├── Dockerfile          # Docker configuration
├── go.mod              # Go module definition
└── README.md           # This file
```

### Testing

Run tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## Deployment

This service can be deployed independently using:
- Docker
- Kubernetes
- Any cloud platform that supports Go applications

## License

Part of the GOX Framework application.
