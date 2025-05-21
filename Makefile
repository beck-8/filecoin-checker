# 使用 bash 作为默认 shell
SHELL=/usr/bin/env bash

# 定义变量
BINARY := filecoin-check
COMMIT := $(shell git rev-parse --short HEAD)
COMMIT_TIMESTAMP := $(shell git log -1 --format=%ct)
VERSION := $(shell git describe --tags --abbrev=0)
GO_BIN := go

# 构建标志
CGO_ENABLED := 0
FLAGS := -trimpath
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.CurrentCommit=$(COMMIT)

# 声明伪目标
.PHONY: all build run gotool clean help linux-amd64 linux-arm64 linux-arm linux-386 windows-amd64 windows-arm64 windows-386 darwin-amd64 darwin-arm64 build-all

# 默认目标：整理代码并编译当前环境
all:  build

# 默认构建：当前环境
build:
	 $(GO_BIN) build -o $(BINARY) $(FLAGS) -ldflags "$(LDFLAGS)"

# 清理
clean:
	@if [ -f $(BINARY) ]; then rm -f $(BINARY); fi

linux:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO_BIN) build -o $(BINARY)_linux_amd64 $(FLAGS) -ldflags "$(LDFLAGS)"

# 帮助信息
help:
	@echo "make              - 编译当前环境的二进制文件"
