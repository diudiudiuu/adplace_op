use ed25519_dalek::{SigningKey, Signature};
use hex;

/// 定义常量私钥
const PRIVATE: &str = "b23dfd4a2ca755e6ddccfcaa31201d1f6ec7fe0290086d53b931e0cd1a6f80fd";

/// 对字符串进行 ed25519 签名
pub fn sign(message: &str) -> Result<String, Box<dyn std::error::Error>> {
    // signing_key 加载私钥
    let private_key_bytes = hex::decode(PRIVATE)?;
    let signing_key = SigningKey::from_bytes(&private_key_bytes)?;

    let message_bytes = message.as_bytes();
    let signature: Signature = signing_key.try_sign(message_bytes)?;

    // 返回签名结果
    Ok(hex::encode(signature.to_bytes()))
}
