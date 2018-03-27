#!/bin/bash
# port-forward elastic search
NAMESPACE=${1:-elastic}
export ELASTIC_POD=$(kubectl get pods --namespace $NAMESPACE -l "app=elasticsearch,component=client,release=elastic" -o jsonpath="{.items[0].metadata.name}")
kubectl port-forward --namespace $NAMESPACE $ELASTIC_POD 9200:9200