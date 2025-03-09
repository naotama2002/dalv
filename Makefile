.PHONY: all build test clean lint

# 変数定義
BINARY_NAME=dalv
BUILD_DIR=./bin
MAIN_PATH=./cmd/dalv
GO=go
GOTEST=$(GO) test
GOLINT=$(GO) run golang.org/x/lint/golint

# デフォルトターゲット
all: clean lint test build

# ビルドターゲット
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# テストターゲット
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...
	@echo "Tests completed"

# リントターゲット
lint:
	@echo "Running linter..."
	$(GO) vet ./...
	$(GOLINT) ./...
	@echo "Linting completed"

# クリーンターゲット
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean completed"

# インストールターゲット
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "Installation completed"

# ヘルプターゲット
help:
	@echo "Available targets:"
	@echo "  all      - Run clean, lint, test, and build"
	@echo "  build    - Build the binary"
	@echo "  test     - Run tests"
	@echo "  lint     - Run linter"
	@echo "  clean    - Remove build artifacts"
	@echo "  install  - Install binary to GOPATH/bin"
	@echo "  help     - Show this help message"
