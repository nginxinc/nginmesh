extern crate ngx_rust;

use ngx_rust::bindings:: { ngx_int_t, ngx_str_t } ;


#[repr(C)]
pub struct ngx_http_collector_main_conf_t {

    pub collector_server: ngx_str_t
}


