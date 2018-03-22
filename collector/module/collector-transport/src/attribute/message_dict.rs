use std::collections::HashMap;
use super::global_dict::GlobalDictionary;

pub struct MessageDictionary {

    global_dict: GlobalDictionary,
    message_words: Vec<String>,
    message_dict:  HashMap<String,i32>
}

// return public index of the message which is negative per API
fn message_dict_index( idx: i32) -> i32 {
    return -(idx + 1);
}

impl MessageDictionary  {

    pub fn new(global_dict: GlobalDictionary) -> MessageDictionary  {

        MessageDictionary {
            global_dict,
            message_words: Vec::new(),
            message_dict: HashMap::new()
        }

    }

    pub fn get_words(&self) -> &Vec<String>  {
        return &self.message_words;
    }

    //find index, try look up in the global, otherwise look up in the local
    pub fn index_of(&mut self, name: &str) -> i32  {

        if let Some(index) = self.global_dict.index_of(name) {
            return *index;
        }

        if let Some(index) = self.message_dict.get(name) {
           return message_dict_index(*index);
        }

        let index  = self.message_words.len() as i32;
        self.message_words.push(String::from(name));
        self.message_dict.insert(String::from(name),index as i32 );

        message_dict_index(index)

    }

    pub fn global_dict_size(&self) -> usize {
        self.global_dict.size()
    }

}


