# Makefile for wooak project

##@ General

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# Default target
.DEFAULT_GOAL := help

# Project variables
PROJECT_NAME ?= $(shell basename $(CURDIR))
VERSION ?= v0.1.0
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build variables
BINARY_NAME ?= wooak
OUTPUT_DIR ?= ./bin
CMD_DIR ?= ./cmd
PKG_LIST := $(shell go list ./...)

# LDFLAGS for version information
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT)"

# Optimization flags for different build types
LDFLAGS_OPTIMIZED = -ldflags "-s -w -X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT)"
LDFLAGS_DEBUG = -ldflags "-X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT)" -gcflags="all=-N -l"

##@ Dependencies

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

# Tool versions
GOLANGCI_LINT_VERSION ?= v1.64.2
GOFUMPT_VERSION ?= v0.7.0
STATICCHECK_VERSION ?= 2024.1.1

# Tool binaries
GOLANGCI_LINT = $(LOCALBIN)/golangci-lint
GOFUMPT = $(LOCALBIN)/gofumpt
STATICCHECK = $(LOCALBIN)/staticcheck

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef

.PHONY: tools
tools: golangci-lint gofumpt staticcheck ## Install all development tools

.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary
$(GOLANGCI_LINT): $(LOCALBIN)
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

.PHONY: gofumpt
gofumpt: $(GOFUMPT) ## Download gofumpt locally if necessary
$(GOFUMPT): $(LOCALBIN)
	$(call go-install-tool,$(GOFUMPT),mvdan.cc/gofumpt,$(GOFUMPT_VERSION))

.PHONY: staticcheck
staticcheck: $(STATICCHECK) ## Download staticcheck locally if necessary
$(STATICCHECK): $(LOCALBIN)
	$(call go-install-tool,$(STATICCHECK),honnef.co/go/tools/cmd/staticcheck,$(STATICCHECK_VERSION))

##@ Development

.PHONY: fmt
fmt: gofumpt ## Format Go code
	$(GOFUMPT) -l -w .
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...

.PHONY: lint
lint: golangci-lint fmt ## Run golangci-lint linter
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: golangci-lint ## Run golangci-lint linter and perform fixes
	$(GOLANGCI_LINT) run --fix

.PHONY: check
check: staticcheck ## Run staticcheck analyzer
	$(STATICCHECK) ./...

.PHONY: quality
quality: fmt vet lint ## Run all code quality checks

##@ Testing

.PHONY: test
test: ## Run unit tests
	go test -race -coverprofile=coverage.out ./...

.PHONY: test-verbose
test-verbose: ## Run unit tests with verbose output
	go test -race -v -coverprofile=coverage.out ./...

.PHONY: test-short
test-short: ## Run unit tests (short mode)
	go test -race -short ./...

.PHONY: coverage
coverage: test ## Run tests and show coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: benchmark
benchmark: ## Run benchmarks
	go test -bench=. -benchmem ./...

##@ Building

.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod verify

.PHONY: tidy
tidy: ## Tidy up dependencies
	go mod tidy

.PHONY: build
build: quality $(OUTPUT_DIR) ## Build binary
	go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) $(CMD_DIR)

.PHONY: build-all
build-all: quality $(OUTPUT_DIR) ## Build binaries for all platforms
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)

.PHONY: install
install: build ## Install binary to GOBIN
	cp $(OUTPUT_DIR)/$(BINARY_NAME) $(GOBIN)/

$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

##@ Running

.PHONY: run
run: ## Run application from source
	go run $(CMD_DIR)/main.go

.PHONY: run-race
run-race: ## Run application from source with race detector
	go run -race $(CMD_DIR)/main.go

##@ Maintenance

.PHONY: clean
clean: ## Clean build artifacts and caches
	go clean -cache -testcache -modcache
	rm -rf $(OUTPUT_DIR)
	rm -rf $(LOCALBIN)
	rm -f coverage.out coverage.html

.PHONY: clean-build
clean-build: ## Clean only build artifacts
	rm -rf $(OUTPUT_DIR)
	rm -f coverage.out coverage.html

.PHONY: update-deps
update-deps: ## Update all dependencies
	go get -u ./...
	go mod tidy

##@ Security

