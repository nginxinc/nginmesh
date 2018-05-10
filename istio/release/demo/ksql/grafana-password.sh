#!/bin/bash
kubectl get secret --namespace kafka grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
