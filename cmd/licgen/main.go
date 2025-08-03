package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/lengxu/golicense/server"
)

func main() {
	var (
		input    = flag.String("i", "", "输入的req.dat文件路径")
		output   = flag.String("o", "license.dat", "输出的license.dat文件路径")
		days     = flag.Int("d", 365, "授权有效期（天数）")
		customer = flag.String("c", "", "客户名称")
		org      = flag.String("org", "", "客户组织")
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
		fmt.Println("  -h    显示帮助信息")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  licgen -i req.dat                                    # 生成1年期license.dat")
		fmt.Println("  licgen -i req.dat -c \"张三\" -org \"ABC公司\"         # 指定客户信息")
		fmt.Println("  licgen -i req.dat -d 30                              # 生成30天期限授权")
		fmt.Println("  licgen -i req.dat -c \"李四\" -d 180 -o custom.dat    # 完整参数")
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

	// 准备客户信息
	customerInfo := server.CustomerInfo{
		Name: *customer,
		Org:  *org,
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
	fmt.Printf("授权有效期: %d 天\n", *days)
	
	if err := server.GenerateLicense(*input, *output, *days, customerInfo); err != nil {
		log.Fatal("生成授权文件失败:", err)
	}

	fmt.Printf("\n✓ 授权文件已生成: %s\n", *output)
	fmt.Println("请将此文件放置到客户端的goweb/bin/目录下")
}