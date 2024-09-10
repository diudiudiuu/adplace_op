
// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

mod api;
mod tools;

use api::server;

fn main() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![
            server::list,
            server::project_info,
            server::exec
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}


