package cmd

import (
	"errors"
	"fmt"

	"github.com/blueseph/cirrus/cfn"
	"github.com/blueseph/cirrus/colors"
	"github.com/blueseph/cirrus/data"
	"github.com/blueseph/cirrus/ui"
	"github.com/urfave/cli/v2"
)

var downFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "stack",
		Aliases:  []string{"s"},
		Usage:    "Specifies stack name",
		Required: true,
	},
}

// DownCommand returns the CLI construct that destroys a CloudFormation stack and watches events
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

// Down manages the stack deletion lifecycle
func Down(stackName string) error {
	err := cfn.VerifyAWSCredentials()
	if err != nil {
		return err
	}

	exists, err := cfn.DetermineIfStackExists(stackName)
	if err != nil {
		return err
	}

	stack, err := cfn.GetStack(stackName)
	if err != nil {
		return err
	}

	info := data.StackInfo{
		StackName: stackName,
		StackID:   *stack.DescribeStacksOutput.Stacks[0].StackId,
	}

	if !exists {
		return errors.New(colors.Error(fmt.Sprintf("Could not find stack %s", stackName)))
	}

	paginator := cfn.GetStackResources(info)
	if err != nil {
		return err
	}

	resources := data.GetResourcesFromPaginator(&paginator)

	err = ui.DisplayDeletes(info, resources)
	if err != nil {
		return err
	}

	return nil
}
