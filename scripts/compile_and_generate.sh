#!/bin/bash

# 以太坊智能合约编译和Go绑定代码生成脚本
# 使用前请确保已安装 solc 和 abigen

echo "=== 以太坊智能合约编译和Go绑定代码生成 ==="

# 检查依赖
if ! command -v solc &> /dev/null; then
    echo "❌ 错误: 未找到 solc 编译器"
    echo "请安装: npm install -g solc"
    exit 1
fi

if ! command -v abigen &> /dev/null; then
    echo "❌ 错误: 未找到 abigen 工具"
    echo "请安装: go install github.com/ethereum/go-ethereum/cmd/abigen@latest"
    exit 1
fi

# 创建输出目录
mkdir -p contracts/compiled

# 编译智能合约
echo "📦 编译智能合约..."
solc --bin --abi --overwrite -o contracts/compiled contracts/SimpleCounter.sol

if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi

echo "✅ 编译成功"

# 生成Go绑定代码
echo "🔧 生成Go绑定代码..."
abigen --bin=contracts/compiled/SimpleCounter.bin --abi=contracts/compiled/SimpleCounter.abi --pkg=counter --out=pkg/eth/counter/counter.go

if [ $? -ne 0 ]; then
    echo "❌ Go绑定代码生成失败"
    exit 1
fi

echo "✅ Go绑定代码生成成功"
echo "📁 生成的绑定文件: pkg/eth/counter/counter.go"

echo ""
echo "=== 任务完成 ==="
echo "1. 智能合约已编译"
echo "2. Go绑定代码已生成"
echo "3. 现在可以运行Go程序与合约交互"