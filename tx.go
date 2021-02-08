package main

import (
	"crypto/sha256"
	"fmt"
	"hash"
)

type Tx struct {
	input  string
	output string
	amount float32
	hash   hash.Hash

	parent *Node
}

func HashTxs(t1 *Tx, t2 *Tx) hash.Hash {
	txData := fmt.Sprintf("%s:%s", t1.hash.Sum(nil), t2.hash.Sum(nil))
	hash := sha256.New()
	hash.Write([]byte(txData))
	return hash
}

func HashTx(tx *Tx) hash.Hash {
	txData := fmt.Sprintf("%s:%s:%f", tx.input, tx.output, tx.amount)
	hash := sha256.New()
	hash.Write([]byte(txData))
	return hash
}

func GenerateTransaction(input, output string, amount float32) *Tx {
	tx := &Tx{
		input:  input,
		output: output,
		amount: amount,
	}

	tx.hash = HashTx(tx)
	return tx
}

func PrintTx(tx *Tx) {
	fmt.Println(fmt.Sprintf("%s -> %s (R$ %f) [%x]", tx.input, tx.output, tx.amount, tx.hash.Sum(nil)))
}
