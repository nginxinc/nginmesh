all:
	echo 'building docker image for oss ${PREFIX}:${TAG}'
	docker build -t $(PREFIX):$(TAG) .
	docker tag $(PREFIX):$(TAG) $(PREFIX):$(DEV_TAG_LATEST)
