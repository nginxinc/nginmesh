#!/usr/bin/env bash
set -X
table_name=$1
TABLE_NAME=`echo $table_name | tr '[a-z]' '[A-Z]'`
echo "=================================================================="
echo "Charting " $TABLE_NAME

## Cleanup existing data

# Elastic
curl -X "DELETE" "http://localhost:9200/""$table_name"

# Connect
curl -X "DELETE" "http://localhost:8083/connectors/es_sink_""$TABLE_NAME"

# Grafana
curl -X "DELETE" "http://localhost:3000/api/datasources/name/""$table_name"   --user admin:admin

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
