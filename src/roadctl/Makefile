
#-include .env

VERSION := 1.0.0
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")
PROJDIR := $(shell pwd)
TARGET := $(PROJECTNAME)

ASSETS := $(PROJDIR)/assets/images
ARTIFACTS := $(PROJDIR)/artifacts
BUILDS := $(PROJDIR)/builds
DOCS := $(PROJDIR)/docs
LOGS := $(PROJDIR)/logs

# Go related variables.
GOBASE := $(shell cd ../../;pwd)
GOPATH := $(GOBASE)
export GOPATH = $(GOBASE)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)
GOLINT := $(shell which golint)
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
GOCOVERAGE := $(ARTIFACTS)/coverage.out
GOLINTREPORT := $(ARTIFACTS)/lint.out
GOSECREPORT := $(ARTIFACTS)/gosec.out
GOVETREPORT := $(ARTIFACTS)/govet.out
GOTESTREPORT := https://sonarcloud.io/dashboard?id=acme_films

GIT_TAG := $(shell git describe)

SHELL := /bin/bash

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -X=main.GitTag=$(GIT_TAG)"

# Make is verbose in Linux. Make it silent.
# MAKEFLAGS += --silent

.PHONY: check build compile sonar-scanner

all: compile check

## compile: Compile the binary.
compile: $(LOGS) $(ARTIFACTS) $(ASSETS) $(DOCS) $(BUILDS)
	@echo "  Compiling"
	@-$(MAKE) -s build

## clean: Remove dep, vendor, binary(s), and executs go clean
clean:
	@echo "  execute go-clean"
	@-rm $(GOBIN)/$(PROJECTNAME)* 2> /dev/null || true
	@-rm -R vendor Gopkt.* 2> /dev/null || true
	@-$(MAKE) go-clean

## build: Build the binary for linux / mac x86 and amd
build: go-get build
	@echo "  >  Building binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME)-$(GOOS)-$(GOARCH) $(GOFILES)
# make this conditional on build GOARCH
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOOS="darwin" GOARCH="amd64" go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME)-"darwin"-"amd64" $(GOFILES)
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOOS="darwin" GOARCH="386" go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME)-"darwin"-"386" $(GOFILES)
	cp $(GOBIN)/$(PROJECTNAME)-$(GOOS)-$(GOARCH) $(BUILDS)/$(PROJECTNAME)-$(GOOS)-$(GOARCH)
	cp $(BUILDS)/$(PROJECTNAME)-$(GOOS)-$(GOARCH) $(PROJECTNAME)
	cp $(GOBIN)/$(PROJECTNAME)-"darwin"-"amd64" $(BUILDS)/$(PROJECTNAME)-"darwin"-"amd64"
	cp $(GOBIN)/$(PROJECTNAME)-"darwin"-"386" $(BUILDS)/$(PROJECTNAME)-"darwin"-"386"


Gopkg.toml:
	@echo "  >  initialize dep support..."
	$(shell (export GOPATH=$(GOPATH);dep init))

go-get: Gopkg.toml get-deps $(ASSETS)
	@echo "  >  Creating dependencies graph png..."
	$(shell (export GOPATH=$(GOPATH);dep status -dot | dot -T png -o $(ASSETS)/$(PROJECTNAME).png))

get-deps:
	@echo "  >  dep ensure..."
	$(shell (GOPATH=$(GOPATH);dep ensure $?))

## install: Install packages or main
install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

# TODO: enable sonar scanner once we get API key
# check: lint sonar-scanner $(ARTIFACTS) $(LOGS) $(ASSETS) $(DOCS)
## check: Start services and execute static code analysis and tests
check: lint $(ARTIFACTS) $(LOGS) $(ASSETS) $(DOCS)
	@echo "  >  running to tests..."
#	go test -coverprofile=$(GOCOVERAGE) -v ./...

sonar-scanner: $(ARTIFACTS)
	sonarcloud.sh

## show-coverage: Show go code coverage in browser
show-coverage:
	go tool cover -html=$(GOCOVERAGE)

## show-test: Show sonarcloud test report
show-test:
	xdg-open $(GOTESTREPORT)

lint: $(GOFILES)
	@echo -n "  >  running lint..."
	@echo $?
	$(GOLINT) ./... > $(GOLINTREPORT)
	@echo "  >  running gosec... > $(GOSECREPORT)"
	$(shell (export GOPATH=$(GOPATH);gosec -fmt=sonarqube -tests -out $(GOSECREPORT) -exclude-dir=.templates ./...))
	@echo "  >  running go vet... > $(GOVETREPORT)"
	$(shell (export GOPATH=$(GOPATH);go vet ./... 2> $(GOVETREPORT)))

## fmt: Run gofmt on all code
fmt: $(GOFILES)
	@gofmt -l -w $?

## simplify: Run gofmt with simplify option
simplify: $(GOFILES)
	@gofmt -s -l -w $?

## help: Print possible commands
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

$(ASSETS):
	@echo "  >  Creating assets directory"
	$(shell mkdir -p $(ASSETS))

$(ARTIFACTS):
	@echo "  >  Creating artifacts directory"
	$(shell mkdir -p $(ARTIFACTS))

$(BUILDS):
	@echo "  >  Creating $(BUILDS) directory"
	$(shell mkdir -p $(BUILDS))

$(DOCS):
	@echo "  >  Creating docs directory"
	$(shell mkdir -p $(DOCS))

$(LOGS):
	@echo "  >  Creating logs directory"
	$(shell mkdir -p $(LOGS))
