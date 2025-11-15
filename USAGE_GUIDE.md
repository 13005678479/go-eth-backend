# 使用指南 - 详细步骤说明

## 概述

本指南提供完成两个核心任务的详细步骤，包含完整的代码注释和操作说明。

## 任务1：区块链读写

### 1.1 环境准备

#### 安装必要工具

```bash
# 1. 安装Go语言环境 (如果未安装)
# 下载地址: https://golang.org/dl/

# 2. 安装Node.js (用于solc编译器)
# 下载地址: https://nodejs.org/

# 3. 安装solc编译器
npm install -g solc

# 4. 安装abigen工具
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# 5. 下载项目依赖
go mod tidy
```

#### 获取Infura API Key

1. 访问 https://infura.io/register 注册账户
2. 登录后创建新项目
3. 选择"Ethereum"网络
4. 选择"Sepolia"测试网络
5. 复制项目ID（API Key）

### 1.2 运行区块链查询程序

#### 代码详解 (`cmd/simple_query.go`)

```go
package main

import (
    "context"
    "fmt"
    "log"
    "math/big"

    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
    // 1. 配置网络连接参数
    infuraAPIKey := "ea33fc8cbc4545d9ac08fba394c5046b" // 替换为您的API Key
    rpcURL := fmt.Sprintf("https://sepolia.infura.io/v3/%s", infuraAPIKey)

    // 2. 创建以太坊客户端连接
    client, err := ethclient.Dial(rpcURL)
    if err != nil {
        log.Fatalf("连接失败: %v", err)
    }
    defer client.Close() // 确保程序退出时关闭连接

    // 3. 查询最新区块号
    latestBlockNumber, err := client.BlockNumber(context.Background())
    if err != nil {
        log.Fatalf("获取最新区块号失败: %v", err)
    }
    fmt.Printf("最新区块号: %d\n", latestBlockNumber)

    // 4. 查询最新区块详细信息
    latestBlock, err := client.BlockByNumber(context.Background(), big.NewInt(int64(latestBlockNumber)))
    if err != nil {
        log.Fatalf("获取区块信息失败: %v", err)
    }

    // 5. 显示区块信息
    fmt.Println("区块信息:")
    fmt.Printf("  区块号: %d\n", latestBlock.Number())
    fmt.Printf("  区块哈希: %s\n", latestBlock.Hash().Hex())
    fmt.Printf("  时间戳: %d\n", latestBlock.Time())
    fmt.Printf("  交易数量: %d\n", len(latestBlock.Transactions()))
    fmt.Printf("  Gas使用量: %d\n", latestBlock.GasUsed())
    fmt.Printf("  矿工地址: %s\n", latestBlock.Coinbase().Hex())
}
```

#### 运行命令

```bash
# 运行区块链查询程序
go run cmd/simple_query.go
```

#### 预期输出

```
=== 区块链查询示例程序 ===
任务目标: 连接到Sepolia测试网络，查询区块信息

🔗 正在连接到Sepolia测试网络...
✅ 成功连接到Sepolia测试网络

📦 查询最新区块信息...
📊 最新区块号: 5689432

🔍 查询最新区块详细信息...
📋 区块信息:
  区块号: 5689432
  区块哈希: 0x1234...abcd
  时间戳: 1700000000 (Unix时间戳)
  交易数量: 145
  Gas使用量: 15000000
  Gas限制: 30000000
  矿工地址: 0x742d35Cc6634C0532925a3b8Ffb8a2B15a3F2F20
  难度: 123456789
  父区块哈希: 0x5678...efgh

=== 区块链查询示例完成 ===
```

### 1.3 发送交易程序

#### 准备工作

1. **获取测试ETH**
   - 访问 https://sepoliafaucet.com/
   - 输入您的以太坊地址
   - 等待测试ETH到账（通常需要几分钟）

2. **配置私钥**
   - 编辑 `cmd/simple_transaction.go` 文件
   - 找到 `privateKeyHex := "YOUR_PRIVATE_KEY_HERE"`
   - 替换为您的测试网络私钥

#### 代码详解 (`cmd/simple_transaction.go`)

```go
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
)

func main() {
    // 1. 配置网络连接
    infuraAPIKey := "ea33fc8cbc4545d9ac08fba394c5046b"
    rpcURL := fmt.Sprintf("https://sepolia.infura.io/v3/%s", infuraAPIKey)
    client, err := ethclient.Dial(rpcURL)
    if err != nil {
        log.Fatalf("连接失败: %v", err)
    }
    defer client.Close()

    // 2. 配置账户信息
    privateKeyHex := "YOUR_PRIVATE_KEY_HERE" // 替换为实际私钥
    privateKey, err := crypto.HexToECDSA(privateKeyHex)
    if err != nil {
        log.Fatalf("私钥解析失败: %v", err)
    }

    // 3. 从私钥获取地址
    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        log.Fatal("无法获取公钥")
    }
    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

    // 4. 检查账户余额
    balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
    if err != nil {
        log.Fatalf("查询余额失败: %v", err)
    }
    etherBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
    fmt.Printf("账户余额: %s ETH\n", etherBalance.String())

    // 5. 获取交易参数
    nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
    gasPrice, err := client.SuggestGasPrice(context.Background())
    chainID, err := client.ChainID(context.Background())

    // 6. 创建交易对象
    toAddress := common.HexToAddress("0x742d35Cc6634C0532925a3b8Ffb8a2B15a3F2F20")
    value := big.NewInt(1000000000000000) // 0.001 ETH
    gasLimit := uint64(21000)
    
    tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

    // 7. 签名交易
    signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)

    // 8. 发送交易
    err = client.SendTransaction(context.Background(), signedTx)

    // 9. 等待交易确认
    receipt, err := bind.WaitMined(context.Background(), client, signedTx)
    if receipt.Status == 1 {
        fmt.Printf("交易确认成功! 区块号: %d\n", receipt.BlockNumber.Uint64())
    }
}
```

