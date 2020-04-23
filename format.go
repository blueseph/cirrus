package main

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"strings"
)

func stackOperationColorize(operation stackOperation) string {
	color := " [green::b]"
	end := "[-]"

	if operation == update {
		color = " [yellow::b]"
	}

	if operation == delete {
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
		return color + changeSetASCII[change] + end
	}

	return color + strings.ToUpper(string(change)) + end
}

func resourceTypeFormat(resourceType string) string {
	replaced := strings.ReplaceAll(resourceType, "::", ".")
	lowered := strings.ToLower(replaced)

	return "[grey::d]" + lowered + "[-]"
}

func getTitleBar(info stackInfo) string {
	var title string
	title += "[white]Stack:     [white::b]" + info.stackName + "\n"
	title += "[white]Changeset: [white::b]" + info.changeSetName + "\n"
	title += "[white]Id:        [white::b]" + info.stackID
	return title
}

func formatChange(change changeScreenRow) string {
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
