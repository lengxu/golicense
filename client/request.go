package client

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

// GenerateRequest 生成授权请求文件req.dat
func GenerateRequest(reqFilePath string) error {
	// 1. 获取硬件指纹
	hardwareID := GetHardwareFingerprint()
	
	// 2. 生成请求ID
	requestID := generateRequestID()
	
	// 3. 构造请求数据
	request := LicenseRequest{
		HardwareID:  hardwareID,
		Timestamp:   time.Now().Unix(),
		Version:     "1.0.0",
		MachineInfo: GetMachineInfo(),
		RequestID:   requestID,
	}

	// 4. 计算请求数据hash
	requestHash, err := SHA256Hash(request)
	if err != nil {
		return fmt.Errorf("failed to calculate request hash: %v", err)
	}

	// 5. 生成AES密钥
	aesKey := GenerateAESKey()

	// 6. AES加密请求数据
	encryptedData, err := AESEncrypt(request, aesKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt request data: %v", err)
	}

	// 7. 用内置公钥加密AES密钥
	publicKey := GetEmbeddedPublicKey()
	encryptedKey, err := RSAEncrypt(aesKey, publicKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt AES key: %v", err)
	}

	// 8. 构造请求文件
	reqFile := RequestFile{
		Data:      base64.StdEncoding.EncodeToString(encryptedData),
		Key:       base64.StdEncoding.EncodeToString(encryptedKey),
		Hash:      hex.EncodeToString(requestHash),
		Timestamp: time.Now().Unix(),
	}

	// 9. 编码为字符串并保存req.dat
	encodedString, err := EncodeToString(reqFile)
	if err != nil {
		return fmt.Errorf("failed to encode request file: %v", err)
	}

	if err := os.WriteFile(reqFilePath, []byte(encodedString), 0644); err != nil {
		return fmt.Errorf("failed to write request file: %v", err)
	}

	fmt.Printf("Request file generated successfully:\n")
	fmt.Printf("  Request ID: %s\n", requestID)
	fmt.Printf("  Hardware ID: %s\n", hardwareID)
	fmt.Printf("  Machine Info: %s\n", request.MachineInfo)
	fmt.Printf("  File: %s\n", reqFilePath)

	return nil
}

// generateRequestID 生成请求唯一标识
func generateRequestID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}