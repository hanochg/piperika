SHELL := /bin/bash

.DEFAULT_GOAL = help

GOCMD ?= go
TEST_TAGS ?= -tags=test
.DEFAULT_GOAL = build

help:				## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

# BUILD:

build: clean			## Build Piperika plugin
	$(GOCMD) build -o ./bin/piperika

build-multi-os: clean	## Build Piperika plugin for different OS types
	./multi-os-build.sh

fmt-fix:			## Gofmt fix errors
	gofmt -w -s .

vet:				## GoVet
	$(GOCMD) vet $(TEST_TAGS) ./...

clean:				## Clean from created bins
	rm -f ./bin/*

run:				## Run the plugin (example: ARGS="pp everything" make run)
	$(GOCMD) run main.go $(ARGS)

# PLUGIN INSTALLATION

install: clean build			## Install the plugin to jfrog cli
	mkdir -p ${HOME}/.jfrog/plugins
	\cp bin/piperika ${HOME}/.jfrog/plugins

uninstall:			## Uninstall the plugin to jfrog cli
	rm -f ${HOME}/.jfrog/plugins/piperika

.PHONY: help build fmt-fix vet clean run test test-list cover install uninstall