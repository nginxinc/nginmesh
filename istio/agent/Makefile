VERSION = 0.6.0
NGINX_VER = 1.13.7
GIT_COMMIT=$(shell git rev-parse --short HEAD)
SIDECAR_IMAGE_NAME=proxy_debug
INIT_IMAGE_NAME=proxy_init
DEV_TAG=$(VERSION)-$(GIT_COMMIT)
DEV_TAG_LATEST=$(VERSION)
REPO=gcr.io/$(GCLOUD_PROJECT)
TAG=$(DEV_TAG)
REPO=gcr.io/$(GCLOUD_PROJECT)
TAG=$(DEV_TAG)
DOCKER_PUSH = gcloud docker -- push
BUILD_DIR=build


release: REPO=nginmesh
release: TAG=$(VERSION)
release: DOCKER_PUSH = docker push

BUILD_IN_CONTAINER = 0

include modules.mk

DOCKER_RUN = docker run --rm -v $(shell pwd):/go/src/github.com/nginmesh/nginx-istio-agent -w /go/src/github.com/nginmesh/nginx-istio-agent/
GOLANG_CONTAINER = golang:1.9

all: helpers dirs modules agent-oss push push-init-container

release: helpers dirs modules agent-oss push push-init-container

debug:	helpers dirs modules agent-oss push push-init-container

dirs: modules-dirs
	mkdir -p $(BUILD_DIR) 

helpers: container-tracing-modules-builder

agent-oss:
ifeq ($(BUILD_IN_CONTAINER),1)
	$(DOCKER_RUN) $(GOLANG_CONTAINER) go build -o $(BUILD_DIR)/agent cmd/agent/*.go
else
	GOOS=linux go build -o $(BUILD_DIR)/agent cmd/agent/*.go
endif

agent-test:
	go test pilot/*.go -v
	go test nginx/*.go -v

container-oss:
	cp -r $(BUILD_DIR) docker-oss
	make PREFIX=$(REPO)/$(SIDECAR_IMAGE_NAME) TAG=$(TAG) DEV_TAG_LATEST=$(DEV_TAG_LATEST) -C docker-oss

container-init:
	make INIT_IMAGE_PREFIX=$(REPO)/$(INIT_IMAGE_NAME) TAG=$(TAG) DEV_TAG_LATEST=$(DEV_TAG_LATEST) -C docker-init

push: container-oss
	$(DOCKER_PUSH) $(REPO)/$(SIDECAR_IMAGE_NAME):$(TAG)
	$(DOCKER_PUSH) $(REPO)/$(SIDECAR_IMAGE_NAME):$(DEV_TAG_LATEST)

push-init-container: container-init
	$(DOCKER_PUSH) $(REPO)/$(INIT_IMAGE_NAME):$(TAG)
	$(DOCKER_PUSH) $(REPO)/$(INIT_IMAGE_NAME):$(DEV_TAG_LATEST)

clean: clean-modules
	rm -rf $(BUILD_DIR)
	rm -rf docker-oss/$(BUILD_DIR)
