package main

import "github.com/urfave/cli/v2"

var downFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "stack",
		Aliases:  []string{"s"},
		Usage:    "Specifies stack name",
		Required: true,
	},
}

// UpCommand returns the CLI construct that destroys a CloudFormation stack and watches events
var DownCommand = &cli.Command{
	Name:   "down",
	Usage:  "Bring down a CloudFormation template and watch stack events",
	Action: downAction,
	Flags:  downFlags,
}

func downAction(c *cli.Context) error {
	err := Down(c.String("stack"))
	if err != nil {
		return err
	}

	return nil
}
