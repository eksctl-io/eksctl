package vpc

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
)

// SetSubnets defines CIDRs for each of the subnets,
// it must be called after SetAvailabilityZones
func SetSubnets(vpc *api.ClusterVPC, availabilityZones []string) error {
	var err error

	vpc.Subnets = &api.ClusterSubnets{
		Private: api.NewAZSubnetMapping(),
		Public:  api.NewAZSubnetMapping(),
	}
	if vpc.CIDR == nil {
		cidr := api.DefaultCIDR()
		vpc.CIDR = &cidr
	}
	prefix, _ := vpc.CIDR.Mask.Size()
	if prefix < 16 || prefix > 24 {
		return errors.New("VPC CIDR prefix must be between /16 and /24")
	}
	zonesTotal := len(availabilityZones)

	var zoneCIDRs []*net.IPNet

	switch subnetsTotal := zonesTotal * 2; {
	case subnetsTotal <= 8:
		zoneCIDRs, err = SplitInto8(&vpc.CIDR.IPNet)
		if err != nil {
			return err
		}
		logger.Debug("VPC CIDR (%s) was divided into 8 subnets %v", vpc.CIDR.String(), zoneCIDRs)
	case subnetsTotal <= 16:
		zoneCIDRs, err = SplitInto16(&vpc.CIDR.IPNet)
		if err != nil {
			return err
		}
		logger.Debug("VPC CIDR (%s) was divided into 16 subnets %v", vpc.CIDR.String(), zoneCIDRs)
	default:
		return fmt.Errorf("cannot create more than 16 subnets, %d requested", subnetsTotal)
	}

	for i, zone := range availabilityZones {
		public := zoneCIDRs[i]
		private := zoneCIDRs[i+zonesTotal]
		vpc.Subnets.Private.SetAZ(zone, api.Network{
			CIDR: &ipnet.IPNet{IPNet: *private},
		})
		vpc.Subnets.Public.SetAZ(zone, api.Network{
			CIDR: &ipnet.IPNet{IPNet: *public},
		})
		logger.Info("subnets for %s - public:%s private:%s", zone, public.String(), private.String())
	}

	return nil
}

func SplitInto16(parent *net.IPNet) ([]*net.IPNet, error) {
	networkLength, _ := parent.Mask.Size()
	networkLength += 4

	var subnets []*net.IPNet
	for i := 0; i < 16; i++ {
		ip4 := parent.IP.To4()
		if ip4 != nil {
			n := binary.BigEndian.Uint32(ip4)
			n += uint32(i) << uint(32-networkLength)
			subnetIP := make(net.IP, len(ip4))
			binary.BigEndian.PutUint32(subnetIP, n)

			subnets = append(subnets, &net.IPNet{
				IP:   subnetIP,
				Mask: net.CIDRMask(networkLength, 32),
			})
		} else {
			return nil, fmt.Errorf("Unexpected IP address type: %s", parent)
		}
	}

	return subnets, nil
}

func SplitInto8(parent *net.IPNet) ([]*net.IPNet, error) {
	networkLength, _ := parent.Mask.Size()
	networkLength += 3

	var subnets []*net.IPNet
	for i := 0; i < 8; i++ {
		ip4 := parent.IP.To4()
		if ip4 != nil {
			n := binary.BigEndian.Uint32(ip4)
			n += uint32(i) << uint(32-networkLength)
			subnetIP := make(net.IP, len(ip4))
			binary.BigEndian.PutUint32(subnetIP, n)

			subnets = append(subnets, &net.IPNet{
				IP:   subnetIP,
				Mask: net.CIDRMask(networkLength, 32),
			})
		} else {
			return nil, fmt.Errorf("Unexpected IP address type: %s", parent)
		}
	}

	return subnets, nil
}

