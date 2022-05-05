package fargate

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

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
func (c *Client) CreateProfile(ctx context.Context, profile *api.FargateProfile, waitForCreation bool) error {
	if profile == nil {
		return errors.New("invalid Fargate profile: nil")
	}
	logger.Debug("Fargate profile: create request input: %#v", profile)
	out, err := c.api.CreateFargateProfile(ctx, createRequest(c.clusterName, profile))
	logger.Debug("Fargate profile: create request: received: %#v", out)
	if err != nil {
		return errors.Wrapf(err, "failed to create Fargate profile %q", profile.Name)
	}
	if waitForCreation {
		return c.waitForCreation(ctx, profile.Name)
	}
	return nil
}

func createRequest(clusterName string, profile *api.FargateProfile) *eks.CreateFargateProfileInput {
	request := &eks.CreateFargateProfileInput{
		ClusterName:         &clusterName,
		FargateProfileName:  &profile.Name,
		Selectors:           toSelectorPointers(profile.Selectors),
		PodExecutionRoleArn: strings.NilIfEmpty(profile.PodExecutionRoleARN),
		Subnets:             profile.Subnets,
		Tags:                profile.Tags,
	}
	logger.Debug("Fargate profile: create request: sending: %#v", request)
	return request
}

func (c *Client) waitForCreation(ctx context.Context, name string) error {
	// Clone this client's policy to ensure this method is re-entrant/thread-safe:
	retryPolicy := c.retryPolicy.Clone()
	for !retryPolicy.Done() {
		out, err := c.api.DescribeFargateProfile(ctx, describeRequest(c.clusterName, name))
		if err != nil {
			return errors.Wrapf(err, "failed while waiting for Fargate profile %q's creation", name)
		}
		logger.Debug("Fargate profile: describe request: received: %#v", out)
		if created(out) {
			return nil
		}
		timer := time.NewTimer(retryPolicy.Duration())
		select {
		case <-timer.C:

		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}

	}
	return fmt.Errorf("timed out while waiting for Fargate profile %q's creation", name)
}

func created(out *eks.DescribeFargateProfileOutput) bool {
	return out != nil &&
		out.FargateProfile != nil &&
		out.FargateProfile.Status == ekstypes.FargateProfileStatusActive
}

func toSelectorPointers(in []api.FargateProfileSelector) []ekstypes.FargateProfileSelector {
	out := make([]ekstypes.FargateProfileSelector, len(in))
	for i, selector := range in {
		out[i] = ekstypes.FargateProfileSelector{
			Namespace: strings.Pointer(selector.Namespace),
			Labels:    selector.Labels,
		}
	}
	return out
}
