#!/bin/bash
# generate and install sidecar
SCRIPTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
set -x
IMAGE_HUB=${1:-docker.io}
KAFKA=${2:-my-kafka}
KAFKA_TOPIC=${3:-nginmesh}
mkdir -p $SCRIPTDIR/generated
$SCRIPTDIR/istio/webhook-create-signed-cert.sh \
    --service istio-sidecar-injector \
    --namespace istio-system \
    --secret sidecar-injector-certs
$SCRIPTDIR/generate-sidecar-config.sh $IMAGE_HUB $KAFKA $KAFKA_TOPIC
kubectl apply -f $SCRIPTDIR/generated/nginmesh-sidecar-injector-configmap.yaml
cat $SCRIPTDIR/istio/istio-sidecar-injector.yaml | \
     $SCRIPTDIR/istio/webhook-patch-ca-bundle.sh > \
     $SCRIPTDIR/generated/istio-sidecar-injector-with-ca-bundle.yaml
kubectl apply -f $SCRIPTDIR/generated/istio-sidecar-injector-with-ca-bundle.yaml
