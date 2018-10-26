package matchers

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/onsi/gomega/types"
	"github.com/weaveworks/eksctl/pkg/testutils/aws"
)

// HaveExistingStack returns a GoMega matcher that will check for the existence of an cloudformation stack
func HaveExistingStack(expectedStackName string) types.GomegaMatcher {
	return &existingStack{expectedStackName: expectedStackName}
}

type existingStack struct {
	expectedStackName string
	stackNotFound     bool
}

func (m *existingStack) Match(actual interface{}) (success bool, err error) {
	if actual == nil {
		return false, errors.New("input is nil")
	}

	if reflect.TypeOf(actual).String() != "*session.Session" {
		return false, errors.New("not a AWS session")
	}

	found, err := aws.StackExists(m.expectedStackName, actual.(*session.Session))

	if err != nil {
		return false, err
	}

	m.stackNotFound = !found
	return found, nil
}

func (m *existingStack) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected to find a Cloudformation stack named %s but it wasn't found", m.expectedStackName)
}

func (m *existingStack) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected NOT to find a Cloudformation stack named %s but it found", m.expectedStackName)
}
