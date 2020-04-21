package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type changeScreenRow struct {
	LogicalResourceID string
	ResourceType      string
	Status            string
	Timestamp         time.Time
	StatusReason      string
}

var (
	xPadding = 3
	yPadding = 1
)

func displayChanges(stackName string, changeSet *cloudformation.DescribeChangeSetResponse, operation stackOperation) (bool, error) {
	changesMap := changeBuilder(changeSet.Changes, operation)

	verification, err := showChanges(changesMap, operation, changeSet)

	return verification, err
}

func changeBuilder(changes []cloudformation.Change, operation stackOperation) map[string]changeScreenRow {
	mapChanges := make(map[string]changeScreenRow)
	for _, change := range changes {
		status := operation
		mapChanges[*change.ResourceChange.LogicalResourceId] = changeScreenRow{
			LogicalResourceID: *change.ResourceChange.LogicalResourceId,
			ResourceType:      *change.ResourceChange.ResourceType,
			Status:            string(status),
		}
	}

	return mapChanges
}

func titleBarDrawFn(changeSet cloudformation.DescribeChangeSetResponse) func(tcell.Screen, int, int, int, int) (int, int, int, int) {
	return func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		tview.Print(screen, "[white]Stack:     [white:b]"+*changeSet.StackName, x+xPadding, y+yPadding, width, tview.AlignLeft, tcell.ColorWhite)
		tview.Print(screen, "[white]Changeset: [white:b]"+*changeSet.ChangeSetName, x+xPadding, y+yPadding+1, width, tview.AlignLeft, tcell.ColorWhite)
		tview.Print(screen, "[white]Id:        [white:b]"+*changeSet.StackId, x+xPadding, y+yPadding+2, width, tview.AlignLeft, tcell.ColorWhite)
		return 0, 0, 0, 0
	}
}

func changesDrawFn(changeSet cloudformation.DescribeChangeSetResponse) func(tcell.Screen, int, int, int, int) (int, int, int, int) {
	return func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		for i, change := range changeSet.Changes {
			msg := formatChange(change)

			tview.Print(screen, msg, x+xPadding, y+yPadding+i, width, tview.AlignLeft, tcell.ColorWhite)
		}

		return 0, 0, 0, 0
	}
}

func formatChange(change cloudformation.Change) string {
	var formatted string
	replacement := change.ResourceChange.Replacement

	formatted += "[" + resourceChangeColorize(change.ResourceChange.Action, true) + "] "
	formatted += "[#00b8ea]" + *change.ResourceChange.LogicalResourceId + " [white]"
	formatted += resourceChangeColorize(change.ResourceChange.Action, false) + " "
	formatted += resourceTypeFormat(*change.ResourceChange.ResourceType)

	if replacement == cloudformation.ReplacementTrue {
		formatted += " [red]Replace[white]"
	}

	if replacement == cloudformation.ReplacementConditional {
		formatted += " [yellow]Replace conditional[white]"
	}

	return formatted
}

func showChanges(changes map[string]changeScreenRow, operation stackOperation, changeSet *cloudformation.DescribeChangeSetResponse) (bool, error) {
	app := tview.NewApplication()

	form := tview.NewForm().
		AddButton("Execute change set", nil).
		AddButton("Decline change set", func() {
			defer fmt.Println(ERROR + "User declined change set")
			app.Stop()
		})

	form.SetButtonsAlign(tview.AlignCenter).SetBorder(true).SetTitle(" Actions ")
	
	changeView := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle(" "+*changeSet.StackName+stackOperationColorize(operation)+" ").SetDrawFunc(titleBarDrawFn(*changeSet)), 5, 0, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle(" Changes ").SetDrawFunc(changesDrawFn(*changeSet)), 0, 3, false).
			AddItem(form, 5, 0, false),
			0, 1, false)

	if err := app.SetRoot(changeView, true).SetFocus(form).Run(); err != nil {
		panic(err)
	}

	// fmt.Printf("%+v\n", changeSet)

	return true, nil
}

func stackOperationColorize(operation stackOperation) string {
	color := " [green::b]"
	end := "[-:-:-]"

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
	end := "[-:-:-]"

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

	return "[grey::d]" + lowered + "[-:-:-]"
}
