SHELL := /bin/sh
BINARY_NAME := kubectl-permissions
REV := $(shell git rev-parse --short HEAD 2> /dev/null || echo 'unknown')
GIT_TREE_STATE := $(shell test -z "`git status --porcelain`" && echo "clean" || echo "dirty")

VERSION ?= $(shell echo "$$(git for-each-ref refs/tags/ --count=1 --sort=-version:refname --format='%(refname:short)' 2>/dev/null)-dev+$(REV)-$(GIT_TREE_STATE)" | sed 's/^v//')

build:
	go build -o $(BINARY_NAME) -trimpath -ldflags "-X github.com/garethjevans/permissions/pkg/version.Version=$(VERSION)" cmd/kubectl-permissions.go

install: build
	sudo cp -f $(BINARY_NAME) /usr/local/bin

lint:
	golangci-lint run

test:
	go test -v -short ./...

.PHONY: integration
integration:
	go test -run Integration ./integration/...
