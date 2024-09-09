
use crate::tools::json;

#[tauri::command]
pub fn list() -> String {
    let data = json::load_json_file();
    data
}

