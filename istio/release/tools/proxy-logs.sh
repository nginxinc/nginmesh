#!/bin/bash
# watch logs for proxy container
# default for productpage
APP=${1:-productpage}
 kubectl logs -f $(kubectl get pods -l app=$APP -o jsonpath='{.items[0].metadata.name}') -c istio-proxy