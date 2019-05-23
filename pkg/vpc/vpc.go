package vpc

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"

	"k8s.io/kops/pkg/util/subnet"
)

// SetSubnets defines CIDRs for each of the subnets,
// it must be called after SetAvailabilityZones
func SetSubnets(spec *api.ClusterConfig) error {
	var err error

	vpc := spec.VPC
	vpc.Subnets = &api.ClusterSubnets{
		Private: map[string]api.Network{},
		Public:  map[string]api.Network{},
	}
	if vpc.CIDR == nil {
		cidr := api.DefaultCIDR()
		vpc.CIDR = &cidr
	}
	prefix, _ := spec.VPC.CIDR.Mask.Size()
	if (prefix < 16) || (prefix > 24) {
		return fmt.Errorf("VPC CIDR prefix must be between /16 and /24")
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
		vpc.Subnets.Private[zone] = api.Network{
			CIDR: &ipnet.IPNet{IPNet: *private},
		}
		vpc.Subnets.Public[zone] = api.Network{
			CIDR: &ipnet.IPNet{IPNet: *public},
		}
		logger.Info("subnets for %s - public:%s private:%s", zone, public.String(), private.String())
	}

	return nil
}

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

func describe(povider api.ClusterProvider, vpcID string) (*ec2.Vpc, error) {
	input := &ec2.DescribeVpcsInput{
		VpcIds: []*string{aws.String(vpcID)},
	}
	output, err := povider.EC2().DescribeVpcs(input)
	if err != nil {
		return nil, err
	}
	return output.Vpcs[0], nil
}

// UseFromCluster retrieves the VPC configuration from an existing cluster
// based on stack outputs
// NOTE: it doesn't expect any fields in spec.VPC to be set, the remote state
// is treated as the source of truth
func UseFromCluster(provider api.ClusterProvider, stack *cfn.Stack, spec *api.ClusterConfig) error {
	if spec.VPC == nil {
		spec.VPC = &api.ClusterVPC{}
	}
	// this call is authoritive, and we can safely override the
	// CIDR, as it can only be set to anything due to defaulting
	spec.VPC.CIDR = nil

	requiredCollectors := map[string]outputs.Collector{
		outputs.ClusterVPC: func(v string) error {
			spec.VPC.ID = v
			return nil
		},
		outputs.ClusterSecurityGroup: func(v string) error {
			spec.VPC.SecurityGroup = v
			return nil
		},
	}

	optionalCollectors := map[string]outputs.Collector{
		outputs.ClusterSharedNodeSecurityGroup: func(v string) error {
			spec.VPC.SharedNodeSecurityGroup = v
			return nil
		},
		outputs.ClusterSubnetsPrivate: func(v string) error {
			return ImportSubnetsFromList(provider, spec, api.SubnetTopologyPrivate, strings.Split(v, ","))
		},
		outputs.ClusterSubnetsPublic: func(v string) error {
			return ImportSubnetsFromList(provider, spec, api.SubnetTopologyPublic, strings.Split(v, ","))
		},
	}

	if !outputs.Exists(*stack, outputs.ClusterSubnetsPublic) &&
		outputs.Exists(*stack, outputs.ClusterSubnetsPublicLegacy) {
		optionalCollectors[outputs.ClusterSubnetsPublicLegacy] = func(v string) error {
			return ImportSubnetsFromList(provider, spec, api.SubnetTopologyPublic, strings.Split(v, ","))
		}
	}

	return outputs.Collect(*stack, requiredCollectors, optionalCollectors)
}

// Import will update spec with VPC ID/CIDR
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func Import(provider api.ClusterProvider, spec *api.ClusterConfig, id string) error {
	vpc, err := describe(provider, id)
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
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func ImportSubnets(provider api.ClusterProvider, spec *api.ClusterConfig, topology api.SubnetTopology, subnets []*ec2.Subnet) error {
	if spec.VPC.ID != "" {
		// ensure VPC gets imported and validated firt, if it's already set
		if err := Import(provider, spec, spec.VPC.ID); err != nil {
			return err
		}
	}
	for _, subnet := range subnets {
		if spec.VPC.ID == "" {
			// if VPC wasn't defined, import it based on VPC of the first
			// subnet that we have
			if err := Import(provider, spec, *subnet.VpcId); err != nil {
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

// ImportSubnetsFromList will update spec with subnets, it will call describeSubnets first,
// then pass resulting subnets to ImportSubnets
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func ImportSubnetsFromList(provider api.ClusterProvider, spec *api.ClusterConfig, topology api.SubnetTopology, subnetIDs []string) error {
	if len(subnetIDs) == 0 {
		return nil
	}
	subnets, err := describeSubnets(provider, subnetIDs...)
	if err != nil {
		return err
	}
	return ImportSubnets(provider, spec, topology, subnets)
}

// ImportAllSubnets will update spec with subnets, it will call describeSubnets first,
// then pass resulting subnets to ImportSubnets
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func ImportAllSubnets(provider api.ClusterProvider, spec *api.ClusterConfig) error {
	if spec.VPC.ID != "" {
		// ensure VPC gets imported and validated firt, if it's already set
		if err := Import(provider, spec, spec.VPC.ID); err != nil {
			return err
		}
	}
	if err := ImportSubnetsFromList(provider, spec, api.SubnetTopologyPrivate, spec.PrivateSubnetIDs()); err != nil {
		return err
	}
	if err := ImportSubnetsFromList(provider, spec, api.SubnetTopologyPublic, spec.PublicSubnetIDs()); err != nil {
		return err
	}

	return nil
}
