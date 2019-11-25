package fargate

import (
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
)

// NewClient returns a new Fargate client.
func NewClient(clusterName string, api eksiface.EKSAPI) *Client {
	return &Client{
		clusterName: clusterName,
		api:         api,
	}
}

// Client wraps around an EKS API client to expose high-level methods.
type Client struct {
	clusterName string
	api         eksiface.EKSAPI
}

// CreateProfile creates the provided Fargate profile.
func (c Client) CreateProfile(profile *api.FargateProfile) error {
	if profile == nil {
		return errors.New("invalid Fargate profile: nil")
	}
	logger.Debug("Fargate profile: create request input: %#v", profile)
	out, err := c.api.CreateFargateProfile(createRequest(c.clusterName, profile))
	if err != nil {
		return errors.Wrapf(err, "failed to create Fargate profile \"%v\" in cluster \"%v\"", profile.Name, c.clusterName)
	}
	logger.Debug("successfully created Fargate profile: %s", out)
	return nil
}

// ReadProfile reads the Fargate profile corresponding to the provided name if
// it exists.
func (c Client) ReadProfile(name string) (*api.FargateProfile, error) {
	out, err := c.api.DescribeFargateProfile(describeRequest(c.clusterName, name))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get EKS cluster \"%s\"'s Fargate profile \"%s\"", c.clusterName, name)
	}
	return toFargateProfile(out.FargateProfile), nil
}

// ReadProfiles reads all existing Fargate profiles.
func (c Client) ReadProfiles() ([]*api.FargateProfile, error) {
	profiles := []*api.FargateProfile{}
	out, err := c.api.ListFargateProfiles(listRequest(c.clusterName))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get EKS cluster \"%s\"'s Fargate profile(s)", c.clusterName)
	}
	logger.Debug("Fargate profile: list request: got %v profile(s): %#v", len(out.FargateProfileNames), out)
	for _, name := range out.FargateProfileNames {
		profile, err := c.ReadProfile(*name)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

// DeleteProfile drains and delete the Fargate profile with the provided name.
func (c Client) DeleteProfile(name string) error {
	if name == "" {
		return errors.New("invalid Fargate profile name: empty")
	}
	_, err := c.api.DeleteFargateProfile(&eks.DeleteFargateProfileInput{
		ClusterName:        &c.clusterName,
		FargateProfileName: &name,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to delete Fargate profile \"%v\" from cluster \"%v\"", name, c.clusterName)
	}
	return nil
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
