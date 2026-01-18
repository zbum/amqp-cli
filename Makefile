.PHONY: build clean test lint fmt run

BINARY_NAME=amqp-cli
DIST_DIR=dist

build:
	@mkdir -p $(DIST_DIR)
	go build -o $(DIST_DIR)/$(BINARY_NAME) .

clean:
	rm -rf $(DIST_DIR)
	go clean

test:
	go test -v ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...

run:
	go run .

# Cross-compilation targets
.PHONY: build-linux build-darwin build-windows

build-linux:
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .

build-darwin:
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .

build-windows:
	@mkdir -p $(DIST_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .

build-all: build-linux build-darwin build-windows

# Release targets for Homebrew
.PHONY: release release-tarball

VERSION ?= 0.1.0

release: build-all release-tarball

release-tarball:
	@mkdir -p $(DIST_DIR)/release
	@cd $(DIST_DIR) && tar -czf release/$(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	@cd $(DIST_DIR) && tar -czf release/$(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@cd $(DIST_DIR) && tar -czf release/$(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	@cd $(DIST_DIR) && tar -czf release/$(BINARY_NAME)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	@cd $(DIST_DIR) && zip release/$(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "\nSHA256 checksums:"
	@cd $(DIST_DIR)/release && shasum -a 256 *
