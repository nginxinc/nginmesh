#!/bin/bash
set -x
helm install --name elastic  incubator/elasticsearch --namespace elastic --set data.persistence.enabled=false,master.persistence.enabled=false --set rbac.create=true
