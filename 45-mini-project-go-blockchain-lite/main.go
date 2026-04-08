package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Difficulty level for Proof of Work (number of leading zeros required)
const difficulty = 4

// Block represents each item in the blockchain
type Block struct {
	Index        int
	Timestamp    string
	Data         string
	PreviousHash string
	Hash         string
	Nonce        int
}

// Blockchain is a collection of blocks
type Blockchain struct {
	Blocks []Block
}

// calculateHash generates a SHA-256 hash for a block
func calculateHash(b Block) string {
	record := strconv.Itoa(b.Index) + b.Timestamp + b.Data + b.PreviousHash + strconv.Itoa(b.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// isHashValid checks if a hash meets the difficulty requirement
func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

// generateBlock creates a new block and mines it (Proof of Work)
func (bc *Blockchain) generateBlock(data string) Block {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().String(),
		Data:         data,
		PreviousHash: prevBlock.Hash,
		Nonce:        0,
	}

	fmt.Printf("⛏️  Mining block %d...\n", newBlock.Index)
	for {
		newBlock.Hash = calculateHash(newBlock)
		if isHashValid(newBlock.Hash, difficulty) {
			break
		}
		newBlock.Nonce++
	}

	fmt.Printf("✅ Block mined! Hash: %s\n", newBlock.Hash)
	return newBlock
}

// isChainValid verifies the integrity of the blockchain
func (bc *Blockchain) isChainValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		// Check hash integrity
		if currentBlock.Hash != calculateHash(currentBlock) {
			fmt.Printf("❌ Block %d has invalid hash!\n", i)
			return false
		}

		// Check link consistency
		if currentBlock.PreviousHash != prevBlock.Hash {
			fmt.Printf("❌ Block %d link inconsistency!\n", i)
			return false
		}
	}
	return true
}

func main() {
	// Initialize blockchain with the Genesis block
	genesisBlock := Block{
		Index:        0,
		Timestamp:    time.Now().String(),
		Data:         "Genesis Block",
		PreviousHash: "",
		Nonce:        0,
	}
	genesisBlock.Hash = calculateHash(genesisBlock)

	bc := &Blockchain{
		Blocks: []Block{genesisBlock},
	}

	fmt.Println("🚀 Go-Blockchain-Lite Started!")
	fmt.Printf("Block #0: %s (Genesis)\n", genesisBlock.Hash)

	// Add new blocks
	bc.Blocks = append(bc.Blocks, bc.generateBlock("Transaction: Alice sends 5 BTC to Bob"))
	bc.Blocks = append(bc.Blocks, bc.generateBlock("Transaction: Bob sends 2 BTC to Charlie"))

	fmt.Println("\n--- Blockchain Status ---")
	if bc.isChainValid() {
		fmt.Println("🔗 Blockchain is VALID and SECURE.")
	} else {
		fmt.Println("⚠️  Blockchain is CORRUPTED!")
	}

	// Demo: Tampering with data
	fmt.Println("\n🔍 Testing Tamper Detection...")
	fmt.Println("Modifying transaction data in block 1...")
	bc.Blocks[1].Data = "Transaction: Alice sends 500 BTC to Bob" // Tamper!

	if bc.isChainValid() {
		fmt.Println("🔗 Blockchain is still valid (Something is wrong!)")
	} else {
		fmt.Println("⚠️  CRITICAL: Tamper detected! Blockchain is invalid now.")
	}
}
