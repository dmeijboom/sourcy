package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "sourcy",
		Usage: "some sort of source control I guess",
		Commands: []*cli.Command{
			{
				Name:   "record",
				Action: recordCmd,
			},
			{
				Name:   "list",
				Action: listCmd,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
