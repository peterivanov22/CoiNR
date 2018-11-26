package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)


// Block struct

type Block struct {

	Index        int
	Timestamp    string
	BPM          int
	Hash         string
	PrevHash string
	Difficulty   int
	//in real bitcoin, nonce is 4 bytes
	//string in golang is pointer (of size 8 bytes)
	Nonce string
}

// validates the block.  Usage is block.validate()

func (b *Block) validate(previous *Block) bool {

	if previous.Index+1 != b.Index {
		return false
	}

	if previous.Hash != b.PrevHash {
		return false
	}

	if b.calculateHash() != b.Hash {
		return false
	}

	return true

}

//calculates the hash for the block.  Usage is block.calculateHash()

func (b *Block) calculateHash() string {

	hashString := strconv.Itoa(b.Index) + b.Timestamp + strconv.Itoa(b.BPM) + b.PrevHash + b.Nonce


	// I looked it up, this is one of two built-in hash functions in the go standard crypto lib.
	// This one seems good.
	hasher := sha256.New()

	hasher.Write([]byte(hashString))

	hashed := hasher.Sum(nil)

	return hex.EncodeToString(hashed)

}


func generateBlock(oldBlock Block, BPM int, difficulty int) Block {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty


	for i := 0; ; i++ {

		hexVal := fmt.Sprintf("%x", i)
		newBlock.Nonce = hexVal

		hash := newBlock.calculateHash()
		if !isHashValid(hash, newBlock.Difficulty) {

			// if someone else has beaten us to this block, make a new block with this data.
			if oldBlock != Blockchain[len(Blockchain) -1] {
				return generateBlock(Blockchain[len(Blockchain)-1], BPM, difficulty)
			}

			fmt.Println(hash, " do more work!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(hash, " work done!")
			newBlock.Hash = hash
			break
		}

	}
	return newBlock
}

//check if hash has correct number of zeros
func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}


