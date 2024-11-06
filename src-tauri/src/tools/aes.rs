use aes::cipher::{KeyIvInit, block_padding::Pkcs7, BlockEncryptMut};
use base64::{Engine as _, engine::general_purpose};
use rand::distributions::{Alphanumeric, DistString};

type Aes256CbcEnc = cbc::Encryptor<aes::Aes256>;

// 加密
pub fn encrypt(key: &str, data: &str) -> (String, String) {
    let mut rng = rand::thread_rng();
    let iv_str = Alphanumeric.sample_string(&mut rng, 16);
    let iv = iv_str.as_bytes();
    
    // 加密数据但不拼接 IV 和密文
    let cipher = Aes256CbcEnc::new(key.as_bytes().into(), iv.into())
        .encrypt_padded_vec_mut::<Pkcs7>(data.as_bytes());

    // 返回 IV 和 Base64 编码的密文
    (iv_str, general_purpose::STANDARD.encode(&cipher))
}
