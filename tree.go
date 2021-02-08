package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
)

type Node struct {
	hash hash.Hash

	parent *Node

	leftNode  *Node
	rigthNode *Node

	leftTx  *Tx
	rightTx *Tx
}

func (n *Node) Validate(txHash hash.Hash) bool {
	if txHash == nil {
		return false
	}

	tx := FindTx(n, txHash)

	if tx == nil {
		fmt.Printf("\t - Could not found transaction: %x\n", txHash.Sum(nil))
		return false
	}

	fmt.Printf("\t - Transaction found: %x\n", tx.hash.Sum(nil))

	// Após encontrar a transação a partir da raiz da árvore
	// é necessário validar se ela não está alterada, para isso
	// vamos comparar ela com a "irmã" dela (que pode estar a direita ou a esquerda)
	// do nó pai
	level := 1
	intermediaryNode := tx.parent
	var intermediaryHash hash.Hash

	compareTx := GenerateTransaction(tx.input, tx.output, tx.amount)

	if intermediaryNode.leftTx == tx {
		if intermediaryNode.rightTx != nil {
			intermediaryHash = HashTxs(compareTx, intermediaryNode.rightTx)
		} else {
			intermediaryHash = HashTx(compareTx)
		}
	} else {
		intermediaryHash = HashTxs(intermediaryNode.leftTx, compareTx)
	}

	if bytes.Compare(intermediaryHash.Sum(nil), intermediaryNode.hash.Sum(nil)) != 0 {
		fmt.Printf("Could not evaluate the transaction %x\n", tx.hash.Sum(nil))
		return false
	}

	fmt.Printf("\t- Evaluation %v done!: %x\n", level, intermediaryNode.hash.Sum(nil))

	// após a validação da transação com sua "irmã"
	// nós vamos proseguir com os nós sucessores até chegar à raiz
	for {
		level++

		if intermediaryNode.parent == nil {
			break
		}

		intermediaryParent := intermediaryNode.parent

		if intermediaryParent.leftNode == intermediaryNode {
			if intermediaryParent.rigthNode != nil {
				intermediaryHash = HashNodes(intermediaryNode, intermediaryParent.rigthNode)
			}
		} else {
			intermediaryHash = HashNodes(intermediaryParent.leftNode, intermediaryParent.rigthNode)
		}

		if bytes.Compare(intermediaryHash.Sum(nil), intermediaryParent.hash.Sum(nil)) != 0 {
			fmt.Printf("Could not evaluate the transaction %x\n", tx.hash.Sum(nil))
			return false
		}

		fmt.Printf("\t- Evaluation %v done!: %x\n", level, intermediaryParent.hash.Sum(nil))

		intermediaryNode = intermediaryParent
	}

	return true
}

func FindTx(n *Node, h hash.Hash) *Tx {
	if n.leftNode == nil && n.rigthNode == nil {
		if n.leftTx.hash == h {
			return n.leftTx
		}

		if n.rightTx.hash == h {
			return n.rightTx
		}
	}

	if n.leftNode != nil {
		leftSide := FindTx(n.leftNode, h)

		if leftSide != nil {
			return leftSide
		}
	}

	if n.rigthNode != nil {
		return FindTx(n.rigthNode, h)
	}

	return nil
}

func HashNodes(n1, n2 *Node) hash.Hash {
	data := fmt.Sprintf("%s:%s", n1.hash.Sum(nil), n2.hash.Sum(nil))
	hash := sha256.New()
	hash.Write([]byte(data))
	return hash
}

func GenereteSingleNode(n1 *Node) *Node {
	n := &Node{
		leftNode: n1,
		hash:     n1.hash,

		rigthNode: nil,
		leftTx:    nil,
		rightTx:   nil,
	}

	n1.parent = n
	return n
}

func GenerateSingleNodeFromTx(tx *Tx) *Node {
	n := &Node{
		leftTx: tx,
		hash:   tx.hash,

		rightTx:   nil,
		rigthNode: nil,
		leftNode:  nil,
	}

	tx.parent = n
	return n
}

func GenerateNode(n1, n2 *Node) *Node {
	n := &Node{
		leftNode:  n1,
		rigthNode: n2,
		hash:      HashNodes(n1, n2),

		leftTx:  nil,
		rightTx: nil,
	}

	n1.parent = n
	n2.parent = n
	return n
}

func GenerateNodeFromTxs(t1 *Tx, t2 *Tx) *Node {
	n := &Node{
		leftNode:  nil,
		rigthNode: nil,
		leftTx:    t1,
		rightTx:   t2,
		hash:      HashTxs(t1, t2),
	}

	t1.parent = n
	t2.parent = n
	return n
}

func GenerateTree(txs []*Tx) (*Node, error) {
	txsLen := len(txs)

	if txsLen < 2 {
		return nil, errors.New("Must be at least 2 transactions")
	}

	nodes := make([]*Node, 0)

	for i := 0; i < txsLen; i += 2 {
		if txsLen-1 == i {
			nodes = append(nodes, GenerateSingleNodeFromTx(txs[i]))
			continue
		}

		nodes = append(nodes, GenerateNodeFromTxs(txs[i], txs[i+1]))
	}

	for {
		intermediaryNodesLen := len(nodes)

		if intermediaryNodesLen == 1 {
			return nodes[0], nil
		}

		intermediaryNodes := make([]*Node, 0)

		for i := 0; i < intermediaryNodesLen; i += 2 {
			if intermediaryNodesLen-1 == i {
				intermediaryNodes = append(intermediaryNodes, GenereteSingleNode(nodes[i]))
				continue
			}

			intermediaryNodes = append(intermediaryNodes, GenerateNode(nodes[i], nodes[i+1]))
		}

		nodes = intermediaryNodes
	}
}

func PrintTree(n *Node) {
	if n == nil {
		return
	}

	if n != nil {
		fmt.Printf("Node: %x\n", n.hash.Sum(nil))
	}

	if n.leftNode == nil && n.rigthNode == nil {
		if n.leftTx != nil {
			fmt.Printf("Tx: %x\n", n.leftTx.hash.Sum(nil))
		}

		if n.rightTx != nil {
			fmt.Printf("Tx: %x\n", n.rightTx.hash.Sum(nil))
		}
	}

	if n.leftNode != nil {
		PrintTree(n.leftNode)
	}

	if n.rigthNode != nil {
		PrintTree(n.rigthNode)
	}
}
