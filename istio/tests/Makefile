KUBE_VERSION=v1.9.4
NGINMESH_VERSION=0.8.0
WRK_VERSION=3.1.0
GIT_COMMIT=$(shell git rev-parse --short HEAD)
TEST_IMAGE_NAME=bookinfo_testing
dev: REPO=gcr.io/$(GCLOUD_PROJECT)
release: REPO=nginmesh
TAG=$(NGINMESH_VERSION)-$(GIT_COMMIT)
TAG_LATEST=$(NGINMESH_VERSION)

build:
	docker build --build-arg KUBE_VERSION=$(KUBE_VERSION)  --build-arg NGINMESH_VERSION=$(NGINMESH_VERSION)  --build-arg WRK_VERSION=$(WRK_VERSION) -t $(REPO)/$(TEST_IMAGE_NAME):$(TAG) -t $(REPO)/$(TEST_IMAGE_NAME):$(TAG_LATEST) -f src/Dockerfile .

release: build

dev: build



