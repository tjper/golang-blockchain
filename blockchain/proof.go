package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const Difficulty = 20

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) Init(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var (
		intHash big.Int
		hash    [32]byte
		nonce   = 0
	)
	for nonce < math.MaxInt64 {
		var data = pow.Init(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		}
		nonce++
	}
	fmt.Println()

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	var (
		data = pow.Init(pow.Block.Nonce)
		hash = sha256.Sum256(data)
	)

	intHash.SetBytes(hash[:])
	return intHash.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte {
	var buf = new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, num); err != nil {
		log.Panic(err)
	}
	return buf.Bytes()
}
