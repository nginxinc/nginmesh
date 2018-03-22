
use nginmesh_collector_transport::attribute::attr_wrapper::AttributeWrapper;

pub trait CollectorConfig {

    // convert and migrate values to istio attributes
    fn process_istio_attr(&self, attr: &mut AttributeWrapper);

}
