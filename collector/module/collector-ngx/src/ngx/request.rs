
use serde_json::{Map,Value};
use chrono::{Utc,TimeZone};
use ngx_rust::bindings:: { ngx_http_request_s, ngx_http_headers_out_t} ;
use nginmesh_collector_transport::attribute::attr_wrapper::AttributeWrapper;
use nginmesh_collector_transport::attribute::global_dict::{ REQUEST_HEADER, REQUEST_HOST, REQUEST_METHOD, REQUEST_PATH,
                                                   REQUEST_REFER, REQUEST_SCHEME, REQUEST_SIZE, REQUEST_TIME, REQUEST_USERAGENT,
                                                   SOURCE_IP, SOURCE_UID, SRC_IP_HEADER, SRC_UID_HEADER,
                                                    RESPONSE_CODE, RESPONSE_SIZE, RESPONSE_HEADERS
};

use super::config::CollectorConfig;


impl CollectorConfig for ngx_http_request_s  {



    fn process_istio_attr(&self, attr: &mut AttributeWrapper )  {

        ngx_http_debug!(self,"send request headers to mixer");

        let headers_in = self.headers_in;


        attr.insert_string_attribute(REQUEST_HOST,  headers_in.host_str());
        attr.insert_string_attribute(REQUEST_METHOD, self.method_name.to_str());
        attr.insert_string_attribute(REQUEST_PATH, self.uri.to_str());
        
        ngx_http_debug!(self,"send referer to mixer");
        let referer = headers_in.referer_str();
        if let Some(ref_str) = referer {
            attr.insert_string_attribute(REQUEST_REFER, ref_str);
        }
        
        ngx_http_debug!(self,"send scheme to mixer");
        //let scheme = request.http_protocol.to_str();

        
        attr.insert_string_attribute(REQUEST_SCHEME, "http"); // hard code now
        attr.insert_int64_attribute(REQUEST_SIZE, self.request_length);

     //   attr.insert_time_stamp_attribute(REQUEST_TIME, Utc.timestamp(self.start_sec,self.start_msec as u32));
        attr.insert_time_stamp_attribute(REQUEST_TIME, Utc::now());
      //  attr.insert_string_attribute(REQUEST_USERAGENT, headers_in.user_agent_str());

        
        // fill in the string value
        let mut map: Map<String,Value> = Map::new();
        {
            for (name,value) in headers_in.headers_iterator()   {
                ngx_http_debug!(self,"in header name: {}, value: {}",&name,&value);

                // TODO: remove header
                match name.as_ref()  {

                    SRC_IP_HEADER  => {
                        ngx_http_debug!(self,"source IP received {}",&value);
                        attr.insert_string_attribute( SOURCE_IP,&value);
                    },

                    SRC_UID_HEADER => {
                        ngx_http_debug!(self,"source UID received {}",&value);
                        attr.insert_string_attribute( SOURCE_UID,&value);
                    },
                    _ => {
                        ngx_http_debug!(self,"other source header {}",&name);
                        map.insert(name,json!(value));
                    }
                }


            }
        }

        attr.insert_value(REQUEST_HEADER, json!(map));
        
    }
}



impl CollectorConfig for ngx_http_headers_out_t {

    fn process_istio_attr(&self, attr: &mut AttributeWrapper, )  {


       // ngx_http_debug!("send request header attribute to mixer");
        attr.insert_int64_attribute(RESPONSE_CODE, self.status as i64);
        attr.insert_int64_attribute(RESPONSE_SIZE, self.content_length_n);

        // fill in the string value
        let mut map: Map<String,Value> = Map::new();
        {
            for (name,value) in self.headers_iterator()   {
         //       ngx_http_debug!("processing out header name: {}, value: {}",&name,&value);

                map.insert(name,json!(value));

            }
        }

        attr.insert_value(RESPONSE_HEADERS, json!(map));
    }
}


