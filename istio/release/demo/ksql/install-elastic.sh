#!/bin/bash
set -x
helm install --name my-es-test2 incubator/elasticsearch --namespace elastic2 --set data.persistence.enabled=false,master.persistence.enabled=false --set rbac.create=true