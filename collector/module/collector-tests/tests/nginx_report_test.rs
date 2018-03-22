extern crate reqwest;
extern crate ngx_mixer_test;

#[macro_use]
extern crate hyper;

use ngx_mixer_test::util::make;
use std::io::Read;


#[test]
fn nginx_report_test()  {

   // let _result = make("test-nginx-only");

    let mut response = reqwest::get("http://localhost:8000/report").unwrap();
    assert!(response.status().is_success(),"nginx test check succedd");

    let mut content = String::new();
    response.read_to_string(&mut content);

    println!("response: {}",content);
    assert_eq!(content,"9100","should return local services");
}