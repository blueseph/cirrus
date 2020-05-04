package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/blueseph/cirrus/cfn"
	"github.com/blueseph/cirrus/colors"
	"github.com/blueseph/cirrus/data"
	"github.com/blueseph/cirrus/utils"
	"github.com/rivo/tview"
)

func declineButtonCallbackFn(app *tview.Application, operation cfn.StackOperation) func() {
	return func() {
		declined := "change set"

		if operation != cfn.StackOperationDelete {
			declined = "delete"
		}

		defer fmt.Println(colors.ERROR + "User declined " + declined)
		app.Stop()
	}
}

func executeButtonCallbackFn(app *tview.Application, displayBox *tview.TextView, form *tview.Form, info data.StackInfo, operation cfn.StackOperation, displayRows map[string]data.DisplayRow, fillDisplayBox func(map[string]data.DisplayRow)) func() {
	return func() {
		resetForm(app, displayBox, form)

		activatedDisplayRows := activateRowsAndRender(displayRows, fillDisplayBox)
		executeOperation(operation, info)

		go handleEventsLoop(app, form, info, activatedDisplayRows, fillDisplayBox)
	}
}

func resetForm(app *tview.Application, displayBox *tview.TextView, form *tview.Form) {
	form.ClearButtons().SetTitle(" Errors ")

	app.SetFocus(displayBox)
	app.SetInputCapture(nil)
}

func activateRowsAndRender(displayRows map[string]data.DisplayRow, fillDisplayBox func(map[string]data.DisplayRow)) map[string]data.DisplayRow {
	activatedDisplayRows := data.ActivateDisplayRows(displayRows)
	fillDisplayBox(activatedDisplayRows)

	return activatedDisplayRows
}

func succeed(app *tview.Application) {
	defer fmt.Println(colors.SUCCESS + "Operation Succeeded")
	app.Stop()
}

func fail(app *tview.Application, errors []cloudformation.StackEvent) {
	errorMsg := colors.ERROR + "Operation failed. The following errors prevented the stack from deploying successfully: \n\n"

	for i, err := range errors {
		errorMsg += colors.Magenta(*err.LogicalResourceId) + " - " + string(*err.ResourceStatusReason)
		if i < len(errors)-1 {
			errorMsg += "\n"
		}
	}

	defer fmt.Println(errorMsg)
	app.Stop()
}

func executeOperation(operation cfn.StackOperation, info data.StackInfo) {
	var err error

	if operation == cfn.StackOperationDelete {
		err = cfn.DeleteStack(info)
	} else {
		err = cfn.ExecuteChangeSet(info)
	}

	if err != nil {
		panic(err)
	}
}

func handleEventsLoop(app *tview.Application, form *tview.Form, info data.StackInfo, activatedDisplayRows map[string]data.DisplayRow, fillDisplayBox func(map[string]data.DisplayRow)) {
	now := time.Now()

	eventIds := make(map[string]bool)
	errors := make([]cloudformation.StackEvent, 0)

	for {
		paginator := cfn.GetStackEvents(info)

		for paginator.Next(context.TODO()) {
			events := paginator.CurrentPage()

			for _, event := range utils.ReverseEvents(events.StackEvents) {
				if event.Timestamp.After(now) {
					if *event.ResourceType == data.CloudformationStackResource {
						if utils.ContainsStackStatus(data.RollbackStackStatus, event.ResourceStatus) {
							addErrorBar(form)
						}

						if !utils.ContainsStackStatus(data.PendingStackStatus, event.ResourceStatus) {
							if len(errors) > 0 {
								fail(app, errors)
							} else {
								succeed(app)
							}
						}
					} else if !eventIds[*event.EventId] {
						activatedDisplayRows[*event.LogicalResourceId] = data.CreateDisplayRowFromEvent(event)

						if utils.ContainsResourceStatus(data.NegativeEventStatus, event.ResourceStatus) {
							errors = append(errors, event)
						}

						eventIds[*event.EventId] = true
					}
				}
			}
		}

		fillDisplayBox(activatedDisplayRows)
		time.Sleep(500 * time.Millisecond)
	}
}
