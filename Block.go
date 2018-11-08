package CoiNR

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

// Block struct

type Block struct {

	Index int
	Timestamp string
	BPM int
	Hash string
	PreviousHash string
}


// validates the block.  Usage is block.validate()

func (b *Block) validate(previous *Block) bool {

	if previous.Index +1 != b.Index {
		return false
	}

	if previous.Hash != b.PreviousHash {
		return false
	}

	if b.calculatehash() != b.Hash {
		return false
	}

	return true

}

//calculates the hash for the block.  Usage is block.calculatehash()

func (b *Block)calculatehash() string {

	hashString := strconv.Itoa(b.Index) + b.Timestamp + strconv.Itoa(b.BPM) + b.PreviousHash

	// I looked it up, this is one of two built-in hash functions in the go standard crypto lib.
	// This one seems good.
	hasher := sha256.New()

	hasher.Write([]byte(hashString))

	hashed := hasher.Sum(nil)

	return hex.EncodeToString(hashed)

}

// generates a new block from the previous block and the BPM.  Usage is var newBlock = generateBlock(*oldBlock, bpm)

func generateBlock(prev *Block, BPM int) Block {

	var newBlock Block

	t := time.Now()

	newBlock.Index = prev.Index +1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PreviousHash = prev.Hash
	newBlock.Hash = newBlock.calculatehash()

	return newBlock

}