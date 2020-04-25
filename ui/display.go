package ui

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/blueseph/cirrus/cfn"
	"github.com/blueseph/cirrus/data"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var (
	executeButtonLabel string = "Execute"
	declineButtonLabel string = "Decline"
)

//DisplayChanges shows the change set in a graphic interface and waits for response. Cancels the command if the user declines, or executes and tails the events log
func DisplayChanges(info data.StackInfo, changeSet *cloudformation.DescribeChangeSetResponse, operation cfn.StackOperation) error {
	displayRows := data.ChangeMap(changeSet.Changes, false)

	err := showScreen(displayRows, operation, info)

	return err
}

func createTitleBar(info data.StackInfo, operation cfn.StackOperation) *tview.TextView {
	textView := tview.NewTextView().SetScrollable(false).SetDynamicColors(true).SetWordWrap(true)

	fmt.Fprintf(textView, "%s ", getTitleBar(info))

	textView.SetBorder(true).SetTitle(" " + info.StackName + stackOperationColorize(operation) + " ")

	return textView
}

func createDisplayRowBox(app *tview.Application) *tview.TextView {
	textView := tview.NewTextView().SetRegions(true).SetScrollable(true).SetDynamicColors(true).SetWrap(false).
		SetChangedFunc(func() {
			app.Draw()
		})

	textView.SetBorder(true).SetTitle(" Changes ")

	return textView
}

func fillDisplayBoxFn(displayBox *tview.TextView) func(map[string]data.DisplayRow) {
	return func(displayRows map[string]data.DisplayRow) {
		displayBox.SetText(ParseDisplayRows(displayRows))
	}
}

func createActionBar(app *tview.Application, displayBox *tview.TextView, info data.StackInfo, displayRows map[string]data.DisplayRow, fillDisplayBox func(map[string]data.DisplayRow)) *tview.Form {
	form := tview.NewForm()

	form.
		AddButton(executeButtonLabel, executeButtonCallbackFn(app, displayBox, form, info, displayRows, fillDisplayBox)).
		AddButton(declineButtonLabel, declineButtonCallbackFn(app))

	form.SetButtonsAlign(tview.AlignCenter).SetBorder(true).SetTitle(" Actions ")

	return form
}

func showScreen(displayRows map[string]data.DisplayRow, operation cfn.StackOperation, info data.StackInfo) error {
	app := tview.NewApplication()

	displayBox := createDisplayRowBox(app)
	fillDisplayBox := fillDisplayBoxFn(displayBox)

	titleBar := createTitleBar(info, operation)
	actionBar := createActionBar(app, displayBox, info, displayRows, fillDisplayBox)

	fillDisplayBox(displayRows)
	// liveBar := createLiveBar

	view := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(titleBar, 5, 0, false).
		AddItem(displayBox, 0, 3, false).
		AddItem(actionBar, 5, 0, false)

	// y u ck
	viewSetInputCapture := viewInputCaptureFn(app, actionBar, displayBox)
	view.SetInputCapture(viewSetInputCapture)

	appSetInputCapture := appSetInputCaptureFn(view)
	app.SetInputCapture(appSetInputCapture)

	if err := app.SetRoot(view, true).SetFocus(displayBox).Run(); err != nil {
		panic(err)
	}

	return nil
}

//hacky workaround
func appSetInputCaptureFn(view *tview.Flex) func(*tcell.EventKey) *tcell.EventKey {
	return func(e *tcell.EventKey) *tcell.EventKey {
		if view.HasFocus() {
			view.GetInputCapture()(e)
		}

		return e
	}
}

//hacky workaround
func viewInputCaptureFn(app *tview.Application, actionBar *tview.Form, displayBox *tview.TextView) func(*tcell.EventKey) *tcell.EventKey {
	return func(e *tcell.EventKey) *tcell.EventKey {
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
	}
}
