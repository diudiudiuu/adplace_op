
use crate::tools::json;
use crate::tools::eddsa;
use reqwest;
use tokio::runtime::Runtime;




#[tauri::command]
// 全部列表
pub fn list() -> String {
    let data = json::load_json_file();
    data
}

#[tauri::command]
// 项目详情
pub fn project_info(project_id: String) -> String {
    // 根据项目ID获取项目信息
    let data = json::get_project_by_id(&project_id);
    data.unwrap_or("{}".to_string())
}

#[tauri::command]
// 执行命令
pub fn exec(project_id: String, sql: String, sql_type: String) -> String {
    // 根据项目ID获取项目信息
    async fn exec(project_id: String, sql: String, sql_type: String) -> String {
        // 根据项目ID获取项目信息
        let data = json::get_project_by_id(&project_id);
        let project = data.unwrap_or("{}".to_string());
        let project: serde_json::Value = serde_json::from_str(&project).unwrap();

        // project_api_url
        let project_api_url = project["project_api_url"].as_str().unwrap();

        // 向这个地址发送post 表单请求 
        let url = format!("{}/dbexec", project_api_url);

        let signature = eddsa::signature(&sql);
        let params = [
            ("sql", sql),
            ("signature", signature),
            ("sql_type", sql_type),
        ];
        let client = reqwest::Client::new();
        let res = client.post(url)
            .form(&params)
            .send()
            .await;
        let body = match res {
            Ok(response) => response.text().await.unwrap(),
            Err(error) => panic!("Error: {}", error),
        };
        
        body
    }
    
    let rt = Runtime::new().unwrap();
    let result = rt.block_on(exec(project_id, sql, sql_type));
    result
}
