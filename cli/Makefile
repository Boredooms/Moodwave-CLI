BINARY     := moodwave
MODULE     := github.com/moodwave/moodwave
CMD        := ./cmd/moodwave
DIST       := dist
VERSION    ?= dev
BUILD_TIME := $(shell powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'" 2>/dev/null || date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS    := -s -w -X $(MODULE)/internal/config.Version=$(VERSION) -X $(MODULE)/internal/config.BuildTime=$(BUILD_TIME)

# Default target
.PHONY: all
all: build

# Build for current platform
.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD)

# Install to GOPATH/bin
.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)" $(CMD)

# Cross-compile for all platforms
.PHONY: build-all
build-all: clean
	@echo "Building for all platforms..."
	GOOS=windows GOARCH=amd64   go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-windows-amd64.exe $(CMD)
	GOOS=windows GOARCH=arm64   go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-windows-arm64.exe $(CMD)
	GOOS=darwin  GOARCH=amd64   go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-darwin-amd64     $(CMD)
	GOOS=darwin  GOARCH=arm64   go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-darwin-arm64      $(CMD)
	GOOS=linux   GOARCH=amd64   go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-linux-amd64       $(CMD)
	GOOS=linux   GOARCH=arm64   go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(BINARY)-linux-arm64       $(CMD)
	@echo "Builds complete. Output in $(DIST)/"

# Windows-friendly cross-compile (PowerShell)
.PHONY: cross
cross: clean
	powershell -Command "scripts/build.ps1 -Version '$(VERSION)'"

# Run tests
.PHONY: test
test:
	go test ./... -v -race -timeout 60s

# Short tests (no race detector, faster)
.PHONY: test-short
test-short:
	go test ./... -short -timeout 30s

# Static analysis
.PHONY: vet
vet:
	go vet ./...

# Format
.PHONY: fmt
fmt:
	go fmt ./...

# Lint (requires golangci-lint)
.PHONY: lint
lint:
	golangci-lint run ./...

# Clean
.PHONY: clean
clean:
	@if exist $(BINARY) del /f $(BINARY) 2>nul; true
	@if exist $(BINARY).exe del /f $(BINARY).exe 2>nul; true
	@if exist $(DIST) rmdir /s /q $(DIST) 2>nul; true
	mkdir $(DIST) 2>nul; true

# Run in development mode
.PHONY: run
run:
	go run $(CMD) $(ARGS)

# Run doctor for diagnostics
.PHONY: doctor
doctor:
	go run $(CMD) doctor

# Show version
.PHONY: version
version:
	go run $(CMD) --version

# Tidy dependencies
.PHONY: tidy
tidy:
	go mod tidy

# Generate go.sum
.PHONY: deps
deps:
	go mod download

.PHONY: help
help:
	@echo "Moodwave CLI — Build Targets"
	@echo ""
	@echo "  make build       Build for current platform"
	@echo "  make build-all   Cross-compile for all platforms"
	@echo "  make install     Install to GOPATH/bin"
	@echo "  make test        Run all tests"
	@echo "  make vet         Run go vet"
	@echo "  make fmt         Format code"
	@echo "  make clean       Remove build artifacts"
	@echo "  make run         Run in development mode"
	@echo "  make doctor      Run diagnostics"
	@echo "  make help        Show this help"
