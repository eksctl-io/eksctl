package vpc

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
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
	if vpc.CIDR == nil {
		cidr := api.DefaultCIDR()
		vpc.CIDR = &cidr
	}
	prefix, _ := spec.VPC.CIDR.Mask.Size()
	if (prefix < 16) || (prefix > 24) {
		return fmt.Errorf("VPC CIDR prefix must be betwee /16 and /24")
	}
	zoneCIDRs, err := subnet.SplitInto8(&spec.VPC.CIDR.IPNet)
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
			CIDR: &ipnet.IPNet{IPNet: *public},
		}
		vpc.Subnets[api.SubnetTopologyPrivate][zone] = api.Network{
			CIDR: &ipnet.IPNet{IPNet: *private},
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

// ImportVPC will update spec with VPC ID/CIDR
func ImportVPC(provider api.ClusterProvider, spec *api.ClusterConfig, id string) error {
	vpc, err := describeVPC(provider, id)
	if err != nil {
		return err
	}
	if spec.VPC.ID == "" {
		spec.VPC.ID = *vpc.VpcId
	} else if spec.VPC.ID != *vpc.VpcId {
		return fmt.Errorf("VPC ID %q is the same as not %q", spec.VPC.ID, *vpc.VpcId)
	}
	if spec.VPC.CIDR == nil {
		spec.VPC.CIDR, err = ipnet.ParseCIDR(*vpc.CidrBlock)
		if err != nil {
			return err
		}
	} else if cidr := spec.VPC.CIDR.String(); cidr != *vpc.CidrBlock {
		return fmt.Errorf("VPC CIDR block %q is not the same as %q",
			cidr,
			*vpc.CidrBlock,
		)
	}

	return nil
}

// ImportSubnets will update spec with subnets, if VPC ID/CIDR is unknown
// it will use provider to call describeVPC based on the VPC ID of the
// first subnet; all subnets must be in the same VPC
func ImportSubnets(provider api.ClusterProvider, spec *api.ClusterConfig, topology api.SubnetTopology, subnets []*ec2.Subnet) error {
	if spec.VPC.ID != "" {
		// ensure VPC gets imported and validated firt, if it's already set
		if err := ImportVPC(provider, spec, spec.VPC.ID); err != nil {
			return err
		}
	}
	for _, subnet := range subnets {
		if spec.VPC.ID == "" {
			// if VPC wasn't defined, import it based on VPC of the first
			// subnet that we have
			if err := ImportVPC(provider, spec, *subnet.VpcId); err != nil {
				return err
			}
		} else if spec.VPC.ID != *subnet.VpcId { // be sure to use the same VPC
			return fmt.Errorf("given %s is in %s, not in %s", *subnet.SubnetId, *subnet.VpcId, spec.VPC.ID)
		}

		if err := spec.ImportSubnet(topology, *subnet.AvailabilityZone, *subnet.SubnetId, *subnet.CidrBlock); err != nil {
			return err
		}
		spec.AppendAvailabilityZone(*subnet.AvailabilityZone)
	}
	return nil
}

// UseSubnetsFromList will update spec with subnets, it will call describeSubnets first,
// then pass resulting subnets to ImportSubnets
func UseSubnetsFromList(provider api.ClusterProvider, spec *api.ClusterConfig, topology api.SubnetTopology, subnetIDs []string) error {
	if len(subnetIDs) == 0 {
		return nil
	}
	subnets, err := describeSubnets(provider, subnetIDs...)
	if err != nil {
		return err
	}
	return ImportSubnets(provider, spec, topology, subnets)
}

// UseSubnets will update spec with subnets, it will call describeSubnets first,
// then pass resulting subnets to ImportSubnets
func UseSubnets(provider api.ClusterProvider, spec *api.ClusterConfig) error {
	if spec.VPC.ID != "" {
		// ensure VPC gets imported and validated firt, if it's already set
		if err := ImportVPC(provider, spec, spec.VPC.ID); err != nil {
			return err
		}
	}
	for topology := range spec.VPC.Subnets {
		subnetIDs := []string{}
		for _, subnet := range spec.VPC.Subnets[topology] {
			subnetIDs = append(subnetIDs, subnet.ID)
		}
		if err := UseSubnetsFromList(provider, spec, topology, subnetIDs); err != nil {
			return err
		}
	}
	return nil
}
