package client

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AutoLicenseCheck 自动授权检查和req.dat生成（用于goweb主控平台）
// 这个函数会在找不到license.dat时自动生成req.dat
func AutoLicenseCheck(module string) error {
	// 获取可执行文件所在目录
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)
	
	licensePath := filepath.Join(exeDir, "license.dat")
	reqPath := filepath.Join(exeDir, "req.dat")

	// 1. 检查license.dat是否存在
	if _, err := os.Stat(licensePath); os.IsNotExist(err) {
		fmt.Printf("⚠️  未找到授权文件: %s\n", licensePath)
		return handleMissingLicense(reqPath)
	}

	// 2. 验证license.dat
	if err := ValidateLicense(licensePath); err != nil {
		fmt.Printf("⚠️  授权验证失败: %v\n", err)
		return handleInvalidLicense(reqPath, err)
	}

	// 3. 检查模块授权
	if module != "" {
		if err := CheckLicenseModule(licensePath, module); err != nil {
			return fmt.Errorf("模块授权检查失败: %v", err)
		}
	}

	// 4. 显示授权信息
	displayLicenseStatus(licensePath)
	return nil
}

// handleMissingLicense 处理缺失授权文件的情况
func handleMissingLicense(reqPath string) error {
	// 检查是否已存在req.dat
	if _, err := os.Stat(reqPath); err == nil {
		fmt.Printf("✓ 授权请求文件已存在: %s\n", reqPath)
		fmt.Println("📋 请将此文件发送给授权服务端获取license.dat")
		return fmt.Errorf("等待授权：请联系管理员获取授权文件")
	}

	fmt.Println("🔄 正在生成授权请求文件...")
	
	// 生成新的req.dat
	if err := GenerateRequest(reqPath); err != nil {
		return fmt.Errorf("生成授权请求失败: %v", err)
	}

	fmt.Println("\n📋 授权请求步骤:")
	fmt.Printf("1. 将 %s 发送给授权服务端\n", reqPath)
	fmt.Println("2. 等待获取 license.dat 文件")
	fmt.Printf("3. 将 license.dat 放入 %s 目录\n", filepath.Dir(reqPath))
	fmt.Println("4. 重新启动程序")

	return fmt.Errorf("等待授权：请按照上述步骤获取授权")
}

// handleInvalidLicense 处理无效授权文件的情况
func handleInvalidLicense(reqPath string, validationErr error) error {
	fmt.Println("🔄 授权文件无效，正在重新生成授权请求...")
	
	// 备份旧的req.dat（如果存在且还没有备份）
	if _, err := os.Stat(reqPath); err == nil {
		// 检查是否已存在备份文件，避免重复备份
		backupPattern := reqPath + ".backup.*"
		matches, _ := filepath.Glob(backupPattern)
		if len(matches) == 0 {
			backupPath := reqPath + ".backup." + fmt.Sprintf("%d", time.Now().Unix())
			os.Rename(reqPath, backupPath)
			fmt.Printf("📦 已备份旧请求文件: %s\n", backupPath)
		} else {
			// 如果已有备份，直接删除当前req.dat
			os.Remove(reqPath)
			fmt.Printf("🗑️  删除旧请求文件 (已存在备份)\n")
		}
	}

	// 生成新的req.dat
	if err := GenerateRequest(reqPath); err != nil {
		return fmt.Errorf("生成授权请求失败: %v", err)
	}

	fmt.Println("\n⚠️  授权失效原因:", validationErr.Error())
	fmt.Println("📋 重新授权步骤:")
	fmt.Printf("1. 将新的 %s 发送给授权服务端\n", reqPath)
	fmt.Println("2. 等待获取新的 license.dat 文件")
	fmt.Printf("3. 替换 bin/license.dat 文件\n")
	fmt.Println("4. 重新启动程序")

	return fmt.Errorf("授权失效：%v", validationErr)
}

// displayLicenseStatus 显示授权状态信息
func displayLicenseStatus(licensePath string) {
	license, err := GetLicenseInfo(licensePath)
	if err != nil {
		return
	}

	fmt.Println("✅ 授权验证成功")
	
	// 计算剩余天数
	remainingDays := int((license.ExpiresAt - time.Now().Unix()) / 86400)
	
	if remainingDays <= 7 {
		fmt.Printf("⚠️  授权即将过期！剩余 %d 天\n", remainingDays)
	} else if remainingDays <= 30 {
		fmt.Printf("🔔 授权剩余 %d 天\n", remainingDays)
	}
	
	if license.CustomerName != "" {
		fmt.Printf("📋 授权用户: %s", license.CustomerName)
		if license.CustomerOrg != "" {
			fmt.Printf(" (%s)", license.CustomerOrg)
		}
		fmt.Println()
	}
}

// ValidateOnlyLicense 仅校验授权（用于goscan/gopasswd扫描工具）
// 这个函数只做授权验证，不会生成req.dat文件
func ValidateOnlyLicense(module string) error {
	// 获取可执行文件所在目录
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)
	
	licensePath := filepath.Join(exeDir, "license.dat")
	
	// 1. 检查license.dat是否存在
	if _, err := os.Stat(licensePath); os.IsNotExist(err) {
		return fmt.Errorf("未找到授权文件 %s，请先通过goweb平台获取授权", licensePath)
	}
	
	// 2. 验证license.dat
	if err := ValidateLicense(licensePath); err != nil {
		return fmt.Errorf("授权验证失败: %v，请重新获取授权", err)
	}
	
	// 3. 检查模块授权
	if module != "" {
		if err := CheckLicenseModule(licensePath, module); err != nil {
			return fmt.Errorf("模块授权检查失败: %v", err)
		}
	}
	
	// 4. 显示授权信息（简化版）
	displayLicenseStatusSimple(licensePath)
	return nil
}

// displayLicenseStatusSimple 显示简化的授权状态信息
func displayLicenseStatusSimple(licensePath string) {
	license, err := GetLicenseInfo(licensePath)
	if err != nil {
		return
	}
	
	// 计算剩余天数
	remainingDays := int((license.ExpiresAt - time.Now().Unix()) / 86400)
	
	if remainingDays <= 7 {
		fmt.Printf("⚠️  授权即将过期！剩余 %d 天\n", remainingDays)
	}
}

// QuickLicenseCheck 快速授权检查（仅验证，不生成文件）
func QuickLicenseCheck(module string) error {
	// 获取可执行文件所在目录
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)
	
	licensePath := filepath.Join(exeDir, "license.dat")
	
	// 检查文件是否存在
	if _, err := os.Stat(licensePath); os.IsNotExist(err) {
		return fmt.Errorf("授权文件不存在，请运行完整的授权检查")
	}

	// 验证授权
	if err := ValidateLicense(licensePath); err != nil {
		return err
	}

	// 检查模块授权
	if module != "" {
		return CheckLicenseModule(licensePath, module)
	}

	return nil
}