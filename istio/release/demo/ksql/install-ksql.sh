#!/bin/bash
set -x
kubectl apply -f ksql.yml
./copy.sh create.sql
