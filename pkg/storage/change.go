package storage

import "github.com/dmeijboom/sourcy/pkg/hash"

type ChangeType uint8

const (
	Unknown ChangeType = 0
	Create  ChangeType = iota + 1
	Update
	Remove
)

type Change struct {
	Hash hash.Hash
	Type ChangeType
	File *File
}

func (change *Change) ComputeHash() hash.Hash {
	source := append(
		[]byte{byte(change.Type)},
		change.File.Hash.Bytes()...,
	)

	return hash.New(source)
}

func NewChange(changeType ChangeType, file *File) *Change {
	change := &Change{
		Type: changeType,
		File: file,
	}

	change.Hash = change.ComputeHash()

	return change
}
