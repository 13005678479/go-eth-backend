package eth

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Block 区块信息结构
type Block struct {
	Number           uint64
	Hash             string
	Timestamp        time.Time
	TransactionCount int
	GasUsed          uint64
	Miner            string
	Size             uint64
	Difficulty       *big.Int
	ExtraData        []byte
}

// BlockHeader 区块头信息
type BlockHeader struct {
	ParentHash string
	GasLimit   uint64
	Difficulty *big.Int
}

// Client 以太坊客户端封装
type Client struct {
	Client *ethclient.Client
}

// NewClient 创建新的以太坊客户端
func NewClient(rpcURL string) (*Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("连接以太坊节点失败: %v", err)
	}

	return &Client{Client: client}, nil
}

// GetRawClient 获取原始以太坊客户端
func (c *Client) GetRawClient() *ethclient.Client {
	return c.Client
}

// Close 关闭客户端连接
func (c *Client) Close() {
	if c.Client != nil {
		c.Client.Close()
	}
}

// GetLatestBlock 获取最新区块
func (c *Client) GetLatestBlock() (*Block, error) {
	header, err := c.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("获取最新区块头失败: %v", err)
	}

	block, err := c.Client.BlockByHash(context.Background(), header.Hash())
	if err != nil {
		return nil, fmt.Errorf("获取区块详情失败: %v", err)
	}

	return &Block{
		Number:           block.Number().Uint64(),
		Hash:             block.Hash().Hex(),
		Timestamp:        time.Unix(int64(block.Time()), 0),
		TransactionCount: len(block.Transactions()),
		GasUsed:          block.GasUsed(),
		Miner:            block.Coinbase().Hex(),
		Size:             block.Size(),
		Difficulty:       block.Difficulty(),
		ExtraData:        block.Extra(),
	}, nil
}

// GetBlockByNumber 根据区块号获取区块
func (c *Client) GetBlockByNumber(number uint64) (*Block, error) {
	block, err := c.Client.BlockByNumber(context.Background(), big.NewInt(int64(number)))
	if err != nil {
		return nil, fmt.Errorf("获取区块 %d 失败: %v", number, err)
	}

	return &Block{
		Number:           block.Number().Uint64(),
		Hash:             block.Hash().Hex(),
		Timestamp:        time.Unix(int64(block.Time()), 0),
		TransactionCount: len(block.Transactions()),
		GasUsed:          block.GasUsed(),
		Miner:            block.Coinbase().Hex(),
		Size:             block.Size(),
		Difficulty:       block.Difficulty(),
		ExtraData:        block.Extra(),
	}, nil
}

// GetBlockHeaderByNumber 根据区块号获取区块头
func (c *Client) GetBlockHeaderByNumber(number uint64) (*BlockHeader, error) {
	header, err := c.Client.HeaderByNumber(context.Background(), big.NewInt(int64(number)))
	if err != nil {
		return nil, fmt.Errorf("获取区块头 %d 失败: %v", number, err)
	}

	return &BlockHeader{
		ParentHash: header.ParentHash.Hex(),
		GasLimit:   header.GasLimit,
		Difficulty: header.Difficulty,
	}, nil
}

// GetBlockTransactionCount 获取区块中的交易数量
func (c *Client) GetBlockTransactionCount(number uint64) (int, error) {
	block, err := c.Client.BlockByNumber(context.Background(), big.NewInt(int64(number)))
	if err != nil {
		return 0, fmt.Errorf("获取区块 %d 失败: %v", number, err)
	}

	return len(block.Transactions()), nil
}

// GetBlockByHash 根据区块哈希获取区块
func (c *Client) GetBlockByHash(hash string) (*Block, error) {
	blockHash := common.HexToHash(hash)
	block, err := c.Client.BlockByHash(context.Background(), blockHash)
	if err != nil {
		return nil, fmt.Errorf("获取区块 %s 失败: %v", hash, err)
	}

	return &Block{
		Number:           block.Number().Uint64(),
		Hash:             block.Hash().Hex(),
		Timestamp:        time.Unix(int64(block.Time()), 0),
		TransactionCount: len(block.Transactions()),
		GasUsed:          block.GasUsed(),
		Miner:            block.Coinbase().Hex(),
		Size:             block.Size(),
		Difficulty:       block.Difficulty(),
		ExtraData:        block.Extra(),
	}, nil
}

// GetTransactionByHash 根据交易哈希获取交易信息
func (c *Client) GetTransactionByHash(hash string) (*types.Transaction, bool, error) {
	txHash := common.HexToHash(hash)
	tx, isPending, err := c.Client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		return nil, false, fmt.Errorf("获取交易 %s 失败: %v", hash, err)
	}

	return tx, isPending, nil
}

// GetTransactionReceipt 获取交易收据
func (c *Client) GetTransactionReceipt(hash string) (*types.Receipt, error) {
	txHash := common.HexToHash(hash)
	receipt, err := c.Client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("获取交易收据 %s 失败: %v", hash, err)
	}

	return receipt, nil
}

// GetBalance 获取账户余额
func (c *Client) GetBalance(address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := c.Client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, fmt.Errorf("获取账户余额失败: %v", err)
	}

	return balance, nil
}

// GetNonce 获取账户nonce
func (c *Client) GetNonce(address string) (uint64, error) {
	account := common.HexToAddress(address)
	nonce, err := c.Client.NonceAt(context.Background(), account, nil)
	if err != nil {
		return 0, fmt.Errorf("获取账户nonce失败: %v", err)
	}

	return nonce, nil
}

// GetGasPrice 获取当前gas价格
func (c *Client) GetGasPrice() (*big.Int, error) {
	gasPrice, err := c.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取gas价格失败: %v", err)
	}

	return gasPrice, nil
}