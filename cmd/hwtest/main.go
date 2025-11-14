package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lengxu/golicense/client"
)

func main() {
	fmt.Println("=== 硬件指纹测试工具 ===")
	fmt.Println()

	// 1. 获取当前硬件指纹
	hwID := client.GetHardwareFingerprint()
	fmt.Printf("✓ 当前硬件指纹:\n  %s\n\n", hwID)

	// 2. 获取机器信息
	machineInfo := client.GetMachineInfo()
	fmt.Printf("✓ 机器信息:\n  %s\n\n", machineInfo)

	// 3. 显示密钥派生结果
	key := client.DeriveKeyFromHardware(hwID)
	fmt.Printf("✓ 派生密钥 (完整32字节):\n  %s\n\n", hex.EncodeToString(key))

	// 4. 检查license.dat文件
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	licensePath := filepath.Join(exeDir, "license.dat")

	fmt.Printf("✓ 授权文件路径:\n  %s\n\n", licensePath)

	if _, err := os.Stat(licensePath); os.IsNotExist(err) {
		fmt.Println("❌ license.dat 不存在")
		return
	}

	// 5. 验证授权文件
	fmt.Println("=== 开始验证授权文件 ===")
	err := client.ValidateLicense(licensePath)
	if err != nil {
		fmt.Printf("❌ 授权验证失败:\n  %v\n\n", err)

		// 尝试读取授权文件中的硬件指纹
		license, _ := client.GetLicenseInfo(licensePath)
		if license != nil {
			fmt.Printf("⚠️  授权文件中的硬件指纹:\n  %s\n\n", license.HardwareID)
			if license.HardwareID != hwID {
				fmt.Println("⚠️  硬件指纹不匹配!")
				fmt.Printf("  当前机器: %s\n", hwID)
				fmt.Printf("  授权文件: %s\n", license.HardwareID)
			}
		}
	} else {
		fmt.Println("✅ 授权验证成功!")

		// 6. 显示授权详细信息
		license, _ := client.GetLicenseInfo(licensePath)
		if license != nil {
			fmt.Println("\n=== 授权详细信息 ===")
			fmt.Printf("  客户名称: %s\n", license.CustomerName)
			fmt.Printf("  客户组织: %s\n", license.CustomerOrg)
			fmt.Printf("  授权版本: %s\n", license.Edition)
			fmt.Printf("  序列号: %s\n", license.SerialNumber)
			fmt.Printf("  硬件指纹: %s\n", license.HardwareID)
			fmt.Printf("  签发时间: %d\n", license.IssuedAt)
			fmt.Printf("  过期时间: %d\n", license.ExpiresAt)
			fmt.Printf("  授权模块: %v\n", license.Modules)
			if len(license.ModulePerms) > 0 {
				fmt.Println("  模块权限:")
				for _, perm := range license.ModulePerms {
					fmt.Printf("    - %s: enabled=%v\n", perm.Module, perm.Enabled)
				}
			}
		}
	}
}
