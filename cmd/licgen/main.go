package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golicense/server"
	"golicense/shared"
)

func main() {
	var (
		input    = flag.String("i", "", "输入的req.dat文件路径")
		output   = flag.String("o", "license.dat", "输出的license.dat文件路径")
		days     = flag.Int("d", 365, "授权有效期（天数）")
		customer = flag.String("c", "", "客户名称")
		org      = flag.String("org", "", "客户组织")
		edition  = flag.String("edition", "enterprise", "授权版本 (basic|enterprise)")
		help     = flag.Bool("h", false, "显示帮助信息")
	)
	flag.Parse()

	if *help {
		fmt.Println("licgen - 授权文件生成工具")
		fmt.Println()
		fmt.Println("用法:")
		fmt.Println("  licgen -i <req.dat> [选项]")
		fmt.Println()
		fmt.Println("选项:")
		fmt.Println("  -i string")
		fmt.Println("        输入的req.dat文件路径 (必需)")
		fmt.Println("  -o string")
		fmt.Println("        输出的license.dat文件路径 (默认 \"license.dat\")")
		fmt.Println("  -d int")
		fmt.Println("        授权有效期天数 (默认 365)")
		fmt.Println("  -c string")
		fmt.Println("        客户名称")
		fmt.Println("  -org string")
		fmt.Println("        客户组织")
		fmt.Println("  -edition string")
		fmt.Println("        授权版本 basic(基础版)|enterprise(旗舰版) (默认 \"enterprise\")")
		fmt.Println("  -h    显示帮助信息")
		fmt.Println()
		fmt.Println("授权版本说明:")
		fmt.Println("  basic      - 基础版: 仅包含准入管理功能")
		fmt.Println("  enterprise - 旗舰版: 包含全部功能(漏洞扫描、弱口令扫描、摄像头扫描)")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  licgen -i req.dat                                           # 生成1年期旗舰版license.dat")
		fmt.Println("  licgen -i req.dat -edition basic                            # 生成基础版授权")
		fmt.Println("  licgen -i req.dat -c \"张三\" -org \"ABC公司\"                # 指定客户信息")
		fmt.Println("  licgen -i req.dat -edition basic -d 30                      # 生成30天期限基础版")
		fmt.Println("  licgen -i req.dat -c \"李四\" -d 180 -o custom.dat           # 完整参数")
		return
	}

	if *input == "" {
		fmt.Println("错误: 必须指定输入文件 (-i)")
		fmt.Println("使用 -h 查看帮助信息")
		os.Exit(1)
	}

	// 检查输入文件是否存在
	if _, err := os.Stat(*input); os.IsNotExist(err) {
		log.Fatal("输入文件不存在:", *input)
	}

	// 确保输出目录存在
	outputDir := filepath.Dir(*output)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatal("创建输出目录失败:", err)
		}
	}

	// 验证天数参数
	if *days <= 0 {
		log.Fatal("授权天数必须大于0")
	}

	// 验证并解析版本参数
	var licenseEdition shared.LicenseEdition
	switch *edition {
	case "basic", "b":
		licenseEdition = shared.EditionBasic
	case "enterprise", "e":
		licenseEdition = shared.EditionEnterprise
	default:
		log.Fatal("无效的授权版本:", *edition, "。请使用 basic 或 enterprise")
	}

	// 准备客户信息
	customerInfo := server.CustomerInfo{
		Name:    *customer,
		Org:     *org,
		Edition: licenseEdition,
	}

	// 生成授权文件
	fmt.Printf("正在处理授权请求: %s\n", *input)
	if *customer != "" {
		fmt.Printf("客户信息: %s", *customer)
		if *org != "" {
			fmt.Printf(" (%s)", *org)
		}
		fmt.Println()
	}
	fmt.Printf("授权版本: %s", licenseEdition)
	switch licenseEdition {
	case shared.EditionBasic:
		fmt.Print(" (基础版 - 准入管理)")
	case shared.EditionEnterprise:
		fmt.Print(" (旗舰版 - 全功能)")
	}
	fmt.Println()
	fmt.Printf("授权有效期: %d 天\n", *days)

	if err := server.GenerateLicenseWithEdition(*input, *output, *days, customerInfo); err != nil {
		log.Fatal("生成授权文件失败:", err)
	}

	// 生成智能文件名
	smartOutput := generateSmartFilename(*output, *input, licenseEdition, *customer)

	// 如果智能文件名与原文件名不同，则重命名
	if smartOutput != *output {
		if err := os.Rename(*output, smartOutput); err != nil {
			log.Printf("重命名文件失败: %v，使用原文件名", err)
			smartOutput = *output
		}
	}

	fmt.Printf("\n✓ 授权文件已生成: %s\n", smartOutput)
	fmt.Println("请将此文件放置到客户端的goweb/bin/目录下")

	// 显示授权包含的模块
	modules := shared.GetModulesForEdition(licenseEdition)
	fmt.Printf("\n授权包含的模块:\n")
	for _, module := range modules {
		switch module {
		case shared.ModuleAdmission:
			fmt.Println("  ✓ 准入管理 - 设备发现、NAC控制、设备管理")
		case shared.ModuleVulnerabilityScan:
			fmt.Println("  ✓ 漏洞扫描 - 网络扫描、端口扫描、漏洞检测")
		case shared.ModulePasswordAudit:
			fmt.Println("  ✓ 弱口令扫描 - 弱口令检测、密码策略审计")
		case shared.ModuleCameraScan:
			fmt.Println("  ✓ 摄像头扫描 - 摄像头发现、ONVIF检测")
		}
	}
}

