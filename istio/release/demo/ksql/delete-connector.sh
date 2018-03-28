#!/bin/bash
set -x
TABLE_NAME=$1
curl -X "DELETE" "http://localhost:8083/connectors/es_sink_""$TABLE_NAME" 
