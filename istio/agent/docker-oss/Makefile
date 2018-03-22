all:
	echo 'building docker image for oss ${PREFIX}:${TAG}'
	cp ../../../collector/module/module/release/ngx_http_collector_module.so build/modules
	docker build -t $(PREFIX):$(TAG) .
	docker tag $(PREFIX):$(TAG) $(PREFIX):$(DEV_TAG_LATEST)
