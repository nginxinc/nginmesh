#!/bin/bash
SCRIPTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
set -x
kubectl delete ns kafka
helm repo remove incubator
helm del --purge my-kafka
kubectl delete -f $SCRIPTDIR/kafka-client.yml
