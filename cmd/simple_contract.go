package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	counter "go-eth-backend/pkg/eth/counter" // 导入生成的绑定代码
)

// 智能合约交互示例程序
// 演示如何部署和调用SimpleCounter合约

func main() {
	fmt.Println("=== 智能合约交互示例程序 ===")
	fmt.Println("任务目标: 部署SimpleCounter合约并进行交互")
	fmt.Println()

	// Step 1: 连接到Sepolia测试网络
	infuraAPIKey := "ea33fc8cbc4545d9ac08fba394c5046b" // 请替换为您的API Key
	rpcURL := fmt.Sprintf("https://sepolia.infura.io/v3/%s", infuraAPIKey)

	fmt.Println("🔗 正在连接到Sepolia测试网络...")
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("❌ 连接失败: %v", err)
	}
	defer client.Close()

	fmt.Println("✅ 成功连接到Sepolia测试网络")
	fmt.Println()

	// Step 2: 配置账户信息
	privateKeyHex := "4323d9e4f879855a70a3c19b732dde4d1bdb0829b0c30be408ad4b8e24e45e60" // 测试网络私钥
	
	if privateKeyHex == "YOUR_PRIVATE_KEY_HERE" {
		showTutorial(client)
		return
	}

	// 解析私钥
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("❌ 私钥解析失败: %v", err)
	}

	// 获取账户地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("❌ 无法获取公钥")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	fmt.Printf("📧 账户地址: %s\n", fromAddress.Hex())

	// Step 3: 检查账户余额
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatalf("❌ 查询余额失败: %v", err)
	}

	etherBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	fmt.Printf("💰 账户余额: %s ETH\n", etherBalance.String())

	if balance.Cmp(big.NewInt(10000000000000000)) == -1 { // 0.01 ETH
		fmt.Println("❌ 余额不足，至少需要0.01 ETH用于合约部署")
		fmt.Println("💡 请从Sepolia水龙头获取测试ETH: https://sepoliafaucet.com/")
		return
	}
	fmt.Println()

	// Step 4: 准备交易授权
	fmt.Println("🔧 准备交易授权...")
	
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取链ID失败: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("❌ 创建交易授权失败: %v", err)
	}

	// 设置Gas价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取Gas价格失败: %v", err)
	}
	auth.GasPrice = gasPrice

	fmt.Printf("⛽ Gas价格: %s wei\n", auth.GasPrice.String())
	fmt.Println()

	// Step 5: 部署合约
	fmt.Println("🚀 部署SimpleCounter合约...")
	
	contractAddress, tx, contractInstance, err := counter.DeployCounter(auth, client)
	if err != nil {
		log.Fatalf("❌ 合约部署失败: %v", err)
	}

	fmt.Printf("✅ 合约部署交易已发送!\n")
	fmt.Printf("📋 交易哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("📋 合约地址: %s\n", contractAddress.Hex())
	fmt.Println()

	// Step 6: 等待合约部署确认
	fmt.Println("⏳ 等待合约部署确认...")
	
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatalf("❌ 等待交易确认失败: %v", err)
	}

	if receipt.Status == 1 {
		fmt.Printf("✅ 合约部署成功!\n")
		fmt.Printf("📋 区块号: %d\n", receipt.BlockNumber.Uint64())
		fmt.Printf("📋 Gas使用量: %d\n", receipt.GasUsed)
	} else {
		log.Fatal("❌ 合约部署失败!")
	}
	fmt.Println()

	// Step 7: 与已部署的合约交互
	fmt.Println("🤖 与合约交互...")
	
	// 方法1: 读取合约状态（不需要Gas）
	fmt.Println("1. 读取合约状态:")
	
	// 获取当前计数值
	currentCount, err := contractInstance.GetCount(nil)
	if err != nil {
		log.Fatalf("❌ 读取计数值失败: %v", err)
	}
	fmt.Printf("  当前计数值: %d\n", currentCount)

	// 方法2: 调用合约方法（需要Gas）
	fmt.Println("\n2. 调用合约方法:")
	
	// 调用increment方法增加计数器
	fmt.Println("  调用increment()方法...")
	
	incrementTx, err := contractInstance.Increment(auth)
	if err != nil {
		log.Fatalf("❌ 调用increment方法失败: %v", err)
	}

	fmt.Printf("  ✅ Increment交易已发送: %s\n", incrementTx.Hash().Hex())

	// 等待交易确认
	fmt.Println("  等待交易确认...")
	
	_, err = bind.WaitMined(context.Background(), client, incrementTx)
	if err != nil {
		log.Fatalf("❌ 等待交易确认失败: %v", err)
	}

	fmt.Println("  ✅ 交易确认成功!")

	// 再次读取计数值验证变化
	fmt.Println("\n3. 验证计数值变化:")
	
	updatedCount, err := contractInstance.GetCount(nil)
	if err != nil {
		log.Fatalf("❌ 读取计数值失败: %v", err)
	}
	fmt.Printf("  更新后的计数值: %d\n", updatedCount)

	if updatedCount.Cmp(currentCount) > 0 {
		fmt.Println("  ✅ 计数器成功增加!")
	} else {
		fmt.Println("  ❌ 计数器没有变化")
	}
	fmt.Println()

	// Step 8: 与现有合约交互（如果不想部署新合约）
	fmt.Println("💡 与现有合约交互示例:")
	
	// 示例：连接到一个已存在的合约
	existingContractAddress := common.HexToAddress("0x...") // 替换为实际的合约地址
	existingContract, err := counter.NewCounter(existingContractAddress, client)
	if err != nil {
		fmt.Println("  ⚠️ 无法连接示例合约（需要实际合约地址）")
	} else {
		existingCount, err := existingContract.GetCount(nil)
		if err != nil {
			fmt.Println("  ⚠️ 无法读取示例合约状态")
		} else {
			fmt.Printf("  示例合约计数值: %d\n", existingCount)
		}
	}

	fmt.Println()
	fmt.Println("=== 智能合约交互示例完成 ===")
	fmt.Println("📝 总结:")
	fmt.Println("1. 成功部署了SimpleCounter合约")
	fmt.Println("2. 调用了合约的increment方法")
	fmt.Println("3. 验证了合约状态的改变")
	fmt.Println()
	fmt.Println("💡 下一步:")
	fmt.Println("- 可以尝试调用其他合约方法")
	fmt.Println("- 可以监听合约事件")
	fmt.Println("- 可以部署更复杂的合约")
}

