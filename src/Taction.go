package main

import (
	ec "crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
)

/*
type Taction struct {
	PrivateKey1 string //Payer
	PrivateKey2 string //Payee
	Amount      float64
	Timestamp   string
}
*/

//yeah.. ill do something more secure later
//const privatekey = "5DB0633DDD43355F97CC15A190660A7453737BE343E503F8882D62D6C927C6DA"

var block_reward = 1

func generateKeys() * ec.PrivateKey {

	privateKey,_ := ec.GenerateKey(elliptic.P256(), rand.Reader)
	return privateKey

}



func getPublicKey (key *ec.PrivateKey)  string {
	pubkey := key.PublicKey.X.String() + key.PublicKey.Y.String()
	return pubkey
}

func getThisPublicKey ()  string {
	pubkey := privateKey.PublicKey.X.String() + privateKey.PublicKey.Y.String()
	return pubkey
}

type tactionOut struct {

	Address string
	Amount  float64
}

type tactionIn struct {

	Address      string
	PrevOutId    string
	PrevOutIndex int
	Signature    string

}

type Taction struct {

	ID   string
	TOut tactionOut
	TIn  tactionIn

}


func (t * Taction) SignTaction() (rs string ){

	hash := sha256.New()
	io.WriteString(hash, t.ID)

	r, _, err := ec.Sign(rand.Reader, privateKey, hash.Sum(nil))

	if err != nil {
		log.Println("bad sign")
		return ""
	}


	return r.String()

}

func (t * Taction)generateTransactionId() {

	temp1 := fmt.Sprintf("%f", t.TIn.PrevOutIndex)
	temp2 := fmt.Sprintf("%f", t.TOut.Amount)

	t.ID = (t.TIn.PrevOutId + temp1 + t.TOut.Address + temp2)
}

type availableCoin struct {

	tactionOutId string
	tactionOutIndex int
	address string
	amount float64
}

var availableCoins []availableCoin


func (	B * Block) updateNewOwnership () {

	for i:=0 ; i< len(B.Transactions); i++{
		newCoin := availableCoin{B.Transactions[i].ID,B.Transactions[i].TIn.PrevOutIndex,
		B.Transactions[i].TOut.Address, B.Transactions[i].TOut.Amount}
		availableCoins = append(availableCoins,newCoin)
	}
}

func (B* Block) deleteOldOwnership () {

	for i:=0 ; i< len(B.Transactions); i++ {

		temp_sum := 0.0


		//newCoin := availableCoin{B.Transactions[i].ID,B.Transactions[i].TOut,
			//B.Transactions[i].TOut.Address, B.Transactions[i].TOut.Amount}

		for j:=0 ; j< len(availableCoins); j++ {

			//so far just assuming everything is in increments of 1
			if (availableCoins[j].address == B.Transactions[i].TIn.Address){
				temp_sum += availableCoins[j].amount
				//dirty way to delete this for now
				availableCoins[j].address = "-1"
			}

			if (temp_sum == B.Transactions[i].TOut.Amount){
				break
			}
			//implement finding right owner and taking away right Amount of coins,
			//how to consolidate multiple availablecoin for same owners
			//can implement validation
		}

	}
}

func  findLastUnspent (address string) availableCoin {



		//newCoin := availableCoin{B.Transactions[i].ID,B.Transactions[i].TOut,
		//B.Transactions[i].TOut.Address, B.Transactions[i].TOut.Amount}

		for j:=0 ; j< len(availableCoins); j++ {

			//so far just assuming everything is in increments of 1
			if (availableCoins[j].address == address){
				//dirty way to delete this for now
				return availableCoins[j]
			}

		}
		return availableCoin{}

}


func  showBalance (address string)  {


	var temp = 0.0
	//newCoin := availableCoin{B.Transactions[i].ID,B.Transactions[i].TOut,
	//B.Transactions[i].TOut.Address, B.Transactions[i].TOut.Amount}

	for j:=0 ; j< len(availableCoins); j++ {

		//so far just assuming everything is in increments of 1
		if (availableCoins[j].address == address){
			//dirty way to delete this for now
			temp+= availableCoins[j].amount
		}

	}
	println("This node has: " , temp)

}


//maybe implement transaction ID
//needed for checks

/*
func (t *Taction) equals(otherTaction Taction) bool {

	if t.PrivateKey1 != otherTaction.PrivateKey1 {
		return false
	}

	if t.PrivateKey2 != otherTaction.PrivateKey2 {
		return false
	}

	if t.Amount != otherTaction.Amount {
		return false
	}

	if t.Timestamp != otherTaction.Timestamp {
		return false
	}

	return true
}
*/
func (t *Taction) equals(otherTaction Taction) bool {

	if t.ID != otherTaction.ID {
		return false
	}


	return true
}


func (t *Taction) isValid(payerWallet float64) bool {

	if (payerWallet-t.TOut.Amount < 0) {
		return false
	}

	return true

}

