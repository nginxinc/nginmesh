#!/bin/bash
# exec into proxy container
# default for productpage
APP=${1:-productpage}
CMD=${2:-/bin/bash}
kubectl exec -it  $(kubectl get pods -l app=$APP -o jsonpath='{.items[0].metadata.name}') -c istio-proxy  $CMD