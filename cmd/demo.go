package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"go-eth-backend/config"
	"go-eth-backend/internal/pkg/eth"
	"go-eth-backend/pkg/eth/counter"
)

func main() {
	fmt.Println("=== 以太坊区块链交互演示程序 ===")
	fmt.Println("使用 Sepolia 测试网络")
	fmt.Println()

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建以太坊客户端
	client, err := eth.NewClient(cfg.Ethereum.Networks["sepolia"].RPCURL)
	if err != nil {
		log.Fatalf("创建以太坊客户端失败: %v", err)
	}
	defer client.Close()

	fmt.Println("✅ 已连接到 Sepolia 测试网络")
	fmt.Println()

	// 演示1: 查询最新区块信息
	fmt.Println("📦 演示1: 查询最新区块信息")
	demoBlockQuery(client)
	fmt.Println()

	// 演示2: 查询指定区块信息
	fmt.Println("🔍 演示2: 查询指定区块信息")
	demoSpecificBlockQuery(client)
	fmt.Println()

	// 演示3: 智能合约交互演示
	fmt.Println("🤖 演示3: 智能合约交互演示")
	demoContractInteraction(client, cfg)
	fmt.Println()

	fmt.Println("=== 演示程序完成 ===")
}

func demoBlockQuery(client *eth.Client) {
	// 获取最新区块
	latestBlock, err := client.GetLatestBlock()
	if err != nil {
		log.Printf("获取最新区块失败: %v", err)
		return
	}

	fmt.Printf("最新区块信息:\n")
	fmt.Printf("  区块号: %d\n", latestBlock.Number)
	fmt.Printf("  区块哈希: %s\n", latestBlock.Hash)
	fmt.Printf("  时间戳: %s\n", latestBlock.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("  交易数量: %d\n", latestBlock.TransactionCount)
	fmt.Printf("  Gas使用量: %d\n", latestBlock.GasUsed)
	fmt.Printf("  矿工地址: %s\n", latestBlock.Miner)
	fmt.Printf("  区块大小: %d bytes\n", latestBlock.Size)

	// 获取区块头信息
	header, err := client.GetBlockHeaderByNumber(latestBlock.Number)
	if err != nil {
		log.Printf("获取区块头失败: %v", err)
		return
	}

	fmt.Printf("\n区块头信息:\n")
	fmt.Printf("  父区块哈希: %s\n", header.ParentHash)
	fmt.Printf("  Gas限制: %d\n", header.GasLimit)
}

func demoSpecificBlockQuery(client *eth.Client) {
	// 查询一个较早的区块（例如区块号 4000000）
	targetBlockNumber := uint64(4000000)
	
	block, err := client.GetBlockByNumber(targetBlockNumber)
	if err != nil {
		log.Printf("获取区块 %d 失败: %v", targetBlockNumber, err)
		return
	}

	fmt.Printf("区块 %d 信息:\n", targetBlockNumber)
	fmt.Printf("  区块哈希: %s\n", block.Hash)
	fmt.Printf("  时间戳: %s\n", block.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("  交易数量: %d\n", block.TransactionCount)
	fmt.Printf("  难度: %s\n", block.Difficulty)
	fmt.Printf("  额外数据: %s\n", block.ExtraData)

	// 查询交易数量
	txCount, err := client.GetBlockTransactionCount(targetBlockNumber)
	if err != nil {
		log.Printf("获取交易数量失败: %v", err)
		return
	}
	fmt.Printf("  验证后的交易数量: %d\n", txCount)
}

func demoContractInteraction(client *eth.Client, cfg *config.Config) {
	// 注意：这里需要实际的合约地址和私钥才能进行部署和调用
	// 这里仅演示如何使用生成的绑定代码

	fmt.Println("📝 智能合约功能演示")
	
	// 部署合约（需要私钥，这里仅展示代码结构）
	fmt.Println("1. 合约部署功能已准备（需要实际私钥）")
	
	// 合约调用示例
	fmt.Println("2. 合约调用功能已准备")
	
	// 演示如何使用生成的绑定代码
	fmt.Println("3. 生成的 Go 绑定代码包含以下方法:")
	fmt.Println("   - counter.NewCounter() - 创建合约实例")
	fmt.Println("   - counter.CounterCaller - 读取合约状态")
	fmt.Println("   - counter.CounterTransactor - 写入合约状态")
	fmt.Println("   - counter.CounterFilterer - 监听合约事件")
	
	// 显示合约 ABI 信息
	fmt.Println("\n📋 合约方法:")
	fmt.Println("   - getCount(): 获取当前计数值")
	fmt.Println("   - increment(): 增加计数值")
	fmt.Println("   - decrement(): 减少计数值")
	fmt.Println("   - reset(): 重置计数器")
	fmt.Println("   - CountUpdated 事件: 计数更新时触发")
	
	fmt.Println("\n⚠️  注意：实际部署和调用需要配置有效的私钥和测试网ETH")
}

// deployContract 演示合约部署（需要实际私钥）
func deployContract(client *eth.Client, privateKey string) (common.Address, error) {
	// 解析私钥
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return common.Address{}, fmt.Errorf("解析私钥失败: %v", err)
	}

	// 创建交易选项
	auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(11155111)) // Sepolia chain ID
	if err != nil {
		return common.Address{}, fmt.Errorf("创建交易授权失败: %v", err)
	}

	// 部署合约
	address, tx, _, err := counter.DeployCounter(auth, client.GetRawClient())
	if err != nil {
		return common.Address{}, fmt.Errorf("部署合约失败: %v", err)
	}

	fmt.Printf("合约部署交易哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("合约地址: %s\n", address.Hex())

	return address, nil
}

// interactWithContract 演示合约交互
func interactWithContract(client *eth.Client, contractAddress common.Address, privateKey string) error {
	// 创建合约实例
	contract, err := counter.NewCounter(contractAddress, client.GetRawClient())
	if err != nil {
		return fmt.Errorf("创建合约实例失败: %v", err)
	}

	// 读取合约状态
	count, err := contract.GetCount(nil)
	if err != nil {
		return fmt.Errorf("读取计数值失败: %v", err)
	}
	fmt.Printf("当前计数值: %d\n", count)

	// 调用合约方法（需要私钥）
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return fmt.Errorf("解析私钥失败: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(11155111))
	if err != nil {
		return fmt.Errorf("创建交易授权失败: %v", err)
	}

	// 增加计数器
	tx, err := contract.Increment(auth)
	if err != nil {
		return fmt.Errorf("调用increment方法失败: %v", err)
	}

	fmt.Printf("Increment交易哈希: %s\n", tx.Hash().Hex())

	return nil
}