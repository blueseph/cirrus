package main

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

var (
	cfnClient      *cloudformation.Client
	changeSetASCII map[cloudformation.ChangeAction]string = map[cloudformation.ChangeAction]string{
		cloudformation.ChangeActionAdd:    "+",
		cloudformation.ChangeActionRemove: "-",
		cloudformation.ChangeActionModify: "↻",
	}
)

type stackOperation string

const (
	stackNotFound   string = "does not exist"
	unknownEndpoint string = "unknown endpoint, could not resolve endpoint"

	update stackOperation = "update"
	create stackOperation = "create"
	delete stackOperation = "delete"
)

func getClient() *cloudformation.Client {
	if cfnClient == nil {
		cfg, err := external.LoadDefaultAWSConfig()
		if err != nil {
			panic(ERROR + "unable to load SDK config, " + err.Error())
		}

		cfnClient = cloudformation.New(cfg)
	}

	return cfnClient
}

func createChanges(stackName string, changeSetID string, template []byte, exists bool) (*cloudformation.DescribeChangeSetResponse, error) {
	err := createChangeSet(stackName, changeSetID, template, exists)
	if err != nil {
		return nil, err
	}

	err = waitForChangeSet(stackName, changeSetID)
	if err != nil {
		return nil, err
	}

	changes, err := describeChangeSet(stackName, changeSetID)

	return changes, err
}

func createChangeSet(stackName string, changeSetID string, template []byte, exists bool) error {
	stringTemplate := string(template)

	changeSetType := cloudformation.ChangeSetTypeCreate
	if exists {
		changeSetType = cloudformation.ChangeSetTypeUpdate
	}

	client := getClient()

	input := cloudformation.CreateChangeSetInput{
		ChangeSetName: &changeSetID,
		StackName:     &stackName,
		TemplateBody:  &stringTemplate,
		ChangeSetType: changeSetType,
	}

	req := client.CreateChangeSetRequest(&input)

	_, err := req.Send(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func waitForChangeSet(stackName string, changeSetID string) error {
	client := getClient()

	input := cloudformation.DescribeChangeSetInput{
		StackName:     &stackName,
		ChangeSetName: &changeSetID,
	}

	err := client.WaitUntilChangeSetCreateComplete(context.Background(), &input)

	if err != nil {
		return err
	}

	return nil
}

func executeChangeSet(stackName string, changeSetID string) error {
	input := cloudformation.ExecuteChangeSetInput{
		StackName:     &stackName,
		ChangeSetName: &changeSetID,
	}

	client := getClient()

	req := client.ExecuteChangeSetRequest(&input)

	_, err := req.Send(context.Background())

	return err
}

func describeChangeSet(stackName string, changeSetID string) (*cloudformation.DescribeChangeSetResponse, error) {
	input := cloudformation.DescribeChangeSetInput{
		StackName:     &stackName,
		ChangeSetName: &changeSetID,
	}

	client := getClient()

	req := client.DescribeChangeSetRequest(&input)

	return req.Send(context.Background())
}

func applyChangeSet(stackName string, changeSetID string) (*cloudformation.DescribeChangeSetResponse, error) {
	err := executeChangeSet(stackName, changeSetID)
	if err != nil {
		return nil, err
	}

	changeSet, err := describeChangeSet(stackName, changeSetID)
	if err != nil {
		return nil, err
	}

	return changeSet, nil
}

func getChanges(stackName string, changeSetID string) ([]cloudformation.Change, error) {
	changeSet, err := describeChangeSet(stackName, changeSetID)

	return changeSet.Changes, err
}

func getStack(stackName string) (*cloudformation.DescribeStacksResponse, error) {
	input := cloudformation.DescribeStacksInput{
		StackName: &stackName,
	}

	client := getClient()

	req := client.DescribeStacksRequest(&input)

	return req.Send(context.Background())
}

func determineIfStackExists(stackName string) (bool, error) {
	stack, err := getStack(stackName)

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

func deleteStack(stackName string) error {
	input := cloudformation.DeleteStackInput{
		StackName: &stackName,
	}

	client := getClient()

	req := client.DeleteStackRequest(&input)

	_, err := req.Send(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// we essentially call a list stack to verify credentials are correctly set up
func verifyAWSCredentials() error {
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
		msg = ERROR + "Unable to verify AWS credentials. Ensure your configuration is correct. \n \n" + DOCS + "https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html"
	}

	return errors.New(msg)
}
