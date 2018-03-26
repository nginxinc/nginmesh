# kubectl run ksql  --image confluentinc/ksql-cli:0.5  --rm -n kafka -ti --command /bin/bash
# then type in cr
# ksql-cli local --bootstrap-server my-kafka-broker.kafka:9092
# create stream request (request_scheme VARCHAR, request_time BIGINT,request_host VARCHAR,request_method VARCHAR,response_duration INTEGER) with ( TIMESTAMP='request_time',VALUE_FORMAT = 'JSON',KAFKA_TOPIC = 'nginmesh');
 