package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

type Blockchain struct {
	Chain               []*Block
	CurrentTransactions []*Transaction
	Nodes               map[string]struct{}
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Chain: []*Block{
			{
				Index:        1,
				Timestamp:    0,
				Transactions: nil,
				Proof:        100,
				PreviousHash: "1",
			},
		},
		CurrentTransactions: nil,
		Nodes:               make(map[string]struct{}),
	}
}

// RegisterNodes Add a new node to the list of nodes
func (b *Blockchain) RegisterNodes(addr string) {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err)
	}
	b.Nodes[addr] = struct{}{}
}

// ValidChain Determine if a given blockchain is valid
func (b *Blockchain) ValidChain(chain []*Block) bool {
	lastBlock := chain[0]
	currentIndex := 1
	for currentIndex < len(chain) {
		block := chain[currentIndex]

		// Check that the hash of the block is correct
		if block.PreviousHash != Hash(*lastBlock) {
			return false
		}
		// Check that the Proof of Work is correct
		if !ValidProof(lastBlock.Proof, block.Proof) {
			return false
		}
		lastBlock = block
		currentIndex += 1
	}

	return true
}

type NodeResponse struct {
	Length int
	Chain  []*Block
}

// ResolveConflicts
// This is our Consensus Algorithm, it resolves conflicts
// by replacing our chain with the longest one in the network.
func (b *Blockchain) ResolveConflicts() bool {
	var newChain []*Block
	neighbors := b.Nodes

	// We're only looking for chains longer than ours
	maxLength := len(b.Chain)

	for key, _ := range neighbors {
		var respdata NodeResponse

		url := fmt.Sprintf("http://%s/chain", key)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			os.Exit(1)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("client: error making http request: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("client: status code: %d\n", res.StatusCode)

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
			os.Exit(1)
		}

		json.Unmarshal(resBody, &respdata)
		fmt.Println(respdata)

		if respdata.Length > maxLength && b.ValidChain(respdata.Chain) {
			maxLength = respdata.Length
			newChain = respdata.Chain
		}
	}

	if len(newChain) > 0 {
		b.Chain = newChain
		return true
	}

	return false
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

type RegisterNodeResponse struct {
	Node string
}

func (b *Blockchain) RegisterNodeHandler(w http.ResponseWriter, r *http.Request) {
	var resp RegisterNodeResponse
	data, error := io.ReadAll(r.Body)
	if error != nil {
		w.Write([]byte("readall error\n"))
		return
	}

	json.Unmarshal(data, &resp)

	b.RegisterNodes(resp.Node)
	w.Write([]byte(fmt.Sprintf("New nodes have been added. Total nodes: %v\n", b.Nodes)))
}

func (b *Blockchain) ResolveHandler(w http.ResponseWriter, r *http.Request) {
	replaced := b.ResolveConflicts()
	if replaced {
		w.Write([]byte(fmt.Sprintf("chain was replaced, %v\n", b.Chain)))
		return
	}
	w.Write([]byte(fmt.Sprintf("our chain is authoritative %v\n", b.Chain)))
}

func main() {
	port := os.Args[1] // :8080

	blockChain := NewBlockchain()
	http.HandleFunc("/mine", blockChain.MineHandler)
	http.HandleFunc("/transactions/new", blockChain.NewTransactionHandler)
	http.HandleFunc("/chain", blockChain.FullChainHandler)
	http.HandleFunc("/nodes/register", blockChain.RegisterNodeHandler)
	http.HandleFunc("/nodes/resolve", blockChain.ResolveHandler)
	fmt.Printf("running on port %s\n", port)
	http.ListenAndServe(port, nil)
}
