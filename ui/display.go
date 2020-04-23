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
	displayRows := data.ChangeMap(changeSet.Changes)

	info := data.StackInfo{
		StackID:       *changeSet.StackId,
		ChangeSetName: *changeSet.ChangeSetName,
		StackName:     *changeSet.StackName,
	}

	err := showScreen(displayRows, operation, info)

	return err
}

func createTitleBar(info data.StackInfo, operation cfn.StackOperation) *tview.TextView {
	textView := tview.NewTextView().SetScrollable(false).SetDynamicColors(true).SetWordWrap(true)

	go func() {
		fmt.Fprintf(textView, "%s ", getTitleBar(info))
	}()

	textView.SetBorder(true).SetTitle(" " + info.StackName + stackOperationColorize(operation) + " ")

	return textView
}

func createDisplayRowBox() *tview.TextView {
	textView := tview.NewTextView().SetRegions(true).SetScrollable(true).SetDynamicColors(true).SetWordWrap(false)

	textView.SetBorder(true).SetTitle(" Changes ")

	return textView
}

func fillDisplayBoxFn(displayBox *tview.TextView) func(map[string]data.DisplayRow) {
	return func(displayRows map[string]data.DisplayRow) {
		msg := ParseDisplayRows(displayRows)
		fmt.Fprintf(displayBox, "%s ", msg)
	}
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

func showScreen(displayRows map[string]data.DisplayRow, operation cfn.StackOperation, info data.StackInfo) error {
	app := tview.NewApplication()
	titleBar := createTitleBar(info, operation)
	actionBar := createActionBar(app)

	displayBox := createDisplayRowBox()
	fillDisplayBox := fillDisplayBoxFn(displayBox)

	fillDisplayBox(displayRows)
	// liveBar := createLiveBar

	view := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titleBar, 5, 0, false).
		AddItem(displayBox, 0, 3, false).
		AddItem(actionBar, 5, 0, false)

	// hacky workarounds
	view.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		executeButton := actionBar.GetButton(0)
		declineButton := actionBar.GetButton(1)

		if e.Key() == tcell.KeyTab {
			switch {
			case displayBox.HasFocus():
				app.SetFocus(executeButton)
			case actionBar.HasFocus():
				if executeButton.GetFocusable().HasFocus() {
					app.SetFocus(declineButton)
				} else {
					app.SetFocus(displayBox)
				}
			}
		}

		if e.Key() == tcell.KeyBacktab {
			switch {
			case displayBox.HasFocus():
				app.SetFocus(declineButton)
			case actionBar.HasFocus():
				if executeButton.GetFocusable().HasFocus() {
					app.SetFocus(displayBox)
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

	if err := app.SetRoot(view, true).SetFocus(displayBox).Run(); err != nil {
		panic(err)
	}

	return nil
}
