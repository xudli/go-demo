package block

import (
	"blockchain/types"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

type Blockchain struct {
	tip []byte   // 最新区块的哈希
	db  *bolt.DB // 数据库连接
}

// 创建新的区块链
func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("没有找到现有区块链，创建一个新的...")
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// 添加区块到区块链
func (bc *Blockchain) AddBlock(block *types.Block) {
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockData := block.Serialize()

		err := b.Put(block.Hash, blockData)
		if err != nil {
			return err
		}

		err = b.Put([]byte("l"), block.Hash)
		if err != nil {
			return err
		}

		bc.tip = block.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// 获取最新区块
func (bc *Blockchain) GetLastBlock() *types.Block {
	var block *types.Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockData := b.Get(bc.tip)
		block = types.DeserializeBlock(blockData)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return block
}

// 关闭数据库
func (bc *Blockchain) Close() {
	bc.db.Close()
}
