#!/bin/bash
kubectl port-forward $(kubectl get pod -l  app=zipkin -n istio-system -o jsonpath='{.items[0].metadata.name}')  -n istio-system 9411:9411
