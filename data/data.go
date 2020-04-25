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
	Source            DisplayRowSource
	Active            bool
}

//StackInfo is a normalized data structure to store identifier properties of a stack/change set
type StackInfo struct {
	StackID       string
	ChangeSetName string
	StackName     string
}

//DisplayRowSource is an enum to determine the origin of the display row
type DisplayRowSource string

const (
	//DisplayRowSourceChangeSet indicates a display row came from a Change Set
	DisplayRowSourceChangeSet DisplayRowSource = "change"

	//DisplayRowSourceEvent indicates a display row came from an Event
	DisplayRowSourceEvent DisplayRowSource = "event"

	//CloudformationStackResource is the string that represents a CloudFormation stack in a template
	CloudformationStackResource string = "AWS::CloudFormation::Stack"
)

var (
	//PositiveEventStatus indicates positive event statuses
	PositiveEventStatus []cloudformation.ResourceStatus = []cloudformation.ResourceStatus{
		cloudformation.ResourceStatusCreateComplete,
		cloudformation.ResourceStatusDeleteComplete,
		cloudformation.ResourceStatusUpdateComplete,
	}

	//NegativeEventStatus indicates negative event statuses
	NegativeEventStatus []cloudformation.ResourceStatus = []cloudformation.ResourceStatus{
		cloudformation.ResourceStatusCreateFailed,
		cloudformation.ResourceStatusDeleteFailed,
		cloudformation.ResourceStatusUpdateFailed,
	}

	//PendingEventStatus indicates an event status that is in a pending state
	PendingEventStatus []cloudformation.ResourceStatus = []cloudformation.ResourceStatus{
		cloudformation.ResourceStatusCreateInProgress,
		cloudformation.ResourceStatusDeleteInProgress,
		cloudformation.ResourceStatusUpdateInProgress,
	}

	//PositiveStackStatus status indicates a stack is in a positive terminal state
	PositiveStackStatus []cloudformation.StackStatus = []cloudformation.StackStatus{
		cloudformation.StackStatusCreateComplete,
		cloudformation.StackStatusDeleteComplete,
		cloudformation.StackStatusUpdateComplete,
		cloudformation.StackStatusRollbackComplete,
	}

	//NegativeStackStatus status indicates a stack is in a negative terminal state
	NegativeStackStatus []cloudformation.StackStatus = []cloudformation.StackStatus{
		cloudformation.StackStatusCreateFailed,
		cloudformation.StackStatusDeleteFailed,
		cloudformation.StackStatusUpdateRollbackComplete,
		cloudformation.StackStatusUpdateRollbackFailed,
		cloudformation.StackStatusRollbackFailed,
	}

	//PendingStackStatus status indicates a stack is not yet in a terminal state
	PendingStackStatus []cloudformation.StackStatus = []cloudformation.StackStatus{
		cloudformation.StackStatusCreateInProgress,
		cloudformation.StackStatusDeleteInProgress,
		cloudformation.StackStatusUpdateInProgress,
		cloudformation.StackStatusReviewInProgress,
		cloudformation.StackStatusUpdateRollbackInProgress,
		cloudformation.StackStatusRollbackInProgress,
	}
)

// ChangeMap normalizes a slice of changes into a map of DisplayRows
func ChangeMap(changes []cloudformation.Change, active bool) map[string]DisplayRow {
	mapChanges := make(map[string]DisplayRow)

	for _, change := range changes {
		mapChanges[*change.ResourceChange.LogicalResourceId] = CreateDisplayRowFromChange(change, active)
	}

	return mapChanges
}

//CreateDisplayRowFromChange normalizes a cloudformation change into a display row
func CreateDisplayRowFromChange(change cloudformation.Change, active bool) DisplayRow {
	return DisplayRow{
		LogicalResourceID: *change.ResourceChange.LogicalResourceId,
		ResourceType:      *change.ResourceChange.ResourceType,
		Replacement:       change.ResourceChange.Replacement,
		Action:            change.ResourceChange.Action,
		Source:            DisplayRowSourceChangeSet,
		Active:            active,
	}
}

// EventMap normalizes a slice of changes into a map of DisplayRows
func EventMap(events []cloudformation.StackEvent) map[string]DisplayRow {
	mapEvents := make(map[string]DisplayRow)

	for _, event := range events {
		mapEvents[*event.LogicalResourceId] = CreateDisplayRowFromEvent(event)
	}

	return mapEvents
}

//CreateDisplayRowFromEvent normalizes a cloudformation event into a display row
func CreateDisplayRowFromEvent(event cloudformation.StackEvent) DisplayRow {
	return DisplayRow{
		LogicalResourceID: *event.LogicalResourceId,
		ResourceType:      *event.ResourceType,
		Status:            event.ResourceStatus,
		Timestamp:         *event.Timestamp,
		Source:            DisplayRowSourceEvent,
	}
}

//ActivateDisplayRows iterates through a display row map and sets the active flag to true
func ActivateDisplayRows(displayRows map[string]DisplayRow) map[string]DisplayRow {
	activatedDisplayRows := make(map[string]DisplayRow)

	for logicalID, event := range displayRows {
		event.Active = true
		activatedDisplayRows[logicalID] = event
	}

	return activatedDisplayRows
}
