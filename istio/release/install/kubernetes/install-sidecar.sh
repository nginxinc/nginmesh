#!/bin/bash
# generate and install sidecar
set -x
IMAGE_HUB=${1:-docker.io}
KAFKA=${2:-my-kafka}
KAFKA_TOPIC=${3:-nginmesh}
mkdir -p generated
./istio/webhook-create-signed-cert.sh \
    --service istio-sidecar-injector \
    --namespace istio-system \
    --secret sidecar-injector-certs
./generate-sidecar-config.sh $IMAGE_HUB $KAFKA $KAFKA_TOPIC    
kubectl apply -f ./generated/nginmesh-sidecar-injector-configmap.yaml
cat ./istio/istio-sidecar-injector.yaml | \
     ./istio/webhook-patch-ca-bundle.sh > \
     ./generated/istio-sidecar-injector-with-ca-bundle.yaml       
kubectl apply -f ./generated/istio-sidecar-injector-with-ca-bundle.yaml