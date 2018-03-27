#!/bin/bash
# update sidecar
# assume istio is installed
SCRIPTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
set -x
mkdir -p generated
NGINMESH_IMAGE_HUB=${1:-docker.io}
KAFKA=${2:-my-kafka}
KAFKA_TOPIC=${3:-nginmesh}
NGX_LOG_LEVEL=${4:-warn}
echo "generating sidecar config using kafka: $KAFKA, topic: $KAFKA_TOPIC"
KAFKA_SERVER=${KAFKA}-kafka.kafka:9092
NGINMESH_CONFIG_NAME=nginmesh-sidecar-injector-configmap.yaml
GEN_NGINMESH=$SCRIPTDIR/generated/$NGINMESH_CONFIG_NAME
echo "# GENERATED FILE. Use with Istio 0.6" > $GEN_NGINMESH
cat $SCRIPTDIR/templates/$NGINMESH_CONFIG_NAME.tmpl >> $GEN_NGINMESH
sed -i .bak "s|{NGINMESH_IMAGE_HUB}|${NGINMESH_IMAGE_HUB}|" $GEN_NGINMESH
sed -i .bak "s|{KAFKA_SERVER}|${KAFKA_SERVER}|" $GEN_NGINMESH
sed -i .bak "s|{KAFKA_TOPIC}|${KAFKA_TOPIC}|" $GEN_NGINMESH
sed -i .bak "s|{NGX_LOG_LEVEL}|${NGX_LOG_LEVEL}|" $GEN_NGINMESH
