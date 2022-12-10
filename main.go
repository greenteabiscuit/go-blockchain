package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

// ProofOfWork Simple Proof of Work Algorithm:
//   - Find a number p' such that hash(pp') contains leading 4 zeroes, where p is the previous p'
//   - p is the previous proof, and p' is the new proof
func (b *Blockchain) ProofOfWork(lastProof int) int {
	proof := 0
	for {
		if ValidProof(lastProof, proof) {
			break
		}
		proof += 1
	}
	return proof
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

// ValidProof Validates the Proof: Does hash(last_proof, proof) contain 4 leading zeroes?
func ValidProof(lastProof, proof int) bool {
	guess := fmt.Sprintf("%d%d", lastProof, proof)
	guessHash := sha256.Sum256([]byte(guess))
	return string(guessHash[:4]) == "0000"
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
