package matchers

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/onsi/gomega/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// HaveCfnStack returns a GoMega matcher that will check for the existence of an cloudformatioin stack
func HaveCfnStack(expectedStackName string) types.GomegaMatcher {
	return &haveCfnStackMatcher{expectedStackName: expectedStackName}
}

type haveCfnStackMatcher struct {
	expectedStackName string
	stackNotFound     bool
}

func (m *haveCfnStackMatcher) Match(actual interface{}) (success bool, err error) {
	if actual == nil {
		return false, errors.New("input is nil")
	}

	if reflect.TypeOf(actual).String() != "*session.Session" {
		return false, errors.New("not a AWS session")
	}

	cfn := cloudformation.New(actual.(*session.Session))

	input := &cloudformation.ListStackResourcesInput{
		StackName: aws.String(m.expectedStackName),
	}
	_, err = cfn.ListStackResources(input)

	//TODO: test for not found

	if err != nil {
		// Check if its a not found error: ResourceNotFoundException
		//if !strings.Contains(err.Error(), awseks.ErrCodeResourceNotFoundException) {
		return false, err
		//}

		//m.clusterNotFound = true
		//return false, nil
	}

	return true, nil
}

func (m *haveCfnStackMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected to find a Cloudformation stack named %s but it wasn't found", m.expectedStackName)
}

func (m *haveCfnStackMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected NOT to find a Cloudformation stack named %s but it found", m.expectedStackName)
}
