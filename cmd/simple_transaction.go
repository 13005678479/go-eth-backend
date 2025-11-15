package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go-eth-backend/internal/pkg/config"
)

// 交易发送示例程序
// 演示如何在Sepolia测试网络上发送以太币交易

func main() {
	fmt.Println("=== 以太币交易发送示例程序 ===")
	fmt.Println("任务目标: 在Sepolia测试网络上发送以太币转账交易")
	fmt.Println()

	// ⚠️ 重要安全提示 ⚠️
	fmt.Println("🚨 安全警告:")
	fmt.Println("- 请确保使用测试网络私钥")
	fmt.Println("- 不要在生产环境中使用硬编码的私钥")
	fmt.Println("- 测试网络ETH没有实际价值")
	fmt.Println()

	// Step 1: 加载配置文件
	fmt.Println("📋 加载配置文件...")
	config := config.LoadConfigOrExit("config.yaml")
	
	// 从配置文件获取网络配置
	sepoliaConfig := config.GetSepoliaConfig()
	rpcURL := sepoliaConfig.RPCURL

	fmt.Println("🔗 正在连接到Sepolia测试网络...")
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("❌ 连接失败: %v", err)
	}
	defer client.Close()

	fmt.Println("✅ 成功连接到Sepolia测试网络")
	fmt.Println()

	// Step 2: 配置账户信息
	// 从配置文件获取私钥
	privateKeyHex := config.GetTestPrivateKey()
	
	fmt.Printf("🔐 已加载测试私钥，账户地址: %s
", getAddressFromPrivateKey(privateKeyHex))
	fmt.Println()

	// 解析私钥
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("❌ 私钥解析失败: %v", err)
	}

	// 从私钥获取公钥和地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("❌ 无法获取公钥")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("📧 发送方地址: %s\n", fromAddress.Hex())

	// 接收方地址（示例地址）
	toAddress := common.HexToAddress("0x742d35Cc6634C0532925a3b8Ffb8a2B15a3F2F20")
	fmt.Printf("📧 接收方地址: %s\n", toAddress.Hex())
	fmt.Println()

	// Step 3: 查询账户余额
	fmt.Println("💰 查询账户余额...")

	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatalf("❌ 查询余额失败: %v", err)
	}

	// 将wei转换为ether
	etherBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	fmt.Printf("📊 发送方余额: %s ETH\n", etherBalance.String())

	// 检查余额是否足够
	if balance.Cmp(big.NewInt(10000000000000000)) == -1 { // 0.01 ETH
		fmt.Println("❌ 余额不足，至少需要0.01 ETH用于测试")
		fmt.Println("💡 请从Sepolia水龙头获取测试ETH: https://sepoliafaucet.com/")
		return
	}
	fmt.Println()

	// Step 4: 获取网络信息
	fmt.Println("🌐 获取网络信息...")

	// 获取链ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取链ID失败: %v", err)
	}
	fmt.Printf("📊 链ID: %s\n", chainID.String())

	// 获取最新区块号
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取区块号失败: %v", err)
	}
	fmt.Printf("📊 最新区块号: %d\n", blockNumber)
	fmt.Println()

	// Step 5: 准备交易参数
	fmt.Println("📝 准备交易参数...")

	// 获取nonce（交易序号）
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("❌ 获取nonce失败: %v", err)
	}
	fmt.Printf("📊 交易序号(nonce): %d\n", nonce)

	// 设置转账金额（0.001 ETH）
	value := big.NewInt(1000000000000000) // 0.001 ETH in wei
	fmt.Printf("💰 转账金额: 0.001 ETH\n")

	// 获取推荐的Gas价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("❌ 获取Gas价格失败: %v", err)
	}
	fmt.Printf("⛽ Gas价格: %s wei\n", gasPrice.String())

	// 设置Gas限制
	gasLimit := uint64(21000) // 标准ETH转账的Gas限制
	fmt.Printf("⛽ Gas限制: %d\n", gasLimit)
	fmt.Println()

	// Step 6: 创建交易对象
	fmt.Println("🔧 创建交易对象...")

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	fmt.Printf("📋 交易已创建:\n")
	fmt.Printf("  发送方: %s\n", fromAddress.Hex())
	fmt.Printf("  接收方: %s\n", toAddress.Hex())
	fmt.Printf("  金额: 0.001 ETH\n")
	fmt.Printf("  Nonce: %d\n", nonce)
	fmt.Println()

	// Step 7: 签名交易
	fmt.Println("✍️  签名交易...")

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		log.Fatalf("❌ 交易签名失败: %v", err)
	}
	fmt.Println("✅ 交易签名成功")
	fmt.Println()

	// Step 8: 发送交易
	fmt.Println("🚀 发送交易到网络...")

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("❌ 发送交易失败: %v", err)
	}

	fmt.Printf("✅ 交易发送成功!\n")
	fmt.Printf("📋 交易哈希: %s\n", signedTx.Hash().Hex())
	fmt.Println()

	// Step 9: 等待交易确认
	fmt.Println("⏳ 等待交易确认...")

	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		log.Fatalf("❌ 等待交易确认失败: %v", err)
	}

	if receipt.Status == 1 {
		fmt.Printf("✅ 交易确认成功!\n")
		fmt.Printf("📋 区块号: %d\n", receipt.BlockNumber.Uint64())
		fmt.Printf("📋 Gas使用量: %d\n", receipt.GasUsed)
	} else {
		fmt.Printf("❌ 交易失败!\n")
	}
	fmt.Println()

	fmt.Println("=== 交易发送示例完成 ===")
}

// 显示教程信息
func showTutorial(client *ethclient.Client) {
	fmt.Println("📚 使用教程:")
	fmt.Println("1. 获取Sepolia测试ETH:")
	fmt.Println("   - 访问: https://sepoliafaucet.com/")
	fmt.Println("   - 输入您的以太坊地址获取测试ETH")
	fmt.Println()

	fmt.Println("2. 配置私钥:")
	fmt.Println("   - 编辑 simple_transaction.go 文件")
	fmt.Println("   - 将 YOUR_PRIVATE_KEY_HERE 替换为您的测试网络私钥")
	fmt.Println("   - 注意: 不要使用主网私钥!")
	fmt.Println()

	fmt.Println("3. 生成测试账户（可选）:")
	fmt.Println("   可以使用以下命令生成测试账户:")
	fmt.Println("   $ openssl ecparam -name secp256k1 -genkey -noout | openssl ec -text -noout")
	fmt.Println()

	fmt.Println("4. 运行程序:")
	fmt.Println("   $ go run cmd/simple_transaction.go")
	fmt.Println()

	// 显示示例账户信息（仅用于演示）
	fmt.Println("💡 示例账户信息（仅用于演示）:")
	fmt.Println("   地址: 0x742d35Cc6634C0532925a3b8Ffb8a2B15a3F2F20")
	fmt.Println("   余额: 可以从水龙头获取测试ETH")
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
