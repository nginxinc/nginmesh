#!/bin/bash
set -x
PROP=connect-distributed.properties
kubectl cp $PROP kafka/kconnect:/etc/kafka/connect-distributed.properties
