# Makefile which works as a simple task system for the project.
# Copyright 2018 Rafal Zajac <rzajac@gmail.com>.

BUILD_DATE = $(shell date -u +%Y%m%d.%H%M%S)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

VERSION_LONG = $(shell git describe --tags --dirty="-dev")
VERSION ?= $(shell git describe --tags --abbrev=0 --dirty="-dev")

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -X github.com/rzajac/iotdet/version.BuildDate=${BUILD_DATE}
LDFLAGS += -X github.com/rzajac/iotdet/version.Version=${GIT_TAG}
LDFLAGS += -X github.com/rzajac/iotdet/version.GitHash=${GIT_SHA}
LDFLAGS += -X github.com/rzajac/iotdet/version.GitTreeState=${GIT_DIRTY}

PWD = $(shell pwd)
DIST = $(PWD)/dist
CONFIG ?= $(PWD)/iotdet.yaml

help:
	@echo "Usage: make TARGET"

build:
	@cd cmd/iotdet && go build -race -ldflags "$(LDFLAGS)" -o $(DIST)/iotdet .

build-linux:
	@cd cmd/iotdet && CGO_ENABLED=0 GOOS=linux go build -race -a -installsuffix cgo -ldflags "$(LDFLAGS)" -o $(DIST)/iotdet .

lint:
	@golint -set_exit_status pkg/iotdet/ cmd/... 2>&1 | tee tmp/lint.txt

vet:
	@go tool vet -shadow=true pkg/ cmd/ 2>&1 | tee tmp/vet.txt

coverage:
	@gocov convert tmp/coverage.out | gocov-xml > tmp/coverage.xml

coverage-html:
	@go tool cover -html=tmp/coverage.out -o tmp/coverage.html

test:
	@(cd pkg/iotdet && go test -v -coverprofile=../../tmp/coverage.out 2>&1 | tee ../../tmp/tests.txt)
	@cat tmp/tests.txt | go2xunit -fail -output tmp/tests.xml

## Development targets.

build-dev:
	@echo "Building development version. DO NOT USE FOR PRODUCTION BUILDS!!!"
	@cd cmd/iotdet && go build -ldflags "$(LDFLAGS)" -o $(DIST)/iotdet .
	@echo "Development binary put in" $(DIST)/iotdet

start: build-dev
	@$(DIST)/iotdet start -c $(CONFIG)

.PHONY: help build-dev build start build-linux
