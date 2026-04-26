BINARY_NAME := mfg-dl
BUILD_DIR := bin

MAIN_PATH := ./cmd/$(BINARY_NAME)/main.go

GO := go
GOFLAGS := -trimpath
LDFLAGS := -ldflags="-s -w"
CGO_ENABLED := 0

.DEFAULT_GOAL := build

.PHONY: deps build build-all build-linux build-windows build-darwin release clean test update

deps:
	@$(GO) mod download

build: $(BUILD_DIR) deps
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

build-all: build-linux build-windows build-darwin

build-linux: $(BUILD_DIR) deps
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-x64 $(MAIN_PATH)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)

build-windows: $(BUILD_DIR) deps
	CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-x64.exe $(MAIN_PATH)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=arm64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe $(MAIN_PATH)

build-darwin: $(BUILD_DIR) deps
	CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-macos-intel $(MAIN_PATH)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-macos-silicon $(MAIN_PATH)

release: GOFLAGS := -trimpath -a
release: build-all

clean:
	rm -rf $(BUILD_DIR)/
	$(GO) clean -cache -testcache

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

test:
	$(GO) test ./test/...

test-verbose:
	$(GO) test -v ./test/...

update:
	$(GO) mod tidy
	@if command -v nix > /dev/null 2>&1; then \
		nix run github:nix-community/gomod2nix -- generate --with-deps; \
	fi
