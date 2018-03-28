#!/bin/bash
# exec into proxy container
# default for productpage
NAMESPACE=${1:-kafka}
kubectl port-forward  -n $NAMESPACE $(kubectl get pods --namespace $NAMESPACE -l "app=grafana-grafana,component=grafana" -o jsonpath="{.items[0].metadata.name}") 3000:3000