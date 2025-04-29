import CryptoJS from 'crypto-js'

export const setProperty = (prop: string, val: any, dom = document.documentElement) => {
    dom.style.setProperty(prop, val);
};

export const mix = (color1: string, color2: string, weight: number = 0.5): string => {
    let color = '#';
    for (let i = 0; i <= 2; i++) {
        const c1 = parseInt(color1.substring(1 + i * 2, 3 + i * 2), 16);
        const c2 = parseInt(color2.substring(1 + i * 2, 3 + i * 2), 16);
        const c = Math.round(c1 * weight + c2 * (1 - weight));
        color += c.toString(16).padStart(2, '0');
    }
    return color;
};


// 解密
export const decryptAes = (data:string ) => {
    // 分割符号
    const splitStr = '/+/'
    // 分割为数组
    const splitArr = data.split(splitStr)
    const content = splitArr[0] || ''
    const key = splitArr[1] || ''
    const iv = splitArr[2] || ''

    const parsedKey = CryptoJS.enc.Hex.parse(key);
    const parsedIv = CryptoJS.enc.Hex.parse(iv);

    // Base64 解码密文
    const encryptedData = CryptoJS.enc.Base64.parse(content);

    // 解密数据
    const decrypted = CryptoJS.AES.decrypt(
        { ciphertext: encryptedData },
        parsedKey,
        {
            iv: parsedIv,
            mode: CryptoJS.mode.CBC,
            padding: CryptoJS.pad.Pkcs7,
        },
    );

    const decryptedText = decrypted.toString(CryptoJS.enc.Utf8);
    return JSON.parse(decryptedText) || {}
}

// 加密
export const encryptAes = (data: string) => {
    // 生成随机密钥和初始向量
    const key = CryptoJS.lib.WordArray.random(16);
    const iv = CryptoJS.lib.WordArray.random(16);

    // 将数据转换为字符串
    const dataStr = JSON.stringify(data);

    // 加密数据
    const encrypted = CryptoJS.AES.encrypt(dataStr, key, {
        iv: iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7,
    });

    // 将密钥、初始向量和加密后的数据拼接成一个字符串
    const encryptedData = `${encrypted.toString()}/+/${key.toString()}/+/${iv.toString()}`;
    return encryptedData;
}
