# Copyright (c) 2018 soren yang
#
# Licensed under the MIT License
# you may not use this file except in complicance with the License.
# You may obtain a copy of the License at
#
#     https://opensource.org/licenses/MIT
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This repo's root import path (under GOPATH).
ROOT := github.com/lsytj0413/fyllo

# Target binaries. You can build multiple binaries for a single project.
TARGETS := fyllo

# Container image prefix and suffix added to targets.
# The final built images are:
#   $[REGISTRY]/$[IMAGE_PREFIX]$[TARGET]$[IMAGE_SUFFIX]:$[VERSION]
# $[REGISTRY] is an item from $[REGISTRIES], $[TARGET] is an item from $[TARGETS].
IMAGE_PREFIX ?= $(strip )
IMAGE_SUFFIX ?= $(strip )

# Container registries.
REGISTRY ?= 

# repository prefix for image.
REGISTRY_PREFIX ?= $(REGISTRY)
ifneq ($(REGISTRY),)
REGISTRY_PREFIX := $(REGISTRY_PREFIX)/
endif

# Container registry for base images.
BASE_REGISTRY ?= 

# image for build executable file.
BASE_BUILD_IMAGE ?= golang:1.12.9-alpine3.10
ifneq ($(BASE_REGISTRY),)
BASE_BUILD_IMAGE := $(BASE_REGISTRY)/$(BASE_BUILD_IMAGE)
endif

#
# These variables should not need tweaking.
#

# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /bin/bash

# Project main package location (can be multiple ones).
CMD_DIR := ./cmd

# Project output directory.
OUTPUT_DIR := ./bin

# Build direcotory.
BUILD_DIR := ./build

# Current version of the project.
VERSION ?= $(shell git describe --tags --always --dirty)

# Available cpus for compiling.
CPUS ?= $(shell sh hack/read_cpus_available.sh)

# Track code version with Docker Label.
DOCKER_LABELS ?= git-describe="$(shell date -u +v%Y%m%d)-$(shell git describe --tags --always --dirty)"

# Golang standard bin directory.
GOPATH ?= $(shell go env GOPATH)
BIN_DIR := $(firstword $(subst :, , $(GOPATH)))/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

#
# Define all targets. At least the following commands are required:
#

# All targets.
.PHONY: lint test build container push

build: build-local

# more info about `GOGC` env: https://github.com/golangci/golangci-lint#memory-usage-of-golangci-lint
lint: $(GOLANGCI_LINT)
	@GOGC=5 $(GOLANGCI_LINT) run

$(GOLANGCI_LINT):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_DIR) v1.16.0

# TODO: if vendor exists skip ensure?
dep: $(GODEP)
	@if [ ! -d ./vendor ]; then                    \
	  	dep ensure;                                \
	else                                           \
	  	echo "vendor exists, skip dep ensure";     \
	fi
$(GODEP):
	go get -u -v github.com/golang/dep/cmd/dep

test:
	@go test -p $(CPUS) $$(go list ./... | grep -v /vendor | grep -v /test) -coverprofile=coverage.out
	@go tool cover -func coverage.out | tail -n 1 | awk '{ print "Total coverage: " $$3 }'

build-local:
	@for target in $(TARGETS); do                                                      \
	  go build -i -v -o $(OUTPUT_DIR)/$${target} -p $(CPUS)                            \
	  -ldflags "-s -w -X $(ROOT)/pkg/version.VERSION=$(VERSION)                        \
	    -X $(ROOT)/pkg/version.REPOROOT=$(ROOT)"                                       \
	  $(CMD_DIR)/$${target};                                                           \
	done

build-linux:
	@for target in $(TARGETS); do                                                      \
	  docker run --rm                                                                  \
	    -v $(PWD):/go/src/$(ROOT)                                                      \
	    -w /go/src/$(ROOT)                                                             \
	    -e GOOS=linux                                                                  \
	    -e GOARCH=amd64                                                                \
	    -e GOPATH=/go                                                                  \
	    $(BASE_BUILD_IMAGE)                                         \
	      go build -i -v -o $(OUTPUT_DIR)/$${target} -p $(CPUS)                        \
	        -ldflags "-s -w -X $(ROOT)/pkg/version.VERSION=$(VERSION)                  \
	          -X $(ROOT)/pkg/version.REPOROOT=$(ROOT)"                                 \
	        $(CMD_DIR)/$${target};                                                     \
	done

container: build-linux
	@for target in $(TARGETS); do                                                      \
	  image=$(IMAGE_PREFIX)$${target}$(IMAGE_SUFFIX);                                  \
	  docker build -t $(REGISTRY)$${image}:$(VERSION)                                 \
	    --label $(DOCKER_LABELS)                                                       \
	    -f $(BUILD_DIR)/$${target}/Dockerfile .;                                       \
	done

push: container
	@for target in $(TARGETS); do                                                      \
	  image=$(IMAGE_PREFIX)$${target}$(IMAGE_SUFFIX);                                  \
	  docker push $(REGISTRY)$${image}:$(VERSION);                                    \
	done

.PHONY: clean
clean:
	@-rm -vrf ${OUTPUT_DIR}