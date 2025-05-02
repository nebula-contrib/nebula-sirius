# Project metadata
MODULE_NAME := nebula-sirius
PKG := ./...
BIN_DIR := bin
COVERAGE_FILE := coverage.out

# Tools
GOLINT := golangci-lint
GOFMT := gofmt
GODOC := godoc

# Tool installation stage
.PHONY: tools
tools:
	@echo "🔧 Checking required tools..."
	@if ! command -v gofmt >/dev/null 2>&1; then \
		echo "Installing gofmt..."; \
		go install golang.org/x/tools/cmd/gofmt@latest; \
	else \
		echo "✔ gofmt found"; \
	fi
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	else \
		echo "✔ golangci-lint found"; \
	fi
	@if ! command -v godoc >/dev/null 2>&1; then \
		echo "Installing godoc..."; \
		go install golang.org/x/tools/cmd/godoc@latest; \
	else \
		echo "✔ godoc found"; \
	fi
	@if ! command -v mockery >/dev/null 2>&1; then \
		echo "Installing mockery..."; \
		go install github.com/vektra/mockery/v2@latest; \
	else \
		echo "✔ mockery found"; \
	fi

.PHONY: fmt
fmt:
	@echo "🧹 Formatting Go code (excluding nebula/ and mocks/)..."
	@find . -type f -name '*.go' \
		! -path './nebula/*' \
		! -path './mocks/*' \
		-exec gofmt -s -w {} +

# Lint code
.PHONY: lint
lint:
	@echo "Linting..."
	@$(GOLINT) run

# Run `go mod tidy`
.PHONY: tidy
tidy:
	@echo "Tidying modules..."
	@go mod tidy

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=$(COVERAGE_FILE) $(PKG)

# Clean generated files
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR) $(COVERAGE_FILE)

# Display coverage in terminal
.PHONY: cover
cover:
	@go tool cover -func=$(COVERAGE_FILE)

# Serve documentation locally
.PHONY: doc
doc:
	@echo "Starting documentation server at http://localhost:6060"
	@godoc -http=:6060


########################################################################
#######  BEGIN: DOCKER COMPOSE PART
.PHONY: docker-compose-up
docker-compose-up:
	cd nebulagraph-light-deployment && docker-compose -f docker-compose-lite.yml up -d

.PHONY: docker-compose-down
docker-compose-down:
	cd nebulagraph-light-deployment && docker-compose -f docker-compose-lite.yml down

.PHONY: docker-compose-up-ssl
docker-compose-up-ssl:
	cd nebulagraph-light-deployment && enable_ssl=true docker-compose -f docker-compose-lite-ssl.yml up -d

.PHONY: docker-compose-down-ssl
docker-compose-down-ssl:
	cd nebulagraph-light-deployment && docker-compose -f docker-compose-lite-ssl.yml down

########################################################################
#######  END: DOCKER COMPOSE PART


########################################################################
#######  BEGIN: THRIFT FILES SYNC PART
########################################################################
THRIFT_FILES_DIR = ./thriftfiles
REMOTE_THRIFT_FILES_URL = https://raw.githubusercontent.com/vesoft-inc/nebula/master/src/interface

.PHONY: download-thrift-files
download-thrift-files: $(THRIFT_FILES_DIR)
		curl -s -o $(THRIFT_FILES_DIR)/common.thrift $(REMOTE_THRIFT_FILES_URL)/common.thrift \
		curl -s -o $(THRIFT_FILES_DIR)/meta.thrift $(REMOTE_THRIFT_FILES_URL)/meta.thrift \
		curl -s -o $(THRIFT_FILES_DIR)/graph.thrift $(REMOTE_THRIFT_FILES_URL)/graph.thrift \
		curl -s -o $(THRIFT_FILES_DIR)/storage.thrift $(REMOTE_THRIFT_FILES_URL)/storage.thrift \
		curl -s -o $(THRIFT_FILES_DIR)/raftex.thrift $(REMOTE_THRIFT_FILES_URL)/raftex.thrift \

########################################################################
#######  END: THRIFT FILES SYNC PART
########################################################################


########################################################################
#######  BEGIN: APACHE THRIFT CODE GENERATION PART
########################################################################
# Variables
THRIFT_DIR := thriftfiles
OUTPUT_DIR := nebula
PACKAGE_PREFIX := nebula

# Find all Thrift files in the directory
THRIFT_FILES := $(wildcard $(THRIFT_DIR)/*.thrift)

.PHONY: generate-from-contracts
generate-from-contracts:
		thrift --gen go:package_prefix=github.com/nebula-contrib/nebula-sirius/ -out  . thriftfiles/graph.thrift \

		thrift --gen go:package_prefix=github.com/nebula-contrib/nebula-sirius/ -out . thriftfiles/meta.thrift \

		thrift --gen go:package_prefix=github.com/nebula-contrib/nebula-sirius/ -out . thriftfiles/storage.thrift \

		thrift --gen go:package_prefix=github.com/nebula-contrib/nebula-sirius/ -out . thriftfiles/common.thrift

########################################################################
#######  END: APACHE THRIFT CODE GENERATION PART
########################################################################


########################################################################
#######  BEGIN: MOCK GENERATION PART
########################################################################

.PHONY: install-mockery
install-mockery:
		GOBIN=$(go env GOPATH)/bin go install github.com/vektra/mockery/v2@latest

# Generate mocks for all interfaces
.PHONY: generate-mocks
generate-mocks:
		$(go env GOPATH)/bin/mockery --all --recursive --output=./mocks --case=underscore

# Generate mock for specific interface thrift.TTransport
.PHONY: generate-mock-thrift-transport
generate-mock-thrift-transport: install-mockery
		$(go env GOPATH)/bin/mockery --name=TTransport --dir=$(go env GOPATH)/pkg/mod/github.com/apache/thrift@v0.21.0/lib/go/thrift --output=./mocks --case=underscore

# Clean generated mocks
.PHONY: clean-mocks
clean-mocks:
		@echo "Cleaning mocks..."
		@rm -rf ./mocks

# Add other targets as needed
########################################################################
#######  END: MOCK GENERATION PART
########################################################################
