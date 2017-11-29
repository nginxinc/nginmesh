#!/bin/bash

kubectl port-forward $(kubectl get pod -l app=grafana -n istio-system -o jsonpath='{.items[0].metadata.name}')  -n istio-system 3000:3000
