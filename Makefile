SHELL := /bin/bash

.PHONY: build
build: deps
	go build ./...

.PHONY: test
test:
	go test -race -v ./...
