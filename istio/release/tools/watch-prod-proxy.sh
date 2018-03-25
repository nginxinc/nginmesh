#!/bin/bash
 kubectl logs -f $(kubectl get pod -l app=productpage  -o jsonpath='{.items[0].metadata.name}') -c istio-proxy 
