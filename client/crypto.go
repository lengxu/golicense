package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
)

// RSAEncrypt RSA加密
func RSAEncrypt(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, nil)
}

// RSAVerify RSA签名验证
func RSAVerify(data interface{}, signature []byte, publicKey *rsa.PublicKey) bool {
	// 序列化数据
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false
	}
	
	// 计算hash
	hashed := sha256.Sum256(jsonData)
	
	// 验证签名
	err = rsa.VerifyPKCS1v15(publicKey, 0, hashed[:], signature)
	return err == nil
}

// GenerateAESKey 生成AES密钥
func GenerateAESKey() []byte {
	key := make([]byte, 32) // AES-256
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}
	return key
}

// AESEncrypt AES加密
func AESEncrypt(data interface{}, key []byte) ([]byte, error) {
	// 序列化数据
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	
	// 创建AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	
	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	
	// 生成nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	
	// 加密
	ciphertext := gcm.Seal(nonce, nonce, jsonData, nil)
	return ciphertext, nil
}

// AESDecrypt AES解密
func AESDecrypt(ciphertext []byte, key []byte, result interface{}) error {
	// 创建AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	
	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return errors.New("ciphertext too short")
	}
	
	// 提取nonce和密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	
	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}
	
	// 反序列化
	return json.Unmarshal(plaintext, result)
}

// SHA256Hash 计算SHA256哈希
func SHA256Hash(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	
	hash := sha256.Sum256(jsonData)
	return hash[:], nil
}

// DeriveKeyFromHardware 从硬件指纹派生AES密钥
func DeriveKeyFromHardware(hardwareID string) []byte {
	// 使用硬件指纹和固定salt生成密钥
	data := hardwareID + "_license_key_salt_2024"
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}