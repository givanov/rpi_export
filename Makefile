.EXPORT_ALL_VARIABLES:

.DEFAULT_GOAL := build

APP_NAME := rpi_exporter

BINDIR := bin

LDFLAGS := -extldflags "-static"

BUILD_PATH = github.com/givanov/rpi_export

GOOS ?= linux
GOARCH ?= arm64

HAS_GIT := $(shell command -v git;)

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

ifdef HAS_GIT
GIT_SHORT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_TAG    := $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

TMP_VERSION := canary

BINARY_VERSION := ""

ifndef VERSION
ifeq ($(GIT_DIRTY), clean)
ifdef GIT_TAG
	TMP_VERSION = $(GIT_TAG)
	BINARY_VERSION = $(GIT_TAG)
endif
endif
else
  BINARY_VERSION = $(VERSION)
endif

endif

VERSION ?= $(TMP_VERSION)

DIST_DIR := _dist
TARGETS ?= darwin/amd64 linux/amd64 windows/amd64 linux/arm64
TARGET_DIRS = find * -type d -exec

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),"")
	LDFLAGS += -X $(BUILD_PATH)/pkg/version.Version=$(VERSION)
	CHART_VERSION = $(VERSION)
endif

LDFLAGS += -X $(BUILD_PATH)/pkg/version.GitCommit=$(GIT_SHORT_COMMIT)

SHELL := /bin/sh

.PHONY: info
info:
	@echo "How are you:       $(GIT_DIRTY)"
	@echo "Version:           $(VERSION)"
	@echo "Git Tag:           $(GIT_TAG)"
	@echo "Git Commit:        $(GIT_SHORT_COMMIT)"
	@echo "binary:            $(BINARY_VERSION)"

build: clean-bin info bootstrap generate tidy fmt 
	@echo "build target..."
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINDIR)/$(APP_NAME) -ldflags '$(LDFLAGS)' ./main.go


.PHONY: clean-bin
clean-bin: 
	@rm -rf $(BINDIR)

.PHONY: clean-dist
clean-dist:
	@rm -rf $(DIST_DIR)

.PHONY: clean
clean: clean-bin clean-dist

.PHONY: tidy
tidy:
	@echo "tidy target..."
	@go mod tidy

.PHONY: generate
generate: bootstrap
	@echo "generate target..."
	@rm -rf ./pkg/mocks
	@go generate ./...

.PHONY: vendor
vendor: tidy
	@echo "vendor target..."
	@go mod vendor

.PHONY: test
test: generate build
	@echo "test target..."
	@go test ./... -v -count=1

.PHONY: bootstrap
bootstrap: 
	@echo "bootstrap target..."

.PHONY: fmt
fmt:
	@echo "fmt target..."
	@gofmt -l -w -s $(SRC)

# Semantic Release
.PHONY: semantic-release-dependencies
semantic-release-dependencies:
	@npm install --save-dev semantic-release
	@npm install @semantic-release/exec conventional-changelog-conventionalcommits -D

.PHONY: semantic-release
semantic-release: semantic-release-dependencies
	@npm ci
	@npx semantic-release

.PHONY: semantic-release-ci
semantic-release-ci: semantic-release-dependencies
	@npx semantic-release

.PHONY: semantic-release-dry-run
semantic-release-dry-run: semantic-release-dependencies
	@npm ci
	@npx semantic-release -d

.PHONY: export-tag-github-actions
export-tag-github-actions:
	@echo "version=$(VERSION)" >> $${GITHUB_OUTPUT}