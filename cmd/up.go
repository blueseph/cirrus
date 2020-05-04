package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/blueseph/cirrus/cfn"
	"github.com/blueseph/cirrus/colors"
	"github.com/blueseph/cirrus/data"
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
		Name:    "parameters",
		Aliases: []string{"p"},
		Value:   "./parameters.json",
		Usage:   "Specifies location of parameters `file`",
	},
	&cli.StringFlag{
		Name:  "tags",
		Value: "./tags.json",
		Usage: "Specifies location of tags `file`",
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
	&cli.BoolFlag{
		Name:    "overwrite",
		Aliases: []string{"o"},
		Usage:   "Overwrites existing empty (0 resource) stacks before updating",
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

	tags, err := data.GetTags(c.String("tags"))
	if err != nil {
		return err
	}

	parameters, err := data.GetParameters(c.String("parameters"))
	if err != nil {
		return err
	}

	stack := c.String("stack")
	overwrite := c.Bool("overwrite")

	err = Up(stack, overwrite, template, tags, parameters)
	if err != nil {
		fmt.Println(colors.Error("Cirrus encountered a fatal error:"))
		return err
	}

	return nil
}

// Up kicks off the stack creation lifecycle, creating a change set, confirming the change set, and tailing the events.
func Up(stackName string, overwrite bool, template []byte, tags []cloudformation.Tag, parameters []cloudformation.Parameter) error {
	changeSetName := stackName + "-" + fmt.Sprint(time.Now().Unix())

	info := data.StackInfo{
		StackName:     stackName,
		ChangeSetName: changeSetName,
	}

	err := cfn.VerifyAWSCredentials()
	if err != nil {
		return err
	}

	exists, err := cfn.DetermineIfStackExists(info.StackName)
	if err != nil {
		return err
	}

	empty := cfn.DetermineIfStackIsEmpty(info)

	if exists && empty {
		err := handleOverwrite(overwrite, exists, info)
		if err != nil {
			return err
		}
	}

	fmt.Println(colors.Status("Creating change set..."))
	changeSet, err := cfn.CreateChanges(info, template, tags, parameters, exists)
	if err != nil {
		return err
	}

	info.StackID = *changeSet.StackId

	operation := cfn.StackOperationCreate
	if exists {
		operation = cfn.StackOperationUpdate
	}

	err = ui.DisplayChanges(info, changeSet, operation)

	return err

}

func askYesNoQuestion(question string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(question)

	for {
		char, _, err := reader.ReadRune()

		if err != nil {
			return false, err
		}

		char = unicode.ToLower(char)

		switch char {
		case 'y':
			return true, nil
		case 'n':
			return false, nil
		default:
			fmt.Println("Please enter Y/N")
		}
	}
}

func handleOverwrite(overwrite bool, exists bool, info data.StackInfo) error {
	var err error
	confirm := overwrite

	if !confirm {
		confirm, err = askYesNoQuestion(colors.Status("Empty stack detected. Overwrite? [Y/N]"))
		if err != nil {
			return err
		}
	}

	if confirm {
		fmt.Println(colors.Status("Deleting stack..."))
		err := cfn.DeleteStackAndWait(info)
		exists = false
		if err != nil {
			return err
		}
	} else {
		fmt.Println(colors.Status("User declined empty stack deletion. Terminating"))
		return nil
	}

	return nil
}
