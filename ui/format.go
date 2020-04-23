package ui

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/blueseph/cirrus/cfn"
	"github.com/blueseph/cirrus/data"
)

func stackOperationColorize(operation cfn.StackOperation) string {
	color := " [green::b]"
	end := "[-]"

	if operation == cfn.StackOperationUpdate {
		color = " [yellow::b]"
	}

	if operation == cfn.StackOperationDelete {
		color = " [red::b]"
	}

	return color + strings.ToUpper(string(operation)) + end
}

func resourceChangeColorize(change cloudformation.ChangeAction, ascii bool) string {
	color := "[green::b]"
	end := "[-]"

	if change == cloudformation.ChangeActionModify {
		color = "[yellow::b]"
	}

	if change == cloudformation.ChangeActionRemove {
		color = "[red::b]"
	}

	if ascii {
		return color + cfn.ChangeSetASCII[change] + end
	}

	return color + strings.ToUpper(string(change)) + end
}

func resourceTypeFormat(resourceType string) string {
	replaced := strings.ReplaceAll(resourceType, "::", ".")
	lowered := strings.ToLower(replaced)

	return "[grey::d]" + lowered + "[-]"
}

func getTitleBar(info data.StackInfo) string {
	var title string
	title += "[white]Stack:     [white::b]" + info.StackName + "\n"
	title += "[white]Changeset: [white::b]" + info.ChangeSetName + "\n"
	title += "[white]Id:        [white::b]" + info.StackID
	return title
}

func formatChange(change data.DisplayRow) string {
	var formatted string
	replacement := change.Replacement

	formatted += "[" + resourceChangeColorize(change.Action, true) + "] "
	formatted += "[#00b8ea]" + change.LogicalResourceID + " [white]"
	formatted += resourceChangeColorize(change.Action, false) + " "
	formatted += resourceTypeFormat(change.ResourceType)

	if replacement == cloudformation.ReplacementTrue {
		formatted += " [red]Replace[white]"
	}

	if replacement == cloudformation.ReplacementConditional {
		formatted += " [yellow]Replace conditional[white]"
	}

	return formatted + "\n"
}