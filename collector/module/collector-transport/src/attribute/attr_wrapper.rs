
use chrono::{Utc,DateTime};
use serde_json::{Map,Value,to_string};

pub struct AttributeWrapper {

    values: Map<String,Value>,       // map of value
}


impl AttributeWrapper  {

    pub fn new() -> AttributeWrapper {
        AttributeWrapper {
            values: Map::new()
        }
    }

    #[allow(dead_code)]
    pub fn key_exists(&self, key: &str) -> bool  {

        self.values.contains_key(key)
    }


    // insert string attributes
    pub fn insert_value(&mut self, key: &str, value: Value) {
        self.values.insert(String::from(key),value);

    }

    pub fn insert_string_attribute(&mut self, key: &str, value: &str) {
        if value.len() > 0 {
            self.insert_value(key, json!(String::from(value)));
        }
    }

    pub fn insert_int64_attribute(&mut self, key: &str, value: i64) {
        self.insert_value(key,json!(value));
    }

    #[allow(dead_code)]
    pub fn insert_f64_attribute(&mut self, key: &str, value: f64) {
        self.insert_value(key,json!(value));
    }

    #[allow(dead_code)]
    pub fn insert_bool_attribute(&mut self, key: &str, value: bool) {
        self.insert_value(key,json!(value));
    }


    pub fn insert_time_stamp_attribute(&mut self, key: &str, value: DateTime<Utc>) {

        self.insert_value(key,json!(value.to_string()));
    }



    pub fn to_string(&self) -> String {
        let result = to_string(&self.values);
        if result.is_ok() {
            result.ok().unwrap()
        } else {
            "Error".to_owned()
        }
    }




}

