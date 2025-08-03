package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lengxu/golicense/client"
)

func main() {
	var (
		license = flag.String("l", "license.dat", "license.dat文件路径")
		module  = flag.String("m", "", "检查特定模块授权")
		help    = flag.Bool("h", false, "显示帮助信息")
	)
	flag.Parse()

	if *help {
		fmt.Println("liccheck - 授权文件检查工具")
		fmt.Println()
		fmt.Println("用法:")
		fmt.Println("  liccheck [选项]")
		fmt.Println()
		fmt.Println("选项:")
		fmt.Println("  -l string")
		fmt.Println("        license.dat文件路径 (默认 \"license.dat\")")
		fmt.Println("  -m string")
		fmt.Println("        检查特定模块授权 (如: goscan, gopasswd, goweb)")
		fmt.Println("  -h    显示帮助信息")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  liccheck                          # 检查默认授权文件")
		fmt.Println("  liccheck -l goweb/bin/license.dat # 检查指定授权文件")
		fmt.Println("  liccheck -m goscan                # 检查goscan模块授权")
		return
	}

	// 检查授权文件是否存在
	if _, err := os.Stat(*license); os.IsNotExist(err) {
		log.Fatal("授权文件不存在:", *license)
	}

	fmt.Printf("正在检查授权文件: %s\n", *license)
	fmt.Printf("当前硬件指纹: %s\n", client.GetHardwareFingerprint())
	fmt.Println()

	// 如果指定了模块，只检查模块授权
	if *module != "" {
		if err := client.CheckLicenseModule(*license, *module); err != nil {
			log.Fatal("模块授权检查失败:", err)
		}
		fmt.Printf("✓ 模块 '%s' 授权有效\n", *module)
		return
	}

	// 验证授权
	if err := client.ValidateLicense(*license); err != nil {
		log.Fatal("授权验证失败:", err)
	}

	// 获取授权详细信息
	licenseInfo, err := client.GetLicenseInfo(*license)
	if err != nil {
		log.Fatal("获取授权信息失败:", err)
	}

	// 显示授权信息
	fmt.Println("✓ 授权验证成功")
	fmt.Println()
	fmt.Println("授权信息:")
	fmt.Printf("  客户ID: %s\n", licenseInfo.CustomerID)
	fmt.Printf("  签发时间: %s\n", time.Unix(licenseInfo.IssuedAt, 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("  过期时间: %s\n", time.Unix(licenseInfo.ExpiresAt, 0).Format("2006-01-02 15:04:05"))
	
	// 计算剩余天数
	remainingDays := int((licenseInfo.ExpiresAt - time.Now().Unix()) / 86400)
	if remainingDays > 0 {
		fmt.Printf("  剩余天数: %d 天\n", remainingDays)
	}
	
	fmt.Printf("  最大扫描次数: %d\n", licenseInfo.MaxScans)
	fmt.Printf("  授权模块: %v\n", licenseInfo.Modules)
	fmt.Printf("  授权功能: %v\n", licenseInfo.Features)
}