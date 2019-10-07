PROJ=dex
ORG_PATH=github.com/dexidp
REPO_PATH=$(ORG_PATH)/$(PROJ)
export PATH := $(PWD)/bin:$(PATH)
THIS_DIRECTORY:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

VERSION ?= $(shell ./scripts/git-version)

DOCKER_REPO=quay.io/dexidp/dex
DOCKER_IMAGE=$(DOCKER_REPO):$(VERSION)

$( shell mkdir -p bin )

user=$(shell id -u -n)
group=$(shell id -g -n)

export GOBIN=$(PWD)/bin

LD_FLAGS="-w -X $(REPO_PATH)/version.Version=$(VERSION)"

build: bin/dex bin/example-app

bin/dex:
	@go install -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/dex

bin/example-app:
	@go install -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/example-app


.PHONY: release-binary
release-binary:
	@go build -o /go/bin/dex -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/dex

.PHONY: revendor
revendor:
	@go mod tidy -v
	@go mod vendor -v
	@go mod verify

test:
	@go test -v ./...

testrace:
	@go test -v --race ./...

vet:
	@go vet ./...

fmt:
	@./scripts/gofmt ./...

lint: bin/golint
	@./bin/golint -set_exit_status $(shell go list ./...)

.PHONY: docker-image
docker-image:
	@sudo docker build -t $(DOCKER_IMAGE) .


bin/golint:
	@go install -v $(THIS_DIRECTORY)/vendor/golang.org/x/lint/golint

clean:
	@rm -rf bin/

testall: testrace vet fmt lint

FORCE:

.PHONY: test testrace vet fmt lint testall
