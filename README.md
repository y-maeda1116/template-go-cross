# Go CLI + Desktop Template

A cross-platform Go application template supporting both CLI and Desktop (Wails) interfaces.

## Features

- **Dual Interface:** CLI (Cobra) and Desktop (Wails + React)
- **Shared Core:** Business logic shared between CLI and Desktop
- **Structured Logging:** zap-based logging with multiple levels
- **Configuration:** Viper-based config with YAML/TOML/JSON support
- **Hot Reload:** Air for CLI, Wails dev for Desktop
- **Testing:** testify + mockgen with 80%+ coverage goal
- **CI/CD:** GitHub Actions for testing and building

## Project Structure

```
.
├── cmd/
│   ├── cli/           # CLI entry point (Cobra)
│   └── desktop/       # Wails desktop entry point
├── internal/
│   ├── core/          # Shared business logic
│   ├── ui/            # Wails UI bindings
│   ├── config/        # Configuration (Viper)
│   ├── logger/        # Logging (zap)
│   └── cli/          # CLI commands
├── frontend/         # Wails frontend (React)
├── pkg/              # Reusable packages
├── test/
│   └── mocks/        # Generated mocks
├── Makefile          # Build targets
├── wails.json        # Wails configuration
└── air.toml          # Hot reload configuration
```

## Getting Started

### Prerequisites

- Go 1.22 or later
- Node.js 20 or later (for Desktop)
- Make

### Installation

```bash
# Clone repository
git clone https://github.com/your-username/your-repo.git
cd your-repo

# Install Go dependencies
go mod download

# Install frontend dependencies (for Desktop)
cd frontend && npm install && cd ..
```

### Configuration

Create environment and config files:

```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your settings
```

## Usage

### CLI

```bash
# Run CLI
make run-cli

# Show help
make run-cli ARGS="--help"

# Say hello
make run-cli ARGS="hello --name World"

# Build CLI
make build-cli
```

### Desktop

```bash
# Run Desktop in dev mode
make run-desktop

# Build Desktop
make build-desktop
```

### Development

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detector
make test-race

# Generate mocks
make mocks

# Format code
make fmt

# Run linter
make lint

# Clean build artifacts
make clean
```

## Architecture

```
Application Layer
├── CLI (Cobra)
└── Desktop (Wails + React)
         ↓
Core Business Logic Layer
└── Shared services
         ↓
Infrastructure Layer
├── Config (Viper)
└── Logger (zap)
```

## License

MIT
