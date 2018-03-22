NGINX_VER = 1.13.7
TAG=dev
RUST_COMPILER_TAG = 1.21.0
DOCKER_REPO=nginmesh
UNAME_S := $(shell uname -s)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
NGX_DEBUG="--with-debug"
export MODULE_DIR=${PWD}
NGX_MODULES = --with-compat  --with-threads --with-http_addition_module \
     --with-http_auth_request_module   --with-http_gunzip_module --with-http_gzip_static_module  \
     --with-http_random_index_module --with-http_realip_module --with-http_secure_link_module \
     --with-http_slice_module  --with-http_stub_status_module --with-http_sub_module \
     --with-stream --with-stream_realip_module --with-stream_ssl_preread_module
ifeq ($(UNAME_S),Linux)
    NGINX_SRC += nginx-linux
    NGX_OPT= $(NGX_MODULES) \
       --with-file-aio
       --with-cc-opt='-g -fstack-protector-strong -Wformat -Werror=format-security -Wp,-D_FORTIFY_SOURCE=2 -fPIC' \
       --with-ld-opt='-Wl,-Bsymbolic-functions -Wl,-z,relro -Wl,-z,now -Wl,--as-needed -pie'
endif
ifeq ($(UNAME_S),Darwin)
    NGINX_SRC += nginx-darwin
    NGX_OPT= $(NGX_MODULES)
endif
DOCKER_BUILD=./docker
DOCKER_MODULE_IMAGE = $(DOCKER_REPO)/${MODULE_NAME}
DOCKER_MODULE_BASE_IMAGE = $(DOCKER_REPO)/${MODULE_NAME}-base
DOCKER_MODULE_NGINX_BUILD_IMAGE = $(DOCKER_REPO)/${MODULE_NAME}-ngx-build
DOCKER_MODULE_NGINX_BASE_IMAGE= $(DOCKER_REPO)/${MODULE_NAME}-ngx-base
DOCKER_RUST_IMAGE = $(DOCKER_REPO)/ngx-rust-tool:${RUST_COMPILER_TAG}
DOCKER_NGIX_IMAGE = $(DOCKER_REPO)/nginx-dev:${NGINX_VER}
DOCKER_MIXER_IMAGE = $(DOCKER_REPO)/ngix-mixer:1.0
MODULE_SO_DIR=nginx/nginx-linux/objs
MODULE_SO_BIN=${MODULE_SO_DIR}/${MODULE_NAME}.so
NGINX_BIN=${MODULE_SO_DIR}/nginx
MODULE_SO_HOST=module/release/${MODULE_NAME}.so
NGINX_SO_HOST=config


DOCKER_BUILD_TOOL=docker run -it --rm -v ${ROOT_DIR}:/src -w /src/${MODULE_PROJ_NAME} ${DOCKER_RUST_IMAGE}
DOCKER_NGINX_TOOL=docker run -it --rm -v ${ROOT_DIR}:/src -w /src/${MODULE_PROJ_NAME} ${DOCKER_NGIX_IMAGE}
DOCKER_NGINX_NAME=nginx-test
DOCKER_NGINX_EXEC=docker exec -it ${DOCKER_NGINX_NAME}
DOCKER_NGINX_EXECD=docker exec -d ${DOCKER_NGINX_NAME}
DOCKER_NGINX_DAEMON=docker run -d -p 8000:8000  --privileged --name  ${DOCKER_NGINX_NAME} \
    --sysctl net.ipv4.ip_nonlocal_bind=1 \
    --sysctl net.ipv4.ip_forward=1 \
	-v ${MODULE_DIR}/module/release:/etc/nginx/modules \
	-v ${MODULE_DIR}:/src  -w /src   ${DOCKER_NGIX_IMAGE}


setup-nginx:
	mkdir -p nginx


nginx-source:	setup-nginx
	rm -rf nginx/${NGINX_SRC}
	wget http://nginx.org/download/nginx-${NGINX_VER}.tar.gz
	tar zxf nginx-${NGINX_VER}.tar.gz
	mv nginx-${NGINX_VER} ${NGINX_SRC}
	mv ${NGINX_SRC} nginx
	rm nginx-${NGINX_VER}.tar.gz*

nginx-configure:
	cd nginx/${NGINX_SRC}; \
    ./configure --add-dynamic-module=../../module $(NGX_OPT)


nginx-setup:	nginx-source nginx-configure


