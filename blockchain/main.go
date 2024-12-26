package main

import (
	"blockchain/block"
	"blockchain/pow"
	"fmt"
)

func main() {
	// 创建新的区块链
	bc := block.NewBlockchain()
	defer bc.Close()

	// 添加一些区块
	transactions := []string{
		"Alice 发送 1 BTC 给 Bob",
		"Bob 发送 2 BTC 给 Charlie",
		"Charlie 发送 0.5 BTC 给 David",
	}

	for _, tx := range transactions {
		// 创建新区块
		prevBlock := bc.GetLastBlock()
		newBlock := block.NewBlock(tx, prevBlock.Hash)

		// 计算工作量证明
		pow := pow.NewProofOfWork(newBlock)
		nonce, hash := pow.Run()

		// 更新区块
		newBlock.Nonce = nonce
		newBlock.Hash = hash

		// 将区块添加到链中
		bc.AddBlock(newBlock)
		fmt.Printf("添加新区块: %x\n", newBlock.Hash)
	}

	fmt.Println("区块链创建成功！")
}
