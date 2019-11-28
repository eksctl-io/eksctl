package fargate

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/retry"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
)

// DefaultWaitTimeout is the default maximum time to wait for long running
// operations.
const DefaultWaitTimeout = 5 * time.Minute

// NewClient returns a new Fargate client.
func NewClient(clusterName string, api eksiface.EKSAPI) *Client {
	return NewClientWithWaitTimeout(clusterName, api, DefaultWaitTimeout)
}

// NewClientWithWaitTimeout returns a new Fargate client configured with the
// provided wait timeout for blocking/waiting operations.
func NewClientWithWaitTimeout(clusterName string, api eksiface.EKSAPI, waitTimeout time.Duration) *Client {
	return NewClientWithRetryPolicy(clusterName, api, &retry.TimingOutExponentialBackoff{
		Timeout:  waitTimeout,
		TimeUnit: time.Second,
	})
}

// NewClientWithRetryPolicy returns a new Fargate client configured with the
// provided retry policy for blocking/waiting operations.
func NewClientWithRetryPolicy(clusterName string, api eksiface.EKSAPI, retryPolicy retry.Policy) *Client {
	return &Client{
		clusterName: clusterName,
		api:         api,
		retryPolicy: retryPolicy,
	}
}

// Client wraps around an EKS API client to expose high-level methods.
type Client struct {
	clusterName string
	api         eksiface.EKSAPI
	retryPolicy retry.Policy
}

// CreateProfile creates the provided Fargate profile.
func (c Client) CreateProfile(profile *api.FargateProfile, waitForCreation bool) error {
	if profile == nil {
		return errors.New("invalid Fargate profile: nil")
	}
	logger.Debug("Fargate profile: create request input: %#v", profile)
	out, err := c.api.CreateFargateProfile(createRequest(c.clusterName, profile))
	logger.Debug("Fargate profile: create request: received: %#v", out)
	if err != nil {
		return errors.Wrapf(err, "failed to create Fargate profile %q in cluster %q", profile.Name, c.clusterName)
	}
	if waitForCreation {
		return c.waitForCreation(profile.Name)
	}
	return nil
}

// ReadProfile reads the Fargate profile corresponding to the provided name if
// it exists.
func (c Client) ReadProfile(name string) (*api.FargateProfile, error) {
	out, err := c.api.DescribeFargateProfile(describeRequest(c.clusterName, name))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get EKS cluster %q's Fargate profile %q", c.clusterName, name)
	}
	logger.Debug("Fargate profile: describe request: received: %#v", out)
	return toFargateProfile(out.FargateProfile), nil
}

