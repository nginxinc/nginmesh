# KSQL Demo

## Install

### Install Elastic Search
```
./install-elastic.sh
```
This will install elastic cluster in the namespace 'elastic'

### Install Grafana
```
./install-grafana.sh
```
This will install grafana in the namespace 'kafka'

### Install Connect
```
./install-connect.sh
```
This will create pod with Kafka Connect.
After pod is successfully created, run following script to copy connect properties
```
./copy-connect.sh
```
Then run following script to shell into connect and start connect
```
./run-connect.sh
```
Then in the shell, run
```
cd /etc/kafka
connect-distributed connect-distributed.properties
```

### Install KSQL
```
./install-ksql.sh
```
After pod is created, run following script
```
./copy-sql.sh
```
Then run following shell to exec into ksql pod
```
./run-ksql.sh
```
run following script in the SQL
``
run script '/tmp/create.sql';
``
### Start following port-forwarding, this is required in order to connect kafka to elastic search to grafana
```
./elastic-portforward.sh
./connect-portforward.sh
./grafana-portforward.sh
```

### Connect following tables to Elastic Search and Set up data source to Grafana

```
./ksql-tables-to-grafana.sh request_path_stat_ts
./ksql-tables-to-grafana.sh request_path_stat

## Visualization

TBD