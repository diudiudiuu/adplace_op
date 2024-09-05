use rusqlite::Result;
use r2d2::Pool;
use r2d2_sqlite::SqliteConnectionManager;
use std::path::PathBuf;
use dirs;

pub type SqlitePool = Pool<SqliteConnectionManager>;


pub fn init() -> Result<SqlitePool> {
    let db_path = db_path();
    if !db_path.exists() {
        print!("db path: {:?}", db_path);
        std::fs::create_dir_all(db_path.parent().unwrap()).unwrap();
    }

    let manager = SqliteConnectionManager::file(db_path);
    let pool = Pool::new(manager).expect("db pool failed");
    let conn = pool.get().expect("db connection failed");
    conn.execute(
        "CREATE TABLE IF NOT EXISTS todo (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            done INTEGER NOT NULL DEFAULT 0
        )",
        [],
    )?;

    Ok(pool)
}

// 获取数据库路径 同级目录下的 db.sqlite
fn db_path() -> PathBuf {
    let mut path = dirs::data_dir().unwrap_or_else(|| PathBuf::from("."));
    path.push("todo");
    
    // 确保目录存在
    std::fs::create_dir_all(&path).expect("Failed to create database directory");
    
    path.push("db.sqlite");
    path
}
