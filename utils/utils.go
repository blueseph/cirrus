package utils

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

// ContainsResourceStatus takes a ResourceStatus slice and looks for a specific value in it
func ContainsResourceStatus(slice []cloudformation.ResourceStatus, val cloudformation.ResourceStatus) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}

	return false
}

// ContainsStackStatus takes a StackStatus slice and looks for a specific value in it
func ContainsStackStatus(slice []cloudformation.StackStatus, val cloudformation.ResourceStatus) bool {
	for _, item := range slice {
		// aws::cloudformation::stack can have a resource status be a stack status
		if string(item) == string(val) {
			return true
		}
	}

	return false
}

// ReverseEvents returns a reversed slice of events without side-effects
func ReverseEvents(s []cloudformation.StackEvent) []cloudformation.StackEvent {
	a := make([]cloudformation.StackEvent, len(s))
	copy(a, s)

	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}
