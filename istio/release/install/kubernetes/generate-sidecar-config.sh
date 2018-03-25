#!/bin/bash
# update sidecar
# assume istio is installed
set -x
mkdir -p generated
KAFKA=${1:-my-kafka}
KAFKA_TOPIC=${2:-nginmesh}
NGX_LOG_LEVEL=${3:-warn}
echo "generating sidecar config using kafka: $KAFKA, topic: $KAFKA_TOPIC"
KAFKA_SERVER=${KAFKA}-kafka.kafka:9092
NGINMESH_CONFIG_NAME=nginmesh-sidecar-injector-configmap.yaml
GEN_NGINMESH=./generated/$NGINMESH_CONFIG_NAME
echo "# GENERATED FILE. Use with Istio 0.6" > $GEN_NGINMESH
cat ./templates/$NGINMESH_CONFIG_NAME.tmpl >> $GEN_NGINMESH
sed -i .bak "s|{KAFKA_SERVER}|${KAFKA_SERVER}|" $GEN_NGINMESH
sed -i .bak "s|{KAFKA_TOPIC}|${KAFKA_TOPIC}|" $GEN_NGINMESH
sed -i .bak "s|{NGX_LOG_LEVEL}|${NGX_LOG_LEVEL}|" $GEN_NGINMESH