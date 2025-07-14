# blog-app

A GOX Framework project created with the GOX CLI.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- GOX CLI installed

### Installation

```bash
# Install dependencies
go mod download

# Start development server
gox dev
```

### Project Structure

```
blog-app/
├── app/            # API Gateway with UI
│   ├── pages/          # Frontend pages (.gox files)
│   ├── components/     # Project-specific components
│   ├── shared/         # Shared UI components
│   ├── middleware/     # Gateway middleware
│   └── routing/        # Page routing
├── services/           # Microservices
├── common/             # Shared code across services
│   ├── middleware/     # Service middleware
│   └── discovery/      # Service discovery
├── infra/              # Infrastructure code
├── tests/              # End-to-end tests
├── docker-compose.yml  # Local development
├── gox.config.yaml     # Main configuration
├── go.mod              # Go dependencies
├── go.work             # Go workspace
└── README.md           # This file
```

## Development

### Running the development server

```bash
gox dev
```

### Generating components

```bash
# Generate a new page
gox generate page dashboard

# Generate a new component
gox generate component user-card

# Generate a new service
gox generate service users
```

## Testing

```bash
# Run all tests
gox test

# Run tests with coverage
gox test --coverage
```

## Building

```bash
# Build for production
gox build

# Build for specific target
gox build --target=linux-amd64
```

## Deployment

```bash
# Deploy to Docker
gox deploy --target=docker

# Deploy to Kubernetes
gox deploy --target=k8s
```

## License

This project is licensed under the MIT License.
