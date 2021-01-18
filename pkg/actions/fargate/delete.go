package fargate

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
)

// DeleteProfile drains and delete the Fargate profile with the provided name.
func (m *Manager) DeleteProfile(name string, waitForDeletion bool) error {
	if name == "" {
		return errors.New("invalid Fargate profile name: empty")
	}
	out, err := m.api.DeleteFargateProfile(deleteRequest(m.clusterName, name))
	logger.Debug("Fargate profile: delete request: received: %#v", out)
	if err != nil {
		return errors.Wrapf(err, "failed to delete Fargate profile %q", name)
	}
	if waitForDeletion {
		return m.waitForDeletion(name)
	}
	return nil
}

func (m *Manager) waitForDeletion(name string) error {
	// Clone this manager's policy to ensure this method is re-entrant/thread-safe:
	retryPolicy := m.retryPolicy.Clone()
	for !retryPolicy.Done() {
		names, err := m.ListProfiles()
		if err != nil {
			return err
		}
		if !contains(names, name) {
			return nil
		}
		time.Sleep(retryPolicy.Duration())
	}
	return fmt.Errorf("timed out while waiting for Fargate profile %q's deletion", name)
}

func deleteRequest(clusterName string, profileName string) *eks.DeleteFargateProfileInput {
	request := &eks.DeleteFargateProfileInput{
		ClusterName:        &clusterName,
		FargateProfileName: &profileName,
	}
	logger.Debug("Fargate profile: delete request: sending: %#v", request)
	return request
}

func contains(array []*string, target string) bool {
	for _, value := range array {
		if value != nil && *value == target {
			return true
		}
	}
	return false
}

// IsUnauthorizedError reports whether the error is an authorization error
// Unauthorized errors are of the form:
//   AccessDeniedException: Account <account> is not authorized to use this service
func IsUnauthorizedError(err error) bool {
	awsErr, ok := errors.Cause(err).(awserr.Error)
	if !ok {
		return false
	}
	return awsErr.Code() == "AccessDeniedException"
}
