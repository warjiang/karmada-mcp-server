# Basic variables
GOOS 					?= $(shell go env GOOS)
GOARCH 					?= $(shell go env GOARCH)
VERSION 				?= $(shell hack/version.sh)

# Default target when just running 'make'
.DEFAULT_GOAL := all

# Build targets
TARGETS := karmada-mcp-server

# Docker image related variables
REGISTRY				?= docker.io/karmada
REGISTRY_USER_NAME  	?= 
REGISTRY_PASSWORD   	?= 
REGISTRY_SERVER_ADDRESS ?= 
IMAGE_TARGET := $(addprefix image-, $(TARGETS))


###################
# Build Targets   #
###################

# Build all binaries (alias for build)
.PHONY: all
all: build

# Build all binaries
.PHONY: build
build: $(TARGETS)

# Build specific binary
.PHONY: $(TARGETS)
$(TARGETS):
	BUILD_PLATFORMS=$(GOOS)/$(GOARCH) hack/build.sh $@

