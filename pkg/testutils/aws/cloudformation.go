package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

const (
	errorMessageTemplate = "Stack with id %s does not exist"
)

// StackExists checks to see if a CloudFormation stack exists
func StackExists(stackName string, session *session.Session) (bool, error) {
	cfn := cloudformation.New(session)

	input := &cloudformation.ListStackResourcesInput{
		StackName: aws.String(stackName),
	}
	_, err := cfn.ListStackResources(input)

	if err != nil {
		// Check if its a not found error
		errorMessage := fmt.Sprintf(errorMessageTemplate, stackName)
		if !strings.Contains(err.Error(), errorMessage) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// DeleteStack deletes a cloudformation stack
func DeleteStack(stackName string, session *session.Session) error {
	cfn := cloudformation.New(session)

	input := &cloudformation.DeleteStackInput{
		StackName: aws.String(stackName),
	}

	_, err := cfn.DeleteStack(input)

	return err

}
