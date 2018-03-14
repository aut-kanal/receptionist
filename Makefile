receptionistd: *.go */*.go */*/*.go Gopkg.lock
	$(GO_VARS) $(GO) build -i -o="receptionistd" -ldflags="$(LD_FLAGS)" $(ROOT)/cmd/receptionist

clean:
	rm -rf receptionistd

help:
	@echo "Please use \`make <ROOT>' where <ROOT> is one of"
	@echo "  receptionistd     to build the main binary for current platform"
	@echo "  clean             to remove generated files"


## Project Vars ##########################################################
ROOT := gitlab.com/kanalbot/receptionist
.PHONY: help clean


## Commons Vars ##########################################################
GO_VARS = ENABLE_CGO=0
GO ?= go
GIT ?= git
COMMIT := $(shell $(GIT) rev-parse HEAD)
VERSION ?= $(shell $(GIT) describe --tags ${COMMIT} 2> /dev/null || echo "$(COMMIT)")
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
ROOT_DIRECTORY := $(shell pwd)
LD_FLAGS := -X $(ROOT).Version=$(VERSION) -X $(ROOT).Commit=$(COMMIT) -X $(ROOT).BuildTime=$(BUILD_TIME) -X $(ROOT).RootDirectory=$(ROOT_DIRECTORY) -X $(ROOT).Title=receptionistd
