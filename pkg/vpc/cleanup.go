package vpc

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
)

func fmtSecurityGroupNameRegexForCluster(name string) string {
	const ourSecurityGroupNameRegexFmt = "^eksctl-%s-(cluster|nodegroup)-.+$"
	return fmt.Sprintf(ourSecurityGroupNameRegexFmt, name)
}

func findDanglingENIs(ctx context.Context, ec2API awsapi.EC2, spec *api.ClusterConfig) ([]string, error) {
	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{spec.VPC.ID},
			},
			{
				Name:   aws.String("status"),
				Values: []string{"available"},
			},
		},
	}

	securityGroupRE, err := regexp.Compile(fmtSecurityGroupNameRegexForCluster(spec.Metadata.Name))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create security group regex")
	}

	var eniIDs []string

	paginator := ec2.NewDescribeNetworkInterfacesPaginator(ec2API, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to list dangling network interfaces in %q: %w", spec.VPC.ID, err)
		}

		for _, eni := range output.NetworkInterfaces {
			id := *eni.NetworkInterfaceId
			for _, sg := range eni.Groups {
				if securityGroupRE.MatchString(*sg.GroupName) {
					logger.Debug("found %q, which belongs to our security group %q (%s)", id, *sg.GroupName, *sg.GroupId)
					eniIDs = append(eniIDs, id)
					break
				}
				logger.Debug("found %q, but it belongs to security group %q (%s), which does not appear to be ours", id, *sg.GroupName, *sg.GroupId)

			}
		}
	}

	return eniIDs, nil
}

// CleanupNetworkInterfaces finds and deletes any dangling ENIs
func CleanupNetworkInterfaces(ctx context.Context, ec2API awsapi.EC2, spec *api.ClusterConfig) error {
	eniIDs, err := findDanglingENIs(ctx, ec2API, spec)
	if err != nil {
		return err
	}
	for _, eniID := range eniIDs {
		input := &ec2.DeleteNetworkInterfaceInput{
			NetworkInterfaceId: &eniID,
		}
		if _, err := ec2API.DeleteNetworkInterface(ctx, input); err != nil {
			return errors.Wrapf(err, "unable to delete network interface %q", eniID)
		}
		logger.Debug("deleted network interface %q", eniID)
	}
	return nil
}
