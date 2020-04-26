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
