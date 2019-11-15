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
	out, err := c.api.CreateFargateProfile(createRequest(c.clusterName, profile))
	if err != nil {
		return errors.Wrapf(err, "failed to create Fargate profile \"%v\" in cluster \"%v\"", profile.Name, c.clusterName)
	}
	logger.Debug("successfully created Fargate profile: %s", out)
	return nil
}

// ReadProfiles reads and returns all existing Fargate profiles.
func (c Client) ReadProfiles() ([]*api.FargateProfile, error) {
	profiles := []*api.FargateProfile{}
	var nextToken *string // used for "pagination" of the retrieval.
	for {
		out, err := c.api.ListFargateProfiles(&eks.ListFargateProfilesInput{
			ClusterName: &c.clusterName,
			NextToken:   nextToken,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get EKS cluster \"%s\"'s Fargate profile(s) (current token: %v)", c.clusterName, nextToken)
		}
		nextToken = out.NextToken
		if out.FargateProfiles == nil || len(out.FargateProfiles) == 0 {
			break
		}
		profiles = append(profiles, toFargateProfiles(out.FargateProfiles)...)
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
		PodExecutionRoleARN: strings.NilIfEmpty(profile.PodExecutionRoleARN),
		Subnets:             strings.NilPointersArrayIfEmpty(strings.ToPointersArray(profile.Subnets)),
	}
}

func toFargateProfiles(in []*eks.FargateProfile) []*api.FargateProfile {
	out := make([]*api.FargateProfile, len(in))
	for i := range in {
		out[i] = toFargateProfile(in[i])
	}
	return out
}

func toFargateProfile(in *eks.FargateProfile) *api.FargateProfile {
	return &api.FargateProfile{
		Name:                *in.FargateProfileName,
		Selectors:           toSelectors(in.Selectors),
		PodExecutionRoleARN: strings.EmptyIfNil(in.PodExecutionRoleARN),
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
