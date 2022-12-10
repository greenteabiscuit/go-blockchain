package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Blockchain struct {
	Chain               []*Block
	CurrentTransactions []*Transaction
}

// NewBlock Adds a new block and adds it to the chain
func (b *Blockchain) NewBlock(proof int, previousHash string) *Block {
	block := &Block{
		Index:        len(b.Chain) + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: b.CurrentTransactions,
		Proof:        proof,
		PreviousHash: previousHash,
	}
	b.CurrentTransactions = []*Transaction{}
	b.Chain = append(b.Chain, block)
	return block
}

// NewTransaction Adds a new transaction to the list of transactions
func (b *Blockchain) NewTransaction(sender, recipient string, amount int) int {
	b.CurrentTransactions = append(b.CurrentTransactions, &Transaction{sender, recipient, amount})
	return b.LastBlock().Index
}

// LastBlock Returns the last Block in the chain
func (b *Blockchain) LastBlock() *Block {
	return b.Chain[len(b.Chain)-1]
}

// Hash creates a SHA-256 hash of a block
func Hash(block Block) string {
	blockString, err := json.Marshal(block)
	if err != nil {
		return ""
	}

	binarySha256 := sha256.Sum256(blockString)

	return hex.EncodeToString(binarySha256[:])
}

type Block struct {
	Index        int
	Timestamp    int64
	Transactions []*Transaction
	Proof        int
	PreviousHash string
}

type Transaction struct {
	Sender    string
	Recipient string
	Amount    int
}
