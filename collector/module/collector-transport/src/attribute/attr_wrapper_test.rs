use protobuf::well_known_types::Timestamp;
use std::collections::HashMap;

use super::global_dict::GlobalDictionary;
use super::message_dict::MessageDictionary;
use super::attr_wrapper::AttributeWrapper;

const TEST_AGENT: &str = "mac123";

#[test]
fn simple_string_mapping() {

    let global_dict = GlobalDictionary::new();
    let mut dict = MessageDictionary::new(global_dict);
    let mut attr_wrapper = AttributeWrapper::new();

    attr_wrapper.insert_string_attribute("source.ip","10.0.0.0");
    attr_wrapper.insert_string_attribute("destination.ip","10.0.0.0");

    let attributes = attr_wrapper.as_attributes(&mut dict);
    let index = attributes.get_strings().get(&0).unwrap();
    assert_eq!(*index,-1);

    let destination_ip_index = dict.index_of("destination.ip");
    let index = attributes.get_strings().get(&destination_ip_index).unwrap();
    assert_eq!(*index,-1);

}



#[test]
fn test_attr_key_exists() {

     let mut attr_wrapper = AttributeWrapper::new();

    attr_wrapper.insert_string_attribute("source.ip","10.0.0.0");
    attr_wrapper.insert_string_attribute("destination.ip","10.0.0.0");

    assert!(attr_wrapper.key_exists("source.ip"));
    assert!(!attr_wrapper.key_exists("source4.ip"));
}


#[test]
fn simple_int64_mapping() {

    let global_dict = GlobalDictionary::new();
    let mut dict = MessageDictionary::new(global_dict);
    let mut attr_wrapper = AttributeWrapper::new();

    attr_wrapper.insert_int64_attribute("response.duration",50);

    let attributes = attr_wrapper.as_attributes(&mut dict);
    let response_dur_index = dict.index_of("response.duration");
    let duration = attributes.get_int64s().get(&response_dur_index).unwrap();
    assert_eq!(*duration,50);
}

#[test]
fn simple_double_mapping() {

    let global_dict = GlobalDictionary::new();
    let mut dict = MessageDictionary::new(global_dict);
    let mut attr_wrapper = AttributeWrapper::new();

    attr_wrapper.insert_f64_attribute("response.duration", 0.5_f64);

    let attributes = attr_wrapper.as_attributes(&mut dict);
    let response_dur_index = dict.index_of("response.duration");
    let duration = attributes.get_doubles().get(&response_dur_index).unwrap();
    assert_eq!(*duration,0.5_f64);
}

#[test]
fn simple_bool_mapping() {

    let global_dict = GlobalDictionary::new();
    let mut dict = MessageDictionary::new(global_dict);
    let mut attr_wrapper = AttributeWrapper::new();

    attr_wrapper.insert_bool_attribute("true", true);

    let attributes = attr_wrapper.as_attributes(&mut dict);
    let response_dur_index = dict.index_of("true");
    let duration = attributes.get_bools().get(&response_dur_index).unwrap();
    assert_eq!(*duration,true);
}

#[test]
fn simple_time_stamp() {

    let global_dict = GlobalDictionary::new();
    let mut dict = MessageDictionary::new(global_dict);
    let mut attr_wrapper = AttributeWrapper::new();

    let mut request_time = Timestamp::new();
    request_time.set_seconds(1000);
    attr_wrapper.insert_time_stamp_attribute("request.time", request_time);

    let attributes = attr_wrapper.as_attributes(&mut dict);
    let response_dur_index = dict.index_of("request.time");
    let duration = attributes.get_timestamps().get(&response_dur_index).unwrap();

    assert_eq!(duration.get_seconds(),1000);
}



#[test]
fn simple_stringmap_mapping() {
    let global_dict = GlobalDictionary::new();
    let mut dict = MessageDictionary::new(global_dict);
    let mut attr_wrapper = AttributeWrapper::new();

    let mut string_map: HashMap<String,String> = HashMap::new();
    string_map.insert(String::from("request.scheme"), String::from("http"));
    string_map.insert(String::from("request.useragent"), String::from(TEST_AGENT));
    attr_wrapper.insert_string_map("request.headers", string_map);

    let attributes = attr_wrapper.as_attributes(&mut dict);

    let str_map = attributes.get_string_maps().get(&dict.index_of("request.headers")).unwrap();
    let str_http_index = str_map.get_entries().get(&dict.index_of("request.scheme")).unwrap();
    assert_eq!(*str_http_index, dict.index_of("http"));

    let mac_index = dict.index_of(TEST_AGENT);
    println!("mac index: {}",mac_index);
    println!("words:  {:?}",attributes.get_words());
    assert_eq!(attributes.get_words().get( (mac_index * -1  -1 ) as usize ).unwrap(),TEST_AGENT);
}
