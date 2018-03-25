#!/bin/bash
# add new topic
# assume testclient has been installed
TOPIC_NAME=$1
CLIENT=testclient
KAFKA_NAME=my-kafka
kubectl -n kafka exec $CLIENT -- /usr/bin/kafka-console-consumer \
    --bootstrap-server $KAFKA_NAME-kafka:9092 --topic $TOPIC_NAME --from-beginning
