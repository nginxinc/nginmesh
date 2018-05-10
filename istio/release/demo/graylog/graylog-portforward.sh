#!/bin/bash
# port-forward graylog
NAMESPACE=${1:-graylog}
export POD=$(kubectl -n $NAMESPACE get pod -l service=graylog -o jsonpath='{.items[0].metadata.name}')
kubectl port-forward --namespace $NAMESPACE $POD 9000:9000
