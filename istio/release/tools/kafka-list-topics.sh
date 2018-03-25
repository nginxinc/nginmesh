#!/bin/bash
# list topic
CLIENT=testclient
KAFKA_NAME=my-kafka
kubectl -n kafka exec -ti $CLIENT -- /usr/bin/kafka-topics --zookeeper $KAFKA_NAME-zookeeper:2181 --list
