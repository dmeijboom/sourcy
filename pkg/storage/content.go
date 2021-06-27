package storage

import "github.com/dmeijboom/sourcy/pkg/hash"

type Content struct {
	Hash hash.Hash
	Data []byte
}

func NewContent(data []byte) *Content {
	return &Content{
		Hash: hash.New(data),
		Data: data,
	}
}
