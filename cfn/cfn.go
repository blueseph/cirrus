package cfn

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/blueseph/cirrus/colors"
	"github.com/blueseph/cirrus/data"
)

var (
	cfnClient *cloudformation.Client

	//ChangeSetASCII is a map to convert a change action to a glyph representing the action. + for Add, - for Remove, ↻ for Modify
	ChangeSetASCII map[cloudformation.ChangeAction]string = map[cloudformation.ChangeAction]string{
		cloudformation.ChangeActionAdd:    "+",
		cloudformation.ChangeActionRemove: "-",
		cloudformation.ChangeActionModify: "↻",
	}
)

//StackOperation is the cloudFormation type of stack operations
type StackOperation string

const (
	stackNotFound   string = "does not exist"
	unknownEndpoint string = "unknown endpoint, could not resolve endpoint"

	//StackOperationUpdate is the enum value for Stack Operation of update
	StackOperationUpdate StackOperation = "update"

	//StackOperationCreate is the enum value for Stack Operation of create
	StackOperationCreate StackOperation = "create"

	//StackOperationDelete is the enum value for Stack Operation of delete
	StackOperationDelete StackOperation = "delete"
)

func getClient() *cloudformation.Client {
	if cfnClient == nil {
		cfg, err := external.LoadDefaultAWSConfig()
		if err != nil {
			panic(colors.ERROR + "unable to load SDK config, " + err.Error())
		}

		cfnClient = cloudformation.New(cfg)
	}

	return cfnClient
}

//CreateChanges creates a change set, waits for it to complete creating, then describes the change set.
func CreateChanges(info data.StackInfo, template []byte, exists bool) (*cloudformation.DescribeChangeSetResponse, error) {
	err := createChangeSet(info, template, exists)
	if err != nil {
		return nil, err
	}

	err = waitForChangeSet(info)
	if err != nil {
		return nil, err
	}

	changes, err := describeChangeSet(info)

	return changes, err
}

func createChangeSet(info data.StackInfo, template []byte, exists bool) error {
	stringTemplate := string(template)
	capabilities := []cloudformation.Capability{
		cloudformation.CapabilityCapabilityAutoExpand,
		cloudformation.CapabilityCapabilityIam,
		cloudformation.CapabilityCapabilityNamedIam,
	}

	changeSetType := cloudformation.ChangeSetTypeCreate
	if exists {
		changeSetType = cloudformation.ChangeSetTypeUpdate
	}

	client := getClient()

	input := cloudformation.CreateChangeSetInput{
		ChangeSetName: &info.ChangeSetName,
		StackName:     &info.StackName,
		TemplateBody:  &stringTemplate,
		ChangeSetType: changeSetType,
		Capabilities:  capabilities,
	}

	req := client.CreateChangeSetRequest(&input)

	_, err := req.Send(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func waitForChangeSet(info data.StackInfo) error {
	client := getClient()

	input := cloudformation.DescribeChangeSetInput{
		StackName:     &info.StackName,
		ChangeSetName: &info.ChangeSetName,
	}

	err := client.WaitUntilChangeSetCreateComplete(context.Background(), &input)

	if err != nil {
		return err
	}

	return nil
}

// ExecuteChangeSet executes the given change set
func ExecuteChangeSet(info data.StackInfo) error {
	input := cloudformation.ExecuteChangeSetInput{
		StackName:     &info.StackName,
		ChangeSetName: &info.ChangeSetName,
	}

	client := getClient()

	req := client.ExecuteChangeSetRequest(&input)

	_, err := req.Send(context.Background())

	return err
}

func describeChangeSet(info data.StackInfo) (*cloudformation.DescribeChangeSetResponse, error) {
	input := cloudformation.DescribeChangeSetInput{
		StackName:     &info.StackName,
		ChangeSetName: &info.ChangeSetName,
	}

	client := getClient()

	req := client.DescribeChangeSetRequest(&input)

	return req.Send(context.Background())
}

func getChanges(info data.StackInfo) ([]cloudformation.Change, error) {
	changeSet, err := describeChangeSet(info)

	return changeSet.Changes, err
}

//GetStack retrieves the information for the given stack name
func GetStack(stackName string) (*cloudformation.DescribeStacksResponse, error) {
	input := cloudformation.DescribeStacksInput{
		StackName: &stackName,
	}

	client := getClient()

	req := client.DescribeStacksRequest(&input)

	stack, err := req.Send(context.Background())
	if err != nil {
		return nil, err
	}

	return stack, err
}

// DetermineIfStackExists pulls a stack via the stackName and determines if it exists. If it is in a "review in progress" state, it counts as not existing
func DetermineIfStackExists(stackName string) (bool, error) {
	stack, err := GetStack(stackName)

	if err != nil {
		s := err.Error()

		if strings.Contains(s, stackNotFound) {
			return false, nil
		}

		return false, err
	}

	exists := true

	if stack.Stacks[0].StackStatus == cloudformation.StackStatusReviewInProgress {
		exists = false
	}

	return exists, nil
}

// DeleteStack deletes the stack given a stack name
func DeleteStack(info data.StackInfo) error {
	input := cloudformation.DeleteStackInput{
		StackName: &info.StackName,
	}

	client := getClient()

	req := client.DeleteStackRequest(&input)

	_, err := req.Send(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// GetStackEvents gets all the events from a particular CloudFormation stack
func GetStackEvents(info data.StackInfo) cloudformation.DescribeStackEventsPaginator {
	input := cloudformation.DescribeStackEventsInput{
		StackName: &info.StackID,
	}

	client := getClient()

	req := client.DescribeStackEventsRequest(&input)

	events := cloudformation.NewDescribeStackEventsPaginator(req)

	return events
}

// GetStackResources get all the resources that exist ina particular CloudFormation stack
func GetStackResources(info data.StackInfo) cloudformation.ListStackResourcesPaginator {
	input := cloudformation.ListStackResourcesInput{
		StackName: &info.StackName,
	}

	client := getClient()

	req := client.ListStackResourcesRequest(&input)

	resources := cloudformation.NewListStackResourcesPaginator(req)

	return resources
}

// VerifyAWSCredentials verifies AWS credentials are properly configured by running a List Stack command and analyzing errors for common issues with credentials
func VerifyAWSCredentials() error {
	input := cloudformation.ListStacksInput{}

	client := getClient()

	req := client.ListStacksRequest(&input)

	_, err := req.Send(context.Background())
	if err != nil {
		err = handleCredentialsError(err)

		return err
	}

	return nil
}

func handleCredentialsError(err error) error {
	strErr := err.Error()
	var msg string

	if strings.Contains(strErr, unknownEndpoint) {
		msg = colors.ERROR + "Unable to verify AWS credentials. Ensure your configuration is correct. \n \n" + colors.DOCS + "https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html"
	}

	return errors.New(msg)
}
