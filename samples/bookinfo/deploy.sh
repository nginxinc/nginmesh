#!/bin/bash

./nginx-inject -f bookinfo.yaml  | kubectl create -f -
