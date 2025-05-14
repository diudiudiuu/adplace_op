use crate::tools::aes;
use crate::tools::json;
use reqwest;
use tokio::runtime::Runtime;

#[tauri::command]
// 全部列表
pub async fn list(authorization: String) -> String {
    let data = json::load_json_file(&authorization).await;
    data
}

#[tauri::command]
// 服务器详情
pub async fn server_info(server_id: String, authorization: String) -> String {
    // 根据服务器ID获取服务器信息
    let data = json::get_server_by_id(&server_id, &authorization).await;
    data.unwrap_or("{}".to_string())
}

#[tauri::command]
// 客户详情
pub async fn project_info(project_id: String, authorization: String) -> String {
    // 根据客户ID获取客户信息
    let data = json::get_project_by_id(&project_id, &authorization).await;
    data.unwrap_or("{}".to_string())
}

#[tauri::command]
// 添加/更新客户
pub async fn project_form(server_id: String, project_info: String, authorization: String) -> String {
    // 添加客户
    let data = json::project_form(&server_id, &project_info, &authorization).await;
    data
}

#[tauri::command]
// 执行命令
pub fn exec(project_id: String, sql: String, sql_type: String, _authorization: String) -> String {
    // 根据客户ID获取客户信息
    async fn exec(project_id: String, sql: String, sql_type: String, _authorization: String) -> String {
        // 根据客户ID获取客户信息
        let data = json::get_project_by_id(&project_id, &_authorization).await;
        let project = data.unwrap_or("{}".to_string());
        let project: serde_json::Value = serde_json::from_str(&project).unwrap();

        // project_api_url
        let project_api_url = project["project_api_url"].as_str().unwrap();
        // 向这个地址发送post 表单请求
        let url = format!("{}/dbexec", project_api_url);

        let signature = aes::encrypt(sql.as_str());
        let params = [("sql_type", sql_type), ("signature", signature)];
        let client = reqwest::Client::new();
        let res = client.post(url).form(&params).send().await;
        let body = match res {
            Ok(response) => response.text().await.unwrap(),
            // 这里错误也要正常返回
            Err(e) => {
                format!("{{\"code\": 500, \"msg\": \"{}\"}}", e)
            }
        };

        body
    }

    let rt = Runtime::new().unwrap();
    let result = rt.block_on(exec(project_id, sql, sql_type, _authorization));
    result
}
