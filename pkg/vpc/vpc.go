package vpc

import (
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"k8s.io/kops/pkg/util/subnet"
)

// SetSubnets defines CIDRs for each of the subnets,
// it must be called after SetAvailabilityZones
func SetSubnets(spec *api.ClusterConfig) error {
	var err error

	vpc := spec.VPC
	vpc.Subnets = map[api.SubnetTopology]map[string]api.Network{
		api.SubnetTopologyPublic:  map[string]api.Network{},
		api.SubnetTopologyPrivate: map[string]api.Network{},
	}
	prefix, _ := spec.VPC.CIDR.Mask.Size()
	if (prefix < 16) || (prefix > 24) {
		return fmt.Errorf("VPC CIDR prefix must be betwee /16 and /24")
	}
	zoneCIDRs, err := subnet.SplitInto8(spec.VPC.CIDR)
	if err != nil {
		return err
	}

	logger.Debug("VPC CIDR (%s) was divided into 8 subnets %v", vpc.CIDR.String(), zoneCIDRs)

	zonesTotal := len(spec.AvailabilityZones)
	if 2*zonesTotal > len(zoneCIDRs) {
		return fmt.Errorf("insufficient number of subnets (have %d, but need %d) for %d availability zones", len(zoneCIDRs), 2*zonesTotal, zonesTotal)
	}

	for i, zone := range spec.AvailabilityZones {
		public := zoneCIDRs[i]
		private := zoneCIDRs[i+zonesTotal]
		vpc.Subnets[api.SubnetTopologyPublic][zone] = api.Network{
			CIDR: public,
		}
		vpc.Subnets[api.SubnetTopologyPrivate][zone] = api.Network{
			CIDR: private,
		}
		logger.Info("subnets for %s - public:%s private:%s", zone, public.String(), private.String())
	}

	return nil
}

func describeSubnets(porvider api.ClusterProvider, subnetIDs ...string) ([]*ec2.Subnet, error) {
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: aws.StringSlice(subnetIDs),
	}
	output, err := porvider.EC2().DescribeSubnets(input)
	if err != nil {
		return nil, err
	}
	return output.Subnets, nil
}

func describeVPC(povider api.ClusterProvider, vpcID string) (*ec2.Vpc, error) {
	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{aws.String(vpcID)},
	}
	output, err := povider.EC2().DescribeVpcs(input)
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
			_, spec.VPC.CIDR, err = net.ParseCIDR(*vpc.CidrBlock)
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