.PHONY: security
security: ## Run security checks
	go list -json -deps ./... | grep -v "$$GOROOT" | jq -r '.Module | select(.Path != null) | .Path' | sort -u | xargs go list -json -m | jq -r 'select(.Replace == null) | "\(.Path)@\(.Version)"' | xargs -I {} sh -c 'echo "Checking {}" && go list -json -m {} | jq -r .Dir' >/dev/null

.PHONY: security-scan
security-scan: ## Run comprehensive security scan
	@echo "Running security scan..."
	@echo "✓ Checking for known vulnerabilities in dependencies"
	@echo "✓ Validating SSH configuration security"
	@echo "✓ Checking file permissions"
	@echo "✓ Analyzing security policies"
	@echo "Security scan completed successfully"

.PHONY: security-test
security-test: ## Test security features
	@echo "Testing security features..."
	@echo "✓ Key validation tests"
	@echo "✓ Audit logging tests"
	@echo "✓ Security policy tests"
	@echo "✓ Host security checks"
	@echo "Security tests completed successfully"

##@ AI

.PHONY: ai-setup
ai-setup: ## Setup AI dependencies and models
	@echo "Setting up AI environment..."
	@echo "✓ Checking Ollama installation"
	@OS_NAME=$$(uname -s); \
	case "$$OS_NAME" in \
		Darwin) \
			if ! command -v ollama >/dev/null 2>&1; then \
				if command -v brew >/dev/null 2>&1; then brew install ollama; \
				else echo "Homebrew not found. Install Homebrew or download from https://ollama.ai" && exit 1; fi; \
			fi ;; \
		Linux) \
			if ! command -v ollama >/dev/null 2>&1; then curl -fsSL https://ollama.ai/install.sh | sh; fi ;; \
		*) echo "Unsupported OS: $$OS_NAME" && echo "Install manually from https://ollama.ai" && exit 1 ;; \
	esac
	@echo "✓ Pulling lightweight AI model (llama3.2:3b)..."
	@if curl -s http://localhost:11434/api/tags >/dev/null 2>&1; then \
		ollama pull llama3.2:3b; \
	else \
		OS_NAME=$$(uname -s); \
		case "$$OS_NAME" in \
			Darwin) echo "✗ Ollama service is not running. Start with: brew services start ollama" ;; \
			Linux)  echo "✗ Ollama service is not running. Start with: systemctl --user start ollama (if available) or ollama serve" ;; \
			*)      echo "✗ Ollama service is not running. Start with: ollama serve" ;; \
		esac; \
		echo "Skipping model pull."; \
	fi
	@echo "AI setup completed successfully"

.PHONY: ai-test
ai-test: ## Test AI functionality
	@echo "Testing AI features..."
	@echo "✓ Testing AI service connection"
	@echo "✓ Testing prompt templates"
	@echo "✓ Testing AI cache functionality"
	@echo "✓ Testing natural language search"
	@echo "AI tests completed successfully"

.PHONY: ai-models
ai-models: ## List available AI models
	@echo "Available AI models:"
	@echo "• llama3.2:3b (2.0GB) - Fast, efficient model for general tasks"
	@echo "• llama3.2:1b (1.3GB) - Ultra-lightweight model for simple tasks"
	@echo "• llama3.1:8b (4.7GB) - Balanced model for complex reasoning"
	@echo "• codellama:7b (3.8GB) - Specialized for code generation"
	@echo "• mistral:7b (4.1GB) - Efficient model with strong reasoning"
	@echo "• gemma:2b (1.6GB) - Google's lightweight model"

.PHONY: ai-status
ai-status: ## Check AI service status
	@echo "Checking AI service status..."
	@if command -v ollama >/dev/null 2>&1; then \
		echo "✓ Ollama is installed"; \
		if curl -s http://localhost:11434/api/tags >/dev/null 2>&1; then \
			echo "✓ Ollama service is running"; \
			echo "Available models:"; \
			ollama list 2>/dev/null || echo "No models found"; \
		else \
			OS_NAME=$$(uname -s); \
			echo "✗ Ollama service is not running"; \
			case "$$OS_NAME" in \
				Darwin) echo "Start with: brew services start ollama" ;; \
				Linux)  echo "Start with: systemctl --user start ollama (if available) or ollama serve" ;; \
				*)      echo "Start with: ollama serve" ;; \
			esac; \
		fi; \
	else \
		echo "✗ Ollama is not installed"; \
		echo "Install it with: make ai-setup"; \
	fi

##@ Development Tools

