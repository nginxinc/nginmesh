#!/bin/bash
# port-forward connect server
NAMESPACE=${1:-kafka}
kubectl port-forward --namespace $NAMESPACE kconnect 28082:28082
