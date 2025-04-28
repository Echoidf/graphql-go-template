# 检测操作系统和架构
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./bin/server

run: build
	./bin/server

test: 
	go test -v ./...