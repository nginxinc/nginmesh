#!/bin/bash
# generate and install sidecar
SCRIPTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
set -x
kubectl delete -f $SCRIPTDIR/generated/nginmesh-sidecar-injector-configmap.yaml    
kubectl delete -f $SCRIPTDIR/generated/istio-sidecar-injector-with-ca-bundle.yaml