// describeSubnets fetches subnet metadata from EC2
// directly using `subnetIDs` (`vpcID` can be empty) or
// indirectly by specifying `cidrBlocks` AND `vpcID`
func describeSubnets(ctx context.Context, ec2API awsapi.EC2, vpcID string, subnetIDs, cidrBlocks, azs []string) ([]ec2types.Subnet, error) {
	var byID []ec2types.Subnet
	if len(subnetIDs) > 0 {
		output, err := ec2API.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
			SubnetIds: subnetIDs,
		})
		if err != nil {
			return nil, err
		}
		byID = output.Subnets
	}
	var byCIDR []ec2types.Subnet
	if len(cidrBlocks) > 0 {
		if vpcID == "" {
			return nil, errors.New("can't describe subnet by CIDR without VPC id")
		}
		input := &ec2.DescribeSubnetsInput{
			Filters: []ec2types.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []string{vpcID},
				},
				{
					Name:   aws.String("cidr-block"),
					Values: cidrBlocks,
				},
			},
		}
		output, err := ec2API.DescribeSubnets(ctx, input)
		if err != nil {
			return nil, err
		}
		byCIDR = output.Subnets
	}
	var byAZ []ec2types.Subnet
	if len(azs) > 0 {
		if vpcID == "" {
			return nil, errors.New("can't describe subnet by AZ without VPC id")
		}
		input := &ec2.DescribeSubnetsInput{
			Filters: []ec2types.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []string{vpcID},
				},
				{
					Name:   aws.String("availability-zone"),
					Values: azs,
				},
			},
		}
		output, err := ec2API.DescribeSubnets(ctx, input)
		if err != nil {
			return nil, err
		}
		byAZ = output.Subnets
	}
	return append(append(byID, byCIDR...), byAZ...), nil
}

func describeVPC(ctx context.Context, ec2API awsapi.EC2, vpcID string) (ec2types.Vpc, error) {
	input := &ec2.DescribeVpcsInput{
		VpcIds: []string{vpcID},
	}
	output, err := ec2API.DescribeVpcs(ctx, input)
	if err != nil {
		return ec2types.Vpc{}, err
	}
	return output.Vpcs[0], nil
}

// UseFromClusterStack retrieves the VPC configuration from an existing cluster
// based on stack outputs
// NOTE: it doesn't expect any fields in spec.VPC to be set, the remote state
// is treated as the source of truth
func UseFromClusterStack(ctx context.Context, provider api.ClusterProvider, stack *types.Stack, spec *api.ClusterConfig) error {
	if spec.VPC == nil {
		spec.VPC = api.NewClusterVPC(spec.IPv6Enabled())
	}
	// this call is authoritative, and we can safely override the
	// CIDR, as it can only be set to anything due to defaulting
	spec.VPC.CIDR = nil

	// Cluster Endpoint Access isn't part of the EKS CloudFormation Cluster stack at this point
	// Retrieve the current configuration via the SDK
	if err := UseEndpointAccessFromCluster(provider, spec); err != nil {
		return err
	}

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
			return ImportSubnetsFromIDList(ctx, provider.EC2(), spec, api.SubnetTopologyPrivate, strings.Split(v, ","))
		},
		outputs.ClusterSubnetsPublic: func(v string) error {
			return ImportSubnetsFromIDList(ctx, provider.EC2(), spec, api.SubnetTopologyPublic, strings.Split(v, ","))
		},
		outputs.ClusterFullyPrivate: func(v string) error {
			spec.PrivateCluster.Enabled = v == "true"
			return nil
		},
	}

	if !outputs.Exists(*stack, outputs.ClusterSubnetsPublic) &&
		outputs.Exists(*stack, outputs.ClusterSubnetsPublicLegacy) {
		optionalCollectors[outputs.ClusterSubnetsPublicLegacy] = func(v string) error {
			return ImportSubnetsFromIDList(ctx, provider.EC2(), spec, api.SubnetTopologyPublic, strings.Split(v, ","))
		}
	}

	return outputs.Collect(*stack, requiredCollectors, optionalCollectors)
}

// importVPC will update spec with VPC ID/CIDR
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func importVPC(ctx context.Context, ec2API awsapi.EC2, spec *api.ClusterConfig, id string) error {
	vpc, err := describeVPC(ctx, ec2API, id)
	if err != nil {
		return err
	}
	if spec.VPC.ID == "" {
		spec.VPC.ID = *vpc.VpcId
	} else if spec.VPC.ID != *vpc.VpcId {
		return fmt.Errorf("VPC ID %q is not the same as %q", spec.VPC.ID, *vpc.VpcId)
	}
	if spec.VPC.CIDR == nil {
		spec.VPC.CIDR, err = ipnet.ParseCIDR(*vpc.CidrBlock)
		if err != nil {
			return err
		}
	} else if cidr := spec.VPC.CIDR.String(); cidr != *vpc.CidrBlock {
		for _, cidrAssoc := range vpc.CidrBlockAssociationSet {
			if aws.ToString(cidrAssoc.CidrBlock) == cidr {
				return nil
			}
		}
		return fmt.Errorf("VPC CIDR block %q not found in VPC", cidr)
	}

	return nil
}

