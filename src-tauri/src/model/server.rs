
use rusqlite::NO_PARAMS;
use serde::{Deserialize, Serialize};
use crate::model::db;


// 查询 数据库中server列表

pub fn list() -> Vec<Server> {
    let conn = db::get_conn();
    let mut stmt = conn
        .prepare("SELECT id, name, host, port, username, password FROM server")
        .unwrap();
    let servers = stmt
        .query_map([], |row| {
            Ok(Server {
                id: row.get(0)?,
                name: row.get(1)?,
                host: row.get(2)?,
                port: row.get(3)?,
                username: row.get(4)?,
                password: row.get(5)?,
            })
        })
        .unwrap()
        .map(|r| r.unwrap())
        .collect();
    servers
}