package main

import (
	"fmt"
	"log"
	"time"

	"go-eth-backend/internal/pkg/eth"
)

func main() {
	fmt.Println("=== Sepolia 测试网络连接测试 ===")
	
	// 测试 Sepolia 网络连接
	client, err := eth.NewClient("sepolia", "")
	if err != nil {
		log.Fatalf("❌ 连接 Sepolia 网络失败: %v", err)
	}
	fmt.Println("✅ Sepolia 网络连接成功")
	
	// 获取最新区块号
	fmt.Println("\n正在获取最新区块信息...")
	blockNumber, err := client.GetLatestBlockNumber()
	if err != nil {
		log.Printf("⚠️  获取最新区块号失败: %v", err)
	} else {
		fmt.Printf("✅ 最新区块号: %d\n", blockNumber)
	}
	
	// 测试区块信息获取
	if blockNumber > 0 {
		fmt.Println("\n正在获取区块详情...")
		block, err := client.GetBlockByNumber(blockNumber)
		if err != nil {
			log.Printf("⚠️  获取区块详情失败: %v", err)
		} else {
			fmt.Printf("✅ 区块哈希: %s\n", block.Hash)
			fmt.Printf("✅ 时间戳: %s\n", block.Timestamp.Format(time.RFC3339))
		}
	}
	
	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("你的 Sepolia 测试网络配置已成功应用！")
}