// ReadProfiles reads all existing Fargate profiles.
func (c Client) ReadProfiles() ([]*api.FargateProfile, error) {
	names, err := c.ListProfiles()
	if err != nil {
		return nil, err
	}
	profiles := []*api.FargateProfile{}
	for _, name := range names {
		profile, err := c.ReadProfile(*name)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

// ListProfiles lists all existing Fargate profiles.
func (c Client) ListProfiles() ([]*string, error) {
	out, err := c.api.ListFargateProfiles(listRequest(c.clusterName))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get EKS cluster %q's Fargate profile(s)", c.clusterName)
	}
	logger.Debug("Fargate profile: list request: got %v profile(s): %#v", len(out.FargateProfileNames), out)
	return out.FargateProfileNames, nil
}

// DeleteProfile drains and delete the Fargate profile with the provided name.
func (c Client) DeleteProfile(name string, waitForDeletion bool) error {
	if name == "" {
		return errors.New("invalid Fargate profile name: empty")
	}
	out, err := c.api.DeleteFargateProfile(deleteRequest(c.clusterName, name))
	logger.Debug("Fargate profile: delete request: received: %#v", out)
	if err != nil {
		return errors.Wrapf(err, "failed to delete Fargate profile %q from cluster %q", name, c.clusterName)
	}
	if waitForDeletion {
		return c.waitForDeletion(name)
	}
	return nil
}

func (c Client) waitForCreation(name string) error {
	// Clone this client's policy to ensure this method is re-entrant/thread-safe:
	retryPolicy := c.retryPolicy.Clone()
	for !retryPolicy.Done() {
		out, err := c.api.DescribeFargateProfile(describeRequest(c.clusterName, name))
		if err != nil {
			return errors.Wrapf(err, "failed while waiting for Fargate profile %q's creation", name)
		}
		logger.Debug("Fargate profile: describe request: received: %#v", out)
		if created(out) {
			return nil
		}
		time.Sleep(retryPolicy.Duration())
	}
	return fmt.Errorf("timed out while waiting for Fargate profile %q's creation", name)
}

func created(out *eks.DescribeFargateProfileOutput) bool {
	return out != nil &&
		out.FargateProfile != nil &&
		out.FargateProfile.Status != nil &&
		*out.FargateProfile.Status == eks.FargateProfileStatusActive
}

func (c Client) waitForDeletion(name string) error {
	// Clone this client's policy to ensure this method is re-entrant/thread-safe:
	retryPolicy := c.retryPolicy.Clone()
	for !retryPolicy.Done() {
		names, err := c.ListProfiles()
		if err != nil {
			return err
		}
		if !contains(names, name) {
			return nil
		}
		time.Sleep(retryPolicy.Duration())
	}
	return fmt.Errorf("deleting of Fargate profile %q timed out", name)
}

func contains(array []*string, target string) bool {
	for _, value := range array {
		if value != nil && *value == target {
			return true
		}
	}
	return false
}

func createRequest(clusterName string, profile *api.FargateProfile) *eks.CreateFargateProfileInput {
	request := &eks.CreateFargateProfileInput{
		ClusterName:         &clusterName,
		FargateProfileName:  &profile.Name,
		Selectors:           toSelectorPointers(profile.Selectors),
		PodExecutionRoleArn: strings.NilIfEmpty(profile.PodExecutionRoleARN),
		Subnets:             strings.NilPointersArrayIfEmpty(strings.ToPointersArray(profile.Subnets)),
	}
	logger.Debug("Fargate profile: create request: sending: %#v", request)
	return request
}

func describeRequest(clusterName string, profileName string) *eks.DescribeFargateProfileInput {
	request := &eks.DescribeFargateProfileInput{
		ClusterName:        &clusterName,
		FargateProfileName: &profileName,
	}
	logger.Debug("Fargate profile: describe request: sending: %#v", request)
	return request
}

func listRequest(clusterName string) *eks.ListFargateProfilesInput {
	request := &eks.ListFargateProfilesInput{
		ClusterName: &clusterName,
	}
	logger.Debug("Fargate profile: list request: sending: %#v", request)
	return request
}

func deleteRequest(clusterName string, profileName string) *eks.DeleteFargateProfileInput {
	request := &eks.DeleteFargateProfileInput{
		ClusterName:        &clusterName,
		FargateProfileName: &profileName,
	}
	logger.Debug("Fargate profile: delete request: sending: %#v", request)
	return request
}

func toFargateProfile(in *eks.FargateProfile) *api.FargateProfile {
	return &api.FargateProfile{
		Name:                *in.FargateProfileName,
		Selectors:           toSelectors(in.Selectors),
		PodExecutionRoleARN: strings.EmptyIfNil(in.PodExecutionRoleArn),
		Subnets:             strings.ToValuesArray(in.Subnets),
	}
}

func toSelectorPointers(in []api.FargateProfileSelector) []*eks.FargateProfileSelector {
	out := make([]*eks.FargateProfileSelector, len(in))
	for i, selector := range in {
		out[i] = &eks.FargateProfileSelector{
			Namespace: strings.Pointer(selector.Namespace),
			Labels:    strings.NilPointersMapIfEmpty(strings.ToPointersMap(selector.Labels)),
		}
	}
	return out
}

func toSelectors(in []*eks.FargateProfileSelector) []api.FargateProfileSelector {
	out := make([]api.FargateProfileSelector, len(in))
	for i, selector := range in {
		out[i] = api.FargateProfileSelector{
			Namespace: *selector.Namespace,
			Labels:    strings.ToValuesMap(selector.Labels),
		}
	}
	return out
}
