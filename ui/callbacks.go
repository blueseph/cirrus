package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/blueseph/cirrus/cfn"
	"github.com/blueseph/cirrus/colors"
	"github.com/blueseph/cirrus/data"
	"github.com/blueseph/cirrus/utils"
	"github.com/rivo/tview"
)

func declineButtonCallbackFn(app *tview.Application) func() {
	return func() {
		defer fmt.Println(colors.ERROR + "User declined change set")
		app.Stop()
	}
}

// this is awful, can we make this better
func executeButtonCallbackFn(app *tview.Application, displayBox *tview.TextView, form *tview.Form, info data.StackInfo, displayRows map[string]data.DisplayRow, fillDisplayBox func(map[string]data.DisplayRow)) func() {
	return func() {
		now := time.Now()
		form.ClearButtons().SetTitle(" Errors ")

		app.SetFocus(displayBox)
		app.SetInputCapture(nil)

		activatedDisplayRows := data.ActivateDisplayRows(displayRows)
		fillDisplayBox(activatedDisplayRows)

		err := cfn.ExecuteChangeSet(info)
		if err != nil {
			panic(err)
		}

		exit := func() {
			defer fmt.Println(colors.SUCCESS + "Operation Succeeded")
			app.Stop()
		}

		go func() {
			eventIds := make(map[string]bool)

			for {
				paginator := cfn.GetStackEvents(info)

				for paginator.Next(context.TODO()) {
					events := paginator.CurrentPage()

					for _, event := range events.StackEvents {

						if *event.ResourceType == data.CloudformationStackResource {
							if !utils.ContainsStackStatus(data.PendingStackStatus, event.ResourceStatus) {
								exit()
							}
							eventIds[*event.EventId] = true
						} else if !eventIds[*event.EventId] {
							if event.Timestamp.After(now) {
								activatedDisplayRows[*event.LogicalResourceId] = data.CreateDisplayRowFromEvent(event)
							}
							eventIds[*event.EventId] = true
						}
					}
				}

				fillDisplayBox(activatedDisplayRows)
				time.Sleep(500 * time.Millisecond)
			}
		}()
	}
}
