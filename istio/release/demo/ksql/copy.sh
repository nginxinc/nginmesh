#!/bin/bash
SCRIPT=$1
kubectl cp $1 kafka/ksql-cli:/tmp/$1
