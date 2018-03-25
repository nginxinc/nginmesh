#!/bin/bash
# exec into proxy container
# default for productpage
APP=${1:-productpage}
CMD=${2:-/bin/bash}
kubectl port-forward  -n istio-system $(kubectl get pods -n istio-system -l istio=pilot -o jsonpath='{.items[0].metadata.name}') 15003:15003