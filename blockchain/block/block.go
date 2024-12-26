package block

import (
	"blockchain/pow"
	"blockchain/types"
	"time"
)

func NewBlock(data string, prevBlockHash []byte) *types.Block {
	block := &types.Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
	}

	// 创建区块时就计算工作量证明
	powInstance := pow.NewProofOfWork(block)
	nonce, hash := powInstance.Run()

	block.Nonce = nonce
	block.Hash = hash

	return block
}

func NewGenesisBlock() *types.Block {
	return NewBlock("Genesis Block", []byte{})
}
