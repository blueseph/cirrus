package ui

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/blueseph/cirrus/cfn"
	"github.com/blueseph/cirrus/colors"
	"github.com/blueseph/cirrus/data"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	executeButtonLabel string = "Execute"
	declineButtonLabel string = "Decline"
)

//DisplayChanges shows the change set in a graphic interface and waits for response. Cancels the command if the user declines, or executes and tails the events log
func DisplayChanges(stackName string, changeSet *cloudformation.DescribeChangeSetResponse, operation cfn.StackOperation) error {
	changesMap := changeBuilder(changeSet.Changes)
	info := data.StackInfo{
		StackID:       *changeSet.StackId,
		ChangeSetName: *changeSet.ChangeSetName,
		StackName:     *changeSet.StackName,
	}

	err := showChanges(changesMap, operation, info)

	return err
}

func changeBuilder(changes []cloudformation.Change) map[string]data.DisplayRow {
	mapChanges := make(map[string]data.DisplayRow)
	for _, change := range changes {
		mapChanges[*change.ResourceChange.LogicalResourceId] = data.DisplayRow{
			LogicalResourceID: *change.ResourceChange.LogicalResourceId,
			ResourceType:      *change.ResourceChange.ResourceType,
			Replacement:       change.ResourceChange.Replacement,
			Action:            change.ResourceChange.Action,
		}
	}

	return mapChanges
}

func eventBuilder(events []cloudformation.StackEvent) map[string]data.DisplayRow {
	mapEvents := make(map[string]data.DisplayRow)

	for _, event := range events {
		mapEvents[*event.LogicalResourceId] = data.DisplayRow{
			LogicalResourceID: *event.LogicalResourceId,
			ResourceType:      *event.ResourceType,
			Status:            event.ResourceStatus,
			Timestamp:         *event.Timestamp,
		}
	}

	return mapEvents
}

func getChangesString(changes map[string]data.DisplayRow) string {
	var allChanges string
	for _, change := range changes {
		msg := formatChange(change)
		allChanges += msg
	}
	return allChanges
}

// func formatEvent(change data.DisplayRow) string {
// 	var formatted string
// 	replacement := change.Replacement
// }

func createTitleBar(info data.StackInfo, operation cfn.StackOperation) *tview.TextView {
	textView := tview.NewTextView().SetScrollable(false).SetDynamicColors(true).SetWordWrap(true)

	go func() {
		fmt.Fprintf(textView, "%s ", getTitleBar(info))
	}()

	textView.SetBorder(true).SetTitle(" " + info.StackName + stackOperationColorize(operation) + " ")

	return textView
}

func createChangesBox(changes map[string]data.DisplayRow) *tview.TextView {
	textView := tview.NewTextView().SetRegions(true).SetScrollable(true).SetDynamicColors(true).SetWordWrap(false)

	go func() {
		msg := getChangesString(changes)
		fmt.Fprintf(textView, "%s ", msg)
	}()

	textView.SetBorder(true).SetTitle(" Changes ")

	return textView
}

func createActionBar(app *tview.Application) *tview.Form {
	form := tview.NewForm().
		AddButton(executeButtonLabel, nil).
		AddButton(declineButtonLabel, func() {
			defer fmt.Println(colors.ERROR + "User declined change set")
			app.Stop()
		})

	form.SetButtonsAlign(tview.AlignCenter).SetBorder(true).SetTitle(" Actions ")

	return form
}

func showChanges(changes map[string]data.DisplayRow, operation cfn.StackOperation, info data.StackInfo) error {
	app := tview.NewApplication()
	titleBar := createTitleBar(info, operation)
	changesBox := createChangesBox(changes)
	actionBar := createActionBar(app)
	// liveBar := createLiveBar

	view := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titleBar, 5, 0, false).
		AddItem(changesBox, 0, 3, false).
		AddItem(actionBar, 5, 0, false)

	// hacky workarounds
	view.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		executeButton := actionBar.GetButton(0)
		declineButton := actionBar.GetButton(1)

		if e.Key() == tcell.KeyTab {
			switch {
			case changesBox.HasFocus():
				app.SetFocus(executeButton)
			case actionBar.HasFocus():
				if executeButton.GetFocusable().HasFocus() {
					app.SetFocus(declineButton)
				} else {
					app.SetFocus(changesBox)
				}
			}
		}

		if e.Key() == tcell.KeyBacktab {
			switch {
			case changesBox.HasFocus():
				app.SetFocus(declineButton)
			case actionBar.HasFocus():
				if executeButton.GetFocusable().HasFocus() {
					app.SetFocus(changesBox)
				} else {
					app.SetFocus(executeButton)
				}
			}
		}

		return e
	})

	app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if view.HasFocus() {
			view.GetInputCapture()(e)
		}

		return e
	})

	if err := app.SetRoot(view, true).SetFocus(changesBox).Run(); err != nil {
		panic(err)
	}

	return nil
}
