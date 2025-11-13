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
	"github.com/ethereum/go-ethereum/ethclient"

	"go-eth-backend/internal/pkg/eth"
	"go-eth-backend/pkg/eth/counter"
)

func main() {
	fmt.Println("=== 智能合约交互示例 ===")
	
	// 初始化以太坊客户端
	client, err := eth.NewClient("sepolia", "")
	if err != nil {
		log.Fatalf("❌ 连接 Sepolia 网络失败: %v", err)
	}
	fmt.Println("✅ Sepolia 网络连接成功")

	// 示例合约地址（这里需要替换为实际部署的合约地址）
	contractAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
	
	// 测试合约交互功能
	fmt.Println("\n1. 创建合约实例...")
	contract, err := counter.NewCounter(contractAddress, client.Client)
	if err != nil {
		log.Printf("⚠️  创建合约实例失败: %v", err)
		fmt.Println("将演示合约部署过程...")
		deployContractDemo(client.Client)
		return
	}
	
	fmt.Println("✅ 合约实例创建成功")
	
	// 测试读取合约数据
	fmt.Println("\n2. 读取合约数据...")
	
	// 获取当前计数器值
	count, err := contract.GetCount(nil)
	if err != nil {
		log.Printf("⚠️  获取计数器值失败: %v", err)
	} else {
		fmt.Printf("✅ 当前计数器值: %s\n", count.String())
	}
	
	// 获取合约所有者
	owner, err := contract.Owner(nil)
	if err != nil {
		log.Printf("⚠️  获取合约所有者失败: %v", err)
	} else {
		fmt.Printf("✅ 合约所有者: %s\n", owner.Hex())
	}
	
	fmt.Println("\n3. 合约交互演示完成")
	fmt.Println("注意：要实际执行写入操作，需要提供有效的私钥和已部署的合约地址")
}

// deployContractDemo 演示合约部署过程
func deployContractDemo(client *ethclient.Client) {
	fmt.Println("\n=== 合约部署演示 ===")
	
	// 创建一个测试私钥（仅用于演示，实际使用时请使用真实私钥）
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Printf("❌ 生成测试私钥失败: %v", err)
		return
	}
	
	// 获取链ID
	chainID, err := client.ChainID(context.Background())
	if err != nil nil {
		log.Printf("❌ 获取链ID失败: %v", err)
		return
	}
	
	// 创建交易选项
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Printf("❌ 创建交易选项失败: %v", err)
		return
	}
	
	// 设置 gas 价格和 gas 限制
	auth.GasPrice = big.NewInt(20000000000) // 20 Gwei
	auth.GasLimit = uint64(3000000)         // 300万 gas
	
	fmt.Println("✅ 交易选项配置完成")
	
	// 部署合约
	fmt.Println("\n开始部署合约...")
	address, tx, contract, err := counter.DeployCounter(auth, client)
	if err != nil {
		log.Printf("❌ 部署合约失败: %v", err)
		return
	}
	
	fmt.Printf("✅ 合约部署交易已提交:\n")
	fmt.Printf("   合约地址: %s\n", address.Hex())
	fmt.Printf("   交易哈希: %s\n", tx.Hash().Hex())
	
	// 等待交易确认
	fmt.Println("\n等待交易确认...")
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Printf("❌ 等待交易确认失败: %v", err)
		return
	}
	
	if receipt.Status == 1 {
		fmt.Println("✅ 合约部署成功!")
		
		// 测试新部署的合约
		fmt.Println("\n测试新部署的合约...")
		
		// 获取初始计数器值
		initialCount, err := contract.GetCount(nil)
		if err != nil {
			log.Printf("⚠️  获取初始计数器值失败: %v", err)
		} else {
			fmt.Printf("✅ 初始计数器值: %s\n", initialCount.String())
		}
		
		// 获取合约所有者
		owner, err := contract.Owner(nil)
		if err != nil {
			log.Printf("⚠️  获取合约所有者失败: %v", err)
		} else {
			fmt.Printf("✅ 合约所有者: %s\n", owner.Hex())
		}
		
		fmt.Println("\n部署演示完成!")
		fmt.Printf("新合约地址: %s\n", address.Hex())
	} else {
		fmt.Println("❌ 合约部署失败")
	}
}

// ContractInteractionDemo 演示完整的合约交互流程
func ContractInteractionDemo() {
	fmt.Println("\n=== 完整的合约交互演示 ===")
	
	// 这个函数展示了如何使用合约进行完整的读写操作
	// 包括：读取状态、写入状态、监听事件等
	
	fmt.Println("功能包括:")
	fmt.Println("1. 读取合约状态 (getCount, getOwner)")
	fmt.Println("2. 写入合约状态 (increment, decrement, reset)")
	fmt.Println("3. 监听合约事件")
	fmt.Println("4. 权限控制 (onlyOwner 修饰器)")
	
	fmt.Println("\n要实际执行这些操作，需要:")
	fmt.Println("1. 有效的 Sepolia 测试网账户和私钥")
	fmt.Println("2. 已部署的合约地址")
	fmt.Println("3. 足够的测试 ETH 用于支付 gas 费用")
}