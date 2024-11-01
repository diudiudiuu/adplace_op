
use std::fs::File;
use std::io::prelude::*;
use std::path::PathBuf;
use serde_json::{Value};



// 加载当前目录下的json文件 返回json 数据
pub fn load_json_file() -> String {
    // 获取当前工作目录
    let current_dir = std::env::current_dir().expect("Failed to get current directory");
    // 设定文件路径
    let file_path = current_dir.join("data.json");

    // 判断文件是否存在 如果不存在则创建
    if !file_path.exists() {
        let mut file = File::create(&file_path).expect("Failed to create file");
        file.write_all(b"[]").expect("Failed to write to file");
    }

    let data = read(file_path);
    data
}

// 根据 server_id 获取服务器信息
pub fn get_server_by_id(server_id: &str) -> Option<String> {
    let data = load_json_file();
    // 将 JSON 字符串转换为 JSON 对象
    let json: Value = serde_json::from_str(&data).ok()?; // 安全地解析 JSON

    // 检查 JSON 是否是一个数组
    if let Some(servers) = json.as_array() {
        // 循环 JSON 数组
        for server in servers {
            // 判断是否是当前服务器
            if let Some(id) = server["server_id"].as_str() {
                if id == server_id {
                    return Some(server.to_string()); // 返回找到的服务器
                }
            }
        }
    }

    None // 如果没有找到，返回 None
}

// 根据 project_id 获取项目信息
pub fn get_project_by_id(project_id: &str) -> Option<String> {
    let data = load_json_file();
    // 将 JSON 字符串转换为 JSON 对象
    let json: Value = serde_json::from_str(&data).ok()?; // 安全地解析 JSON

    // 检查 JSON 是否是一个数组
    if let Some(servers) = json.as_array() {
        // 循环 JSON 数组
        for server in servers {
            // 检查 "project_list" 是否是一个数组
            if let Some(project_list) = server["project_list"].as_array() {
                // 循环 project_list 数组
                for project in project_list {
                    // 判断是否是当前项目
                    if let Some(id) = project["project_id"].as_str() {
                        if id == project_id {
                            return Some(project.to_string()); // 返回找到的项目
                        }
                    }
                }
            }
        }
    }

    None // 如果没有找到，返回 None
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