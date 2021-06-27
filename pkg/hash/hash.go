package hash

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

const Len = 32

type Hash [Len]byte

func New(source []byte) Hash {
	return sha256.Sum256(source)
}

func FromBytesString(hashBytes []byte) (Hash, error) {
	return FromString(string(hashBytes))
}

func FromString(hashString string) (Hash, error) {
	var hash Hash

	hashBytes := []byte{}

	_, err := fmt.Sscanf(hashString, "%x", &hashBytes)
	if err != nil {
		return hash, err
	}

	if len(hashBytes) != Len {
		return hash, fmt.Errorf("invalif hash format: %s", hashString)
	}

	copy(hash[:], hashBytes)

	return hash, nil
}

func (hash Hash) Eq(other Hash) bool {
	return bytes.Equal(hash[:], other[:])
}

func (hash Hash) Bytes() []byte {
	return hash[:]
}

func (hash Hash) String() string {
	return fmt.Sprintf("%x", hash[:])
}

func (hash Hash) Short() string {
	return hash.String()[:10]
}
