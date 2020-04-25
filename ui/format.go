package ui

import (
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/blueseph/cirrus/cfn"
	"github.com/blueseph/cirrus/data"
	"github.com/blueseph/cirrus/utils"
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

func colorizeAction(change cloudformation.ChangeAction, ascii bool) string {
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

func colorizeStatus(status cloudformation.ResourceStatus) string {
	color := "[green::b]"
	end := "[-]"

	if utils.ContainsStatus(data.PendingEventStatus, status) {
		color = "[yellow::b]"
	}

	if utils.ContainsStatus(data.NegativeEventStatus, status) {
		color = "[red::b]"
	}

	return color + strings.ToUpper(string(status)) + end
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

func parseDisplayRow(row data.DisplayRow) string {
	if row.Source == data.DisplayRowSourceChangeSet {
		return parseChangeRow(row)
	}

	return parseEventRow(row)
}

func parseChangeRow(row data.DisplayRow) string {
	var formatted string
	replacement := row.Replacement

	if row.Active {
		formatted += "[[grey]PENDING_" + strings.ToUpper(string(row.Action)) + "[-]]"
	} else {
		formatted += "[" + colorizeAction(row.Action, true) + "] "
	}

	formatted += "[#00b8ea]" + row.LogicalResourceID + " [white]"
	if !row.Active {
		formatted += colorizeAction(row.Action, false) + " "
	}
	formatted += resourceTypeFormat(row.ResourceType)

	if !row.Active {
		if replacement == cloudformation.ReplacementTrue {
			formatted += " [red]Replace[white]"
		}

		if replacement == cloudformation.ReplacementConditional {
			formatted += " [yellow]Replace conditional[white]"
		}
	}

	return formatted + "\n"
}

func parseEventRow(row data.DisplayRow) string {
	var formatted string

	formatted += "[" + colorizeStatus(row.Status) + "]"
	formatted += "[#00b8ea]" + row.LogicalResourceID + " [white]"
	formatted += resourceTypeFormat(row.ResourceType)

	return formatted + "\n"
}

//ParseDisplayRows parses and sorts the map of display rows and returns a tview.TextBox consumable string
func ParseDisplayRows(displayRows map[string]data.DisplayRow) string {
	keys := make([]string, 0)
	var allChanges string

	for key := range displayRows {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		msg := parseDisplayRow(displayRows[key])
		allChanges += msg
	}
	return allChanges
}
