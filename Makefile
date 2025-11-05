.PHONY: build test test-verbose test-coverage clean tidy fmt lint help docker-build docker-build-platform

# Default target
.DEFAULT_GOAL := help

# Docker/Podman configuration
DOCKER_BIN ?= docker
PLATFORM ?= linux/amd64
RELEASE ?= latest
IMAGE_TAG ?= hrexed/otelcol-securityevents
IMAGE_FULL ?= $(IMAGE_TAG):$(RELEASE)

# Build the processor
build:
	@echo "Building Security Event Processor..."
	go build ./...

# Run all tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test ./... -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out
	@echo "Coverage report generated: coverage.out"
	@echo "Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "HTML coverage report generated: coverage.html"

# Run tests for a specific package
test-package:
	@echo "Running tests for specific package..."
	@if [ -z "$(PKG)" ]; then \
		echo "Usage: make test-package PKG=./processor/securityevent/internal/openreports"; \
		exit 1; \
	fi
	go test -v $(PKG)

# Run tests with race detection
test-race:
	@echo "Running tests with race detector..."
	go test ./... -race -v

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	go clean ./...
	rm -f coverage.out coverage.html
	@echo "Clean complete"

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy
	@echo "Dependencies updated"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi
	golangci-lint run ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	@echo "Dependencies installed"

# Verify dependencies
verify:
	@echo "Verifying dependencies..."
	go mod verify

# Build collector using OCB
ocb-build:
	@echo "Building collector with OCB..."
	@chmod +x ocb/build.sh
	@./ocb/build.sh

# Build collector with simple manifest
ocb-build-simple:
	@echo "Building collector with simple manifest..."
	@chmod +x ocb/build.sh
	@./ocb/build.sh manifest-simple.yaml

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@echo "  Binary: $(DOCKER_BIN)"
	@echo "  Platform: $(PLATFORM)"
	@echo "  Release: $(RELEASE)"
	@echo "  Image: $(IMAGE_FULL)"
	@if ! command -v $(DOCKER_BIN) &> /dev/null; then \
		echo "Error: $(DOCKER_BIN) not found. Please install $(DOCKER_BIN) or set DOCKER_BIN=podman"; \
		exit 1; \
	fi
	$(DOCKER_BIN) build \
		--platform $(PLATFORM) \
		-t $(IMAGE_FULL) \
		-t $(IMAGE_TAG):latest \
		-f docker/Dockerfile \
		.
	@echo "✓ Docker image built successfully: $(IMAGE_FULL)"

# Build Docker image for multiple platforms
docker-build-platform:
	@echo "Building Docker image for platform $(PLATFORM)..."
	@echo "  Binary: $(DOCKER_BIN)"
	@echo "  Platform: $(PLATFORM)"
	@echo "  Release: $(RELEASE)"
	@echo "  Image: $(IMAGE_FULL)"
	@if ! command -v $(DOCKER_BIN) &> /dev/null; then \
		echo "Error: $(DOCKER_BIN) not found. Please install $(DOCKER_BIN) or set DOCKER_BIN=podman"; \
		exit 1; \
	fi
	@if [ "$(DOCKER_BIN)" = "podman" ]; then \
		echo "Note: Podman does not support --platform flag in build. Building for current platform."; \
		$(DOCKER_BIN) build \
			-t $(IMAGE_FULL) \
			-t $(IMAGE_TAG):latest \
			-f docker/Dockerfile \
			.; \
	else \
		$(DOCKER_BIN) build \
			--platform $(PLATFORM) \
			-t $(IMAGE_FULL) \
			-t $(IMAGE_TAG):latest \
			-f docker/Dockerfile \
			.; \
	fi
	@echo "✓ Docker image built successfully: $(IMAGE_FULL)"

# Build Docker image using docker-compose
docker-compose-build:
	@echo "Building Docker image with docker-compose..."
	@echo "  Release: $(RELEASE)"
	@echo "  Image: $(IMAGE_FULL)"
	@if ! command -v docker-compose &> /dev/null && ! command -v docker &> /dev/null; then \
		echo "Error: docker-compose or docker not found"; \
		exit 1; \
	fi
	@if command -v docker-compose &> /dev/null; then \
		IMAGE_TAG=$(IMAGE_FULL) docker-compose -f docker/docker-compose.yml build; \
	else \
		IMAGE_TAG=$(IMAGE_FULL) docker compose -f docker/docker-compose.yml build; \
	fi

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	@echo "  Binary: $(DOCKER_BIN)"
	@echo "  Image: $(IMAGE_FULL)"
	$(DOCKER_BIN) run -d \
		--name otelcol-securityevent \
		-p 4317:4317 \
		-p 4318:4318 \
		-p 55679:55679 \
		-p 13133:13133 \
		-v $(PWD)/ocb/config.yaml:/etc/otelcol/config.yaml:ro \
		$(IMAGE_FULL)

