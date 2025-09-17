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

// CheckLicenseModule 检查模块授权（兼容旧版本）
func CheckLicenseModule(licensePath string, module string) error {
	license, err := GetLicenseInfo(licensePath)
	if err != nil {
		return err
	}

	// 新版本授权：检查详细模块权限
	if len(license.ModulePerms) > 0 {
		for _, perm := range license.ModulePerms {
			if string(perm.Module) == module && perm.Enabled {
				return nil
			}
		}
		return fmt.Errorf("模块 '%s' 未授权或已禁用", module)
	}

	// 旧版本授权：检查简单模块列表
	for _, allowedModule := range license.Modules {
		if string(allowedModule) == module {
			return nil
		}
	}

	return fmt.Errorf("模块 '%s' 未授权", module)
}

// CheckModulePermission 检查模块的详细权限
func CheckModulePermission(licensePath string, module LicenseModule) (*ModulePermissions, error) {
	license, err := GetLicenseInfo(licensePath)
	if err != nil {
		return nil, err
	}

	// 查找对应模块的权限配置
	for _, perm := range license.ModulePerms {
		if perm.Module == module {
			if !perm.Enabled {
				return nil, fmt.Errorf("模块 '%s' 已被禁用", module)
			}
			return &perm, nil
		}
	}

	return nil, fmt.Errorf("模块 '%s' 未找到授权配置", module)
}

// IsModuleEnabled 检查模块是否启用
func IsModuleEnabled(licensePath string, module LicenseModule) bool {
	_, err := CheckModulePermission(licensePath, module)
	return err == nil
}

// GetLicenseEdition 获取授权版本
func GetLicenseEdition(licensePath string) (LicenseEdition, error) {
	license, err := GetLicenseInfo(licensePath)
	if err != nil {
		return "", err
	}

	// 新版本直接返回版本信息
	if license.Edition != "" {
		return license.Edition, nil
	}

	// 旧版本通过模块数量推断版本
	moduleCount := len(license.Modules)
	if moduleCount <= 1 {
		return EditionBasic, nil
	}
	return EditionEnterprise, nil
}

// ValidateOnlyLicense 仅验证授权（不生成req.dat） - 兼容旧接口
func ValidateOnlyLicense(appName string) error {
	// 查找license.dat文件
	licensePath := "license.dat"

	// 尝试多个可能的路径
	possiblePaths := []string{
		"license.dat",
		"./license.dat",
		"../license.dat",
		"bin/license.dat",
		"./bin/license.dat",
	}

	var foundPath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			foundPath = path
			break
		}
	}

	if foundPath == "" {
		return errors.New("license.dat file not found")
	}

	// 验证授权
	if err := ValidateLicense(foundPath); err != nil {
		return fmt.Errorf("授权验证失败: %v", err)
	}

	// 检查应用模块权限
	return CheckAppModulePermission(foundPath, appName)
}

// CheckAppModulePermission 检查应用的模块权限
func CheckAppModulePermission(licensePath string, appName string) error {
	// 根据应用名称映射到模块
	var moduleToCheck LicenseModule
	switch appName {
	case "goscan":
		moduleToCheck = ModuleVulnerabilityScan
	case "gopasswd":
		moduleToCheck = ModulePasswordAudit
	case "goonvif", "onvif":
		moduleToCheck = ModuleCameraScan
	case "goweb":
		// goweb需要准入管理权限
		moduleToCheck = ModuleAdmission
	default:
		// 默认检查准入管理权限
		moduleToCheck = ModuleAdmission
	}

	// 检查模块权限
	_, err := CheckModulePermission(licensePath, moduleToCheck)
	return err
}

// GetAvailableModules 获取可用的模块列表
func GetAvailableModules(licensePath string) ([]LicenseModule, error) {
	license, err := GetLicenseInfo(licensePath)
	if err != nil {
		return nil, err
	}

	var availableModules []LicenseModule

	// 新版本：从ModulePerms获取启用的模块
	if len(license.ModulePerms) > 0 {
		for _, perm := range license.ModulePerms {
			if perm.Enabled {
				availableModules = append(availableModules, perm.Module)
			}
		}
		return availableModules, nil
	}

	// 旧版本：返回所有模块
	return license.Modules, nil
}