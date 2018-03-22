extern crate ngx_mixer_transport;
extern crate grpc;
extern crate futures;

use std::collections::HashMap;
use ngx_mixer_transport::istio_client::mixer_client_wrapper::MixerClientWrapper ;
use ngx_mixer_transport::attribute::attr_wrapper::AttributeWrapper;
use ngx_mixer_transport::transport::mixer_grpc::GrpcTransport;
use ngx_mixer_transport::transport::server_info::MixerInfo;
use ngx_mixer_transport::transport::status::{ StatusCodeEnum};
use futures::future::Future;

// run integration test, in order to run this, mixer server should be running https://github.com/istio/mixer/blob/master/doc/dev/development.md
// ./bazel-bin//cmd/server/mixs server --logtostderr --configStore2URL=fs://$(pwd)/testdata/config --configStoreURL=fs://$(pwd)/testdata/configroot  -v=4


#[test]
fn intg_check_empty_request() {


    let info = MixerInfo { server_name: String::from("localhost"), server_port: 9091};
    let attributes = AttributeWrapper::new();

    let transport = GrpcTransport::new(info,attributes);

    let client = MixerClientWrapper::new();

    let result = client.check(transport).wait();

    println!("result, {:?}",result);

    match result  {
        Ok(_) =>  assert!(true,"succeed"),
        Err(_)  => assert!(false,"failed check")
    }
}


#[test]
fn intg_check_deny() {


    let info = MixerInfo { server_name: String::from("localhost"), server_port: 9091};
    let mut attributes = AttributeWrapper::new();
    attributes.insert_string_attribute("destination.service","abc.ns.svc.cluster.local");
    attributes.insert_string_attribute("source.name","myservice");
    attributes.insert_string_attribute("source.port","8080");

    let mut string_map: HashMap<String,String> = HashMap::new();
    string_map.insert("clnt".to_string(),"abc".to_string());
    string_map.insert("source".to_string(),"abcd".to_string());
    string_map.insert("destination.labels".to_string(),"app:ratings".to_string());
    string_map.insert("labels".to_string(),"version:v2".to_string());

    attributes.insert_string_map("request.headers",string_map);

    //  --string_attributes destination.service=abc.ns.svc.cluster.local,source.name=myservice,target.port=8080 --stringmap_attributes "request.headers=clnt:abc;source:abcd,destination.labels=app:ratings,source.labels=version:v2"   --timestamp_attributes request.time="2017-07-04T00:01:10Z" --bytes_attributes source.ip=c0:0:0:2
    //2017/10/31 15:21:18 grpc: addrConn.resetTrans

    let transport = GrpcTransport::new(info,attributes);

    let client = MixerClientWrapper::new();

    let result = client.check(transport).wait();

    println!("result, {:?}",result);

    match result  {
        Ok(_) =>  assert!(false,"should not have succeed"),
        Err(error)  => assert_eq!(error.get_error_code(),StatusCodeEnum::PERMISSION_DENIED,"permission denied expected")
    }
}


