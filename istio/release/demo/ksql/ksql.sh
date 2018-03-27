#!/bin/bash
kubectl exec -it ksql-cli -n kafka -- /usr/bin/ksql-cli local --bootstrap-server my-kafka-kafka.kafka:9092
