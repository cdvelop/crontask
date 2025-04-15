# Variables
BINARY_NAME=crontask
CMD_PATH=./cmd/crontask
OUTPUT_DIR=./bin

# Default target: Build for the current platform
all: build

# Build the binary for the current platform
build: clean
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(OUTPUT_DIR)
	@go build -ldflags="-s -w" -o $(OUTPUT_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Build complete: $(OUTPUT_DIR)/$(BINARY_NAME)"

# Cross-compile for Linux, Windows, macOS (amd64)
build-all: clean
	@echo "Starting cross-compilation..."
	@mkdir -p $(OUTPUT_DIR)
	@echo "Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(OUTPUT_DIR)/$(BINARY_NAME)_linux_amd64 $(CMD_PATH)
	@echo "Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(OUTPUT_DIR)/$(BINARY_NAME)_windows_amd64.exe $(CMD_PATH)
	@echo "Building for macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(OUTPUT_DIR)/$(BINARY_NAME)_darwin_amd64 $(CMD_PATH)
	@echo "Cross-compilation complete. Binaries are in $(OUTPUT_DIR)/"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)
	@echo "Clean complete."

# Phony targets (targets that don't represent files)
.PHONY: all build build-all clean
