package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

// 区块链查询示例程序
// 连接到Sepolia测试网络并查询区块信息

func main() {
	fmt.Println("=== 区块链查询示例程序 ===")
	fmt.Println("任务目标: 连接到Sepolia测试网络，查询区块信息")
	fmt.Println()

	// Step 1: 连接到Sepolia测试网络
	// 使用Infura的免费API端点（需要替换为您的API Key）
	// 获取Infura API Key: https://infura.io/register
	infuraAPIKey := "ea33fc8cbc4545d9ac08fba394c5046b" // 使用正确的Infura API Key
	rpcURL := fmt.Sprintf("https://sepolia.infura.io/v3/%s", infuraAPIKey)

	fmt.Println("🔗 正在连接到Sepolia测试网络...")
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("❌ 连接失败: %v", err)
	}
	defer client.Close()

	fmt.Println("✅ 成功连接到Sepolia测试网络")
	fmt.Println()

	// Step 2: 查询最新区块号
	fmt.Println("📦 查询最新区块信息...")

	// 获取最新区块号
	latestBlockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取最新区块号失败: %v", err)
	}

	fmt.Printf("📊 最新区块号: %d\n", latestBlockNumber)
	fmt.Println()

	// Step 3: 查询最新区块详细信息
	fmt.Println("🔍 查询最新区块详细信息...")

	// 获取最新区块完整信息
	latestBlock, err := client.BlockByNumber(context.Background(), big.NewInt(int64(latestBlockNumber)))
	if err != nil {
		log.Fatalf("❌ 获取区块信息失败: %v", err)
	}

	// 显示区块信息
	fmt.Println("📋 区块信息:")
	fmt.Printf("  区块号: %d\n", latestBlock.Number())
	fmt.Printf("  区块哈希: %s\n", latestBlock.Hash().Hex())
	fmt.Printf("  时间戳: %d (Unix时间戳)\n", latestBlock.Time())
	fmt.Printf("  交易数量: %d\n", len(latestBlock.Transactions()))
	fmt.Printf("  Gas使用量: %d\n", latestBlock.GasUsed())
	fmt.Printf("  Gas限制: %d\n", latestBlock.GasLimit())
	fmt.Printf("  矿工地址: %s\n", latestBlock.Coinbase().Hex())
	fmt.Printf("  难度: %s\n", latestBlock.Difficulty().String())
	fmt.Printf("  父区块哈希: %s\n", latestBlock.ParentHash().Hex())
	fmt.Println()

	// Step 4: 查询指定区块信息（示例：查询第1000个区块）
	fmt.Println("🔍 查询指定区块信息（示例：区块号 1000）...")

	targetBlockNumber := uint64(1000)
	targetBlock, err := client.BlockByNumber(context.Background(), big.NewInt(int64(targetBlockNumber)))
	if err != nil {
		fmt.Printf("⚠️  查询区块 %d 失败: %v\n", targetBlockNumber, err)
		fmt.Println("    (可能是区块号不存在或同步问题)")
	} else {
		fmt.Printf("📋 区块 %d 信息:\n", targetBlockNumber)
		fmt.Printf("  区块哈希: %s\n", targetBlock.Hash().Hex())
		fmt.Printf("  时间戳: %d\n", targetBlock.Time())
		fmt.Printf("  交易数量: %d\n", len(targetBlock.Transactions()))
	}
	fmt.Println()

	// Step 5: 网络状态检查
	fmt.Println("🌐 网络状态检查...")

	// 检查网络同步状态
	syncProgress, err := client.SyncProgress(context.Background())
	if err != nil {
		fmt.Printf("⚠️  获取同步状态失败: %v\n", err)
	} else if syncProgress != nil {
		fmt.Printf("📊 同步进度: %d/%d (%.2f%%)\n",
			syncProgress.CurrentBlock, syncProgress.HighestBlock,
			float64(syncProgress.CurrentBlock)/float64(syncProgress.HighestBlock)*100)
	} else {
		fmt.Println("✅ 网络已完全同步")
	}

	fmt.Println()
	fmt.Println("=== 区块链查询示例完成 ===")
	fmt.Println("📝 总结:")
	fmt.Println("1. 成功连接到Sepolia测试网络")
	fmt.Println("2. 查询了最新区块信息")
	fmt.Println("3. 展示了区块的基本属性")
	fmt.Println("4. 验证了网络连接状态")
	fmt.Println()
	fmt.Println("💡 下一步:")
	fmt.Println("- 可以尝试查询交易信息")
	fmt.Println("- 可以查询账户余额")
	fmt.Println("- 可以发送测试交易")
}

// 查询交易数量（辅助函数）
func getTransactionCount(client *ethclient.Client, blockNumber uint64) (int, error) {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		return 0, err
	}
	return len(block.Transactions()), nil
}
