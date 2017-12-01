#!/bin/bash

# based on https://github.com/opentracing-contrib/nginx-opentracing/blob/master/docker/Dockerfile

set -ex

NGINX_VERSION=1.13.7

TEMP_DIR="$(mktemp -d)" 

cd $TEMP_DIR

## Build opentracing-cpp
git clone https://github.com/opentracing/opentracing-cpp.git
cd opentracing-cpp
mkdir .build && cd .build
cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_TESTING=OFF ..
make && make install

cd $TEMP_DIR

## Build zipkin-cpp-opentracing
git clone https://github.com/rnburn/zipkin-cpp-opentracing.git
cd zipkin-cpp-opentracing
mkdir .build && cd .build
cmake -DBUILD_SHARED_LIBS=1 -DCMAKE_BUILD_TYPE=Release -DBUILD_TESTING=OFF ..
make && make install 

cd $TEMP_DIR

### Get nginx-opentracing source
git clone https://github.com/opentracing-contrib/nginx-opentracing.git 

## Build nginx-opentracing modules
wget -O nginx-release-$NGINX_VERSION.tar.gz https://github.com/nginx/nginx/archive/release-$NGINX_VERSION.tar.gz
tar zxf nginx-release-$NGINX_VERSION.tar.gz
cd nginx-release-$NGINX_VERSION
auto/configure \
        --with-compat \
        --add-dynamic-module=${TEMP_DIR}/nginx-opentracing/opentracing \
        --add-dynamic-module=${TEMP_DIR}/nginx-opentracing/zipkin
make modules 

mkdir -p /build/modules
cp objs/ngx_http_opentracing_module.so objs/ngx_http_zipkin_module.so /build/modules

mkdir -p /build/libs
cp -L /usr/local/lib/libopentracing.so /usr/local/lib/libzipkin.so /usr/local/lib/libzipkin_opentracing.so /build/libs 