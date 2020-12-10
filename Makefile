#!/usr/bin/env make

build:
	go build -o ./bin/domaininfo ./cmd/domaininfo

fmt:
	go fmt ./...

lint:
	golangci-lint run

test:
	CGO_ENABLED=1 go test -race ./... -count 1

tidy:
	go mod tidy

vet:
	go vet ./...

.PHONY: build fmt lint test
