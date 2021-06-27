package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/dmeijboom/sourcy/pkg/hash"
)

type TransactionHeader struct {
	Hash   string  `json:"hash"`
	Parent *string `json:"parent"`
}

type Tree struct {
	dir  string
	Root *Transaction
}

func OpenTree(dir string) (*Tree, error) {
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		return &Tree{dir: dir}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open tree: %w", err)
	}

	content, err := ioutil.ReadFile(fmt.Sprintf("%s/root", dir))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, TransactionDoesNotExist
		}

		return nil, fmt.Errorf("failed to read transaction: %w", err)
	}

	rootHash, err := hash.FromBytesString(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction hash: %w", err)
	}

	tree := &Tree{dir: dir}

	transaction, err := tree.ReadTransaction(rootHash)
	if err != nil {
		return nil, fmt.Errorf("failed to read root transaction: %w", err)
	}

	tree.Root = transaction

	return tree, nil
}

func (tree *Tree) Append(transaction *Transaction) error {
	if tree.Root == nil {
		tree.Root = transaction
	} else {
		if tree.Root.Hash.Eq(transaction.Hash) {
			return fmt.Errorf("unable to add empty transaction")
		}

		transaction.Parent = tree.Root
		tree.Root = transaction
	}

	// write transaction to disk
	filename := fmt.Sprintf("%s/tr/%s", tree.dir, transaction.Hash.String())
	dirname := path.Dir(filename)

	_, err := os.Stat(dirname)

	if os.IsNotExist(err) {
		if err := os.MkdirAll(dirname, 0700); err != nil {
			return fmt.Errorf("failed to create transaction dir: %w", err)
		}
	} else if err != nil {
		return err
	}

	if err := tree.StoreTransaction(transaction); err != nil {
		return fmt.Errorf("failed to store transaction %s: %w", transaction.Hash.String(), err)
	}

	// update the root transaction reference
	rootFilename := fmt.Sprintf("%s/root", tree.dir)

	if err := ioutil.WriteFile(rootFilename, []byte(transaction.Hash.String()), 0644); err != nil {
		return fmt.Errorf("failed to update tree root: %w", err)
	}

	return nil
}

func (tree *Tree) ReadTransaction(txHash hash.Hash) (*Transaction, error) {
	filename := fmt.Sprintf("%s/tr/%s", tree.dir, txHash.String())
	header := &TransactionHeader{}

	if err := decode(filename, header); err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %w", err)
	}

	var parent *Transaction

	if header.Parent != nil {
		var err error

		parentHash, err := hash.FromString(*header.Parent)
		if err != nil {
			return nil, fmt.Errorf("invalid hash format: %s", string(*header.Parent))
		}

		parent, err = tree.ReadTransaction(parentHash)
		if err != nil {
			return nil, err
		}
	}

	actualHash, err := hash.FromString(header.Hash)
	if err != nil {
		return nil, fmt.Errorf("invalid hash format: %s", string(header.Hash))
	}

	if !actualHash.Eq(txHash) {
		return nil, fmt.Errorf("failed to read corrupt transaction")
	}

	return &Transaction{
		Hash:   actualHash,
		Parent: parent,
	}, nil
}

func (tree *Tree) StoreTransaction(transaction *Transaction) error {
	header := &TransactionHeader{
		Hash: transaction.Hash.String(),
	}

	if transaction.Parent != nil {
		parentHashString := transaction.Parent.Hash.String()
		header.Parent = &parentHashString
	}

	data, err := json.Marshal(header)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/tr/%s", tree.dir, transaction.Hash)

	return ioutil.WriteFile(filename, data, 0644)
}
