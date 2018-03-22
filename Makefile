## Project Vars ##########################################################
GO_VARS = ENABLE_CGO=0 GOOS=darwin GOARCH=amd64
GO ?= go
GIT ?= git
COMMIT := $(shell $(GIT) rev-parse HEAD)
VERSION ?= $(shell $(GIT) describe --tags ${COMMIT} 2> /dev/null || echo "$(COMMIT)")
BUILD_TIME := $(shell LANG=en_US date +"%F_%T_%z")
ROOT := gitlab.com/kanalbot/receptionist
ROOT_DIRECTORY := $(shell pwd)
LD_FLAGS := -X $(ROOT).Version=$(VERSION) -X $(ROOT).Commit=$(COMMIT) -X $(ROOT).BuildTime=$(BUILD_TIME) -X $(ROOT).Title=receptionistd
DOCKER_IMAGE := registry.gitlab.com/kanalbot/receptionist

.PHONY: help clean

receptionistd: *.go */*.go */*/*.go Gopkg.lock
	$(GO_VARS) $(GO) build -i -o="receptionistd" -ldflags="$(LD_FLAGS)" $(ROOT)/cmd/receptionist

clean:
	rm -rf receptionistd

help:
	@echo "Please use \`make <ROOT>' where <ROOT> is one of"
	@echo "  receptionistd     to build the main binary for current platform"
	@echo "  clean             to remove generated files"

docker: receptionistd Dockerfile
	docker build -t $(DOCKER_IMAGE):$(VERSION) .
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest

push:
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):latest
