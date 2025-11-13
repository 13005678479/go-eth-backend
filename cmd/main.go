package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go-eth-backend/internal/pkg/eth"
)

func main() {
	// 命令行参数
	network := flag.String("network", "mainnet", "Ethereum network (mainnet, ropsten, etc)")
	rpcURL := flag.String("rpc", "", "Ethereum RPC URL")
	flag.Parse()

	// 初始化以太坊客户端
	client, err := eth.NewClient(*network, *rpcURL)
	if err != nil {
		log.Fatalf("Failed to create Ethereum client: %v", err)
	}

	fmt.Println("=== Go Ethereum Backend ===")
	fmt.Printf("Network: %s\n", *network)
	
	// 获取最新区块号
	blockNumber, err := client.GetLatestBlockNumber()
	if err != nil {
		log.Printf("Failed to get latest block number: %v", err)
	} else {
		fmt.Printf("Latest block number: %d\n", blockNumber)
	}

	// 启动服务
	fmt.Println("\nStarting Ethereum backend service...")
	fmt.Println("Press Ctrl+C to exit")

	// 简单的等待循环
	select {}
}