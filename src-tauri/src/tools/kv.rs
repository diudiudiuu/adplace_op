use reqwest::{Client, Method};
use serde::{Deserialize, Serialize};
use std::error::Error;

const KV_BASE_URL: &str = "https://kv.adswds.com/1ep2d8wb";

#[derive(Debug, Deserialize, Serialize)]
pub struct KvData {
    pub key: String,
    pub value: String,  // 这里直接是 String
}


// 响应体结构体
#[derive(Debug, Deserialize, Serialize)]
pub struct ApiResponse {
    pub code: u16,
    pub msg: String,
    pub data: Option<KvData>,
}

// 公共请求方法
async fn request(
    method: Method,
    authorization: &str,
    key: &str,
    value: Option<&str>,
) -> Result<ApiResponse, Box<dyn Error>> {
    let client = Client::new();

    let mut url = format!("{}?key={}", KV_BASE_URL, key);
    if let Some(v) = value {
        url = format!("{}&value={}", url, v);
    }

    let resp = client
        .request(method, &url)
        .header("Authorization", format!("Bearer {}", authorization))
        .send()
        .await?
        .json::<ApiResponse>()
        .await?;

    Ok(resp)
}

// 新增数据（POST）
pub async fn create_key(key: &str, value: &str, authorization: &str) -> ApiResponse {
    let res = request(Method::POST, authorization, key, Some(value)).await.unwrap();
    println!("POST result: {:?}", res);
    res
}

// 查询数据（GET）
pub async fn get_key(key: &str, authorization: &str) -> ApiResponse {
    let res = request(Method::GET, authorization, key, None).await.unwrap();
    println!("GET result: {:?}", res);
    res
}

// 更新数据（PUT）
pub async fn update_key(key: &str, value: &str, authorization: &str) -> ApiResponse {
    let res = request(Method::PUT, authorization, key, Some(value)).await.unwrap();
    println!("PUT result: {:?}", res);
    res
}

