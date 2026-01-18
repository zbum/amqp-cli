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
.PHONY: release release-tarball publish update-formula

VERSION ?= 0.1.0
HOMEBREW_TAP_DIR ?= /tmp/homebrew-tap

release: build-all release-tarball

release-tarball:
	@mkdir -p $(DIST_DIR)/release
	@cd $(DIST_DIR) && cp $(BINARY_NAME)-darwin-amd64 $(BINARY_NAME) && tar -czf release/$(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME) && rm $(BINARY_NAME)
	@cd $(DIST_DIR) && cp $(BINARY_NAME)-darwin-arm64 $(BINARY_NAME) && tar -czf release/$(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME) && rm $(BINARY_NAME)
	@cd $(DIST_DIR) && cp $(BINARY_NAME)-linux-amd64 $(BINARY_NAME) && tar -czf release/$(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME) && rm $(BINARY_NAME)
	@cd $(DIST_DIR) && cp $(BINARY_NAME)-linux-arm64 $(BINARY_NAME) && tar -czf release/$(BINARY_NAME)-linux-arm64.tar.gz $(BINARY_NAME) && rm $(BINARY_NAME)
	@cd $(DIST_DIR) && zip release/$(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "\n=== SHA256 checksums ==="
	@cd $(DIST_DIR)/release && shasum -a 256 *
	@echo "\nUpdate Formula/$(BINARY_NAME).rb with these SHA256 values"

# Create GitHub release and upload assets
gh-release:
	@echo "Creating GitHub release v$(VERSION)..."
	gh release create v$(VERSION) $(DIST_DIR)/release/* --title "v$(VERSION)" --notes "Release v$(VERSION)"

# Update homebrew-tap repo
update-tap:
	@echo "Updating homebrew-tap..."
	@rm -rf $(HOMEBREW_TAP_DIR)
	@git clone https://github.com/zbum/homebrew-tap.git $(HOMEBREW_TAP_DIR)
	@cp Formula/$(BINARY_NAME).rb $(HOMEBREW_TAP_DIR)/Formula/
	@cd $(HOMEBREW_TAP_DIR) && git add Formula/$(BINARY_NAME).rb && git commit -m "Update $(BINARY_NAME) to v$(VERSION)" && git push
	@echo "Done! homebrew-tap updated."

# Full publish: release + gh-release + update-tap
# Usage: make publish VERSION=0.2.0
publish: release gh-release update-tap
	@echo "\n=== Published v$(VERSION) ==="
	@echo "1. GitHub release: https://github.com/zbum/amqp-cli/releases/tag/v$(VERSION)"
	@echo "2. Homebrew: brew upgrade amqp-cli"
