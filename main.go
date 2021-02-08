package main

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func genStr(l int) string {
	buff := make([]byte, l)
	cryptoRand.Read(buff)
	str := base64.StdEncoding.EncodeToString(buff)
	return str[:l]
}

func genAmount() float32 {
	return rand.Float32()
}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	fmt.Println("====== Generate transactions ======")

	txs := make([]*Tx, 0)

	for i := 0; i < 5; i++ {
		tx := GenerateTransaction(genStr(10), genStr(10), genAmount())
		PrintTx(tx)
		txs = append(txs, tx)
	}

	fmt.Println("\n\n====== Put all transactions on the tree ======")

	root, err := GenerateTree(txs)

	if err != nil {
		log.Fatalf("Problems while generating tree: %v", err)
	}

	fmt.Printf("Tree root: %x\n", root.hash.Sum(nil))

	fmt.Println("\n\n====== Print all transactions ======")

	PrintTree(root)

	fmt.Println("\n\n====== Validate transactions tree ======")

	for _, tx := range txs {
		fmt.Println(fmt.Sprintf("Validate transaction: %s -> %s (R$ %f) [%x]", tx.input, tx.output, tx.amount, tx.hash.Sum(nil)))

		if !root.Validate(tx.hash) {
			fmt.Println("WARNING - The tree is tampered!")
			continue
		}
	}

	fmt.Println("\n\n====== Change one transactions ======")

	txs[3].output = "Tome now"
	txs[3].amount = 20000.00

	txs[0].output = "Tome now"
	txs[0].amount = 20000.00

	fmt.Println("\n\n====== Validate transactions tree ======")

	for _, tx := range txs {
		fmt.Println(fmt.Sprintf("Validate transaction: %s -> %s (R$ %f) [%x]", tx.input, tx.output, tx.amount, tx.hash.Sum(nil)))

		if !root.Validate(tx.hash) {
			fmt.Println("WARNING - The tree is tampered!")
			continue
		}
	}
}