// ImportSubnets will update spec with subnets, if VPC ID/CIDR is unknown
// it will use provider to call describeVPC based on the VPC ID of the
// first subnet; all subnets must be in the same VPC
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func ImportSubnets(ctx context.Context, ec2API awsapi.EC2, spec *api.ClusterConfig, topology api.SubnetTopology, subnets []ec2types.Subnet) error {
	if spec.VPC.ID != "" {
		// ensure managed NAT is disabled
		// if we are importing an existing VPC/subnets, the expectation is that the user has
		// already setup NAT, routing, etc. for these subnets
		disable := api.ClusterDisableNAT
		spec.VPC.NAT = &api.ClusterNAT{
			Gateway: &disable,
		}

		// ensure VPC gets imported and validated first, if it's already set
		if err := importVPC(ctx, ec2API, spec, spec.VPC.ID); err != nil {
			return err
		}
	}

	for _, sn := range subnets {
		if spec.VPC.ID == "" {
			// if VPC wasn't defined, import it based on VPC of the first
			// subnet that we have
			if err := importVPC(ctx, ec2API, spec, *sn.VpcId); err != nil {
				return err
			}
		} else if spec.VPC.ID != *sn.VpcId { // be sure to use the same VPC
			return fmt.Errorf("given %s is in %s, not in %s", *sn.SubnetId, *sn.VpcId, spec.VPC.ID)
		}

		if err := spec.ImportSubnet(topology, *sn.AvailabilityZone, *sn.SubnetId, *sn.CidrBlock); err != nil {
			return err
		}
		spec.AppendAvailabilityZone(*sn.AvailabilityZone)
	}
	return nil
}

// ImportSubnetsFromList will update spec with subnets, it will call describeSubnets first,
// then pass resulting subnets to ImportSubnets
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func importSubnetsFromList(ctx context.Context, ec2API awsapi.EC2, spec *api.ClusterConfig, topology api.SubnetTopology, subnetIDs, cidrs, azs []string) error {
	subnets, err := describeSubnets(ctx, ec2API, spec.VPC.ID, subnetIDs, cidrs, azs)
	if err != nil {
		return err
	}

	return ImportSubnets(ctx, ec2API, spec, topology, subnets)
}

// importSubnetsForTopology will update spec with subnets, it will call describeSubnets first,
// then pass resulting subnets to ImportSubnets
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func importSubnetsForTopology(ctx context.Context, ec2API awsapi.EC2, spec *api.ClusterConfig, topology api.SubnetTopology) error {
	var subnetMapping api.AZSubnetMapping
	if spec.VPC.Subnets != nil {
		switch topology {
		case api.SubnetTopologyPrivate:
			subnetMapping = spec.VPC.Subnets.Private
		case api.SubnetTopologyPublic:
			subnetMapping = spec.VPC.Subnets.Public
		default:
			panic(fmt.Sprintf("unexpected subnet topology: %s", topology))
		}
	}

	subnetIDs := subnetMapping.WithIDs()
	cidrs := subnetMapping.WithCIDRs()
	azs := subnetMapping.WithAZs()

	subnets, err := describeSubnets(ctx, ec2API, spec.VPC.ID, subnetIDs, cidrs, azs)
	if err != nil {
		return err
	}

	return ImportSubnets(ctx, ec2API, spec, topology, subnets)
}

// ImportSubnetsFromIDList will update cluster config with subnets _only specified by ID_
// then pass resulting subnets to ImportSubnets
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func ImportSubnetsFromIDList(ctx context.Context, ec2API awsapi.EC2, spec *api.ClusterConfig, topology api.SubnetTopology, subnetIDs []string) error {
	return importSubnetsFromList(ctx, ec2API, spec, topology, subnetIDs, []string{}, []string{})
}