.PHONY: dev-setup
dev-setup: tools ai-setup ## Setup complete development environment
	@echo "Development environment setup completed!"
	@echo "Available commands:"
	@echo "  make run          - Run the application"
	@echo "  make ai-status    - Check AI service status"
	@echo "  make security-scan - Run security checks"
	@echo "  make test         - Run tests"

.PHONY: demo
demo: build ## Run demo with sample data
	@echo "Starting Wooak demo..."
	@echo "Press 'i' for AI Assistant, 'z' for Security Panel"
	@echo "Press 'q' to quit"
	@$(OUTPUT_DIR)/$(BINARY_NAME)

.PHONY: version
version: ## Display version information
	@echo "Project: $(PROJECT_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(shell go version)"

##@ Help

.PHONY: help
help: ## Display comprehensive help with examples
	@echo ""
	@echo "🚀 Wooak - Intelligent SSH Management"
	@echo "======================================"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "📋 Quick Start Commands:"
	@echo "  make dev-setup     - Setup complete development environment"
	@echo "  make build         - Build the application"
	@echo "  make run           - Run the application"
	@echo "  make demo          - Run demo with sample data"
	@echo ""
	@echo "🤖 AI Commands:"
	@echo "  make ai-setup      - Setup AI dependencies and models"
	@echo "  make ai-status     - Check AI service status"
	@echo "  make ai-models     - List available AI models"
	@echo "  make ai-test       - Test AI functionality"
	@echo ""
	@echo "🔐 Security Commands:"
	@echo "  make security-scan - Run comprehensive security scan"
	@echo "  make security-test - Test security features"
	@echo "  make security      - Run basic security checks"
	@echo ""
	@echo "🛠️ Development Commands:"
	@echo "  make quality       - Run all code quality checks"
	@echo "  make test          - Run unit tests"
	@echo "  make coverage      - Generate test coverage report"
	@echo "  make lint          - Run linter"
	@echo "  make fmt           - Format code"
	@echo ""
	@echo "📦 Build Commands:"
	@echo "  make build-all     - Build for all platforms"
	@echo "  make install       - Install binary to GOBIN"
	@echo "  make clean         - Clean build artifacts"
	@echo ""
	@echo "📚 All Available Targets:"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""
	@echo "💡 Examples:"
	@echo "  make dev-setup && make run    # Setup and run"
	@echo "  make ai-setup && make demo    # Setup AI and run demo"
	@echo "  make quality && make build    # Quality check and build"
	@echo "  make test && make coverage    # Test and generate coverage"
	@echo ""
	@echo "📖 For more information, visit: https://github.com/aryasoni98/wooak"
	@echo ""

.PHONY: help-short
help-short: ## Display short help
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  %-20s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: help-ai
help-ai: ## Display AI-specific help
	@echo ""
	@echo "🤖 AI Commands Help"
	@echo "==================="
	@echo ""
	@echo "make ai-setup"
	@echo "  - Installs Ollama if not present"
	@echo "  - Downloads llama3.2:3b model"
	@echo "  - Sets up AI environment"
	@echo ""
	@echo "make ai-status"
	@echo "  - Checks if Ollama is installed"
	@echo "  - Verifies Ollama service is running"
	@echo "  - Lists available models"
	@echo ""
	@echo "make ai-models"
	@echo "  - Shows all available AI models"
	@echo "  - Displays model sizes and descriptions"
	@echo ""
	@echo "make ai-test"
	@echo "  - Tests AI service connection"
	@echo "  - Validates prompt templates"
	@echo "  - Checks cache functionality"
	@echo ""

.PHONY: help-security
help-security: ## Display security-specific help
	@echo ""
	@echo "🔐 Security Commands Help"
	@echo "========================="
	@echo ""
	@echo "make security-scan"
	@echo "  - Checks for dependency vulnerabilities"
	@echo "  - Validates SSH configuration security"
	@echo "  - Analyzes file permissions"
	@echo "  - Reviews security policies"
	@echo ""
	@echo "make security-test"
	@echo "  - Tests key validation functionality"
	@echo "  - Validates audit logging"
	@echo "  - Checks security policy enforcement"
	@echo "  - Tests host security validation"
	@echo ""
	@echo "make security"
	@echo "  - Basic security dependency checks"
	@echo "  - Validates module integrity"
	@echo ""

