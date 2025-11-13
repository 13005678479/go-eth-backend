package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Transaction 交易信息结构体
type Transaction struct {
	Hash        string   `json:"hash"`
	From        string   `json:"from"`
	To          string   `json:"to"`
	Value       *big.Int `json:"value"`
	GasPrice    *big.Int `json:"gasPrice"`
	GasLimit    uint64   `json:"gasLimit"`
	Nonce       uint64   `json:"nonce"`
	Data        []byte   `json:"data"`
	ChainID     *big.Int `json:"chainId"`
}

// SendTransaction 发送以太币交易
func (c *Client) SendTransaction(fromPrivateKey, toAddress string, amount *big.Int) (string, error) {
	// 解析私钥
	privateKey, err := crypto.HexToECDSA(fromPrivateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	// 获取公钥地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取nonce
	nonce, err := c.client.PendingNonceAt(c.ctx, fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	// 获取gas价格
gasPrice, err := c.client.SuggestGasPrice(c.ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %v", err)
	}

	// 获取链ID
	chainID, err := c.client.ChainID(c.ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID: %v", err)
	}

	// 构建交易
	tx := types.NewTransaction(
		nonce,
		common.HexToAddress(toAddress),
		amount,
		21000, // 标准的以太币转账gas限制
		gasPrice,
		nil,
	)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 发送交易
	err = c.client.SendTransaction(c.ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx.Hash().Hex(), nil
}

// GetTransactionReceipt 获取交易收据
func (c *Client) GetTransactionReceipt(txHash string) (*types.Receipt, error) {
	hash := common.HexToHash(txHash)
	receipt, err := c.client.TransactionReceipt(c.ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction receipt: %v", err)
	}
	return receipt, nil
}

// GetTransactionByHash 根据哈希获取交易信息
func (c *Client) GetTransactionByHash(txHash string) (*types.Transaction, error) {
	hash := common.HexToHash(txHash)
	tx, isPending, err := c.client.TransactionByHash(c.ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}
	if isPending {
		return nil, fmt.Errorf("transaction is still pending")
	}
	return tx, nil
}

// WaitForTransactionReceipt 等待交易确认
func (c *Client) WaitForTransactionReceipt(txHash string) (*types.Receipt, error) {
	hash := common.HexToHash(txHash)
	
	// 设置超时上下文
	ctx, cancel := context.WithTimeout(c.ctx, 300*time.Second)
	defer cancel()
	
	// 等待交易确认
	receipt, err := bind.WaitMined(ctx, c.client, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction receipt: %v", err)
	}
	
	return receipt, nil
}