use std::process::Command;
use std::process::Output;
use std::env;
use std::io::Result;


// perform make with argument
pub fn make(arg: &str) -> Result<Output> {
    let current_path = env::current_dir().unwrap();
    let make_path = current_path.parent().unwrap();
    let path_name = format!("{}",make_path.display());
    println!("executing make command at {}",path_name);
    let result =  Command::new("/usr/bin/make")
        .args(&[arg])
        .current_dir(path_name)
        .output();

    match result  {
        Err(e)  =>  {
            println!("make error: {}", e);
            return Err(e);
        },

        Ok(output) => {
            println!("status: {}", output.status);
            println!("stdout: {}", String::from_utf8_lossy(&output.stdout));
            println!("stderr: {}", String::from_utf8_lossy(&output.stderr));
            return Ok(output);
        }
    }
}
