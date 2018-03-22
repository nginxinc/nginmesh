use std::collections::HashMap;
use base64::{encode, decode};
use std::str;
use std::vec::Vec;

/**
 * convert istio headers into a single string that can be sent as http header
*/
#[allow(dead_code)]
pub fn encode_istio_header(headers: &Vec<(&str,&str)>)  -> String {

    // for now we do simple serialize, we can convert to same format as envoy for version 0.2 with new mixer API

    let mut out = String::from("");

    for  &(key, value ) in headers  {
        out.push_str(key);
        out.push_str("@");
        out.push_str(value);
        out.push_str("!");
    }

    return encode(&out);
}

// decode istio header and convert to map
#[allow(dead_code)]
pub fn decode_istio_header(encoded_string: &str) -> HashMap<String,String>   {

    let decode_bytes = &decode(encoded_string).unwrap()[..];
    let decode_value = str::from_utf8(decode_bytes).unwrap();


    let mut out = HashMap::new();

    let attrs_list = decode_value.split("!");

    for attr in attrs_list  {
        let tokens = attr.split("@").collect::<Vec<&str>>();
        if tokens.len() == 2  {
            out.insert(String::from(tokens[0]),String::from(tokens[1]));
        }
    }

    return out;



}