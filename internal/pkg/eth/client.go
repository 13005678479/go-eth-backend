package eth

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

// Client 以太坊客户端包装器
type Client struct {
	client *ethclient.Client
	ctx    context.Context
}

// NewClient 创建新的以太坊客户端
func NewClient(network, rpcURL string) (*Client, error) {
	// 如果没有提供 RPC URL，使用默认的
	if rpcURL == "" {
		rpcURL = getDefaultRPCURL(network)
	}

	ctx := context.Background()
	
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %v", err)
	}

	return &Client{
		client: client,
		ctx:    ctx,
	}, nil
}

// GetLatestBlockNumber 获取最新的区块号
func (c *Client) GetLatestBlockNumber() (uint64, error) {
	header, err := c.client.HeaderByNumber(c.ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block: %v", err)
	}

	return header.Number.Uint64(), nil
}

// GetBlockByNumber 根据区块号获取区块信息
func (c *Client) GetBlockByNumber(blockNumber uint64) (*Block, error) {
	number := big.NewInt(int64(blockNumber))
	block, err := c.client.BlockByNumber(c.ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get block %d: %v", blockNumber, err)
	}

	return ParseBlock(block), nil
}

// GetBlockByHash 根据区块哈希获取区块信息
func (c *Client) GetBlockByHash(blockHash string) (*Block, error) {
	hash := common.HexToHash(blockHash)
	block, err := c.client.BlockByHash(c.ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by hash %s: %v", blockHash, err)
	}

	return ParseBlock(block), nil
}

// GetBlockHeaderByNumber 根据区块号获取区块头信息
func (c *Client) GetBlockHeaderByNumber(blockNumber uint64) (*BlockHeader, error) {
	number := big.NewInt(int64(blockNumber))
	header, err := c.client.HeaderByNumber(c.ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get block header %d: %v", blockNumber, err)
	}

	return &BlockHeader{
		Number:     header.Number.Uint64(),
		Hash:       header.Hash().Hex(),
		Timestamp:  time.Unix(int64(header.Time), 0),
		ParentHash: header.ParentHash.Hex(),
		GasLimit:   header.GasLimit,
		GasUsed:    header.GasUsed,
		Miner:      header.Coinbase.Hex(),
	}, nil
}

// GetLatestBlock 获取最新区块信息
func (c *Client) GetLatestBlock() (*Block, error) {
	header, err := c.client.HeaderByNumber(c.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block header: %v", err)
	}

	return c.GetBlockByNumber(header.Number.Uint64())
}

// GetBlockTransactionCount 获取区块中的交易数量
func (c *Client) GetBlockTransactionCount(blockNumber uint64) (int, error) {
	number := big.NewInt(int64(blockNumber))
	txCount, err := c.client.TransactionCount(c.ctx, number)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction count for block %d: %v", blockNumber, err)
	}
	return int(txCount), nil
}

// GetBalance 获取地址余额
func (c *Client) GetBalance(address string) (*big.Int, error) {
	return c.client.BalanceAt(c.ctx, parseAddress(address), nil)
}

// getDefaultRPCURL 根据网络获取默认的 RPC URL
func getDefaultRPCURL(network string) string {
	switch network {
	case "mainnet":
		return "https://mainnet.infura.io/v3/YOUR-PROJECT-ID"
	case "sepolia":
		return "https://sepolia.infura.io/v3/ea33fc8cbc4545d9ac08fba394c5046b"
	case "goerli":
		return "https://goerli.infura.io/v3/YOUR-PROJECT-ID"
	case "ropsten":
		return "https://ropsten.infura.io/v3/YOUR-PROJECT-ID"
	case "rinkeby":
		return "https://rinkeby.infura.io/v3/YOUR-PROJECT-ID"
	default:
		// 默认使用 Sepolia 测试网络
		return "https://sepolia.infura.io/v3/ea33fc8cbc4545d9ac08fba394c5046b"
	}
}

// parseAddress 解析以太坊地址
func parseAddress(address string) string {
	// 这里可以添加地址验证逻辑
	return address
}