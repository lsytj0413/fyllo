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

# Current version of the project
VERSION ?= 1.0.0
GIT_SHA=$(shell git rev-parse --short HEAD)
TAGS=$(GIT_SHA)

# This repo's root import path (under GOPATH)
ROOT := github.com/lsytj0413/fyllo

# Target binaries. You can build multiple binaries for a single project
TARGETS := fyllo

# A list of all packages
PKGS := $(shell go list ./... | grep -v /vendor | grep -v /test)

# Project main package location (can be multiple ones)
CMD_DIR := ./cmd

# Project output directory
OUTPUT_DIR := ./bin

# Build directory
BUILD_DIR := ./build

# Git commit sha
COMMIT := $(shell git rev-parse --short HEAD)

# Golang standard bin directory
BIN_DIR := $(firstword $(subst :, ,$(GOPATH)))/bin
GOMETALINTER := $(BIN_DIR)/gometalinter
GODEP := $(BIN_DIR)/dep

#
# all targets
#

.PHONY: clean lint test build dep

all: test build

# TODO: if vendor exists skip ensure?
dep: $(GODEP)
	@if [ ! -d ./vendor ]; then                    \
	  	dep ensure;                               \
	else                                          \
	  	echo "vendor exists, skip dep ensure";    \
	fi
$(GODEP):
	go get -u -v github.com/golang/dep/cmd/dep
	
test: dep
#	go test $(PKGS)
	@for pkg in $(PKGS); do             \
	  go test $${pkg};                  \
	done

build: build-local

build-local: dep
	@for target in $(TARGETS); do                                     \
	  go build -i -v -o $(OUTPUT_DIR)/$${target}                      \
	   -ldflags "-s -w -X $(ROOT)/pkg/version.Version=$(VERSION)      \
	            -X $(ROOT)/pkg/version.Commit=$(COMMIT)"              \
	   $(CMD_DIR)/$${target};                                         \
	done		

build-docker: 
	@for target in $(TARGETS); do                                 \
	  docker build --force-rm -t $${target}:$(VERSION)                       \
	    -f $(BUILD_DIR)/$${target}/Dockerfile .;                  \
	done
	# docker rmi -f $(shell docker images -q --filter label=stage=intermediate)
	docker image prune -f

lint: $(GOMETALINTER)
	gometalinter ./... --vendor
$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

clean:
	-rm -vrf ${OUTPUT_DIR}