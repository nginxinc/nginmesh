#!/bin/bash
SCRIPT=create.sql
kubectl cp $SCRIPT kafka/ksql-cli:/tmp/$1
