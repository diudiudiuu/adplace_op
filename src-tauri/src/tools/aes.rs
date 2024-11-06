use aes::cipher::{KeyIvInit, block_padding::Pkcs7, BlockEncryptMut};
use base64::{Engine as _, engine::general_purpose};
use bytebuffer::ByteBuffer;
use rand::distributions::{Alphanumeric, DistString};
use std::time::{SystemTime, UNIX_EPOCH};

type Aes256CbcEnc = cbc::Encryptor<aes::Aes256>;

// 定义常量密钥（32字节长度）
const KEY: &str = "3HXV4P8dizvATG5EjLIsUKxSreyghDMB";

// 加密函数
pub fn encrypt(data: &str) -> String {
    let mut rng = rand::thread_rng();
    let iv_str = Alphanumeric.sample_string(&mut rng, 16);
    let iv = iv_str.as_bytes();

    // 获取当前 UTC 时间戳（秒）
    let timestamp = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_secs();
    
    // 加密数据
    let cipher = Aes256CbcEnc::new(KEY.as_bytes().into(), iv.into())
        .encrypt_padded_vec_mut::<Pkcs7>(data.as_bytes());
    
    // 拼接时间戳和 IV + 密文
    let mut buffer = ByteBuffer::new();
    buffer.write_u64(timestamp); // 写入时间戳
    buffer.write_bytes(iv);       // 写入 IV
    buffer.write_bytes(&cipher);  // 写入加密后的数据
    
    // 返回 Base64 编码的结果
    general_purpose::STANDARD.encode(buffer.as_bytes())
}
