
all:
	docker build -t $(INIT_IMAGE_PREFIX):$(TAG) .
	docker tag $(INIT_IMAGE_PREFIX):$(TAG) $(INIT_IMAGE_PREFIX):$(DEV_TAG_LATEST)