// generateSmartFilename 生成智能文件名
func generateSmartFilename(originalOutput, inputFile string, edition shared.LicenseEdition, customer string) string {
	// 如果用户明确指定了输出文件名（不是默认的license.dat），则保持用户指定的名称
	if originalOutput != "license.dat" {
		return originalOutput
	}

	// 从输入文件名提取硬件标识和日期
	inputBase := filepath.Base(inputFile)
	inputName := strings.TrimSuffix(inputBase, filepath.Ext(inputBase))

	// 尝试从req文件名中提取硬件标识（格式：req_[serialNumber]_[hwID]_[date].dat）
	parts := strings.Split(inputName, "_")
	var hwID, date string

	if len(parts) >= 4 && parts[0] == "req" {
		// 新格式：req_[serialNumber]_[hwID]_[date]
		hwID = parts[2]
		date = parts[3]
	} else if len(parts) >= 3 && parts[0] == "req" {
		// 旧格式：req_[hwID]_[date] 或 req_unknown_[hwID]_[date]
		if parts[1] == "unknown" && len(parts) >= 4 {
			hwID = parts[2]
			date = parts[3]
		} else {
			hwID = parts[1]
			date = parts[2]
		}
	} else {
		// 无法解析，使用当前时间
		hwID = "unknown"
		date = time.Now().Format("20060102")
	}

	// 生成版本前缀
	var editionPrefix string
	switch edition {
	case shared.EditionBasic:
		editionPrefix = "NSB"
	case shared.EditionEnterprise:
		editionPrefix = "NSE"
	default:
		editionPrefix = "NSC"
	}

	// 清理客户名称（移除特殊字符）
	cleanCustomer := strings.ReplaceAll(customer, " ", "")
	cleanCustomer = strings.ReplaceAll(cleanCustomer, "/", "")
	cleanCustomer = strings.ReplaceAll(cleanCustomer, "\\", "")
	if cleanCustomer == "" {
		cleanCustomer = "customer"
	}

	// 生成文件名：license_[版本]_[客户名]_[硬件ID]_[日期].dat
	filename := fmt.Sprintf("license_%s_%s_%s_%s.dat",
		editionPrefix, cleanCustomer, hwID, date)

	return filename
}