# Stop Docker container
docker-stop:
	@echo "Stopping Docker container..."
	$(DOCKER_BIN) stop otelcol-securityevent || true
	$(DOCKER_BIN) rm otelcol-securityevent || true

# Clean Docker artifacts
docker-clean:
	@echo "Cleaning Docker artifacts..."
	$(DOCKER_BIN) rmi $(IMAGE_FULL) || true
	$(DOCKER_BIN) rmi $(IMAGE_TAG):latest || true
	$(DOCKER_BIN) rmi otelcol-securityevent:latest || true

# Clean OCB build artifacts
ocb-clean:
	@echo "Cleaning OCB build artifacts..."
	rm -rf dist build

# Install OCB
ocb-install:
	@echo "Installing OCB..."
	go install go.opentelemetry.io/collector/cmd/builder@v0.138.0
	@echo "OCB installed. Add to PATH: export PATH=\$$PATH:\$$(go env GOPATH)/bin"

# Run all checks (format, test, build)
check: fmt test build
	@echo "All checks passed"

# Full build (processor + collector + docker)
# Note: docker-build includes OCB build inside Docker, so ocb-build is optional
all: build docker-build
	@echo "Full build complete"

# Documentation targets
docs-install:
	@echo "Installing documentation dependencies..."
	pip install -r docs/requirements.txt

docs-serve:
	@echo "Starting MkDocs development server..."
	mkdocs serve

docs-build:
	@echo "Building documentation..."
	mkdocs build

docs-deploy:
	@echo "Deploying documentation to GitHub Pages..."
	mkdocs gh-deploy

# Clean everything
clean-all: clean ocb-clean docker-clean
	@echo "All artifacts cleaned"

# Show help
help:
	@echo "Security Event Processor - Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build              - Build the processor"
	@echo "  test               - Run all tests"
	@echo "  test-verbose       - Run tests with verbose output"
	@echo "  test-coverage      - Run tests and generate coverage report"
	@echo "  test-package PKG=  - Run tests for a specific package"
	@echo "  test-race          - Run tests with race detector"
	@echo "  clean              - Clean build artifacts"
	@echo "  tidy               - Tidy Go module dependencies"
	@echo "  fmt                - Format code"
	@echo "  lint               - Lint code (requires golangci-lint)"
	@echo "  deps               - Install dependencies"
	@echo "  verify             - Verify dependencies"
	@echo "  check              - Run format, test, and build"
	@echo ""
	@echo "OCB (OpenTelemetry Collector Builder) targets:"
	@echo "  ocb-install        - Install OCB tool"
	@echo "  ocb-build          - Build collector with OCB (full manifest)"
	@echo "  ocb-build-simple   - Build collector with OCB (simple manifest)"
	@echo "  ocb-clean          - Clean OCB build artifacts"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-build       - Build Docker image (default: docker, linux/amd64, hrexed/otelcol-securityevents:latest)"
	@echo "  docker-build-platform - Build Docker image with platform support"
	@echo "  docker-compose-build - Build Docker image with docker-compose"
	@echo "  docker-run         - Run Docker container"
	@echo "  docker-stop        - Stop Docker container"
	@echo "  docker-clean       - Clean Docker images"
	@echo ""
	@echo "Docker variables (can be overridden):"
	@echo "  DOCKER_BIN         - Container binary (default: docker, use 'podman' for Podman)"
	@echo "  PLATFORM           - Target platform (default: linux/amd64)"
	@echo "  RELEASE            - Release version/tag (default: latest)"
	@echo "  IMAGE_TAG          - Image repository (default: hrexed/otelcol-securityevents)"
	@echo ""
	@echo "Documentation targets:"
	@echo "  docs-install       - Install documentation dependencies"
	@echo "  docs-serve         - Start MkDocs development server"
	@echo "  docs-build         - Build documentation"
	@echo "  docs-deploy        - Deploy documentation to GitHub Pages"
	@echo ""
	@echo "Combined targets:"
	@echo "  all                - Build processor + collector + Docker"
	@echo "  clean-all          - Clean all artifacts"
	@echo "  help               - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build                    - Build the processor"
	@echo "  make ocb-build                - Build collector with OCB"
	@echo "  make docker-build             - Build Docker image (default settings)"
	@echo "  make docker-build DOCKER_BIN=podman - Build with Podman"
	@echo "  make docker-build RELEASE=0.1.0 - Build with specific release version"
	@echo "  make docker-build PLATFORM=linux/arm64 - Build for ARM64"
	@echo "  make docker-build IMAGE_TAG=myregistry/otelcol-securityevents - Custom image tag"
	@echo "  make docker-build DOCKER_BIN=podman RELEASE=0.1.0 PLATFORM=linux/amd64"
	@echo "  make all                      - Full build (processor + collector + docker)"
	@echo "  make test-coverage            - Generate coverage report"
	@echo "  make test-package PKG=./processor/securityevent/internal/openreports"
	@echo "  make check                    - Run all checks"

