package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golicense/client"
)

func main() {
	var (
		output = flag.String("o", "req.dat", "输出文件路径")
		help   = flag.Bool("h", false, "显示帮助信息")
	)
	flag.Parse()

	if *help {
		fmt.Println("reqgen - 授权请求生成工具")
		fmt.Println()
		fmt.Println("用法:")
		fmt.Println("  reqgen [选项]")
		fmt.Println()
		fmt.Println("选项:")
		fmt.Println("  -o string")
		fmt.Println("        输出文件路径 (默认 \"req.dat\")")
		fmt.Println("  -h    显示帮助信息")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  reqgen                    # 生成 req.dat")
		fmt.Println("  reqgen -o request.dat     # 生成 request.dat")
		return
	}

	// 确保输出目录存在
	outputDir := filepath.Dir(*output)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Fatal("创建输出目录失败:", err)
		}
	}

	// 生成授权请求
	fmt.Println("正在获取硬件指纹...")
	fmt.Printf("硬件指纹: %s\n", client.GetHardwareFingerprint())
	fmt.Printf("机器信息: %s\n", client.GetMachineInfo())
	
	fmt.Println("正在生成授权请求文件...")
	if err := client.GenerateRequest(*output); err != nil {
		log.Fatal("生成请求文件失败:", err)
	}

	fmt.Printf("\n✓ 授权请求文件已生成: %s\n", *output)
	fmt.Println("请将此文件发送给授权服务端以获取license.dat")
}