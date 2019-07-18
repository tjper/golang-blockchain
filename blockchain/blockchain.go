package blockchain

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/dgraph-io/badger"
)

const dbPath = "./tmp/blocks"

var lastHashKey = []byte("lh")

type BlockChain struct {
	Database *badger.DB
	LastHash []byte
}

func New(options ...BlockChainOption) *BlockChain {
	var bc = new(BlockChain)
	for _, option := range options {
		option(bc)
	}
	return bc
}

type BlockChainOption func(*BlockChain)

func WithBadgerDB(db *badger.DB) BlockChainOption {
	return func(bc *BlockChain) {
		bc.Database = db
	}
}

func (bc *BlockChain) Init() {
	if err := bc.Database.Update(func(tx *badger.Txn) error {
		item, err := tx.Get(lastHashKey)
		if err == badger.ErrKeyNotFound {
			bc.AddBlock("Genesis")
			return nil
		}
		if err != nil {
			return err
		}

		if err := item.Value(func(val []byte) error {
			bc.LastHash = val
			return nil
		}); err != nil {
			log.Panic(err)
		}
		return nil
	}); err != nil {
		log.Panic(err)
	}
}

func (bc *BlockChain) AddBlock(data string) {
	var newBlock = CreateBlock(data, bc.LastHash)

	if err := bc.Database.Update(func(tx *badger.Txn) error {
		if err := tx.Set(newBlock.Hash, newBlock.Serialize()); err != nil {
			return err
		}
		if err := tx.Set(lastHashKey, newBlock.Hash); err != nil {
			return err
		}
		bc.LastHash = newBlock.Hash
		return nil
	}); err != nil {
		log.Panic(err)
	}
}

func (bc *BlockChain) Print() {
	if err := bc.Database.View(func(tx *badger.Txn) error {
		item, err := tx.Get(lastHashKey)
		if err != nil {
			return err
		}
		var lastHash []byte
		if err := item.Value(func(val []byte) error {
			lastHash = val
			return nil
		}); err != nil {
			return err
		}

		var it = tx.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		var key = lastHash
		for {
			it.Seek(key)
			if err := it.Item().Value(func(val []byte) error {
				var block = Deserialize(val)
				key = block.PrevHash

				fmt.Printf("Previous Hash: %x\n", block.PrevHash)
				fmt.Printf("Data in Block: %s\n", block.Data)
				fmt.Printf("Hash: %x\n", block.Hash)

				var pow = NewProof(block)
				fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
				fmt.Println()

				return nil
			}); err != nil {
				return err
			}
			if bytes.Equal(key, []byte{}) {
				break
			}
		}
		return nil
	}); err != nil {
		log.Panic(err)
	}
}
