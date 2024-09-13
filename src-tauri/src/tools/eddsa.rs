

use ed25519_dalek::SigningKey;
use ed25519_dalek::Signature;
use ed25519_dalek::SECRET_KEY_LENGTH;
use hex::FromHex;
use rand::rngs::OsRng;
use ed25519_dalek::Signer;
use sha2::Sha512;

// 信息签名
pub fn signature(message: &str) -> String {
    
    // let sec_bytes = hex!("b23dfd4a2ca755e6ddccfcaa31201d1f6ec7fe0290086d53b931e0cd1a6f80fd");
    // let signing_key = SigningKey::from_bytes(&sec_bytes);

    // let mut prehash_for_signing = Sha512::default();
    // prehash_for_signing.update(message.as_bytes());

    // let signature: Signature = signing_key.sign_prehashed(prehash_for_signing, None).unwrap();
    // hex::encode(signature.to_bytes())


    let sec_bytes: Vec<u8> = FromHex::from_hex("302e020100300506032b6570042204201aeb9c77c50e6b877ebfc99492f792798ad18bcc36163d85abfb66eae17bd51").unwrap();
    let sec_bytes = &sec_bytes[..SECRET_KEY_LENGTH].try_into().unwrap();
    let signing_key = SigningKey::from_bytes(sec_bytes);
    let msg_bytes = message.as_bytes();
    let signature: Signature = signing_key.sign(msg_bytes);
    hex::encode(signature.to_bytes())
}

// 生成密钥对
pub fn generate() -> (String, String) {
    let mut csprng = OsRng;
    let signing_key = SigningKey::generate(&mut csprng);
   
   "123".to_string();
    let verifying_key = signing_key.verify_key();
    (hex::encode(signing_key.to_bytes()), hex::encode(verifying_key.to_bytes()))
}