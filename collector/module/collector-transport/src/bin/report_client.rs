/**
  test sending of report
 */


extern crate nginmesh_collector_transport;


use std::collections::HashMap;
// use ngx_mixer_module::mixer::service_grpc::MixerClient;
// use ngx_mixer_module::mixer::report::ReportRequest;
// use ngx_mixer_module::mixer::attributes::Attributes;
//use ngx_mixer_module::mixer::service_grpc::Mixer;


static REQUEST_HEADER: i32 = 0;
static TARGET_SERVICE: i32 = 1;
static REQUEST_HOST: i32 = 2;


fn main() {

    println!("creating mixer client at local host");

    //let client = MixerClient::new_plain("localhost", 9091, Default::default()).expect("init");


  //  let mut rf = RepeatedField::default();
   // let attr = Attributes::new();
    //attr.set_string_attributes("")

    let mut dict_values: HashMap<i32,String> = HashMap::new();
    dict_values.insert(REQUEST_HEADER,String::from("request.headers"));
    dict_values.insert(TARGET_SERVICE,String::from("target.service"));
    dict_values.insert(REQUEST_HOST,String::from("request.host"));
 //   attr.set_dictionary(dict_values);

    let mut string_values: HashMap<i32,String> = HashMap::new();
    string_values.insert(TARGET_SERVICE,String::from("reviews.default.svc.cluster.local"));
    string_values.insert(REQUEST_HOST,String::from("test.com"));
  //  attr.set_string_attributes(string_values);


//    req.set_attribute_update(attr);


 //   rf.push(req);


   // let resp = client.report(grpc::RequestOptions::new(), req);
    ;


  //  println!("send report request: {} count",resp.wait_drop_metadata().count());
}
