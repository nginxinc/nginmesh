extern crate ngx_rust;

use ngx_rust::bindings:: { ngx_str_t, ngx_flag_t } ;
use nginmesh_collector_transport::attribute::attr_wrapper::AttributeWrapper;
use nginmesh_collector_transport::attribute::global_dict::{ DESTINATION_SERVICE };

use ngx::config::CollectorConfig;

#[repr(C)]
pub struct ngx_http_collector_loc_conf_t {
    pub enable_report: ngx_flag_t,              // for every location, we need flag to enable/disable mixer
    pub enable_check: ngx_flag_t,
    pub destination_service: ngx_str_t
}

impl CollectorConfig for ngx_http_collector_loc_conf_t  {


    fn process_istio_attr(&self,attr: &mut AttributeWrapper) {

        attr.insert_string_attribute( DESTINATION_SERVICE,self.destination_service.to_str());

    }


}


