package eth

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Block 区块信息结构体
type Block struct {
	Number          uint64        `json:"number"`
	Hash            string        `json:"hash"`
	Timestamp       time.Time     `json:"timestamp"`
	ParentHash      string        `json:"parentHash"`
	Difficulty      string        `json:"difficulty"`
	GasLimit        uint64        `json:"gasLimit"`
	GasUsed         uint64        `json:"gasUsed"`
	Miner           string        `json:"miner"`
	ExtraData       string        `json:"extraData"`
	Transactions    []string      `json:"transactions"`
	TransactionCount int          `json:"transactionCount"`
	Size            uint64        `json:"size"`
}

// ParseBlock 从原生区块类型解析为自定义区块结构
func ParseBlock(block *types.Block) *Block {
	transactions := make([]string, len(block.Transactions()))
	for i, tx := range block.Transactions() {
		transactions[i] = tx.Hash().Hex()
	}

	return &Block{
		Number:          block.Number().Uint64(),
		Hash:            block.Hash().Hex(),
		Timestamp:       time.Unix(int64(block.Time()), 0),
		ParentHash:      block.ParentHash().Hex(),
		Difficulty:      block.Difficulty().String(),
		GasLimit:        block.GasLimit(),
		GasUsed:         block.GasUsed(),
		Miner:           block.Coinbase().Hex(),
		ExtraData:       fmt.Sprintf("%x", block.Extra()),
		Transactions:    transactions,
		TransactionCount: len(transactions),
		Size:            block.Size(),
	}
}

// BlockHeader 区块头信息结构体
type BlockHeader struct {
	Number        uint64    `json:"number"`
	Hash          string    `json:"hash"`
	Timestamp     time.Time `json:"timestamp"`
	ParentHash    string    `json:"parentHash"`
	GasLimit      uint64    `json:"gasLimit"`
	GasUsed       uint64    `json:"gasUsed"`
	Miner         string    `json:"miner"`
}