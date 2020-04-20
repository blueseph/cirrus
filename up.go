package main

import (
	"io/ioutil"

	"github.com/urfave/cli/v2"
)

var upFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "template",
		Aliases: []string{"t"},
		Value:   "./template.yaml",
		Usage:   "Specifies location of template `file`",
	},
	&cli.StringFlag{
		Name:     "stack",
		Aliases:  []string{"s"},
		Usage:    "Specifies `stack name`",
		Required: true,
	},
	&cli.BoolFlag{
		Name:    "skip-lint",
		Aliases: []string{"sl"},
		Usage:   "Skips linting (not recommended)",
	},
}

// UpCommand returns the CLI construct that uploads a template to CloudFormation and watches the response
var UpCommand = &cli.Command{
	Name:   "up",
	Usage:  "Deploy a CloudFormation template and watch stack events",
	Action: upAction,
	Flags:  upFlags,
}

func upAction(c *cli.Context) error {
	template, err := ioutil.ReadFile(c.String("template"))
	if err != nil {
		return err
	}

	err = Up(c.String("stack"), template)
	if err != nil {
		return err
	}

	return nil
}