// 显示教程信息
func showTutorial(client *ethclient.Client) {
	fmt.Println("📚 使用教程:")
	fmt.Println("1. 准备工作:")
	fmt.Println("   - 获取Sepolia测试ETH: https://sepoliafaucet.com/")
	fmt.Println("   - 配置私钥到 simple_contract.go 文件")
	fmt.Println()
	
	fmt.Println("2. 生成Go绑定代码:")
	fmt.Println("   $ chmod +x scripts/compile_and_generate.sh")
	fmt.Println("   $ ./scripts/compile_and_generate.sh")
	fmt.Println()
	
	fmt.Println("3. 运行程序:")
	fmt.Println("   $ go run cmd/simple_contract.go")
	fmt.Println()
	
	fmt.Println("4. 智能合约详情:")
	fmt.Println("   合约名称: SimpleCounter")
	fmt.Println("   合约方法:")
	fmt.Println("     - increment(): 增加计数器")
	fmt.Println("     - getCount(): 获取当前计数值")
	fmt.Println("   合约事件: CountIncremented")
	fmt.Println()
	
	// 检查网络连接
	fmt.Println("🌐 网络连接状态:")
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		fmt.Println("❌ 网络连接失败")
	} else {
		fmt.Printf("✅ 网络连接正常，最新区块: %d\n", blockNumber)
	}
}

// 查询合约事件（辅助函数）
func queryContractEvents(contractInstance *counter.Counter) {
	fmt.Println("📊 查询合约事件...")
	
	// 这里可以添加事件查询逻辑
	// 实际使用时需要实现事件过滤和监听
	fmt.Println("  事件查询功能待实现...")
}