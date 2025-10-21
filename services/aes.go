package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"time"
)

const KEY = "3HXV4P8dizvATG5EjLIsUKxSreyghDMB" // 32字节密钥

// AesService AES加密服务
type AesService struct{}

// NewAesService 创建AES服务实例
func NewAesService() *AesService {
	return &AesService{}
}

// Encrypt AES加密
func (s *AesService) Encrypt(data string) (string, error) {
	// 生成随机IV
	iv := make([]byte, 16)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	// 获取当前时间戳
	timestamp := time.Now().Unix()

	// 创建AES加密器
	block, err := aes.NewCipher([]byte(KEY))
	if err != nil {
		return "", err
	}

	// PKCS7填充
	plaintext := []byte(data)
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := make([]byte, len(plaintext)+padding)
	copy(padtext, plaintext)
	for i := len(plaintext); i < len(padtext); i++ {
		padtext[i] = byte(padding)
	}

	// CBC模式加密
	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(padtext))
	mode.CryptBlocks(ciphertext, padtext)

	// 组合时间戳 + IV + 密文
	result := make([]byte, 8+16+len(ciphertext))
	binary.BigEndian.PutUint64(result[0:8], uint64(timestamp))
	copy(result[8:24], iv)
	copy(result[24:], ciphertext)

	// Base64编码
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt AES解密 (如果需要的话)
func (s *AesService) Decrypt(encryptedData string) (string, error) {
	// Base64解码
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	if len(data) < 24 {
		return "", fmt.Errorf("invalid encrypted data length")
	}

	// 提取时间戳、IV和密文
	// timestamp := binary.BigEndian.Uint64(data[0:8])
	iv := data[8:24]
	ciphertext := data[24:]

	// 创建AES解密器
	block, err := aes.NewCipher([]byte(KEY))
	if err != nil {
		return "", err
	}

	// CBC模式解密
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除PKCS7填充
	padding := int(plaintext[len(plaintext)-1])
	if padding > aes.BlockSize || padding == 0 {
		return "", fmt.Errorf("invalid padding")
	}

	for i := len(plaintext) - padding; i < len(plaintext); i++ {
		if plaintext[i] != byte(padding) {
			return "", fmt.Errorf("invalid padding")
		}
	}

	return string(plaintext[:len(plaintext)-padding]), nil
}
