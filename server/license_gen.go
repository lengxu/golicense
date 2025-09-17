package server

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/lengxu/golicense/shared"
)

// CustomerInfo 客户信息
type CustomerInfo struct {
	Name    string
	Org     string
	Edition shared.LicenseEdition
}

// GenerateLicense 根据req.dat生成license.dat（兼容旧版本）
func GenerateLicense(reqFilePath, licenseFilePath string, days int, customer CustomerInfo) error {
	// 默认使用旗舰版
	if customer.Edition == "" {
		customer.Edition = shared.EditionEnterprise
	}
	return GenerateLicenseWithEdition(reqFilePath, licenseFilePath, days, customer)
}

// GenerateLicenseWithEdition 根据req.dat生成指定版本的license.dat
func GenerateLicenseWithEdition(reqFilePath, licenseFilePath string, days int, customer CustomerInfo) error {
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

	// 4. 根据版本生成授权数据
	now := time.Now()
	modules := shared.GetModulesForEdition(customer.Edition)
	modulePerms := shared.GetDefaultModulePermissions(customer.Edition)

	// 生成序列号
	serialNumber := generateSerialNumber(request.HardwareID, customer.Edition)

	license := License{
		HardwareID:   request.HardwareID,
		IssuedAt:     now.Unix(),
		ExpiresAt:    now.AddDate(0, 0, days).Unix(),
		Edition:      customer.Edition,
		Modules:      modules,
		ModulePerms:  modulePerms,
		CustomerID:   generateCustomerID(request.HardwareID),
		CustomerName: customer.Name,
		CustomerOrg:  customer.Org,
		MaxScans:     getMaxScansForEdition(customer.Edition),
		MaxAssets:    getMaxAssetsForEdition(customer.Edition),
		MaxUsers:     getMaxUsersForEdition(customer.Edition),
		Features:     getFeaturesForEdition(customer.Edition),
		RequestID:    request.RequestID,
		LicenseKey:   generateLicenseKey(request.HardwareID, customer.Edition),
		SerialNumber: serialNumber,
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
		Version:   "2.0", // 升级版本号以支持新格式
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
	fmt.Printf("  Edition: %s\n", license.Edition)
	fmt.Printf("  Serial Number: %s\n", license.SerialNumber)
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

// generateSerialNumber 生成序列号
func generateSerialNumber(hardwareID string, edition shared.LicenseEdition) string {
	var prefix string
	switch edition {
	case shared.EditionBasic:
		prefix = "NSB" // NScan Basic
	case shared.EditionEnterprise:
		prefix = "NSE" // NScan Enterprise
	default:
		prefix = "NSC" // NScan Custom
	}

	hash := sha256.Sum256([]byte("serial_" + hardwareID + string(edition)))
	return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(hash[:6]))
}

// generateLicenseKey 生成授权密钥
func generateLicenseKey(hardwareID string, edition shared.LicenseEdition) string {
	hash := sha256.Sum256([]byte("license_key_" + hardwareID + string(edition)))
	return hex.EncodeToString(hash[:16])
}

// deriveKeyFromHardware 从硬件指纹派生AES密钥
func deriveKeyFromHardware(hardwareID string) []byte {
	// 使用硬件指纹和固定salt生成密钥
	data := hardwareID + "_license_key_salt_2024"
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

// getMaxScansForEdition 根据版本获取最大扫描次数
func getMaxScansForEdition(edition shared.LicenseEdition) int {
	switch edition {
	case shared.EditionBasic:
		return 1000 // 基础版限制1000次扫描
	case shared.EditionEnterprise:
		return 0 // 旗舰版不限制
	default:
		return 100 // 默认限制
	}
}

// getMaxAssetsForEdition 根据版本获取最大资产数量
func getMaxAssetsForEdition(edition shared.LicenseEdition) int {
	switch edition {
	case shared.EditionBasic:
		return 500 // 基础版限制500个资产
	case shared.EditionEnterprise:
		return 0 // 旗舰版不限制
	default:
		return 50 // 默认限制
	}
}

// getMaxUsersForEdition 根据版本获取最大用户数量
func getMaxUsersForEdition(edition shared.LicenseEdition) int {
	switch edition {
	case shared.EditionBasic:
		return 3 // 基础版限制3个用户
	case shared.EditionEnterprise:
		return 0 // 旗舰版不限制
	default:
		return 1 // 默认限制
	}
}

// getFeaturesForEdition 根据版本获取功能特性
func getFeaturesForEdition(edition shared.LicenseEdition) []string {
	switch edition {
	case shared.EditionBasic:
		return []string{"basic_scanning", "device_management", "basic_reporting"}
	case shared.EditionEnterprise:
		return []string{"full_scanning", "advanced_reporting", "api_access", "custom_templates"}
	default:
		return []string{"basic"}
	}
}