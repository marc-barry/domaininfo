#!/usr/bin/env make

build:
	go build ./

fmt:
	go fmt ./...

lint:
	golangci-lint run

test:
	CGO_ENABLED=1 go test -race ./... -count 1

vet:
	go vet ./...

.PHONY: build fmt lint test
