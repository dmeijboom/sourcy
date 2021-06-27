package storage

import (
	"errors"

	"github.com/dmeijboom/sourcy/pkg/hash"
)

var TransactionDoesNotExist = errors.New("transaction doesn't exists")

type Transaction struct {
	Hash    hash.Hash
	Parent  *Transaction
	Changes []*Change
}

func NewTransaction(parent *Transaction, changes []*Change) *Transaction {
	transaction := &Transaction{
		Parent:  parent,
		Changes: changes,
	}

	transaction.Hash = transaction.ComputeHash()

	return transaction
}

func (transaction *Transaction) ComputeHash() hash.Hash {
	source := []byte{}

	if transaction.Parent != nil {
		source = append(source, transaction.Parent.Hash[:]...)
	}

	for _, change := range transaction.Changes {
		source = append(source, change.Hash[:]...)
	}

	return hash.New(source)
}
