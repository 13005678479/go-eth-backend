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
	ctx := context.Background()
	
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
	nonce, err := c.Client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	// 获取gas价格
gasPrice, err := c.Client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %v", err)
	}

	// 获取链ID
	chainID, err := c.Client.ChainID(ctx)
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
	err = c.Client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx.Hash().Hex(), nil
}



// WaitForTransactionReceipt 等待交易确认
func (c *Client) WaitForTransactionReceipt(txHash string) (*types.Receipt, error) {
	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	
	// 等待交易确认
	tx, _, err := c.GetTransactionByHash(txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %v", err)
	}
	
	receipt, err := bind.WaitMined(ctx, c.Client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction receipt: %v", err)
	}
	
	return receipt, nil
}