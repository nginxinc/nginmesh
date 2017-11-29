
# Service Mesh with Istio and NGINX

This repo provides an implementation of sidecar proxy for Istio using NGINX open source version.

## What is Service Mesh and Istio?

Please see https://istio.io for detail explanation of service mesh provided by Istio.  
Combination of Nginx and Istio provides best service mesh for deploying micro-services.

## Production Status

This version of nginmesh work with with 0.2.12 release of Istio.
Please see below for Istio features we support.  Nginmesh is not production ready yet.  


<TBD>

## Architecture

Please see diagram below to see how Nginx Sidecar Proxy is implemented as of 0.16 version.
The sidecar run NGINX open source version with custom module to interface to Istio Mixer.

![Alt text](/images/nginx_sidecar.png?raw=true "Nginx Sidecar")

## The Demo

Please see README.md in samples/bookinfo directory.  

The sample app is copied from Istio project without modification.  
We only support deployment using Kubernetes initializer.  

You can use following kubectl snippet to go inside NGNX proxy container.  Running ps show that is is running NGINX that you know and love.

```bash
kubectl exec -it  $(kubectl get pod -l app=productpage  -o jsonpath='{.items[0].metadata.name}') -c proxy  /bin/bash

nginx@productpage-v1-2380123074-j0hx9:/$ ps -ef                                                                                            
UID        PID  PPID  C STIME TTY          TIME CMD
nginx        1     0  0 Sep01 ?        00:00:08 /agent
nginx        8     1  0 Sep01 ?        00:00:00 nginx: master process nginx -g daemon off;
nginx      107     8  0 Sep01 ?        00:00:00 nginx: worker process
nginx      145     0  0 00:12 ?        00:00:00 /bin/bash
nginx      149   145  0 00:12 ?        00:00:00 ps -ef
nginx@productpage-v1-2380123074-j0hx9:/

```

We added custom module to extend NGINX. The Mixer module send telemetry data to Mixer service and enforce SLA.  
The Mixer module is written in Rust language.

The top level NGINX conf for sample app looks like this:

```bash
load_module modules/ngx_http_istio_mixer_module.so;
user  nginx;
worker_processes  auto;

error_log  /var/log/nginx/error.log warn;
pid        /tmp/nginx.pid;


events {
    worker_connections  1024;
}


http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for" '
                      'scr_ip="$http_x_istio_src_ip" src_uid="$http_x_istio_src_uid" host="$host"';

    access_log  /var/log/nginx/access.log  main;

    sendfile        on;
    #tcp_nopush     on;

    keepalive_timeout  65;

    #gzip  on;

    server_names_hash_bucket_size 128;
    variables_hash_bucket_size 128;

    # mixer configuration
    mixer_server istio-mixer;
    mixer_port   9091;
    mixer_target_ip 10.16.2.68;
    mixer_target_uid kubernetes://productpage-v1-2380123074-j0hx9.default;
    mixer_target_service "productpage.default.svc.cluster.local";
    
    include /etc/nginx/conf.d/*.conf;
````

For every target services for each service, we generate virtual host configuration.  For sample app, this looks below

```bash
more /etc/nginx/conf.d/9080.conf

upstream out.9885ad2b8f51ba5a8ef1dbd5c35677a57e298b14 {
	
	server 10.16.0.43:15001;
}
upstream out.c7c42a0ffbf55f227c7c7c70637160bebc86e072 {
	
	server 10.16.3.53:15001;
}


             
server {
    listen 15001;

    server_name  details:9080 details details.default:9080 details.default details.default.svc:9080 details.default.svc details.default.svc.cluster:9080 details.default.svc.cluster details.default.svc.cluster.local:9080 details.default.svc.cluster.local 10.19.253.23:9080 10.19.253.23;

    
	location / {
        
        mixer on;

        

        proxy_pass http://out.2f50c73ed6c7ec0db7e7658c9112bebcfefd9d5c;
        proxy_set_header Host $host;
        proxy_set_header X-ISTIO-SRC-IP 10.16.2.68;
j0hx9.default;set_header X-ISTIO-SRC-UID kubernetes://productpage-v1-2380123074---More--(29%)
    }
}

server {
    listen 15001;

    server_name  reviews:9080 reviews reviews.default:9080 reviews.default reviews.default.svc:9080 reviews.default.svc reviews.default.svc.cluster:9080 reviews.default.svc.cluster reviews.default.svc.cluster.local:9080 reviews.default.svc.cluster.local 10.19.254.107:9080 10.19.254.107;

    
	location / {
        
        mixer on;
             
        
        
        if ($http_cookie ~ '^(.*?;)?(user=jason)(;.*)?$') {
            set $r0_0 't';
        }
        
        set $r0 '${r0_0}';
        if ($r0 != 't') {
            rewrite ^/(.*)$ /1/$1 last;
        }
        

        proxy_pass http://out.9885ad2b8f51ba5a8ef1dbd5c35677a57e298b14;
        proxy_set_header Host $host;
        proxy_set_header X-ISTIO-SRC-IP 10.16.2.68;
        proxy_set_header X-ISTIO-SRC-UID kubernetes://productpage-v1-2380123074-j0hx9.default;
    }
	location /1 {
        internal;
        mixer on;

             

        proxy_pass http://out.c7c42a0ffbf55f227c7c7c70637160bebc86e072$request_uri;
        proxy_set_header Host $host;
        proxy_set_header X-ISTIO-SRC-IP 10.16.2.68;
        proxy_set_header X-ISTIO-SRC-UID kubernetes://productpage-v1-2380123074-j0hx9.default;
    }
}

```
