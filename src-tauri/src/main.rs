
// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

mod model;

use model::db::init;
use tauri::Manager;

fn main() {
    let pool = init().expect("init db failed");

    tauri::Builder::default()
        .setup(move |app| {
            app.manage(pool);
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}


