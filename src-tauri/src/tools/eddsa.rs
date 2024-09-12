

use ed25519_dalek::SigningKey;
use ed25519_dalek::Signature;
use sha2::{Sha512, Digest};
use hex_literal::hex;


pub fn sign(message: &str) -> String {
    
    let sec_bytes = hex!("b23dfd4a2ca755e6ddccfcaa31201d1f6ec7fe0290086d53b931e0cd1a6f80fd");
    let signing_key = SigningKey::from_bytes(&sec_bytes);

    let mut prehash_for_signing = Sha512::default();
    prehash_for_signing.update(message.as_bytes());

    let signature: Signature = signing_key.sign_prehashed(prehash_for_signing, None).unwrap();
    hex::encode(signature.to_bytes())
}
