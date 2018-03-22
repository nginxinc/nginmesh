

use std::thread ;

use ngx_rust::bindings::ngx_int_t;
use ngx_rust::bindings::NGX_OK;


// start background activities
#[no_mangle]
pub extern fn nginmesh_collector_init() -> ngx_int_t {
    /*
    ngx_event_debug!("init mixer start ");
    thread::spawn(|| {
        ngx_event_debug!("starting mixer report background task");
        collector_report_background();
    });


    ngx_event_debug!("init mixer end ");
    */
    return NGX_OK as ngx_int_t;
}

#[no_mangle]
pub extern fn nginmesh_collector_exit() {

    ngx_event_debug!("mixer exit ");
}

