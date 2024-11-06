use aes::cipher::{KeyIvInit, block_padding::Pkcs7, BlockEncryptMut};
use base64::{Engine as _, engine::general_purpose};
use bytebuffer::ByteBuffer;
use rand::distributions::{Alphanumeric, DistString};

type Aes256CbcEnc = cbc::Encryptor<aes::Aes256>;

// 加密
pub fn encrypt(key: &str, data: &str) -> String {
    let mut rng = rand::thread_rng();
    let iv_str = Alphanumeric.sample_string(&mut rng, 16);
    let iv = iv_str.as_bytes();
    
    // 加密数据
    let cipher = Aes256CbcEnc::new(key.as_bytes().into(), iv.into())
        .encrypt_padded_vec_mut::<Pkcs7>(data.as_bytes());
    
    // 将 IV 和密文拼接
    let mut buffer = ByteBuffer::from_bytes(iv);
    buffer.write_bytes(&cipher);
    
    // 返回 Base64 编码后的拼接结果
    general_purpose::STANDARD.encode(buffer.as_bytes())
}
