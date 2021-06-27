package storage

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"os"

	"github.com/dmeijboom/sourcy/pkg/hash"
)

type File struct {
	Hash    hash.Hash
	Name    string
	Mode    os.FileMode
	Content *Content
}

func NewFile(name string, mode os.FileMode, reader io.Reader) (*File, error) {
	buff := bytes.Buffer{}

	writer := zlib.NewWriter(&buff)
	defer writer.Close()

	if _, err := io.Copy(writer, reader); err != nil {
		return nil, err
	}

	content := NewContent(buff.Bytes())

	file := &File{
		Name:    name,
		Mode:    mode,
		Content: content,
	}

	file.Hash = file.ComputeHash()

	return file, nil
}

func (file *File) ComputeHash() hash.Hash {
	modeBytes := make([]byte, 4)
	binary.PutUvarint(modeBytes[:], uint64(file.Mode))

	source := []byte(file.Name)
	source = append(source, modeBytes[:]...)
	source = append(source, file.Content.Hash[:]...)

	return hash.New(source)
}
