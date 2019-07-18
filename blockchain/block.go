package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{
		Hash:     []byte{},
		Data:     []byte(data),
		PrevHash: prevHash,
		Nonce:    0,
	}

	var pow = NewProof(block)
	var nonce, hash = pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func (b *Block) Serialize() []byte {
	var buf = new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(b); err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&block); err != nil {
		log.Panic(err)
	}
	return &block
}
