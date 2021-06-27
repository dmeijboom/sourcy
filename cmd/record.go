package main

import (
	"bytes"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"

	"github.com/dmeijboom/sourcy/pkg/storage"
)

func recordCmd(context *cli.Context) error {
	tree, err := storage.OpenTree(".srcy")
	if err != nil {
		log.Fatal(err)
	}

	file, err := storage.NewFile("test", 0644, bytes.NewReader([]byte("hello world")))
	if err != nil {
		log.Fatal(err)
	}

	err = tree.Append(storage.NewTransaction(nil, []*storage.Change{
		storage.NewChange(
			storage.Create,
			file,
		),
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tree)

	return nil
}
