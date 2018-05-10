#!/bin/bash
set -x
helm install --name grafana --namespace=kafka  stable/grafana 
