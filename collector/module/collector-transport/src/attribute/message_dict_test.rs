use super::global_dict::GlobalDictionary;
use super::message_dict::MessageDictionary;

// test for accessing global dictionary
#[test]
fn test_message_dict_global() {
    let global_dict = GlobalDictionary::new();
    let mut dict = MessageDictionary::new(global_dict);
    let index = dict.index_of("source.port");
    assert_eq!(index, 1, "check source port");
}

// test if we are adding new message word
#[test]
fn test_message_dict_local_first() {
    let global_dict = GlobalDictionary::new();
    let mut dict = MessageDictionary::new(global_dict);
    assert_eq!(dict.index_of("unknown.x"),-1,"check new");
    assert_eq!(dict.index_of("unknown.x"),-1,"check new");      // existing
    assert_eq!(dict.index_of("unknown.y"),-2,"check new");      // new
}
