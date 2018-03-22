extern crate ngx_rust;

use ngx_rust::bindings:: { ngx_uint_t, ngx_str_t } ;

use nginmesh_collector_transport::attribute::attr_wrapper::AttributeWrapper;
use nginmesh_collector_transport::attribute::global_dict::{ SOURCE_IP, SOURCE_UID, SOURCE_SERVICE, SOURCE_PORT,
                    DESTINATION_SERVICE,DESTINATION_IP,DESTINATION_UID
};

use super::config::CollectorConfig;

#[repr(C)]
pub struct ngx_http_collector_srv_conf_t {

    pub destination_service:    ngx_str_t,
    pub destination_uid:        ngx_str_t,
    pub destination_ip:         ngx_str_t,
    pub source_ip:              ngx_str_t,
    pub source_uid:             ngx_str_t,
    pub source_service:         ngx_str_t,
    pub source_port:            ngx_uint_t
}

impl CollectorConfig for  ngx_http_collector_srv_conf_t  {

    fn process_istio_attr(&self, attr: &mut AttributeWrapper) {

        attr.insert_string_attribute( DESTINATION_SERVICE, self.destination_service.to_str());
        attr.insert_string_attribute( DESTINATION_UID, self.destination_uid.to_str());
        attr.insert_string_attribute( DESTINATION_IP, self.destination_ip.to_str());
        attr.insert_string_attribute( SOURCE_IP,self.source_ip.to_str());
        attr.insert_string_attribute(SOURCE_UID,self.source_uid.to_str());
        attr.insert_string_attribute(SOURCE_SERVICE,self.source_service.to_str());
        attr.insert_int64_attribute(SOURCE_PORT,self.source_port as i64);

    }


}


