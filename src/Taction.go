package main

import(
	"crypto"
	ec "crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"math/big"
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

var privateKey = new(ec.PrivateKey)

func generateKeys() {

	privateKey, err := ec.GenerateKey(elliptic.P256(), rand.Reader)


}



type tactionOut struct {

	address string
	amount float64
}

type tactionIn struct {

	prevOutId string
	prevOutIndex int
	signature string

}

type Taction struct {

	id string
	tOut tactionOut
	tIn tactionIn

}


func (t * Taction) SignTaction() (r *big.Int, s *big.Int ){


	hash := sha256.New()


	io.WriteString(hash, t.id)

	r, s, err := ec.Sign(rand.Reader, privateKey, hash.Sum(nil))

	return r,s

}

func (t * Taction)generateTransactionId() {

	temp1 := fmt.Sprintf("%f", t.tIn.prevOutIndex)
	temp2 := fmt.Sprintf("%f", t.tOut.amount)

	t.id = (t.tIn.prevOutId + temp1 + t.tOut.address + temp2)
}



//maybe implement transaction id
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

func (t *Taction) isValid(payerWallet float64) bool {

	if (payerWallet-t.tOut.amount < 0) {
		return false
	}

	return true

}

