.DEFAULT_GOAL := help

.PHONY: help deps build run build-all build-linux build-windows build-darwin release clean clean-cache test test-verbose fmt vet lint check update

## Show this help
help:
	@ grep -h -A 1 "^## " $(MAKEFILE_LIST) \
	| grep -v "^--$$" \
	| awk -F "## " '/^## / { desc = $$2; next } { sub(/:.*/, "", $$1); printf "\033[36m%-20s\033[0m %s\n", $$1, desc }'

GO ?= go

GOFLAGS ?= -trimpath
LDFLAGS ?= -ldflags="-s -w"

CGO_ENABLED ?= 0

BINARY_NAME := mfg-dl
BUILD_DIR := build

MAIN_PATH := ./cmd/$(BINARY_NAME)

## downloads go module dependencies
deps:
	@$(GO) mod download

## builds binary for current platform
build:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

RUN_ARGS ?=
## builds and runs binary with RUN_ARGS
run: build
	$(BUILD_DIR)/$(BINARY_NAME) $(RUN_ARGS)

## builds binaries for linux, windows and macos
build-all: build-linux build-windows build-darwin

build-linux:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-x64 $(MAIN_PATH)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)

build-windows:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-x64.exe $(MAIN_PATH)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=arm64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe $(MAIN_PATH)

build-darwin:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-macos-intel $(MAIN_PATH)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-macos-silicon $(MAIN_PATH)

## cleans, builds all platforms
release:
	@$(MAKE) clean
	@$(MAKE) build-all

## removes the build dir
clean:
	rm -rf $(BUILD_DIR)/

## cleans global go build and test cache
clean-cache:
	$(GO) clean -cache -testcache

## go mod tidy
update:
	$(GO) mod tidy

## runs all tests
test:
	$(GO) test ./...

## verbose test
test-verbose:
	$(GO) test -v ./...

## formats code
fmt:
	$(GO) fmt ./...

## runs go vet
vet:
	$(GO) vet ./...

## runs golangci-lint
lint:
	golangci-lint run ./...

## runs fmt, vet and lint
check: fmt vet lint
