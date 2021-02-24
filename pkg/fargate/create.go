package fargate

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/names"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
)

// CreateOptions groups the parameters required to create a Fargate profile.
type CreateOptions struct {
	ProfileName              string
	ProfileSelectorNamespace string
	// +optional
	ProfileSelectorLabels map[string]string
	// +optional
	Tags map[string]string
}

// Validate validates this Options object's fields.
func (o *CreateOptions) Validate() error {
	if strings.HasPrefix(o.ProfileName, api.ReservedProfileNamePrefix) {
		return fmt.Errorf("invalid Fargate profile: name should NOT start with %q", api.ReservedProfileNamePrefix)
	}
	if o.ProfileSelectorNamespace == "" {
		return errors.New("invalid Fargate profile: empty selector namespace")
	}
	return nil
}

// ToFargateProfile creates a FargateProfile object from this Options object.
func (o CreateOptions) ToFargateProfile() *api.FargateProfile {
	return &api.FargateProfile{
		Name: names.ForFargateProfile(o.ProfileName),
		Selectors: []api.FargateProfileSelector{
			{
				Namespace: o.ProfileSelectorNamespace,
				Labels:    o.ProfileSelectorLabels,
			},
		},
		Tags: o.Tags,
	}
}

// CreateProfile creates the provided Fargate profile.
func (c *Client) CreateProfile(profile *api.FargateProfile, waitForCreation bool) error {
	if profile == nil {
		return errors.New("invalid Fargate profile: nil")
	}
	logger.Debug("Fargate profile: create request input: %#v", profile)
	out, err := c.api.CreateFargateProfile(createRequest(c.clusterName, profile))
	logger.Debug("Fargate profile: create request: received: %#v", out)
	if err != nil {
		return errors.Wrapf(err, "failed to create Fargate profile %q", profile.Name)
	}
	if waitForCreation {
		return c.waitForCreation(profile.Name)
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
		Tags:                strings.NilPointersMapIfEmpty(strings.ToPointersMap(profile.Tags)),
	}
	logger.Debug("Fargate profile: create request: sending: %#v", request)
	return request
}

func (c *Client) waitForCreation(name string) error {
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
