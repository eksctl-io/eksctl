package vpc

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func fmtSecurityGroupNameRegexForCluster(name string) string {
	const ourSecurityGroupNameRegexFmt = "^eksctl-%s-(cluster|nodegroup)-.+$"
	return fmt.Sprintf(ourSecurityGroupNameRegexFmt, name)
}

func forEachDanglingENI(ec2API ec2iface.EC2API, spec *api.ClusterConfig, fn func(eniID string) error) error {
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

	securityGroupRE, err := regexp.Compile(fmtSecurityGroupNameRegexForCluster(spec.Metadata.Name))
	if err != nil {
		return errors.Wrap(err, "failed to create security group regex")
	}

	var fnErr error

	err = ec2API.DescribeNetworkInterfacesPages(input, func(output *ec2.DescribeNetworkInterfacesOutput, lastPage bool) bool {
		for _, eni := range output.NetworkInterfaces {
			id := *eni.NetworkInterfaceId
			for _, sg := range eni.Groups {
				if !securityGroupRE.MatchString(*sg.GroupName) {
					logger.Debug("found %q, but it belongs to security group %q (%s), which does not appear to be ours", id, *sg.GroupName, *sg.GroupId)
					continue
				}

				logger.Debug("found %q, which belongs to our security group %q (%s)", id, *sg.GroupName, *sg.GroupId)
				if err := fn(id); err != nil {
					fnErr = err
					return false
				}
				break
			}
		}
		return !lastPage
	})

	if err != nil {
		return errors.Wrapf(err, "unable to list dangling network interfaces in %q", spec.VPC.ID)
	}

	return fnErr
}

// CleanupNetworkInterfaces finds and deletes any dangling ENIs
func CleanupNetworkInterfaces(ec2API ec2iface.EC2API, spec *api.ClusterConfig) error {
	return forEachDanglingENI(ec2API, spec, func(eniID string) error {
		input := &ec2.DeleteNetworkInterfaceInput{
			NetworkInterfaceId: &eniID,
		}
		if _, err := ec2API.DeleteNetworkInterface(input); err != nil {
			return errors.Wrapf(err, "unable to delete network interface %q", eniID)
		}
		logger.Debug("deleted network interface %q", eniID)
		return nil
	})
}
