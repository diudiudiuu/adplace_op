
// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

mod api;
mod tools;

use api::server;

fn main() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![
            server::list, // 服务器列表
            server::server_info, // 服务器信息
            server::project_info, // 项目信息
            server::exec // 远程服务器执行命令
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}