func ValidateLegacySubnetsForNodeGroups(ctx context.Context, spec *api.ClusterConfig, provider api.ClusterProvider) error {
	subnetsToValidate := sets.NewString()

	selectSubnets := func(ng *api.NodeGroupBase) error {
		if len(ng.AvailabilityZones) > 0 || len(ng.Subnets) > 0 {
			// Check only the public subnets that this ng has
			subnetIDs, err := SelectNodeGroupSubnets(ctx, ng.AvailabilityZones, ng.Subnets, spec.VPC.Subnets.Public, provider.EC2(), spec.VPC.ID)
			if err != nil {
				return errors.Wrap(err, "couldn't find public subnets")
			}
			subnetsToValidate.Insert(subnetIDs...)
		} else {
			// This ng doesn't have AZs defined so we need to check all public subnets
			for _, subnet := range spec.VPC.Subnets.Public {
				subnetsToValidate.Insert(subnet.ID)
			}
		}
		return nil
	}

	for _, ng := range spec.NodeGroups {
		if ng.PrivateNetworking {
			continue
		}
		err := selectSubnets(ng.NodeGroupBase)
		if err != nil {
			return err
		}
	}

	for _, ng := range spec.ManagedNodeGroups {
		if ng.PrivateNetworking {
			continue
		}
		err := selectSubnets(ng.NodeGroupBase)
		if err != nil {
			return err
		}
	}

	if err := ValidateExistingPublicSubnets(ctx, provider, spec.VPC.ID, subnetsToValidate.List()); err != nil {
		// If the cluster endpoint is reachable from the VPC nodes might still be able to join
		if spec.HasPrivateEndpointAccess() {
			logger.Warning("public subnets for one or more nodegroups have %q disabled. This means that nodes won't "+
				"get public IP addresses. If they can't reach the cluster through the private endpoint they won't be "+
				"able to join the cluster", "MapPublicIpOnLaunch")
			return nil
		}

		logger.Critical(err.Error())
		return errors.Errorf("subnets for one or more new nodegroups don't meet requirements. "+
			"To fix this, please run `eksctl utils update-legacy-subnet-settings --cluster %s`",
			spec.Metadata.Name)
	}
	return nil
}

// ValidateExistingPublicSubnets makes sure that subnets have the property MapPublicIpOnLaunch enabled
func ValidateExistingPublicSubnets(ctx context.Context, provider api.ClusterProvider, vpcID string, subnetIDs []string) error {
	if len(subnetIDs) == 0 {
		return nil
	}
	subnets, err := describeSubnets(ctx, provider.EC2(), vpcID, subnetIDs, []string{}, []string{})
	if err != nil {
		return err
	}
	return validatePublicSubnet(subnets)
}

// EnsureMapPublicIPOnLaunchEnabled will enable MapPublicIpOnLaunch in EC2 for all given subnet IDs
func EnsureMapPublicIPOnLaunchEnabled(ctx context.Context, ec2API awsapi.EC2, subnetIDs []string) error {
	if len(subnetIDs) == 0 {
		logger.Debug("no subnets to update")
		return nil
	}

	for _, s := range subnetIDs {
		input := &ec2.ModifySubnetAttributeInput{
			SubnetId: aws.String(s),
			MapPublicIpOnLaunch: &ec2types.AttributeBooleanValue{
				Value: aws.Bool(true),
			},
		}

		logger.Debug("enabling MapPublicIpOnLaunch for subnet %q", s)
		_, err := ec2API.ModifySubnetAttribute(ctx, input)
		if err != nil {
			return errors.Wrapf(err, "unable to set MapPublicIpOnLaunch attribute to true for subnet %q", s)
		}
	}
	return nil
}

// ImportSubnetsFromSpec will update spec with subnets, it will call describeSubnets first,
// then pass resulting subnets to ImportSubnets
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func ImportSubnetsFromSpec(ctx context.Context, provider api.ClusterProvider, spec *api.ClusterConfig) error {
	if spec.VPC.ID != "" {
		// ensure VPC gets imported and validated first, if it's already set
		if err := importVPC(ctx, provider.EC2(), spec, spec.VPC.ID); err != nil {
			return err
		}
	}
	if err := importSubnetsForTopology(ctx, provider.EC2(), spec, api.SubnetTopologyPrivate); err != nil {
		return err
	}
	if err := importSubnetsForTopology(ctx, provider.EC2(), spec, api.SubnetTopologyPublic); err != nil {
		return err
	}
	// to clean up invalid subnets based on AZ after importing both private and public subnets
	cleanupSubnets(spec)
	return nil
}

//UseEndpointAccessFromCluster retrieves the Cluster's endpoint access configuration via the SDK
// as the CloudFormation Stack doesn't support that configuration currently
func UseEndpointAccessFromCluster(provider api.ClusterProvider, spec *api.ClusterConfig) error {
	input := &awseks.DescribeClusterInput{
		Name: &spec.Metadata.Name,
	}
	output, err := provider.EKS().DescribeCluster(input)
	if err != nil {
		return err
	}
	if spec.VPC.ClusterEndpoints == nil {
		spec.VPC.ClusterEndpoints = &api.ClusterEndpoints{}
	}
	spec.VPC.ClusterEndpoints.PublicAccess = output.Cluster.ResourcesVpcConfig.EndpointPublicAccess
	spec.VPC.ClusterEndpoints.PrivateAccess = output.Cluster.ResourcesVpcConfig.EndpointPrivateAccess
	return nil
}

