.PHONY: generate-from-contracts, fmt, lint, unit, docker-compose-up, docker-compose-down

fmt:
	go fmt ./...

lint:
	@test -z `gofmt -l *.go` || (echo "Please run 'make fmt' to format Go code" && exit 1)

unit:
	go mod tidy
	go test -v -race --covermode=atomic --coverprofile coverage.out -cover -json
	go tool cover -func=coverage.out


########################################################################
#######  BEGIN: DOCKER COMPOSE PART
docker-compose-up:
	cd nebulagraph-light-deployment && docker-compose -f docker-compose-lite.yml up -d

docker-compose-down:
	cd nebulagraph-light-deployment && docker-compose -f docker-compose-lite.yml down

docker-compose-up-ssl:
	cd nebulagraph-light-deployment && docker-compose -f docker-compose-lite-ssl.yml up -d

docker-compose-down-ssl:
	cd nebulagraph-light-deployment && docker-compose -f docker-compose-lite-ssl.yml down

########################################################################
#######  END: DOCKER COMPOSE PART


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
		thrift --gen go:package_prefix=github.com/egasimov/nebula-go-sdk/ -out  . thriftfiles/graph.thrift \

		thrift --gen go:package_prefix=github.com/egasimov/nebula-go-sdk/ -out . thriftfiles/meta.thrift \

		thrift --gen go:package_prefix=github.com/egasimov/nebula-go-sdk/ -out . thriftfiles/storage.thrift \

		thrift --gen go:package_prefix=github.com/egasimov/nebula-go-sdk/ -out . thriftfiles/common.thrift

########################################################################
#######  END: APACHE THRIFT CODE GENERATION PART
########################################################################
