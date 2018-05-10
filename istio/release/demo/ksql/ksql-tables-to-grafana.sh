#!/usr/bin/env bash
#set -X
table_name=$1
grafana_password=`kubectl get secret --namespace kafka grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo`
TABLE_NAME=`echo $table_name | tr '[a-z]' '[A-Z]'`
echo "=================================================================="
echo "Charting " $TABLE_NAME

## Cleanup existing data

# Elastic
echo "removing elastic table"
curl -X "DELETE" "http://localhost:9200/""$table_name"

# Connect
echo "cleaning connect elastic connector"
curl -X "DELETE" "http://localhost:28082/connectors/es_sink_""$TABLE_NAME"

# Grafana
echo "cleaning grafana"
curl -X "DELETE" "http://localhost:3000/api/datasources/name/""$table_name"   --user admin:$grafana_password

# Wire in the new connection path
echo "\n\nConnecting KSQL->Elastic->Grafana " "$table_name"
./ksql-connect-es-grafana.sh "$table_name"


echo "Navigate to http://localhost:3000/dashboard/db/click-stream-analysis"


# ========================
#   REST API Notes
# ========================
#
# Extract datasources from grafana
# curl -s "http://localhost:3000/api/datasources"  -u admin:admin|jq -c -M '.[]'
#
# Delete a Grafana DataSource
# curl -X "DELETE" "http://localhost:3000/api/datasources/name/

# List confluent connectors
# curl -X "GET" "http://localhost:8083/connectors"
#
# Delete a Confluent-Connector
# curl -X "DELETE" "http://localhost:8083/connectors/es_sink_PER_USER_KBYTES_TS"
#
# Delete an Elastic Index
# curl -X "DELETE" "http://localhost:9200/per_user_kbytes_ts"
#
