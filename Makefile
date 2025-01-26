.PHONY: generate-from-contracts

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
