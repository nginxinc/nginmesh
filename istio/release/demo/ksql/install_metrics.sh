echo -e "run script 'create.sql';\nexit" | ksql
./ksql-tables-to-grafana.sh request_per_min
./ksql-tables-to-grafana.sh errors_per_min_alert
./ksql-tables-to-grafana.sh request_activity
