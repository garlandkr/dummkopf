#!/usr/bin/make -f

.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

name := dummkopf
ecr := 920858517028.dkr.ecr.us-east-1.amazonaws.com
script_dir := script
OUT := $(name)
VERSION := 0.0.1

## Builds the application
$(OUT): $(wildcard *.go)
	CGO_ENABLED=0 go build -o $@

.PHONY: clean
## Cleans build artifacts
clean:
	rm -rf $(OUT)

.PHONY: container
container:
	DOCKER_BUILDKIT=1 docker build --tag $(name):$(VERSION) .
	docker tag $(name):$(VERSION) $(ecr)/$(name):$(VERSION)
	docker push $(ecr)/$(name):$(VERSION)

.PHONY: test
## Runs tests
test:
	go test

.PHONY: lint
## Runs linter
lint:
	# TODO

.PHONY: help
## Shows help (you are here!)
help:
	@echo "Usage:"
	@echo "  make [<target>] [VERSION=<version>]"
	@echo
	@echo "Targets:"
	@awk -f "$(script_dir)"/make-help.awk ./Makefile
	@echo
	@echo "Options:"
	@echo "  VERSION		Build version"
