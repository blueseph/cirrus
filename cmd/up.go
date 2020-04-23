package cmd

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/blueseph/cirrus/cfn"
	"github.com/blueseph/cirrus/colors"
	"github.com/blueseph/cirrus/ui"
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

// Up kicks off the stack creation lifecycle, creating a change set, confirming the change set, and tailing the events.
func Up(stackName string, template []byte) error {
	err := cfn.VerifyAWSCredentials()
	if err != nil {
		return err
	}

	changeSetID := stackName + "-" + fmt.Sprint(time.Now().Unix())
	exists, err := cfn.DetermineIfStackExists(stackName)
	if err != nil {
		return err
	}

	fmt.Println(colors.STATUS + "Creating change set...")
	changeSet, err := cfn.CreateChanges(stackName, changeSetID, template, exists)
	if err != nil {
		return err
	}

	operation := cfn.StackOperationCreate
	if exists {
		operation = cfn.StackOperationUpdate
	}

	err = ui.DisplayChanges(stackName, changeSet, operation)

	return err
}