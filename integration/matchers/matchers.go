package matchers

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/onsi/gomega/types"
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

	if v := reflect.TypeOf(actual).String(); v != "aws.Config" {
		return false, fmt.Errorf("%s was not of type aws.Config", v)
	}

	found, err := stackExists(m.expectedStackName, actual.(aws.Config))
	if err != nil {
		return false, err
	}

	m.stackNotFound = !found
	return found, nil
}

func (m *existingStack) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected to find a Cloudformation stack named %s but it was NOT found", m.expectedStackName)
}

func (m *existingStack) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected NOT to find a Cloudformation stack named %s but it was found", m.expectedStackName)
}

// HaveExistingCluster returns a GoMega matcher that will check for the existence of an EKS cluster
func HaveExistingCluster(expectedName string, expectedStatus string, expectedVersion string) types.GomegaMatcher {
	return &existingCluster{expectedName: expectedName, expectedStatus: expectedStatus, expectedVersion: expectedVersion}
}

type existingCluster struct {
	expectedName    string
	expectedStatus  string
	expectedVersion string

	clusterNotFound bool
	versionMismatch bool
	statusMismatch  bool

	actualVersion string
	actualStatus  string
}

func (m *existingCluster) Match(actual interface{}) (success bool, err error) {
	if actual == nil {
		return false, errors.New("input is nil")
	}

	if v := reflect.TypeOf(actual).String(); v != "aws.Config" {
		return false, fmt.Errorf("%s was not of type aws.Config", v)
	}

	client := eks.NewFromConfig(actual.(aws.Config))

	input := &eks.DescribeClusterInput{
		Name: aws.String(m.expectedName),
	}
	output, err := client.DescribeCluster(context.Background(), input)

	if err != nil {
		// Check if it's a not found error: ResourceNotFoundException
		var notFoundErr *ekstypes.ResourceNotFoundException
		if !errors.As(err, &notFoundErr) {
			return false, err
		}

		m.clusterNotFound = true
		return false, nil
	}

	m.actualStatus = string(output.Cluster.Status)
	if m.actualStatus != m.expectedStatus {
		m.statusMismatch = true
		return false, nil
	}

	m.actualVersion = *output.Cluster.Version
	if m.actualVersion != m.expectedVersion {
		m.versionMismatch = true
		return false, nil
	}

	return true, nil
}

func (m *existingCluster) FailureMessage(actual interface{}) (message string) {
	if m.statusMismatch {
		return fmt.Sprintf("Expected EKS cluster status: %s to equal actual EKS cluster status: %s", m.expectedStatus, m.actualStatus)
	}
	if m.versionMismatch {
		return fmt.Sprintf("Expected EKS cluster version: %s to equal actual EKS cluster version: %s", m.expectedVersion, m.actualVersion)
	}

	return fmt.Sprintf("Expected to find a cluster named %s but it was NOT found", m.expectedName)
}

func (m *existingCluster) NegatedFailureMessage(_ interface{}) (message string) {
	if m.statusMismatch {
		return fmt.Sprintf("Expected EKS cluster status: %s NOT to equal actual EKS cluster status: %s", m.expectedStatus, m.actualStatus)
	}
	if m.versionMismatch {
		return fmt.Sprintf("Expected EKS cluster version: %s NOT to equal actual EKS cluster version: %s", m.expectedVersion, m.actualVersion)
	}

	return fmt.Sprintf("Expected NOT to find a cluster named %s but it was found", m.expectedName)
}
