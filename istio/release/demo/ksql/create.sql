create stream mesh_stream (request_scheme VARCHAR, request_size INTEGER, request_time BIGINT,request_host VARCHAR,request_path VARCHAR,request_method VARCHAR,response_duration 
INTEGER,response_code INTEGER, source_port INTEGER) with ( TIMESTAMP='request_time',VALUE_FORMAT = 'JSON',KAFKA_TOPIC = 'nginmesh');
 
 create table request_path_stat  as select request_path,count(*) as events  from mesh_stream WINDOW TUMBLING ( size 10 seconds) group by request_path;
create table request_path_stat_ts  as select rowTime as event_ts, * FROM request_path_stat;
