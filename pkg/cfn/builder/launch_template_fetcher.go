package builder

import (
	"context"

	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type launchTemplateFetcher interface {
	DescribeLaunchTemplateVersions(ctx context.Context, params *ec2.DescribeLaunchTemplateVersionsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeLaunchTemplateVersionsOutput, error)
}

// LaunchTemplateFetcher fetches launch template data
type LaunchTemplateFetcher struct {
	fetcher launchTemplateFetcher
}

// NewLaunchTemplateFetcher creates a new LaunchTemplateFetcher
func NewLaunchTemplateFetcher(fetcher launchTemplateFetcher) *LaunchTemplateFetcher {
	return &LaunchTemplateFetcher{fetcher: fetcher}
}

// Fetch fetches the specified launch template
func (l *LaunchTemplateFetcher) Fetch(ctx context.Context, launchTemplate *api.LaunchTemplate) (*ec2types.ResponseLaunchTemplateData, error) {
	input := &ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateId: aws.String(launchTemplate.ID),
	}
	if version := launchTemplate.Version; version != nil {
		input.Versions = []string{*version}
	} else {
		input.Versions = []string{"$Default"}
	}

	output, err := l.fetcher.DescribeLaunchTemplateVersions(ctx, input)
	if err != nil {
		return nil, err
	}
	if len(output.LaunchTemplateVersions) != 1 {
		return nil, errors.Errorf("failed to find launch template with ID %q", launchTemplate.ID)
	}

	return output.LaunchTemplateVersions[0].LaunchTemplateData, nil
}
