package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Blockchain struct {
	Chain               []*Block
	CurrentTransactions []*Transaction
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Chain: []*Block{
			{
				Index:        0,
				Timestamp:    0,
				Transactions: nil,
				Proof:        100,
				PreviousHash: "1",
			},
		},
		CurrentTransactions: nil,
	}
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
	guessHex := hex.EncodeToString(guessHash[:])
	return string(guessHex[:4]) == "0000"
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

func (b *Blockchain) MineHandler(w http.ResponseWriter, r *http.Request) {
	lastBlock := b.LastBlock()
	lastProof := lastBlock.Proof
	proof := b.ProofOfWork(lastProof)

	// We must receive a reward for finding the proof.
	//The sender is "0" to signify that this node has mined a new coin.

	b.NewTransaction("0", "NodeIdentifier", 1)
	previousHash := Hash(*lastBlock)
	newBlock := b.NewBlock(proof, previousHash)
	w.Write([]byte(fmt.Sprintf("New block forged, block index %d, transactions %v proof %d, previousHash %s \n",
		newBlock.Index, newBlock.Transactions, newBlock.Proof, newBlock.PreviousHash)))
}

func (b *Blockchain) NewTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction
	data, error := io.ReadAll(r.Body)
	if error != nil {
		return
	}
	json.Unmarshal(data, &transaction)
	if transaction.Sender == "" || transaction.Recipient == "" {
		w.Write([]byte("field missing"))
	}
	index := b.NewTransaction(transaction.Sender, transaction.Recipient, transaction.Amount)
	w.Write([]byte(fmt.Sprintf("added transaction to Block %d\n", index)))
}

func (b *Blockchain) FullChainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("length %d\nchain: %v\n", len(b.Chain), b.Chain)))
}

func main() {
	blockChain := NewBlockchain()
	http.HandleFunc("/mine", blockChain.MineHandler)
	http.HandleFunc("/transactions/new", blockChain.NewTransactionHandler)
	http.HandleFunc("/chain", blockChain.FullChainHandler)
	http.ListenAndServe(":8080", nil)
}
