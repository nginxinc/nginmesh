extern crate reqwest;
extern crate ngx_mixer_test;
#[macro_use]
extern crate hyper;

use ngx_mixer_test::util::make;
use std::io::Read;
use hyper::header::Headers;
use reqwest::StatusCode;

header! { (Clnt, "clnt") => [String] }

// in order run these test, please run following make step:  'make test-nginx-only'
#[test]
fn nginx_check_basic_test()  {

  //  let _result = make("test-nginx-only");

    let mut response = reqwest::get("http://localhost:8000/check").unwrap();
    assert!(response.status().is_success(),"nginx test check succedd");

    let mut content = String::new();
    response.read_to_string(&mut content);

    println!("response: {}",content);
    assert_eq!(content,"9100","should return local services");
}


// force deny

#[test]
fn nginx_check_deny_test()  {


   // let _result = make("test-nginx-only");

    let client = reqwest::Client::new();
    let mut headers = Headers::new();
    headers.set(Clnt("abc".to_owned()));
    let mut response = client.get("http://localhost:8000/check")
        .headers(headers)
        .send()
        .unwrap();


    assert_eq!(response.status(),StatusCode::Unauthorized,"expected 401");

}
