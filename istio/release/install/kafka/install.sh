#!/bin/bash
SCRIPTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
set -x
kubectl create ns kafka
helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator
helm install --name my-kafka --namespace kafka incubator/kafka --set imageTag=4.1.0 
kubectl apply -f $SCRIPTDIR/kafka-client.yml
