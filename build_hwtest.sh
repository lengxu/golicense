#!/bin/bash
# 编译硬件指纹测试工具 - ARM64版本

echo "开始编译硬件指纹测试工具..."

cd "$(dirname "$0")"

# 编译ARM64版本
echo "编译 ARM64 版本..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o hwtest_linux_arm64 ./cmd/hwtest

if [ $? -eq 0 ]; then
    echo "✓ ARM64版本编译成功: hwtest_linux_arm64"
    ls -lh hwtest_linux_arm64
else
    echo "❌ ARM64版本编译失败"
    exit 1
fi

echo ""
echo "使用方法:"
echo "1. 将 hwtest_linux_arm64 上传到 rk3568 机器"
echo "2. chmod +x hwtest_linux_arm64"
echo "3. ./hwtest_linux_arm64"
echo ""
echo "如果要测试授权文件，将 license.dat 放在同一目录下"
