#!/bin/bash
# generate and install sidecar
set -x
kubectl delete -f ./generated/nginmesh-sidecar-injector-configmap.yaml    
kubectl delete -f ./generated/istio-sidecar-injector-with-ca-bundle.yaml