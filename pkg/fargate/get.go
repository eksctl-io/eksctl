package fargate

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
)

// ReadProfile reads the Fargate profile corresponding to the provided name if
// it exists.
func (c *Client) ReadProfile(ctx context.Context, name string) (*api.FargateProfile, error) {
	out, err := c.api.DescribeFargateProfile(ctx, describeRequest(c.clusterName, name))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get Fargate profile %q", name)
	}
	logger.Debug("Fargate profile: describe request: received: %#v", out)
	return toFargateProfile(out.FargateProfile), nil
}

// ReadProfiles reads all existing Fargate profiles.
func (c *Client) ReadProfiles(ctx context.Context) ([]*api.FargateProfile, error) {
	names, err := c.ListProfiles(ctx)
	if err != nil {
		return nil, err
	}
	profiles := []*api.FargateProfile{}
	for _, name := range names {
		profile, err := c.ReadProfile(ctx, name)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

// ListProfiles lists all existing Fargate profiles.
func (c *Client) ListProfiles(ctx context.Context) (fargateProfileNames []string, err error) {
	request := &eks.ListFargateProfilesInput{
		ClusterName: &c.clusterName,
	}
	var out *eks.ListFargateProfilesOutput
	for out, err = c.sendListRequestOnce(ctx, request); err == nil; out, err = c.sendListRequestOnce(ctx, request) {
		fargateProfileNames = append(fargateProfileNames, out.FargateProfileNames...)

		if request.NextToken = out.NextToken; request.NextToken == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	logger.Debug("Fargate profile: list request: received a total of %v profile(s)", len(fargateProfileNames))
	return fargateProfileNames, nil
}

func (c *Client) sendListRequestOnce(ctx context.Context, request *eks.ListFargateProfilesInput) (*eks.ListFargateProfilesOutput, error) {
	logger.Debug("Fargate profile: list request: sending: %#v", request)
	out, err := c.api.ListFargateProfiles(ctx, request)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get Fargate profile(s) for cluster %q", c.clusterName)
	}
	logger.Debug("Fargate profile: list request: received %v profile(s): %#v", len(out.FargateProfileNames), out)
	return out, nil
}

func toFargateProfile(in *ekstypes.FargateProfile) *api.FargateProfile {
	return &api.FargateProfile{
		Name:                *in.FargateProfileName,
		Selectors:           toSelectors(in.Selectors),
		PodExecutionRoleARN: strings.EmptyIfNil(in.PodExecutionRoleArn),
		Subnets:             in.Subnets,
		Tags:                in.Tags,
		Status:              string(in.Status),
	}
}
