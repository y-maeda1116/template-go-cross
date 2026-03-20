# Go Template Project

A cross-platform Go application template for Mac, Windows, and Linux.

## Features

- Standard Go project structure (`cmd/`, `internal/`)
- Environment variable configuration with `.env` support
- Graceful shutdown with signal handling
- Cross-platform builds (Windows, macOS arm64, Linux)
- GitHub Actions CI/CD
- Makefile for common operations

## Getting Started

### Prerequisites

- Go 1.22 or later
- Make

### Installation

```bash
# Clone the repository
git clone https://github.com/your-username/your-repo.git
cd your-repo

# Install dependencies
go mod download
```

### Configuration

Create a `.env` file in the project root:

```bash
cp env.example .env
# Edit .env with your configuration
```

## Usage

```bash
# Run the application
make run

# Build for current OS
make build

# Build for Windows
make build-win

# Build for Mac (Apple Silicon)
make build-mac

# Clean build artifacts
make clean
```

## Project Structure

```
.
├── .github/workflows/    # GitHub Actions workflows
├── cmd/                  # Application entry points
│   └── app/
│       └── main.go       # Main application
├── internal/             # Internal packages
│   └── config/          # Configuration
├── env.example          # Environment variable template
├── Makefile             # Build commands
└── go.mod               # Go module definition
```

## Development

### Running Tests

```bash
go test -v ./...
```

### Linting

```bash
go vet ./...
```

## License

MIT
