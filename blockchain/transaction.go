package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	Value  int
	PubKey string
}

type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer

	if err := gob.NewEncoder(&encoded).Encode(tx); err != nil {
		log.Panic(err)
	}

	var hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	var (
		txIn = TxInput{
			ID:  []byte{},
			Out: -1,
			Sig: data,
		}
		txOut = TxOutput{
			Value:  100,
			PubKey: to,
		}
		tx = Transaction{
			ID:      nil,
			Inputs:  []TxInput{txIn},
			Outputs: []TxOutput{txOut},
		}
	)
	tx.SetID()

	return &tx
}

func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var (
		inputs  []TxInput
		outputs []TxOutput
	)
}
