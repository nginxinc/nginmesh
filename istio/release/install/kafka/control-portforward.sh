#!/bin/bash
# port-forward connect server
NAMESPACE=${1:-kafka}
kubectl port-forward --namespace $NAMESPACE kcontrol 9021:9021
