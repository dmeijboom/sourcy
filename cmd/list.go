package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"

	"github.com/dmeijboom/sourcy/pkg/storage"
)

func listCmd(context *cli.Context) error {
	tree, err := storage.OpenTree(".srcy")
	if err != nil {
		log.Fatal(err)

		return err
	}

	if tree.Root == nil {
		fmt.Println("empty tree")
		return nil
	}

	node := tree.Root

	for node != nil {
		rootMarker := ""

		if tree.Root.Hash.Eq(node.Hash) {
			rootMarker = "*"
		}

		fmt.Printf("%s%s: %d changes\n", rootMarker, node.Hash.Short(), len(node.Changes))

		node = node.Parent
	}

	return nil
}
