
use std::fs::File;
use std::io::prelude::*;
use std::path::PathBuf;


// 加载当前目录下的json文件 返回json 数据
pub fn load_json_file() -> String {
    // 获取当前工作目录
    let current_dir = std::env::current_dir().expect("Failed to get current directory");
    // 设定文件路径
    let file_path = current_dir.join("data.json");
    let data = read(file_path);
    data
}

// 读取文件
fn read(file_path: PathBuf) -> String {
    let display = file_path.display();

    let mut file = match File::open(&file_path) {
        Err(why) => panic!("couldn't open {}: {}", display, why),
        Ok(file) => file,
    };

    let mut s = String::new();
    match file.read_to_string(&mut s) {
        Err(why) => panic!("couldn't read {}: {}", display, why),
        Ok(_) => println!("{} read success", display),
    }
    s
}