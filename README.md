# 广告投放管理面板

## 启动项目

    npm run tauri dev  


## 打包项目

    npm run tauri build



## openssl EdDSA

    # 生成私钥
    openssl genpkey -algorithm ed25519 -out ed25519_private.pem

    # 从私钥生成公钥
    openssl pkey -in ed25519_private.pem -pubout -out ed25519_public.pem

    # 提取私钥和公钥的原始字节（不是PEM格式，直接用于编码）
    openssl pkey -in ed25519_private.pem -outform DER -out ed25519_private.der
    openssl pkey -in ed25519_private.pem -pubout -outform DER -out ed25519_public.der