#### 运行命令

```bash
# 运行交易发送程序
go run cmd/simple_transaction.go
```

## 任务2：合约代码生成

### 2.1 编译智能合约

#### 智能合约代码 (`contracts/SimpleCounter.sol`)

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title SimpleCounter
 * @dev 一个极简的计数器智能合约，仅包含核心功能
 */
contract SimpleCounter {
    uint256 public count; // 公共状态变量，存储计数值
    
    // 事件：当计数器增加时触发
    event CountIncremented(uint256 newCount);
    
    // 构造函数：初始化计数器为0
    constructor() {
        count = 0;
    }
    
    /**
     * @dev 增加计数器值（每次增加1）
     */
    function increment() public {
        count += 1;
        emit CountIncremented(count); // 触发事件
    }
    
    /**
     * @dev 获取当前计数器值
     * @return 当前计数值
     */
    function getCount() public view returns (uint256) {
        return count;
    }
}
```

#### 编译脚本详解 (`scripts/compile_and_generate.sh`)

```bash
#!/bin/bash

# 以太坊智能合约编译和Go绑定代码生成脚本

echo "=== 以太坊智能合约编译和Go绑定代码生成 ==="

# 检查依赖工具是否安装
if ! command -v solc &> /dev/null; then
    echo "错误: 未找到 solc 编译器"
    echo "请安装: npm install -g solc"
    exit 1
fi

if ! command -v abigen &> /dev/null; then
    echo "错误: 未找到 abigen 工具"
    echo "请安装: go install github.com/ethereum/go-ethereum/cmd/abigen@latest"
    exit 1
fi

# 创建输出目录
mkdir -p contracts/compiled

# 编译智能合约生成ABI和字节码
echo "编译智能合约..."
solc --bin --abi --overwrite -o contracts/compiled contracts/SimpleCounter.sol

# 生成Go绑定代码
echo "生成Go绑定代码..."
abigen --bin=contracts/compiled/SimpleCounter.bin --abi=contracts/compiled/SimpleCounter.abi --pkg=counter --out=pkg/eth/counter/counter.go

echo "任务完成"
```

#### 运行编译脚本

```bash
# 给脚本执行权限
chmod +x scripts/compile_and_generate.sh

# 运行编译脚本
./scripts/compile_and_generate.sh
```

### 2.2 与合约交互

#### 代码详解 (`cmd/simple_contract.go`)

```go
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

func main() {
    // 1. 连接到网络
    client, err := ethclient.Dial("https://sepolia.infura.io/v3/ea33fc8cbc4545d9ac08fba394c5046b")
    if err != nil {
        log.Fatalf("连接失败: %v", err)
    }
    defer client.Close()

    // 2. 配置账户
    privateKey, err := crypto.HexToECDSA("YOUR_PRIVATE_KEY_HERE")
    chainID, err := client.ChainID(context.Background())
    auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)

    // 3. 部署合约
    contractAddress, tx, contractInstance, err := counter.DeployCounter(auth, client)
    if err != nil {
        log.Fatalf("合约部署失败: %v", err)
    }
    fmt.Printf("合约地址: %s\n", contractAddress.Hex())

    // 4. 等待部署确认
    receipt, err := bind.WaitMined(context.Background(), client, tx)
    if receipt.Status == 1 {
        fmt.Println("合约部署成功!")
    }

    // 5. 读取合约状态（不需要Gas）
    currentCount, err := contractInstance.GetCount(nil)
    fmt.Printf("当前计数值: %d\n", currentCount)

    // 6. 调用合约方法（需要Gas）
    incrementTx, err := contractInstance.Increment(auth)
    if err != nil {
        log.Fatalf("调用increment方法失败: %v", err)
    }

    // 7. 等待交易确认
    _, err = bind.WaitMined(context.Background(), client, incrementTx)

    // 8. 验证结果
    updatedCount, err := contractInstance.GetCount(nil)
    fmt.Printf("更新后的计数值: %d\n", updatedCount)
}
```

#### 运行命令

```bash
# 运行合约交互程序
go run cmd/simple_contract.go
```

## 常见问题解答

### Q: 连接失败怎么办？
**A:** 检查以下内容：
1. 网络连接是否正常
2. Infura API Key是否正确
3. 是否选择了正确的网络（Sepolia）

### Q: 余额不足怎么办？
**A:** 从Sepolia水龙头获取测试ETH：
1. 访问 https://sepoliafaucet.com/
2. 输入您的以太坊地址
3. 等待几分钟ETH到账

### Q: 合约编译失败怎么办？
**A:** 检查以下内容：
1. solc编译器是否正确安装：`solc --version`
2. 合约语法是否正确
3. 是否使用了正确的Solidity版本

### Q: 生成的绑定代码无法导入怎么办？
**A:** 确保：
1. 绑定代码生成成功
2. 文件路径正确：`pkg/eth/counter/counter.go`
3. 导入语句正确：`counter "go-eth-backend/pkg/eth/counter"`

## 安全注意事项

1. **私钥安全**：永远不要将主网私钥用于测试网络
2. **测试网络**：所有操作都在Sepolia测试网络进行
3. **敏感信息**：不要将API Key和私钥提交到版本控制
4. **生产环境**：本示例代码仅用于学习，生产环境需要额外安全措施

## 下一步学习

完成本示例后，可以尝试：
1. 编写更复杂的智能合约
2. 实现合约事件监听
3. 构建完整的DApp前端
4. 学习合约安全最佳实践

---

**提示：** 本指南中的所有代码都包含详细的注释，建议仔细阅读代码理解每个步骤的作用。