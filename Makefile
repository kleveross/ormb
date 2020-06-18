# Copyright 2019 The Caicloud Authors.
#
# The old school Makefile, following are required targets. The Makefile is written
# to allow building multiple binaries. You are free to add more targets or change
# existing implementations, as long as the semantics are preserved.
#
#   make              - default to 'build' target
#   make lint         - code analysis
#   make test         - run unit test (or plus integration test)
#   make build        - alias to build-local target
#   make build-local  - build local binary targets
#   make build-linux  - build linux binary targets
#   make container    - build containers
#   $ docker login registry -u username -p xxxxx
#   make push         - push containers
#   make clean        - clean up targets
#
# Not included but recommended targets:
#   make e2e-test
#
# The makefile is also responsible to populate project version information.
#

#
# Tweak the variables based on your project.
#

# This repo's root import path (under GOPATH).
ROOT := github.com/caicloud/ormb

# Target binaries. You can build multiple binaries for a single project.
TARGETS := ormb ormb-storage-initializer

# Container image prefix and suffix added to targets.
# The final built images are:
#   $[REGISTRY]/$[IMAGE_PREFIX]$[TARGET]$[IMAGE_SUFFIX]:$[VERSION]
# $[REGISTRY] is an item from $[REGISTRIES], $[TARGET] is an item from $[TARGETS].
IMAGE_PREFIX ?= $(strip clever-)
IMAGE_SUFFIX ?= $(strip )

# Container registries.
REGISTRY ?= cleveross

# Container registry for base images.
BASE_REGISTRY ?= docker.io

#
# These variables should not need tweaking.
#

# It's necessary to set this because some environments don't link sh -> bash.
export SHELL := /bin/bash

# It's necessary to set the errexit flags for the bash shell.
export SHELLOPTS := errexit

# Project main package location (can be multiple ones).
CMD_DIR := ./cmd

# Project output directory.
OUTPUT_DIR := ./bin

# Build direcotory.
BUILD_DIR := ./build

# Current version of the project.
VERSION ?= $(shell git describe --tags --always --dirty)
GITSHA ?= $(shell git rev-parse --short HEAD)

# Available cpus for compiling, please refer to https://github.com/caicloud/engineering/issues/8186#issuecomment-518656946 for more information.
CPUS ?= $(shell /bin/bash hack/read_cpus_available.sh)

# Track code version with Docker Label.
DOCKER_LABELS ?= git-describe="$(shell date -u +v%Y%m%d)-$(shell git describe --tags --always --dirty)"

# Golang standard bin directory.
GOPATH ?= $(shell go env GOPATH)
BIN_DIR := $(GOPATH)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

# Default golang flags used in build and test
# -mod=vendor: force go to use the vendor files instead of using the `$GOPATH/pkg/mod`
# -p: the number of programs that can be run in parallel
# -race: enable data race detection
# -count: run each test and benchmark 1 times. Set this flag to disable test cache
export GOFLAGS ?= -mod=vendor -p=$(CPUS) -race -count=1

#
# Define all targets. At least the following commands are required:
#

# All targets.
.PHONY: lint test build container push

build: build-local

# more info about `GOGC` env: https://github.com/golangci/golangci-lint#memory-usage-of-golangci-lint
lint: $(GOLANGCI_LINT)
	@$(GOLANGCI_LINT) run

$(GOLANGCI_LINT):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_DIR) v1.23.6

generate:
	@mockgen -source pkg/oras/interface.go -destination pkg/oras/mock/mock.go -package mock
	@mockgen -source pkg/oras/cache/interface.go -destination pkg/oras/cache/mock/mock.go -package mock
	@mockgen -source pkg/oras/orasclient/interface.go -destination pkg/oras/orasclient/mock/mock.go -package mock
	@mockgen -source pkg/saver/interface.go -destination pkg/saver/mock/mock.go -package mock
	@mockgen -source pkg/exporter/interface.go -destination pkg/exporter/mock/mock.go -package mock
	@mockgen -source vendor/github.com/deislabs/oras/pkg/auth/client.go -destination pkg/oras/mock/mock_auth.go -package mock
	@mockgen -source vendor/github.com/containerd/containerd/remotes/resolver.go -destination pkg/oras/mock/mock_resolver.go -package mock

test: generate
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func coverage.out | tail -n 1 | awk '{ print "Total coverage: " $$3 }'

build-local:
	@for target in $(TARGETS); do                                                      \
	  go build -v -o $(OUTPUT_DIR)/$${target}                                          \
	    -ldflags "-s -w -X $(ROOT)/pkg/version.VERSION=$(VERSION)                      \
		  -X $(ROOT)/pkg/version.COMMIT=$(GITSHA)                                      \
	      -X $(ROOT)/pkg/version.REPOROOT=$(ROOT)"                                     \
	    $(CMD_DIR)/$${target};                                                         \
	done

build-linux:
	@docker run --rm -t                                                                \
	  -v $(PWD):/go/src/$(ROOT)                                                        \
	  -w /go/src/$(ROOT)                                                               \
	  -e GOOS=linux                                                                    \
	  -e GOARCH=amd64                                                                  \
	  -e GOPATH=/go                                                                    \
	  -e GOFLAGS="$(GOFLAGS)"                                                          \
	  -e SHELLOPTS="$(SHELLOPTS)"                                                      \
	  $(BASE_REGISTRY)/golang:1.13.9-stretch                                           \
	    /bin/bash -c 'for target in $(TARGETS); do                                     \
	      go build -v -o $(OUTPUT_DIR)/$${target}                                      \
	        -ldflags "-s -w -X $(ROOT)/pkg/version.VERSION=$(VERSION)                  \
	          -X $(ROOT)/pkg/version.REPOROOT=$(ROOT)"                                 \
	        $(CMD_DIR)/$${target};                                                     \
	    done'

container-buildlocal: build-local container-skipbuild

container: build-linux container-skipbuild

push: container
	@for target in $(TARGETS); do                                                      \
	  image=$(IMAGE_PREFIX)$${target}$(IMAGE_SUFFIX);                                  \
	  docker push $(REGISTRY)/$${image}:$(VERSION);                                    \
	done

container-skipbuild:
	@for target in $(TARGETS); do                                                      \
	  image=$(IMAGE_PREFIX)$${target}$(IMAGE_SUFFIX);                                  \
	  docker build -t $(REGISTRY)/$${image}:$(VERSION)                                 \
	    --label $(DOCKER_LABELS)                                                       \
	    -f $(BUILD_DIR)/$${target}/Dockerfile .;                                       \
	done

.PHONY: clean
clean:
	@-rm -vrf ${OUTPUT_DIR}