nginx-module:
	cd nginx/${NGINX_SRC}; \
	make modules; \
	strip objs/*.so



copy-module:
	docker rm -v ngx-copy || true
	docker create --name ngx-copy ${DOCKER_MODULE_IMAGE}:${TAG}
	docker cp ngx-copy:/etc/nginx/modules/${MODULE_NAME}.so ${MODULE_SO_HOST}
	docker rm -v ngx-copy

# build module using docker
# we copy only necessary context to docker daemon (src and module directory)
build-module-docker:
	rm -rf $(DOCKER_BUILD)/context
	mkdir $(DOCKER_BUILD)/context
	cp $(DOCKER_BUILD)/Dockerfile.module $(DOCKER_BUILD)/context
	cp -r collector-ngx $(DOCKER_BUILD)/context
	cp -r collector-transport $(DOCKER_BUILD)/context
	cp -r collector-tests $(DOCKER_BUILD)/context
	cp -r module $(DOCKER_BUILD)/context
	cp -r test $(DOCKER_BUILD)/context
	docker build -f $(DOCKER_BUILD)/context/Dockerfile.module -t ${DOCKER_MODULE_IMAGE}:${GIT_COMMIT} $(DOCKER_BUILD)/context
	docker tag ${DOCKER_MODULE_IMAGE}:${GIT_COMMIT} ${DOCKER_MODULE_IMAGE}:${TAG}
	

# build module and deposit in the module directory
build-module: build-module-docker

# build base container image that pre-compiles rust and nginx modules
build-base:
	docker build -f $(DOCKER_BUILD)/Dockerfile.base -t ${DOCKER_MODULE_BASE_IMAGE}:${GIT_COMMIT} .
	docker tag ${DOCKER_MODULE_BASE_IMAGE}:${GIT_COMMIT} ${DOCKER_MODULE_BASE_IMAGE}:${TAG}


run-base-image:
	docker run -it --rm  ${DOCKER_MODULE_BASE_IMAGE}:dev /bin/bash


run-module-image:
	docker run -it --rm  ${DOCKER_MODULE_IMAGE}:dev /bin/bash



watch-mixer:
	 kubectl logs -f $(kubectl get pod -l istio=mixer -n istio-system -o jsonpath='{.items[0].metadata.name}')  -n istio-system -c mixer	




# setup nginx container for testing
# copies the configuration and modules
# start test services
test-nginx-setup:
	test/deploy.sh


# run integrated test
test-intg:
	cargo +stable test --color=always intg -- --nocapture


test-unit:
	cargo test --lib


# remove nginx container
test-nginx-clean:
	docker rm -f  ${DOCKER_NGINX_NAME} || true


test-nginx-only: test-nginx-clean
	$(DOCKER_NGINX_DAEMON)
	$(DOCKER_NGINX_EXECD) make test-nginx-setup > make.out
	sleep 1

test-docker-only:
	docker rm -f nginx || true
	docker run --name nginx -d ${DOCKER_MODULE_IMAGE}:${TAG} 

test-k8-nginx-clean:
	kubectl delete deployment nginx-test || true

test-k8-only:	
	./scripts/k8-test.sh nginx-test ${DOCKER_NGINX_NAME} ${DOCKER_MODULE_IMAGE} ${TAG}

test-nginx-log:
	docker logs -f nginx-test

test-k8-setup:
	kubectl exec -it nginx-test-57df6c6988-d6wnf /bin/bash
	kubectl port-forward nginx-test-57df6c6988-d6wnf 8000:8000 &

show-k8-logs:
	kubectl logs $(kubectl get pod -l app=nginmesh -o jsonpath='{.items[0].metadata.name}')


kafka-install:
	kubectl create ns kafka
	helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator
	helm install --name my-kafka --namespace kafka incubator/kafka
	kubectl apply -f test/kafka-client.yml

kafka-add-test-topic:
	kubectl -n kafka exec testclient -- /usr/bin/kafka-topics --zookeeper my-kafka-zookeeper:2181 --topic test --create --partitions 1 --replication-factor 1	
		

kafka-list-message:
	kubectl -n kafka exec -ti testclient -- /usr/bin/kafka-console-consumer --bootstrap-server my-kafka-kafka:9092 --topic test --from-beginning

test-nginx-full:	build-module test-nginx-only

# invoke http report
test-http-report:
	curl localhost:8000/report