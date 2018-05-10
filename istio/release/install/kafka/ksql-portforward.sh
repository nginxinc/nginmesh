#!/bin/bash
# port-forward ksql server
NAMESPACE=${1:-kafka}
kubectl port-forward --namespace $NAMESPACE ksql-cli 8088:8088