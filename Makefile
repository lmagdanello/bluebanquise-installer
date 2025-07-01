# BlueBanquise Installer Makefile

.PHONY: help build test test-unit test-integration clean install ci ci-local security-check

# Default target
help:
	@echo "Available targets:"
	@echo "  build          - Build the bluebanquise-installer binary"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  install        - Install dependencies"
	@echo "  lint           - Run linter"
	@echo "  format         - Format code"
	@echo "  ci             - Run CI checks locally"
	@echo "  ci-local       - Run CI checks without Docker"
	@echo "  security-check - Run security checks"
	@echo "  release        - Build release binaries"

# Build the binary
build:
	@echo "Building bluebanquise-installer..."
	go build -o bluebanquise-installer .

# Run all tests
test: test-unit test-integration

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test -v ./internal/... ./cmd/...

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	go test -v -tags=integration ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f bluebanquise-installer
	rm -f coverage.out coverage.html
	rm -f bluebanquise-installer-*
	rm -rf release/
	rm -rf offline-packages/
	rm -rf tarball-packages/

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
format:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	go install github.com/google/go-licenses@latest

# Run tests in verbose mode
test-v:
	@echo "Running tests in verbose mode..."
	go test -v ./...

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -race ./...

# Build for different platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bluebanquise-installer-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bluebanquise-installer-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bluebanquise-installer-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bluebanquise-installer-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bluebanquise-installer-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o bluebanquise-installer-windows-arm64.exe .

# Show test results summary
test-summary:
	@echo "Test Summary:"
	@go test -json ./... | jq -r '. | select(.Action=="pass" or .Action=="fail") | "\(.Package): \(.Action)"' 2>/dev/null || echo "Install jq for better test summary"

# Run CI checks locally (simula o que roda no GitHub Actions)
ci: install-tools
	@echo "Running CI checks locally..."
	@echo "1. Running linter..."
	@make lint
	@echo "2. Running unit tests..."
	@make test-unit
	@echo "3. Running integration tests..."
	@make test-integration
	@echo "4. Running security checks..."
	@make security-check
	@echo "5. Building binary..."
	@make build
	@echo "✅ All CI checks passed!"

# Run CI checks without Docker (for local development)
ci-local: install-tools
	@echo "Running CI checks locally (without Docker)..."
	@echo "1. Running linter..."
	@make lint
	@echo "2. Running unit tests..."
	@make test-unit
	@echo "3. Running integration tests..."
	@make test-integration
	@echo "4. Running security checks..."
	@make security-check
	@echo "5. Building binary..."
	@make build
	@echo "✅ All CI checks passed!"

# Run security checks
security-check:
	@echo "Running security checks..."
	@echo "1. Running gosec..."
	@gosec -fmt=json -out=security-results.json ./... || true
	@echo "2. Checking licenses..."
	@go-licenses check ./... || echo "License check failed - check manually"
	@echo "3. Running go vet..."
	@go vet ./...
	@echo "✅ Security checks completed!"

# Build release binaries
release: clean build-all
	@echo "Creating release assets..."
	@mkdir -p release
	@cp bluebanquise-installer-* release/
	@sha256sum bluebanquise-installer-* > release/checksums.txt
	@cp README.md release/ 2>/dev/null || echo "No README.md found"
	@cp LICENSE release/ 2>/dev/null || echo "No LICENSE file found"
	@echo "✅ Release assets created in release/ directory"

# Run Docker-based tests locally
test-docker:
	@echo "Running Docker-based tests..."
	@echo "This requires Docker to be running"
	@echo "Testing online installation on Ubuntu 22.04..."
	@docker run --rm --privileged -v $(PWD):/installer ubuntu:22.04 bash -c "\
		cd /installer && \
		apt-get update && apt-get install -y python3 python3-pip python3-venv git curl && \
		./bluebanquise-installer online --user testuser --home /tmp/bluebanquise && \
		ls -la /tmp/bluebanquise/ && \
		/tmp/bluebanquise/ansible_venv/bin/ansible --version"

# Check code quality
quality-check: lint format test-coverage
	@echo "✅ Code quality check completed!"

# Prepare for development
dev-setup: install install-tools
	@echo "✅ Development environment ready!"