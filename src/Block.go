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
	Transactions []Taction
	Hash         string
	PrevHash     string
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

	hashString := strconv.Itoa(b.Index) + b.Timestamp

	for _, tact := range b.Transactions {

		hashString += tact.id
	}

	hashString += b.PrevHash + b.Nonce

	// I looked it up, this is one of two built-in hash functions in the go standard crypto lib.
	// This one seems good.
	hasher := sha256.New()

	hasher.Write([]byte(hashString))

	hashed := hasher.Sum(nil)

	return hex.EncodeToString(hashed)

}

func (b *Block) hasTransaction(t Taction) bool {

	for _, bTaction := range b.Transactions {
		if bTaction.equals(t) {
			return true
		}
	}

	return false
}

func (b *Block) equals(otherBlock Block) bool {

	if b.Index != otherBlock.Index {
		return false
	}

	if b.Timestamp != otherBlock.Timestamp {
		return false
	}

	if b.PrevHash != otherBlock.PrevHash {
		return false
	}

	if b.Hash != otherBlock.Hash {
		return false
	}

	if b.Difficulty != otherBlock.Difficulty {
		return false
	}

	if b.Nonce != otherBlock.Nonce {
		return false
	}

	for _, bTaction := range b.Transactions {

		if !otherBlock.hasTransaction(bTaction) {
			return false
		}
	}

	return true
}

func generateBlock(oldBlock Block, tactions []Taction, difficulty int) Block {
	var newBlock Block

	t := time.Now()

	//privKey := getPrivateKey(publicKey)

	newTOut := tactionOut{
		getThisPublicKey(),
		1}

	newTIn := tactionIn{
		"",
		"",
		0,
		"" }


	newTrans := Taction{
		"",
		newTOut,
		newTIn }

	(*Taction).generateTransactionId(&newTrans)
	newTrans.tIn.signature = (*Taction).SignTaction(&newTrans)

	coinbaseTaction := newTrans

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Transactions = append([]Taction{coinbaseTaction}, tactions...)
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty

	for i := 0; ; i++ {

		hexVal := fmt.Sprintf("%x", i)
		newBlock.Nonce = hexVal

		hash := newBlock.calculateHash()
		if !isHashValid(hash, newBlock.Difficulty) {

			// if someone else has beaten us to this block, make a new block with this data.

			if !oldBlock.equals(Blockchain[len(Blockchain)-1]) {

				missingTactions := filterCommittedTactions(tactions)

				pendingTransactions = append(pendingTransactions, missingTactions...)

				return generateBlock(Blockchain[len(Blockchain)-1], pendingTransactions, difficulty)
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

/*
func (b *Block) getWalletAmt(privKey string) float64 {

	var wallet float64

	wallet = 0

	for _, trans := range b.Transactions {
		if trans.PrivateKey2 == privKey {
			wallet += trans.Amount
		}

		if trans.PrivateKey1 == privKey {
			wallet -= trans.Amount
		}

	}

	return wallet

}
*/
