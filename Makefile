.PHONY: all
all: audit test lint build

audit:
	go list -json -m all | nancy sleuth
.PHONY: audit

build:
	go build ./...
.PHONY: build

test:
	go test -race -cover ./...
.PHONY: test

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2
	golangci-lint run ./...
.PHONY: lint
