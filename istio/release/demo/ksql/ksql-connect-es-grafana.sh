#!/usr/bin/env bash
set -x
## An "all-in-once script" to load up a new table and connect all of the relevant parts to allow data to pipe through from KSQL.KafkaTopic->Connect->Elastic->Grafana[DataSource]
## Accepts a KSQL TABLE_NAME where the data is to be sourced from.


if [ "$#" -ne 1 ]; then
    echo "Usage: ksql-connect-es-grafana.sh <TABLENAME>"
    exit;
fi

table_name=$1
KAFKA=my-kafka-kafka.kafka
TABLE_NAME=`echo $1 | tr '[a-z]' '[A-Z]'`


echo "Connecting:" $table_name

## Load the _TS dynamic template into ELASTIC
./elastic-dynamic-template.sh

# Tell Kafka to send this Table-Topic to Elastic
echo "Adding Elastic Source\n\n"


curl -X "POST" "http://localhost:8083/connectors/" \
     -H "Content-Type: application/json" \
     -d $'{
  "name": "es_sink_'$TABLE_NAME'",
  "config": {
    "schema.ignore": "true",
    "topics": "'$TABLE_NAME'",
    "key.converter": "org.apache.kafka.connect.storage.StringConverter",
    "value.converter.schemas.enable": false,
    "connector.class": "io.confluent.connect.elasticsearch.ElasticsearchSinkConnector",
    "key.ignore": "true",
    "value.converter": "org.apache.kafka.connect.json.JsonConverter",
    "type.name": "type.name=kafkaconnect",
    "topic.index.map": "'$TABLE_NAME':'$table_name'",
     "connection.url": "http://elastic-elasticsearch-client.elastic:9200",
    "transforms": "FilterNulls",
    "transforms.FilterNulls.type": "io.confluent.transforms.NullFilter"
  }
}'

## Add the Elastic DataSource into Grafana
curl -X "POST" "http://localhost:3000/api/datasources" \
	    -H "Content-Type: application/json" \
	     --user admin:admin \
	     -d $'{"id":1,"orgId":1,"name":"'$table_name'","type":"elasticsearch","typeLogoUrl":"public/app/plugins/datasource/elasticsearch/img/elasticsearch.svg","access":"proxy","url":"http://elastic-elasticsearch-client.elastic:9200","password":"","user":"","database":"'$table_name'","basicAuth":false,"isDefault":false,"jsonData":{"timeField":"EVENT_TS"}}'

