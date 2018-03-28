#!/bin/bash
set -x
helm install --name grafana --namespace=kafka  stable/grafana --set server.adminPassword=admin --set server.image=grafana/grafana:master
