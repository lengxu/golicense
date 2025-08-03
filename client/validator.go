package client

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"
)

// ValidateLicense 验证授权文件
func ValidateLicense(licenseFilePath string) error {
	// 1. 检查license.dat是否存在
	if _, err := os.Stat(licenseFilePath); os.IsNotExist(err) {
		return errors.New("license file not found")
	}

	// 2. 读取license.dat
	licenseData, err := os.ReadFile(licenseFilePath)
	if err != nil {
		return fmt.Errorf("failed to read license file: %v", err)
	}

	var licenseFile LicenseFile
	if err := DecodeFromString(string(licenseData), &licenseFile); err != nil {
		return fmt.Errorf("failed to decode license file: %v", err)
	}

	// 3. 获取当前硬件指纹
	currentHW := GetHardwareFingerprint()

	// 4. 用硬件指纹派生的密钥解密授权数据
	licenseKey := DeriveKeyFromHardware(currentHW)
	
	// 验证密钥hash
	keyHashArray := sha256.Sum256(licenseKey)
	expectedKeyHash := hex.EncodeToString(keyHashArray[:])
	if licenseFile.Key != expectedKeyHash {
		return errors.New("license key mismatch - hardware fingerprint changed")
	}

	// 5. 解密授权数据
	encryptedData, err := base64.StdEncoding.DecodeString(licenseFile.Data)
	if err != nil {
		return fmt.Errorf("failed to decode license data: %v", err)
	}

	var license License
	if err := AESDecrypt(encryptedData, licenseKey, &license); err != nil {
		return fmt.Errorf("failed to decrypt license data: %v", err)
	}

	// 6. 验证硬件指纹绑定
	if license.HardwareID != currentHW {
		return errors.New("hardware fingerprint mismatch")
	}

	// 7. 验证时间
	now := time.Now().Unix()
	if now < license.IssuedAt {
		return errors.New("license not yet valid")
	}
	if now > license.ExpiresAt {
		return fmt.Errorf("license expired on %s", time.Unix(license.ExpiresAt, 0).Format("2006-01-02 15:04:05"))
	}

	// 8. 验证RSA签名
	signature, err := base64.StdEncoding.DecodeString(licenseFile.Signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %v", err)
	}

	publicKey := GetEmbeddedPublicKey()
	if !RSAVerify(license, signature, publicKey) {
		return errors.New("invalid license signature")
	}

	return nil
}

// GetLicenseInfo 获取授权信息
func GetLicenseInfo(licenseFilePath string) (*License, error) {
	// 先验证授权
	if err := ValidateLicense(licenseFilePath); err != nil {
		return nil, err
	}

	// 读取并解密授权信息
	licenseData, err := os.ReadFile(licenseFilePath)
	if err != nil {
		return nil, err
	}

	var licenseFile LicenseFile
	if err := DecodeFromString(string(licenseData), &licenseFile); err != nil {
		return nil, err
	}

	currentHW := GetHardwareFingerprint()
	licenseKey := DeriveKeyFromHardware(currentHW)

	encryptedData, err := base64.StdEncoding.DecodeString(licenseFile.Data)
	if err != nil {
		return nil, err
	}

	var license License
	if err := AESDecrypt(encryptedData, licenseKey, &license); err != nil {
		return nil, err
	}

	return &license, nil
}

// CheckLicenseModule 检查模块授权
func CheckLicenseModule(licenseFilePath string, moduleName string) error {
	license, err := GetLicenseInfo(licenseFilePath)
	if err != nil {
		return err
	}

	// 检查模块是否在授权列表中
	for _, module := range license.Modules {
		if module == moduleName {
			return nil
		}
	}

	return fmt.Errorf("module '%s' not authorized", moduleName)
}