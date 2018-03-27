#!/bin/bash
set -x
# https://github.com/kubernetes/charts/tree/master/incubator/kafka
kubectl create ns kafka
helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator
helm install --name my-kafka --namespace kafka incubator/kafka
kubectl apply -f ./kafka-client.yml
#helm install --name registry --namespace=kafka  incubator/schema-registry 