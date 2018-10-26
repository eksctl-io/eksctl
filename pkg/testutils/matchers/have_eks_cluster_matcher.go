package matchers

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/onsi/gomega/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awseks "github.com/aws/aws-sdk-go/service/eks"
)

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

	region string
}

func (m *existingCluster) Match(actual interface{}) (success bool, err error) {
	if actual == nil {
		return false, errors.New("input is nil")
	}

	if reflect.TypeOf(actual).String() != "*session.Session" {
		return false, errors.New("not a AWS session")
	}

	eks := awseks.New(actual.(*session.Session))

	input := &awseks.DescribeClusterInput{
		Name: aws.String(m.expectedName),
	}
	output, err := eks.DescribeCluster(input)

	if err != nil {
		// Check if its a not found error: ResourceNotFoundException
		if !strings.Contains(err.Error(), awseks.ErrCodeResourceNotFoundException) {
			return false, err
		}

		m.clusterNotFound = true
		return false, nil
	}

	m.actualStatus = *output.Cluster.Status
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

	return fmt.Sprintf("Expected to find a cluster named %s but it wasn't found", m.expectedName)
}

func (m *existingCluster) NegatedFailureMessage(_ interface{}) (message string) {
	if m.statusMismatch {
		return fmt.Sprintf("Expected EKS cluster status: %s NOT to equal actual EKS cluster status: %s", m.expectedStatus, m.actualStatus)
	}
	if m.versionMismatch {
		return fmt.Sprintf("Expected EKS cluster version: %s NOT to equal actual EKS cluster version: %s", m.expectedVersion, m.actualVersion)
	}

	return fmt.Sprintf("Expected NOT to find a cluster named %s but it found", m.expectedName)
}