// cleanupSubnets clean up subnet entries having invalid AZ
func cleanupSubnets(spec *api.ClusterConfig) {
	availabilityZones := make(map[string]struct{})
	for _, az := range spec.AvailabilityZones {
		availabilityZones[az] = struct{}{}
	}

	cleanup := func(subnets *api.AZSubnetMapping) {
		for name, subnet := range *subnets {
			if _, ok := availabilityZones[subnet.AZ]; !ok {
				delete(*subnets, name)
			}
		}
	}

	cleanup(&spec.VPC.Subnets.Private)
	cleanup(&spec.VPC.Subnets.Public)
}

func validatePublicSubnet(subnets []ec2types.Subnet) error {
	legacySubnets := make([]string, 0)
	for _, sn := range subnets {
		if sn.MapPublicIpOnLaunch == nil || !*sn.MapPublicIpOnLaunch {
			legacySubnets = append(legacySubnets, *sn.SubnetId)
		}
	}
	if len(legacySubnets) > 0 {
		return fmt.Errorf("found mis-configured or non-public subnets %q. Expected public subnets with property "+
			"\"MapPublicIpOnLaunch\" enabled. Without it new nodes won't get an IP assigned", legacySubnets)
	}

	return nil
}

// getSubnetByID returns a subnet based on an ID.
func getSubnetByID(ctx context.Context, ec2API awsapi.EC2, id string) (ec2types.Subnet, error) {
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: []string{id},
	}
	output, err := ec2API.DescribeSubnets(ctx, input)
	if err != nil {
		return ec2types.Subnet{}, err
	}
	if len(output.Subnets) != 1 {
		return ec2types.Subnet{}, fmt.Errorf("subnet with id %q not found", id)
	}
	return output.Subnets[0], nil
}

func SelectNodeGroupSubnets(ctx context.Context, nodegroupAZs, nodegroupSubnets []string, subnets api.AZSubnetMapping, ec2API awsapi.EC2, vpcID string) ([]string, error) {
	// We have validated that either azs are provided or subnets are provided
	numNodeGroupsAZs := len(nodegroupAZs)
	numNodeGroupsSubnets := len(nodegroupSubnets)
	if numNodeGroupsAZs == 0 && numNodeGroupsSubnets == 0 {
		return nil, nil
	}

	makeErrorDesc := func() string {
		return fmt.Sprintf("(allSubnets=%#v AZs=%#v subnets=%#v)", subnets, nodegroupAZs, nodegroupSubnets)
	}
	if len(subnets) < numNodeGroupsAZs || len(subnets) < numNodeGroupsSubnets {
		return nil, fmt.Errorf("mapping doesn't have enough subnets: %s", makeErrorDesc())
	}
	subnetIDs := []string{}
	// We validate previously that either AZs or subnets is set
	for _, az := range nodegroupAZs {
		azSubnetIDs := []string{}
		for _, s := range subnets {
			if s.AZ == az {
				azSubnetIDs = append(azSubnetIDs, s.ID)
			}
		}
		if len(azSubnetIDs) == 0 {
			return nil, fmt.Errorf("mapping doesn't have subnet with AZ %s: %s", az, makeErrorDesc())
		}
		subnetIDs = append(subnetIDs, azSubnetIDs...)
	}
	for _, subnetName := range nodegroupSubnets {
		var subnetID string
		if subnet, ok := subnets[subnetName]; !ok {
			for _, s := range subnets {
				if s.ID != subnetName {
					continue
				}
				subnetID = s.ID
			}
		} else {
			subnetID = subnet.ID
		}
		if subnetID == "" {
			subnet, err := getSubnetByID(ctx, ec2API, subnetName)
			if err != nil {
				return nil, err
			}
			if subnet.VpcId != nil && *subnet.VpcId != vpcID {
				return nil, fmt.Errorf("subnet with id %q is not in the attached vpc with id %q", *subnet.SubnetId, vpcID)
			}
			subnetID = *subnet.SubnetId
		}
		subnetIDs = append(subnetIDs, subnetID)
	}
	return subnetIDs, nil
}
