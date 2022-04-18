# Binary
TAG ?= dev-local
BUILD_HASH := $(shell git rev-parse HEAD)
BUILD_HASH_SHORT = $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u +%Y%m%d.%H%M%S)
LDFLAGS += -X "github.com/mattermost/ponos/function.BuildDate=$(BUILD_DATE)"
LDFLAGS += -X "github.com/mattermost/ponos/function.BuildHash=$(BUILD_HASH)"
LDFLAGS += -X "github.com/mattermost/ponos/function.BuildHashShort=$(BUILD_HASH_SHORT)"

## Golang
GO ?= go
GO_TEST_FLAGS ?= -race

## Binaries.
TOOLS_BIN_DIR := $(abspath bin)

GO_IMAGE_LINT = "golangci/golangci-lint:v1.45.0"

OUTDATED_VER := master
OUTDATED_BIN := go-mod-outdated
OUTDATED_GEN := $(TOOLS_BIN_DIR)/$(OUTDATED_BIN)

## Docker
DOCKER_IMAGE ?= mattermost/ponos-service:${TAG}

## Docker Build Versions
DOCKER_BUILD_IMAGE = golang:1.17.9@sha256:55759506d9a7b33b28117977170db9713903799c110db99af2c15a3f603602af
DOCKER_BASE_IMAGE = gcr.io/distroless/static@sha256:d6fa9db9548b5772860fecddb11d84f9ebd7e0321c0cb3c02870402680cc315f

## Tools version
TERRAFORM_VERSION=1.1.7
################################################################################

.PHONY: all
## all: builds and pushes the docker image
all: lint test docker-build docker-push

.PHONY: build
## build: builds linux binary
build:
	@echo Building binary for linux for ${BUILD_SERVICE}
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -ldflags '$(LDFLAGS)' -trimpath -o build/_output/bin/ponos-${BUILD_SERVICE}-linux-amd64 ./cmd/${BUILD_SERVICE}

.PHONY: build-local
## build-local: build the executable for local usage
build-local:
	@echo Building for local use only
	$(GO) build -o build/_output/bin/ponos-${BUILD_SERVICE} ./cmd/${BUILD_SERVICE}

.PHONY: docker-build
## docker-build: builds the docker image
docker-build:
	@echo Building Ponos Service Docker Image
	docker build \
	--build-arg DOCKER_BUILD_IMAGE=$(DOCKER_BUILD_IMAGE) \
	--build-arg DOCKER_BASE_IMAGE=$(DOCKER_BASE_IMAGE) \
	. -f build/Dockerfile -t $(DOCKER_IMAGE)

.PHONY: docker-login
## docker-login: docker login to registry
docker-login:
	@echo Docker login
	echo $(DOCKER_PASSWORD) | docker login --username $(DOCKER_USERNAME) --password-stdin

.PHONY: docker-push
## docker-push: push the docker image
docker-push: docker-login
	@echo Taging Pushing Ponos Serve Docker Image
	docker push mattermost/mattermost-ponos-service:${TAG}

.PHONY: check-modules
## check-modules: Check outdated modules
check-modules:
	@echo Checking outdated modules
	$(GO) list -u -m -json all | docker run --rm -i psampaz/go-mod-outdated

.PHONY: check-style
## check-style: Runs govet and gofmt against all packages.
check-style: govet lint
	@echo Checking for style guide compliance

.PHONY: clean
## clean: deletes all
clean:
	rm -rf build/_output/bin/

.PHONY: dist
## dist-aws: creates the bundle file for AWS Lambda deployments
dist: build
	@echo Building dist for AWS Lambda
	cp -r static dist
	cp manifest.json dist/
	cp build/_output/bin/ponos-app-linux-amd64 dist/ponos
	cd dist/; zip -qr go-function ponos; zip -r bundle.zip go-function.zip manifest.json static

.PHONY: get-terraform
## get-terraform: download terraform only if it's not available. Used in the docker build
get-terraform: 
	@if [ ! -f build/terraform ]; then \
		curl -Lo build/terraform.zip https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && cd build && unzip terraform.zip &&\
		chmod +x terraform && rm terraform.zip;\
	fi

.PHONY: vet
## govet: Runs govet against all packages.
govet:
	@echo Running govet
	$(GO) vet ./...
	@echo Govet success

.PHONY: lint
## lint: Run golangci-lint on codebase
lint:
	@echo Running lint with GolangCI
	docker run --rm -v $(PWD):/app -w /app ${GO_IMAGE_LINT} golangci-lint run --timeout=5m

.PHONY: run
## run: runs the service
run: build-local
	@echo Running chaos engine with debug
	build/_output/bin/ponos-${BUILD_SERVICE}

.PHONY: test
## test: tests all packages
test:
	@echo Running tests
	$(GO) test $(GO_TEST_FLAGS) ./...

.PHONY: vendor
## vendor: create a vendor folder
vendor:
	@echo Running vendor
	go mod vendor

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
