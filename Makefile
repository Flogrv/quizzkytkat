.PHONY: build run clean test docker-build docker-run docker-stop install dev

# Build the application
build:
	@echo "🔨 Building quiz-server..."
	@go build -o quiz-server .
	@echo "✅ Build complete!"

# Run the application
run: build
	@echo "🚀 Starting quiz-server on port 2222..."
	@./quiz-server

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed!"

# Development mode with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "🔥 Starting development server..."
	@air

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -f quiz-server
	@rm -rf data/
	@echo "✅ Clean complete!"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t cybersec-quiz:latest .
	@echo "✅ Docker image built!"

# Run Docker container
docker-run: docker-build
	@echo "🐳 Starting Docker container..."
	@docker-compose up -d
	@echo "✅ Container started on port 2222!"
	@echo "Connect with: ssh -p 2222 localhost"

# Stop Docker container
docker-stop:
	@echo "🛑 Stopping Docker container..."
	@docker-compose down
	@echo "✅ Container stopped!"

# Show logs
logs:
	@docker-compose logs -f

# Generate SSH keys (if needed)
keygen:
	@echo "🔑 Generating SSH host key..."
	@mkdir -p .ssh
	@ssh-keygen -t ed25519 -f .ssh/id_ed25519 -N ""
	@echo "✅ SSH key generated!"

# Help
help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Build and run locally"
	@echo "  make install      - Install dependencies"
	@echo "  make dev          - Run with hot reload (requires air)"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run with Docker Compose"
	@echo "  make docker-stop  - Stop Docker container"
	@echo "  make logs         - Show Docker logs"
	@echo "  make keygen       - Generate SSH host key"
