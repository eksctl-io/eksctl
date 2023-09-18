package outposts

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/kris-nova/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/outposts"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

// A ClusterExtender extends a cluster with resources required to support nodegroups on Outposts.
type ClusterExtender struct {
	StackUpdater stackUpdater
	EC2API       awsapi.EC2
	OutpostsAPI  awsapi.Outposts
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o fakes . stackUpdater
type stackUpdater interface {
	AppendNewClusterStackResource(ctx context.Context, extendForOutposts, plan bool) (bool, error)
}

// ClusterToExtend represents a cluster that needs to be extended.
//
//counterfeiter:generate -o fakes . ClusterToExtend
type ClusterToExtend interface {
	// IsControlPlaneOnOutposts returns true if the control plane is on Outposts.
	IsControlPlaneOnOutposts() bool

	// FindNodeGroupOutpostARN checks whether any nodegroups are on Outposts and returns the Outpost ARN.
	FindNodeGroupOutpostARN() (outpostARN string, found bool)
}

// ExtendWithOutpostSubnetsIfRequired extends a cluster's stack with Outpost subnets if required.
func (e *ClusterExtender) ExtendWithOutpostSubnetsIfRequired(ctx context.Context, cluster ClusterToExtend, clusterVPC *api.ClusterVPC) error {
	if cluster.IsControlPlaneOnOutposts() {
		return nil
	}
	nodeGroupOutpostARN, found := cluster.FindNodeGroupOutpostARN()
	if !found {
		return nil
	}

	subnetsOutpostARN, found := clusterVPC.FindOutpostSubnetsARN()
	if found {
		if subnetsOutpostARN != nodeGroupOutpostARN {
			return fmt.Errorf("cannot extend a cluster with two different Outposts; found subnets on Outpost %q but nodegroup is using %q", subnetsOutpostARN, nodeGroupOutpostARN)
		}
		return nil
	}

	logger.Info("extending cluster with subnets for Outposts")

	outpost, err := e.OutpostsAPI.GetOutpost(ctx, &outposts.GetOutpostInput{
		OutpostId: aws.String(nodeGroupOutpostARN),
	})
	if err != nil {
		return fmt.Errorf("error getting Outpost details: %w", err)
	}

	newSubnets, err := vpc.ExtendWithOutpostSubnets(clusterVPC.CIDR.IPNet, len(clusterVPC.Subnets.Public)+len(clusterVPC.Subnets.Private), nodeGroupOutpostARN, *outpost.Outpost.AvailabilityZone)
	if err != nil {
		return fmt.Errorf("error extending cluster with Outpost subnets: %w", err)
	}

	existingSubnets, err := describeVPCSubnets(ctx, e.EC2API, clusterVPC.ID)
	if err != nil {
		return err
	}

	if err := validateNewSubnets(existingSubnets, newSubnets, clusterVPC.Subnets); err != nil {
		return err
	}

	if err := e.extendWithOutpostSubnets(ctx, newSubnets, clusterVPC); err != nil {
		return fmt.Errorf("error adding new Outpost subnets: %w", err)
	}
	logger.Info("cluster has been extended with Outpost subnets")
	return nil
}

func (e *ClusterExtender) extendWithOutpostSubnets(ctx context.Context, newSubnets *vpc.SubnetPair, clusterVPC *api.ClusterVPC) error {
	if err := addNewSubnets(newSubnets.Public, clusterVPC.Subnets.Public); err != nil {
		return err
	}
	if err := addNewSubnets(newSubnets.Private, clusterVPC.Subnets.Private); err != nil {
		return err
	}
	_, err := e.StackUpdater.AppendNewClusterStackResource(ctx, true, false)
	if err != nil {
		return fmt.Errorf("error updating cluster stack with Outpost resources: %w", err)
	}
	return nil
}

func addNewSubnets(newSubnets []api.AZSubnetSpec, existingSubnetsMap api.AZSubnetMapping) error {
	for i, newSubnet := range newSubnets {
		subnetAlias := vpc.MakeExtendedSubnetAlias(newSubnet.AZ, i+1)
		if _, ok := existingSubnetsMap[subnetAlias]; ok {
			return fmt.Errorf("unexpected error adding new Outpost subnets: subnet alias %q generated for new Outpost subnet already exists", subnetAlias)
		}
		existingSubnetsMap[subnetAlias] = newSubnet
	}
	return nil
}

func validateNewSubnets(existingSubnets []ec2types.Subnet, newSubnets *vpc.SubnetPair, clusterSubnets *api.ClusterSubnets) error {
	var newSubnetPrefixes []netip.Prefix
	for _, subnets := range [][]api.AZSubnetSpec{newSubnets.Public, newSubnets.Private} {
		for _, s := range subnets {
			addr, ok := netip.AddrFromSlice(s.CIDR.IP)
			if !ok {
				return fmt.Errorf("unexpected error creating a netip.Addr from subnet CIDR %q", s.CIDR)
			}
			ones, _ := s.CIDR.Mask.Size()
			prefix := netip.PrefixFrom(addr, ones)
			newSubnetPrefixes = append(newSubnetPrefixes, prefix)
		}
	}

	for _, subnet := range existingSubnets {
		subnetPrefix, err := netip.ParsePrefix(*subnet.CidrBlock)
		if err != nil {
			return fmt.Errorf("unexpected error parsing subnet CIDR %q: %w", *subnet.CidrBlock, err)
		}
		if err := validateCIDROverlap(newSubnetPrefixes, subnetPrefix, subnet, clusterSubnets); err != nil {
			return err
		}
	}
	return nil
}

func isExternalSubnet(subnets *api.ClusterSubnets, subnet ec2types.Subnet) bool {
	for _, subnetMap := range []api.AZSubnetMapping{subnets.Public, subnets.Private} {
		for _, s := range subnetMap {
			if s.ID == *subnet.SubnetId {
				return false
			}
		}
	}
	return true
}

func validateCIDROverlap(newSubnetPrefixes []netip.Prefix, subnetPrefix netip.Prefix, subnet ec2types.Subnet, clusterSubnets *api.ClusterSubnets) error {
	for _, newSubnetPrefix := range newSubnetPrefixes {
		if subnetPrefix.Overlaps(newSubnetPrefix) {
			if isExternalSubnet(clusterSubnets, subnet) {
				return fmt.Errorf("cannot create subnets on Outpost; subnet CIDR %q (ID: %s) created outside of eksctl overlaps with new CIDR %q", *subnet.CidrBlock, *subnet.SubnetId, newSubnetPrefix)
			}
			return fmt.Errorf("unexpected error calculating subnet CIDRs for Outposts: new CIDR %q overlaps with existing CIDR %q (ID: %s)", newSubnetPrefix, *subnet.CidrBlock, *subnet.SubnetId)
		}
	}
	return nil
}

func describeVPCSubnets(ctx context.Context, ec2API awsapi.EC2, vpcID string) ([]ec2types.Subnet, error) {
	subnetsPaginator := ec2.NewDescribeSubnetsPaginator(ec2API, &ec2.DescribeSubnetsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcID},
			},
		},
	})
	var ret []ec2types.Subnet
	for subnetsPaginator.HasMorePages() {
		output, err := subnetsPaginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("error describing subnets: %w", err)
		}
		ret = append(ret, output.Subnets...)
	}
	return ret, nil
}
