package main

import (
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/rivo/tview"
)

type changeScreenRow struct {
	LogicalResourceID string
	ResourceType      string
	Status            string
	Timestamp         time.Time
	StatusReason      string
}

func displayChanges(stackName string, changeSet *cloudformation.DescribeChangeSetResponse, operation stackOperation) (bool, error) {
	changesMap := changeBuilder(changeSet.Changes, operation)

	verification, err := showChanges(changesMap, stackName, operation, changeSet)

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

func showChanges(changes map[string]changeScreenRow, stackName string, operation stackOperation, changeSet *cloudformation.DescribeChangeSetResponse) (bool, error) {
	app := tview.NewApplication()
	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle(stackName+stackOperationColorize(operation)), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	return true, nil
}

func testBox() {
	app := tview.NewApplication()
	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("bruce"+stackOperationColorize(update)), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Changes"), 0, 3, false),
			0, 2, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func stackOperationColorize(operation stackOperation) string {
	color := " [green]"
	end := "[white]"

	if operation == update {
		color = " [yellow]"
	}

	if operation == delete {
		color = " [red]"
	}

	return color + strings.ToUpper(string(operation)) + end
}
