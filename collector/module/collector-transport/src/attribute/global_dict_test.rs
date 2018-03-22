use super::global_dict::GlobalDictionary;

#[test]
fn test_global_simple() {

    let dict = GlobalDictionary::new();
    let index = dict.index_of("source.port").unwrap();
    assert_eq!(*index,1,"check source port");
    let index = dict.index_of("source.service").unwrap();
    assert_eq!(*index,158,"check source service");
}
