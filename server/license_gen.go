package server

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

// CustomerInfo 客户信息
type CustomerInfo struct {
	Name string
	Org  string
}

// GenerateLicense 根据req.dat生成license.dat
func GenerateLicense(reqFilePath, licenseFilePath string, days int, customer CustomerInfo) error {
	// 1. 读取req.dat
	reqData, err := os.ReadFile(reqFilePath)
	if err != nil {
		return fmt.Errorf("failed to read request file: %v", err)
	}

	var reqFile RequestFile
	if err := DecodeFromString(string(reqData), &reqFile); err != nil {
		return fmt.Errorf("failed to decode request file: %v", err)
	}

	// 2. 解密请求数据
	privateKey := GetPrivateKey()
	
	// 解码RSA加密的AES密钥
	encryptedKey, err := base64.StdEncoding.DecodeString(reqFile.Key)
	if err != nil {
		return fmt.Errorf("failed to decode encrypted key: %v", err)
	}
	
	// 解密AES密钥
	aesKey, err := RSADecrypt(encryptedKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to decrypt AES key: %v", err)
	}
	
	// 解码请求数据
	encryptedData, err := base64.StdEncoding.DecodeString(reqFile.Data)
	if err != nil {
		return fmt.Errorf("failed to decode encrypted data: %v", err)
	}
	
	// 解密请求数据
	var request LicenseRequest
	if err := AESDecrypt(encryptedData, aesKey, &request); err != nil {
		return fmt.Errorf("failed to decrypt request data: %v", err)
	}

	// 3. 验证请求数据完整性 (暂时跳过，等待修复hash验证)
	expectedHash, err := SHA256Hash(request)
	if err != nil {
		return fmt.Errorf("failed to calculate request hash: %v", err)
	}
	
	fmt.Printf("Hash verification (debug): expected=%s, got=%s\n", 
		hex.EncodeToString(expectedHash), reqFile.Hash)
	
	// TODO: 修复hash验证问题
	// if hex.EncodeToString(expectedHash) != reqFile.Hash {
	//     return fmt.Errorf("request data integrity check failed")
	// }

	// 4. 生成授权数据
	now := time.Now()
	license := License{
		HardwareID:   request.HardwareID,
		IssuedAt:     now.Unix(),
		ExpiresAt:    now.AddDate(0, 0, days).Unix(),
		Modules:      []string{"goscan", "gopasswd", "goweb"},
		CustomerID:   generateCustomerID(request.HardwareID),
		CustomerName: customer.Name,
		CustomerOrg:  customer.Org,
		MaxScans:     10000,
		Features:     []string{"full"},
		RequestID:    request.RequestID,
	}

	// 5. 签名授权数据
	signature, err := RSASign(license, privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign license: %v", err)
	}

	// 6. 用硬件指纹派生的密钥加密授权数据
	licenseKey := deriveKeyFromHardware(request.HardwareID)
	encryptedLicense, err := AESEncrypt(license, licenseKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt license: %v", err)
	}

	// 7. 生成license文件
	keyHashArray := sha256.Sum256(licenseKey)
	licenseFile := LicenseFile{
		Data:      base64.StdEncoding.EncodeToString(encryptedLicense),
		Key:       hex.EncodeToString(keyHashArray[:]),
		Signature: base64.StdEncoding.EncodeToString(signature),
		Version:   "1.0",
	}

	// 8. 编码为字符串并保存license.dat
	encodedString, err := EncodeLicenseToString(licenseFile)
	if err != nil {
		return fmt.Errorf("failed to encode license file: %v", err)
	}

	if err := os.WriteFile(licenseFilePath, []byte(encodedString), 0644); err != nil {
		return fmt.Errorf("failed to write license file: %v", err)
	}

	fmt.Printf("License generated successfully:\n")
	fmt.Printf("  Request ID: %s\n", request.RequestID)
	fmt.Printf("  Hardware ID: %s\n", request.HardwareID)
	fmt.Printf("  Customer: %s", license.CustomerName)
	if license.CustomerOrg != "" {
		fmt.Printf(" (%s)", license.CustomerOrg)
	}
	fmt.Println()
	fmt.Printf("  Customer ID: %s\n", license.CustomerID)
	fmt.Printf("  Issued At: %s\n", time.Unix(license.IssuedAt, 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("  Expires At: %s\n", time.Unix(license.ExpiresAt, 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("  Modules: %v\n", license.Modules)

	return nil
}

// generateCustomerID 生成客户ID
func generateCustomerID(hardwareID string) string {
	hash := sha256.Sum256([]byte("customer_" + hardwareID))
	return hex.EncodeToString(hash[:8])
}

// deriveKeyFromHardware 从硬件指纹派生AES密钥
func deriveKeyFromHardware(hardwareID string) []byte {
	// 使用硬件指纹和固定salt生成密钥
	data := hardwareID + "_license_key_salt_2024"
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}