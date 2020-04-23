package data

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

//DisplayRow is a normalized data structure to store change/event data to display
type DisplayRow struct {
	LogicalResourceID string
	ResourceType      string
	Status            cloudformation.ResourceStatus
	Timestamp         time.Time
	StatusReason      string
	Replacement       cloudformation.Replacement
	Action            cloudformation.ChangeAction
}

//StackInfo is a normalized data structure to store identifier properties of a stack/change set
type StackInfo struct {
	StackID       string
	ChangeSetName string
	StackName     string
}

// ChangeMap normalizes a slice of changes into a map of DisplayRows
func ChangeMap(changes []cloudformation.Change) map[string]DisplayRow {
	mapChanges := make(map[string]DisplayRow)

	for _, change := range changes {
		mapChanges[*change.ResourceChange.LogicalResourceId] = DisplayRow{
			LogicalResourceID: *change.ResourceChange.LogicalResourceId,
			ResourceType:      *change.ResourceChange.ResourceType,
			Replacement:       change.ResourceChange.Replacement,
			Action:            change.ResourceChange.Action,
		}
	}

	return mapChanges
}

// EventMap normalizes a slice of changes into a map of DisplayRows
func EventMap(events []cloudformation.StackEvent) map[string]DisplayRow {
	mapEvents := make(map[string]DisplayRow)

	for _, event := range events {
		mapEvents[*event.LogicalResourceId] = DisplayRow{
			LogicalResourceID: *event.LogicalResourceId,
			ResourceType:      *event.ResourceType,
			Status:            event.ResourceStatus,
			Timestamp:         *event.Timestamp,
		}
	}

	return mapEvents
}
