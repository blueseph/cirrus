package main

import (
	"log"
	"os"

	"github.com/blueseph/cirrus/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	log.SetFlags(0)

	app := &cli.App{
		Commands: []*cli.Command{
			cmd.UpCommand,
			cmd.DownCommand,
		},
	}

	app.EnableBashCompletion = true

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
