# 广告投放管理面板

## 启动客户

    npm run tauri dev  


## 打包客户

    npm run tauri build



## openssl EdDSA

    # 生成私钥
    openssl genpkey -algorithm ed25519 -out private.pem

    # 从私钥生成公钥
    openssl pkey -in private.pem -pubout -out public.pem

    # 提取私钥和公钥的原始字节（不是PEM格式，直接用于编码）
    openssl pkey -in private.pem -outform DER -out private.der
    openssl pkey -in private.pem -pubout -outform DER -out ed25519_public.der
