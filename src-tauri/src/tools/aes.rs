use aes::Aes256;
use aes::cipher::{
    BlockEncrypt, KeyInit,
    generic_array::GenericArray,
};
use rand::Rng;

// 进行字符串加密,返回加密后的字符串以及加密后的key
pub fn encrypt(data: &str) -> (String, String) {
    // 生成一个随机的key (32 bytes for Aes256)
    let key: [u8; 32] = rand::thread_rng().gen();
    let key_array = GenericArray::from_slice(&key);

    // 生成一个加密器
    let cipher = Aes256::new(&key_array);

    // Ensure data length is 16 bytes by padding (AES block size is 16 bytes)
    let mut padded_data = data.as_bytes().to_vec();
    padded_data.resize(16, 0); // pad with zeros if less than 16 bytes

    // Convert to GenericArray for encryption
    let mut block = GenericArray::clone_from_slice(&padded_data);

    // 加密
    cipher.encrypt_block(&mut block);

    // 返回加密后的字符串以及加密的key（转换为Hex编码，避免UTF-8错误）
    (hex::encode(block.to_vec()), hex::encode(key))
}