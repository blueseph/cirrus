package utils

import "github.com/aws/aws-sdk-go-v2/service/cloudformation"

// ContainsStatus takes a slice and looks for a specific value in it
func ContainsStatus(slice []cloudformation.ResourceStatus, val cloudformation.ResourceStatus) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}

	return false
}
