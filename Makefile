# Use bash as the default shell
SHELL=/usr/bin/env bash

# Define variables
BINARY := filecoin-check
COMMIT := $(shell git rev-parse --short HEAD)
COMMIT_TIMESTAMP := $(shell git log -1 --format=%ct)
VERSION := $(shell git describe --tags --abbrev=0)
GO_BIN := go

# Build flags
CGO_ENABLED := 0
FLAGS := -trimpath
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.CurrentCommit=$(COMMIT)

# Declare phony targets
.PHONY: all build run gotool clean help linux-amd64 linux-arm64 linux-arm linux-386 windows-amd64 windows-arm64 windows-386 darwin-amd64 darwin-arm64 build-all

# Default target: tidy code and compile for current environment
all:  build

# Default build: current environment
build:
	 $(GO_BIN) build -o $(BINARY) $(FLAGS) -ldflags "$(LDFLAGS)"

# Clean
clean:
	@if [ -f $(BINARY) ]; then rm -f $(BINARY); fi

linux:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO_BIN) build -o $(BINARY)_linux_amd64 $(FLAGS) -ldflags "$(LDFLAGS)"

# Help information
help:
	@echo "make              - Compile binary for current environment"
