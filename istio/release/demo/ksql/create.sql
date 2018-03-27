create stream request (request_scheme VARCHAR, request_size INTEGER, request_time BIGINT,request_host VARCHAR,request_path VARCHAR,request_method VARCHAR,response_duration 
INTEGER,response_code INTEGER, source_port INTEGER) with ( TIMESTAMP='request_time',VALUE_FORMAT = 'JSON',KAFKA_TOPIC = 'nginmesh');
 