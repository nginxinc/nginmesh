#!/bin/bash
set -x
DEPLOYMENT_NAME=$1
DOCKER_NGINX_NAME=$2
DOCKER_MODULE_IMAGE=$3
TAG=$4
kubectl delete deployment $DEPLOYMENT_NAME || true
pkill -f 8000 | true
kubectl run $DOCKER_NGINX_NAME --labels='app=nginmesh' --image $DOCKER_MODULE_IMAGE:$TAG
sleep 10
kubectl port-forward $(kubectl get pod -l app=nginmesh -o jsonpath='{.items[0].metadata.name}') 8000:8000 &
kubectl exec $(kubectl get pod -l app=nginmesh -o jsonpath='{.items[0].metadata.name}') /bin/bash /test/deploy.sh
