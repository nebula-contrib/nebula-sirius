.PHONY: download-thrift-files, generate-from-contracts, fmt, lint, unit, docker-compose-up, docker-compose-down, install-mockery, generate-mocks, generate-mock-thrift-transport, clean-mocks

fmt:
	go fmt ./...

lint:
	@test -z `gofmt -l *.go` || (echo "Please run 'make fmt' to format Go code" && exit 1)

m1:
	grep -vE '.*/mocks/.*|.*/nebulagraph-light-deployment/.*|.*/nebula/.*|.*/examples/.*' coverage.out > coverage.filtered.out
	go tool cover -func=coverage.filtered.out
unit:
	go mod tidy
	go test ./... -v -race --covermode=atomic --coverprofile coverage.out -cover -json
	go tool cover -func=coverage.out
	#go tool cover -html=coverage.out -o coverage.html


########################################################################
#######  BEGIN: DOCKER COMPOSE PART
docker-compose-up:
	cd nebulagraph-light-deployment && docker-compose -f docker-compose-lite.yml up -d

docker-compose-down:
	cd nebulagraph-light-deployment && docker-compose -f docker-compose-lite.yml down

docker-compose-up-ssl:
	cd nebulagraph-light-deployment && enable_ssl=true docker-compose -f docker-compose-lite-ssl.yml up -d

docker-compose-down-ssl:
	cd nebulagraph-light-deployment && docker-compose -f docker-compose-lite-ssl.yml down

########################################################################
#######  END: DOCKER COMPOSE PART


########################################################################
#######  BEGIN: THRIFT FILES SYNC PART
########################################################################
THRIFT_FILES_DIR = ./thriftfiles
REMOTE_THRIFT_FILES_URL = https://raw.githubusercontent.com/vesoft-inc/nebula/master/src/interface

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

# Default target
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
# Install mockery if not already installed
install-mockery:
		GOBIN=$(go env GOPATH)/bin go install github.com/vektra/mockery/v2@latest

# Generate mocks for all interfaces
generate-mocks: install-mockery
		$(go env GOPATH)/bin/mockery --all --recursive --output=./mocks --case=underscore

# Generate mock for specific interface thrift.TTransport
generate-mock-thrift-transport: install-mockery
		$(go env GOPATH)/bin/mockery --name=TTransport --dir=$(go env GOPATH)/pkg/mod/github.com/apache/thrift@v0.21.0/lib/go/thrift --output=./mocks --case=underscore

# Clean generated mocks
clean-mocks:
		rm -rf ./mocks

# Add other targets as needed
########################################################################
#######  END: MOCK GENERATION PART
########################################################################
