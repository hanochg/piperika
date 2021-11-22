SHELL := /bin/bash

.DEFAULT_GOAL = help

GOCMD ?= go
TEST_TAGS ?= -tags=test
.DEFAULT_GOAL = build

help:				## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

# BUILD:

build: clean			## Build Piperika plugin
	$(GOCMD) build -o ./build/piprika

fmt-fix:			## Gofmt fix errors
	gofmt -w -s .

vet:				## GoVet
	$(GOCMD) vet $(TEST_TAGS) ./...

clean:				## Clean from created bins
	rm -f ./build/*

run:				## Run the plugin
	$(GOCMD) run main.go


# TEST EXECUTION

test:				## Run all tests
	time $(GOCMD) test ./... $(TEST_TAGS) -count=1

test-list:			## List all tests
	$(GOCMD) list ./...

cover:				## Shows coverage details
	$(GOCMD) test ./... $(TEST_TAGS) -count=1 -coverprofile=coverage


# PLUGIN INSTALLATION

install:			## Install the plugin to jfrog cli
	jfrog plugin install forest

uninstall:			## Uninstall the plugin to jfrog cli
	jfrog plugin uninstall forest

.PHONY: help build fmt-fix vet clean run test test-list cover install uninstall