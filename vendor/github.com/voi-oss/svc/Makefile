PKGS ?= ./...

PHONY: all
all: lint test

PHONY: lint
lint:
	golangci-lint run --config .golangci.yaml

.PHONY: test
test:
	go test -race $(PKGS) -coverprofile=coverage.txt -covermode=atomic

.PHONY: tools
tools:
	# Install golangci for linting. Installer copied from project page.
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.21.0

.PHONY: vendor
vendor:
	go mod vendor
