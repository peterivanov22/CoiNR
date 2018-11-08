package CoiNR

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type Block struct {

	Index int
	Timestamp string
	BPM int
	Hash string
	PreviousHash string
}


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

func (b *Block)calculatehash() string {

	hashString := strconv.Itoa(b.Index) + b.Timestamp + strconv.Itoa(b.BPM) + b.PreviousHash

	hasher := sha256.New()

	hasher.Write([]byte(hashString))

	hashed := hasher.Sum(nil)

	return hex.EncodeToString(hashed)

}

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