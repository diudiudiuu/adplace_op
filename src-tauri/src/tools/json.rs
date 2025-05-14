use serde_json::Value;
use crate::tools::kv;

const KV_KEY: &str = "client_json";

// 加载 KV 中 JSON 字符串
pub async fn load_json_file(authorization: &str) -> String {
    let json = kv::get_key(KV_KEY, authorization).await;

    if json.code == 200 {
        if let Some(data) = json.data {
            let value_str = data.value;
            // value_str 是个 JSON 字符串，再把它反序列化成 serde_json::Value 看看
            let value_json: serde_json::Value = serde_json::from_str(&value_str).unwrap_or_else(|_| serde_json::json!([]));
            return value_json.to_string();
        } else {
            return "[]".to_string();
        }
    } else {
        let empty_json = serde_json::json!([]);
        let empty_json_str = empty_json.to_string();
        kv::create_key(KV_KEY, &empty_json_str, authorization).await;
        empty_json_str
    }
}

// 根据 server_id 获取服务器信息
pub async fn get_server_by_id(server_id: &str, authorization: &str) -> Option<String> {
    let data = load_json_file(authorization).await;
    let json: Value = serde_json::from_str(&data).ok()?;

    if let Some(servers) = json.as_array() {
        for server in servers {
            if let Some(id) = server["server_id"].as_str() {
                if id == server_id {
                    return Some(server.to_string());
                }
            }
        }
    }
    None
}

// 根据 project_id 获取客户信息
pub async fn get_project_by_id(project_id: &str, authorization: &str) -> Option<String> {
    let data = load_json_file(authorization).await;
    let json: Value = serde_json::from_str(&data).ok()?;

    if let Some(servers) = json.as_array() {
        for server in servers {
            if let Some(project_list) = server["project_list"].as_array() {
                for project in project_list {
                    if let Some(id) = project["project_id"].as_str() {
                        if id == project_id {
                            return Some(project.to_string());
                        }
                    }
                }
            }
        }
    }
    None
}

// 添加/更新客户
pub async fn project_form(server_id: &str, project_info: &str, authorization: &str) -> String {
    let data = load_json_file(authorization).await;
    let mut json: Value = serde_json::from_str(&data).unwrap_or_else(|_| serde_json::json!([]));

    if let Some(servers) = json.as_array_mut() {
        for server in servers {
            if let Some(id) = server["server_id"].as_str() {
                if id == server_id {
                    let project_json: Value = serde_json::from_str(project_info).unwrap_or_else(|_| serde_json::json!({}));

                    if let Some(project_list) = server["project_list"].as_array_mut() {
                        let mut project_found = false;

                        for project in project_list.iter_mut() {
                            if let Some(project_id) = project["project_id"].as_str() {
                                if project_id == project_json["project_id"].as_str().unwrap_or("") {
                                    *project = project_json.clone();
                                    project_found = true;
                                    break;
                                }
                            }
                        }

                        if !project_found {
                            project_list.push(project_json);
                        }

                        // 保存到 KV
                        kv::update_key(KV_KEY, &json.to_string(), authorization).await;

                        return "{}".to_string();
                    }
                }
            }
        }
    }

    "{}".to_string()
}

