package main

type Taction struct {
	PrivateKey1 string //Payer
	PrivateKey2 string //Payee
	Amount      float64
	Timestamp   string
}

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

func (t *Taction) isValid(payerWallet float64) bool {

	if (payerWallet-t.Amount < 0) {
		return false
	}

	return true

}