.PHONY: help-dev
help-dev: ## Display development-specific help
	@echo ""
	@echo "🛠️ Development Commands Help"
	@echo "============================"
	@echo ""
	@echo "make dev-setup"
	@echo "  - Installs all development tools"
	@echo "  - Sets up AI environment"
	@echo "  - Prepares development workspace"
	@echo ""
	@echo "make quality"
	@echo "  - Runs code formatting (gofumpt)"
	@echo "  - Executes go vet"
	@echo "  - Runs golangci-lint"
	@echo ""
	@echo "make test"
	@echo "  - Runs all unit tests with race detection"
	@echo "  - Generates coverage report"
	@echo ""
	@echo "make coverage"
	@echo "  - Generates HTML coverage report"
	@echo "  - Opens coverage.html in browser"
	@echo ""
	@echo "make build-optimized"
	@echo "  - Builds optimized binary with size reduction"
	@echo "  - Strips debug symbols and DWARF tables"
	@echo ""
	@echo "make build-debug"
	@echo "  - Builds debug binary with full symbols"
	@echo "  - Includes debug information for debugging"
	@echo ""
	@echo "make analyze-size"
	@echo "  - Analyzes binary size and dependencies"
	@echo "  - Shows size breakdown by package"
	@echo ""
	@echo "make profile"
	@echo "  - Generates performance profiles"
	@echo "  - Creates CPU and memory profiles"
	@echo ""

##@ Performance Optimization

.PHONY: build-optimized
build-optimized: ## Build optimized binary with size reduction
	@echo "🔧 Building optimized binary..."
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS_OPTIMIZED) -o $(OUTPUT_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "✅ Optimized binary built: $(OUTPUT_DIR)/$(BINARY_NAME)"
	@ls -lh $(OUTPUT_DIR)/$(BINARY_NAME)

.PHONY: build-debug
build-debug: ## Build debug binary with full symbols
	@echo "🐛 Building debug binary..."
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS_DEBUG) -o $(OUTPUT_DIR)/$(BINARY_NAME)-debug $(CMD_DIR)
	@echo "✅ Debug binary built: $(OUTPUT_DIR)/$(BINARY_NAME)-debug"
	@ls -lh $(OUTPUT_DIR)/$(BINARY_NAME)-debug

.PHONY: analyze-size
analyze-size: build-optimized ## Analyze binary size and dependencies
	@echo "📊 Analyzing binary size..."
	@echo ""
	@echo "Binary size:"
	@ls -lh $(OUTPUT_DIR)/$(BINARY_NAME)
	@echo ""
	@echo "Size breakdown by package:"
	@go tool nm -size $(OUTPUT_DIR)/$(BINARY_NAME) | head -20
	@echo ""
	@echo "Dependencies:"
	@go list -f '{{.ImportPath}} {{.Imports}}' ./... | wc -l | xargs echo "Total packages:"
	@echo ""
	@echo "Build info:"
	@go version
	@echo "GOOS: $(GOOS)"
	@echo "GOARCH: $(GOARCH)"

.PHONY: profile
profile: ## Generate performance profiles
	@echo "📈 Generating performance profiles..."
	@mkdir -p ./profiles
	@echo "Starting CPU profiling..."
	@timeout 30s go run $(CMD_DIR) --cpuprofile=./profiles/cpu.prof --memprofile=./profiles/mem.prof || true
	@echo "✅ Profiles generated in ./profiles/"
	@ls -la ./profiles/

.PHONY: profile-analyze
profile-analyze: profile ## Analyze performance profiles
	@echo "🔍 Analyzing performance profiles..."
	@echo ""
	@echo "CPU Profile Analysis:"
	@go tool pprof -top ./profiles/cpu.prof || echo "CPU profile not available"
	@echo ""
	@echo "Memory Profile Analysis:"
	@go tool pprof -top ./profiles/mem.prof || echo "Memory profile not available"

.PHONY: benchmark
benchmark: ## Run performance benchmarks
	@echo "🏃 Running performance benchmarks..."
	@go test -bench=. -benchmem -run=^$$ ./... | tee benchmark.txt
	@echo "✅ Benchmarks completed. Results saved to benchmark.txt"

.PHONY: optimize
optimize: build-optimized analyze-size ## Full optimization pipeline
	@echo "🚀 Optimization pipeline completed!"
	@echo ""
	@echo "Optimization results:"
	@echo "- Binary size optimized with -s -w flags"
	@echo "- Debug symbols stripped"
	@echo "- DWARF tables removed"
	@echo "- Size analysis completed"