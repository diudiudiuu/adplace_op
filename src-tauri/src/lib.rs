#[cfg_attr(mobile, tauri::mobile_entry_point)]

mod api;
mod tools;

use api::server;


pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_fs::init())
        .invoke_handler(tauri::generate_handler![
            server::list,         // 服务器列表
            server::server_info,  // 服务器信息
            server::project_info, // 客户信息
            server::project_form, // 添加/更新客户
            server::exec          // 远程服务器执行命令
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
