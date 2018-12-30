package vpc

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha3"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
)

func describeSubnets(provider api.ClusterProvider, subnetIDs ...string) ([]*ec2.Subnet, error) {
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: aws.StringSlice(subnetIDs),
	}
	output, err := provider.EC2().DescribeSubnets(input)
	if err != nil {
		return nil, err
	}
	return output.Subnets, nil
}

func describeVPC(provider api.ClusterProvider, vpcID string) (*ec2.Vpc, error) {
	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{aws.String(vpcID)},
	}
	output, err := provider.EC2().DescribeVpcs(input)
	if err != nil {
		return nil, err
	}
	return output.Vpcs[0], nil
}

// ImportSubnets will update spec with subnets, if VPC ID/CIDR is unknown
// it will use provider to call describeVPC based on the VPC ID of the
// first subnet; all subnets must be in the same VPC
func ImportSubnets(provider api.ClusterProvider, spec *api.ClusterConfig, topology api.SubnetTopology, subnets []*ec2.Subnet) error {
	for _, subnet := range subnets {
		if spec.VPC.ID == "" {
			vpc, err := describeVPC(provider, *subnet.VpcId)
			if err != nil {
				return err
			}
			spec.VPC.ID = *vpc.VpcId
			spec.VPC.CIDR, err = ipnet.ParseCIDR(*vpc.CidrBlock)
			if err != nil {
				return err
			}
		} else if spec.VPC.ID != *subnet.VpcId {
			return fmt.Errorf("given %s is in %s, not in %s", *subnet.SubnetId, *subnet.VpcId, spec.VPC.ID)
		}

		spec.ImportSubnet(topology, *subnet.AvailabilityZone, *subnet.SubnetId, *subnet.CidrBlock)
		spec.AppendAvailabilityZone(*subnet.AvailabilityZone)
	}
	return nil
}

// UseSubnets will update spec with subnets, it will call describeSubnets first,
// then pass resulting subnets to ImportSubnets
func UseSubnets(provider api.ClusterProvider, spec *api.ClusterConfig, topology api.SubnetTopology, subnetIDs []string) error {
	if len(subnetIDs) == 0 {
		return nil
	}

	subnets, err := describeSubnets(provider, subnetIDs...)
	if err != nil {
		return err
	}

	return ImportSubnets(provider, spec, topology, subnets)
}
