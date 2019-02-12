package vpc

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
)

func fmtSecurityGroupNameRegexForCluster(name string) string {
	const ourSecurityGroupNameRegexFmt = "^eksctl-%s-(cluster|nodegroup)-.+$"
	return fmt.Sprintf(ourSecurityGroupNameRegexFmt, name)
}

func findAvailableNetworkInterfaces(provider api.ClusterProvider, spec *api.ClusterConfig) ([]string, error) {
	input := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{&spec.VPC.ID},
			},
			{
				Name:   aws.String("status"),
				Values: []*string{aws.String("available")},
			},
		},
	}
	output, err := provider.EC2().DescribeNetworkInterfaces(input)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to list dangling network interfaces in %q", spec.VPC.ID)
	}
	securityGroupRE, err := regexp.Compile(fmtSecurityGroupNameRegexForCluster(spec.Metadata.Name))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to list dangling network interfaces in %q", spec.VPC.ID)
	}
	networkInterfaces := []string{}
	for _, eni := range output.NetworkInterfaces {
		id := *eni.NetworkInterfaceId
		for _, sg := range eni.Groups {
			if securityGroupRE.MatchString(*sg.GroupName) {
				logger.Debug("found %q, which belongs to our security group %q (%s)", id, *sg.GroupName, *sg.GroupId)
				networkInterfaces = append(networkInterfaces, id)
				break
			} else {
				logger.Debug("found %q, but it belongs to security group %q (%s), which does not appear to be ours", id, *sg.GroupName, *sg.GroupId)
				break
			}
		}
	}
	return networkInterfaces, nil
}

func deleteNetworkInterfaces(provider api.ClusterProvider, networkInterfaces []string) error {
	for _, eni := range networkInterfaces {
		input := &ec2.DeleteNetworkInterfaceInput{
			NetworkInterfaceId: &eni,
		}
		if _, err := provider.EC2().DeleteNetworkInterface(input); err != nil {
			return errors.Wrapf(err, "unable to delete network interface %q", eni)
		}
		logger.Debug("deleted %q", eni)
	}
	return nil
}

// CleanupNetworkInterfaces find and deletes any dangling ENIs
func CleanupNetworkInterfaces(provider api.ClusterProvider, spec *api.ClusterConfig) error {
	networkInterfaces, err := findAvailableNetworkInterfaces(provider, spec)
	if err != nil {
		return err
	}
	return deleteNetworkInterfaces(provider, networkInterfaces)
}
