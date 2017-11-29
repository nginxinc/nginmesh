#!/bin/bash
kubectl port-forward $(kubectl get pod -l  istio=mixer -n istio-system -o jsonpath='{.items[0].metadata.name}')  -n istio-system 9091:9091
