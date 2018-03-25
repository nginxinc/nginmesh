use std::sync::mpsc::{ Sender,Receiver};
use std::sync::Mutex;


pub struct Channels<T> {
    pub tx: Mutex<Sender<T>>,
    pub rx: Mutex<Receiver<T>>
}
