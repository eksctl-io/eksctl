package vpc

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/kris-nova/logger"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/az"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
	"github.com/weaveworks/eksctl/pkg/utils/nodes"
)

// StackDriftError represents a stack drift error.
type StackDriftError struct {
	Msg string
}

// Error implements the error interface.
func (s *StackDriftError) Error() string {
	return s.Msg
}

// SetSubnets defines CIDRs for each of the subnets,
// it must be called after SetAvailabilityZones.
func SetSubnets(vpc *api.ClusterVPC, availabilityZones, localZones []string) error {
	if err := validateVPCCIDR(vpc); err != nil {
		return err
	}

	zonesTotal := len(availabilityZones) + len(localZones)
	subnetsTotal := zonesTotal * 2

	subnetSize, networkLength, err := getSubnetNetworkSize(vpc.CIDR.IPNet, subnetsTotal)
	if err != nil {
		return err
	}

	zoneCIDRs, err := SplitInto(&vpc.CIDR.IPNet, subnetSize, networkLength)
	if err != nil {
		return err
	}

	logger.Debug("VPC CIDR (%s) was divided into %d subnets %v", vpc.CIDR.String(), len(zoneCIDRs), zoneCIDRs)

	vpc.Subnets = &api.ClusterSubnets{
		Private: api.NewAZSubnetMapping(),
		Public:  api.NewAZSubnetMapping(),
	}
	vpc.LocalZoneSubnets = &api.ClusterSubnets{
		Private: api.NewAZSubnetMapping(),
		Public:  api.NewAZSubnetMapping(),
	}

	setSubnets := func(zones []string, startIndex int, subnets *api.ClusterSubnets) {
		for i, zone := range zones {
			publicCIDRIndex := startIndex + i
			privateCIDRIndex := publicCIDRIndex + zonesTotal

			publicCIDR := zoneCIDRs[publicCIDRIndex]
			privateCIDR := zoneCIDRs[privateCIDRIndex]

			subnets.Private[zone] = api.AZSubnetSpec{
				AZ:        zone,
				CIDR:      &ipnet.IPNet{IPNet: *privateCIDR},
				CIDRIndex: privateCIDRIndex,
			}
			subnets.Public[zone] = api.AZSubnetSpec{
				AZ:        zone,
				CIDR:      &ipnet.IPNet{IPNet: *publicCIDR},
				CIDRIndex: publicCIDRIndex,
			}

			logger.Info("subnets for %s - public:%s private:%s", zone, publicCIDR, privateCIDR)
		}
	}

	setSubnets(availabilityZones, 0, vpc.Subnets)
	setSubnets(localZones, len(availabilityZones), vpc.LocalZoneSubnets)
	return nil
}

// A SubnetPair represents a pair of public and private subnets.
type SubnetPair struct {
	Public  []api.AZSubnetSpec
	Private []api.AZSubnetSpec
}

// ExtendWithOutpostSubnets extends the VPC by returning public and private subnet CIDRs for Outposts.
func ExtendWithOutpostSubnets(vpcCIDR net.IPNet, existingSubnetsCount int, outpostARN, outpostAZ string) (*SubnetPair, error) {
	subnetSize, networkLength, err := getSubnetNetworkSize(vpcCIDR, existingSubnetsCount)
	if err != nil {
		return nil, err
	}
	cidrs, err := SplitInto(&vpcCIDR, subnetSize, networkLength)
	if err != nil {
		return nil, err
	}
	if len(cidrs) < existingSubnetsCount {
		return nil, errors.New("unexpected error calculating new subnet CIDRs")
	}

	newCIDRs := cidrs[existingSubnetsCount:]
	if len(newCIDRs) < 2 {
		return nil, fmt.Errorf("VPC cannot be extended with more subnets: expected to find at least two free CIDRs in VPC; got %d", len(newCIDRs))
	}

	makeAZSubnetSpec := func(cidr *net.IPNet, cidrIndex int) api.AZSubnetSpec {
		return api.AZSubnetSpec{
			AZ: outpostAZ,
			CIDR: &ipnet.IPNet{
				IPNet: *cidr,
			},
			OutpostARN: outpostARN,
			CIDRIndex:  cidrIndex,
		}
	}

	publicCIDR, privateCIDR := newCIDRs[0], newCIDRs[1]
	return &SubnetPair{
		Public:  []api.AZSubnetSpec{makeAZSubnetSpec(publicCIDR, existingSubnetsCount+1)},
		Private: []api.AZSubnetSpec{makeAZSubnetSpec(privateCIDR, existingSubnetsCount+2)},
	}, nil
}

func getSubnetNetworkSize(vpcCIDR net.IPNet, subnetsTotal int) (subnetSize, networkLength int, err error) {
	switch maskSize, _ := vpcCIDR.Mask.Size(); {
	case subnetsTotal == 2:
		subnetSize = 2
		networkLength = maskSize + 3
	case subnetsTotal <= 8:
		subnetSize = 8
		networkLength = maskSize + 3
	case subnetsTotal <= 16:
		subnetSize = 16
		networkLength = maskSize + 4
	default:
		return 0, 0, fmt.Errorf("cannot create more than 16 subnets, %d requested", subnetsTotal)
	}
	return subnetSize, networkLength, nil
}

func validateVPCCIDR(vpc *api.ClusterVPC) error {
	if vpc.CIDR == nil {
		cidr := api.DefaultCIDR()
		vpc.CIDR = &cidr
	}

	if prefix, _ := vpc.CIDR.Mask.Size(); prefix < 16 || prefix > 24 {
		return errors.New("VPC CIDR prefix must be between /16 and /24")
	}
	return nil
}

func SplitInto(parent *net.IPNet, size, networkLength int) ([]*net.IPNet, error) {
	if networkLength < 16 || networkLength > 28 {
		return nil, errors.New("CIDR block size must be between a /16 netmask and /28 netmask")
	}
	var subnets []*net.IPNet
	for i := 0; i < size; i++ {
		ip4 := parent.IP.To4()
		if ip4 == nil {
			return nil, fmt.Errorf("unexpected IP address type: %s", parent)
		}

		n := binary.BigEndian.Uint32(ip4)
		n += uint32(i) << uint(32-networkLength)
		subnetIP := make(net.IP, len(ip4))
		binary.BigEndian.PutUint32(subnetIP, n)

		subnets = append(subnets, &net.IPNet{
			IP:   subnetIP,
			Mask: net.CIDRMask(networkLength, 32),
		})
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
			return nil, errors.New("can't describe subnet by CIDR without VPC ID")
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
			return nil, errors.New("can't describe subnet by AZ without VPC ID")
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
func UseFromClusterStack(ctx context.Context, provider api.ClusterProvider, stack *types.Stack, spec *api.ClusterConfig, ignoreDrift bool) error {
	if spec.VPC == nil {
		spec.VPC = api.NewClusterVPC(spec.IPv6Enabled())
	}
	if spec.VPC.Subnets == nil {
		spec.VPC.Subnets = &api.ClusterSubnets{
			Public:  api.NewAZSubnetMapping(),
			Private: api.NewAZSubnetMapping(),
		}
	}
	if spec.VPC.LocalZoneSubnets == nil {
		spec.VPC.LocalZoneSubnets = &api.ClusterSubnets{
			Public:  api.NewAZSubnetMapping(),
			Private: api.NewAZSubnetMapping(),
		}
	}
	// this call is authoritative, and we can safely override the
	// CIDR, as it can only be set to anything due to defaulting
	spec.VPC.CIDR = nil

	// Cluster Endpoint Access isn't part of the EKS CloudFormation Cluster stack at this point
	// Retrieve the current configuration via the SDK
	if err := UseEndpointAccessFromCluster(ctx, provider, spec); err != nil {
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

	splitOutputValue := func(v string) []string {
		return strings.Split(v, ",")
	}
	importSubnetsFromIDList := func(subnetMapping api.AZSubnetMapping, value string) error {
		var (
			vpcSubnets   []string
			stackSubnets []string
			toBeImported []string
		)
		// collect VPC subnets as returned by CFN stack outputs
		stackSubnets = splitOutputValue(value)

		// collect VPC subnets as returned by EC2 API
		ec2API := provider.EC2()
		output, err := ec2API.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
			Filters: []ec2types.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []string{spec.VPC.ID},
				},
			},
		})
		if err != nil {
			return err
		}
		for _, o := range output.Subnets {
			vpcSubnets = append(vpcSubnets, *o.SubnetId)
		}

		// if a subnet is present on the stack outputs, but actually missing from VPC
		// e.g. it was manually deleted by the user using AWS CLI/Console
		// than log a warning and don't import it into cluster spec
		stackDriftFound := false
		for _, ssID := range stackSubnets {
			if !slices.Contains(vpcSubnets, ssID) {
				msg := fmt.Sprintf("%s was found in cluster's CloudFormation stack outputs, but has been removed from VPC %s outside of eksctl", ssID, spec.VPC.ID)
				if !ignoreDrift {
					return &StackDriftError{
						Msg: msg,
					}
				}
				stackDriftFound = true
				logger.Warning(msg)
				continue
			}
			toBeImported = append(toBeImported, ssID)
		}
		if stackDriftFound {
			logger.Warning("VPC %s contains the following subnets: %s", spec.VPC.ID, strings.Join(vpcSubnets, ","))
		}
		return ImportSubnetsFromIDList(ctx, provider.EC2(), spec, subnetMapping, toBeImported)
	}

	optionalCollectors := map[string]outputs.Collector{
		outputs.ClusterSharedNodeSecurityGroup: func(v string) error {
			spec.VPC.SharedNodeSecurityGroup = v
			return nil
		},
		outputs.ClusterSubnetsPrivate: func(v string) error {
			return importSubnetsFromIDList(spec.VPC.Subnets.Private, v)
		},
		outputs.ClusterSubnetsPublic: func(v string) error {
			return importSubnetsFromIDList(spec.VPC.Subnets.Public, v)
		},
		outputs.ClusterSubnetsPrivateLocal: func(v string) error {
			return importSubnetsFromIDList(spec.VPC.LocalZoneSubnets.Private, v)
		},
		outputs.ClusterSubnetsPublicLocal: func(v string) error {
			return importSubnetsFromIDList(spec.VPC.LocalZoneSubnets.Public, v)
		},
		outputs.ClusterSubnetsPrivateExtended: func(v string) error {
			return ImportSubnetsByIDsWithAlias(ctx, provider.EC2(), spec, spec.VPC.Subnets.Private, splitOutputValue(v), MakeExtendedSubnetAliasFunc())
		},
		outputs.ClusterSubnetsPublicExtended: func(v string) error {
			return ImportSubnetsByIDsWithAlias(ctx, provider.EC2(), spec, spec.VPC.Subnets.Public, splitOutputValue(v), MakeExtendedSubnetAliasFunc())
		},
		outputs.ClusterFullyPrivate: func(v string) error {
			spec.PrivateCluster.Enabled = v == "true"
			return nil
		},
		outputs.ClusterFeatureNATMode: func(v string) error {
			spec.VPC.NAT = &api.ClusterNAT{
				Gateway: aws.String(v),
			}
			return nil
		},
	}

	if !outputs.Exists(*stack, outputs.ClusterSubnetsPublic) &&
		outputs.Exists(*stack, outputs.ClusterSubnetsPublicLegacy) {
		optionalCollectors[outputs.ClusterSubnetsPublicLegacy] = func(v string) error {
			return importSubnetsFromIDList(spec.VPC.Subnets.Public, v)
		}
	}

	if err := outputs.Collect(*stack, requiredCollectors, optionalCollectors); err != nil {
		return err
	}
	// to clean up invalid subnets based on AZ after importing valid subnets from stack
	cleanupSubnets(spec)
	return nil
}

// MakeExtendedSubnetAliasFunc returns a function for creating an alias for a subnet that was added as part of extending
// the VPC with Outpost subnets.
func MakeExtendedSubnetAliasFunc() MakeSubnetAlias {
	subnetsCount := 0
	return func(subnet *ec2types.Subnet) string {
		subnetsCount++
		return MakeExtendedSubnetAlias(*subnet.AvailabilityZone, subnetsCount)
	}
}

// MakeExtendedSubnetAlias generates an alias for a subnet that was added as part of extending the VPC
// with Outpost subnets.
func MakeExtendedSubnetAlias(az string, ordinal int) string {
	return fmt.Sprintf("outpost-%s-%d", az, ordinal)
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

type MakeSubnetAlias func(*ec2types.Subnet) string

// ImportSubnets will update spec with subnets, if VPC ID/CIDR is unknown
// it will use provider to call describeVPC based on the VPC ID of the
// first subnet; all subnets must be in the same VPC.
// It imports the specified subnets into ClusterConfig and sets the AZs and local zones used by those subnets.
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func ImportSubnets(ctx context.Context, ec2API awsapi.EC2, spec *api.ClusterConfig, subnetMapping api.AZSubnetMapping, subnets []ec2types.Subnet, makeSubnetAlias MakeSubnetAlias) error {
	if subnetMapping == nil {
		return nil
	}
	if spec.VPC.ID != "" {
		// ensure managed NAT is disabled
		// if we are importing an existing VPC/subnets, the expectation is that the user has
		// already setup NAT, routing, etc. for these subnets
		if spec.VPC.NAT == nil {
			disable := api.ClusterDisableNAT
			spec.VPC.NAT = &api.ClusterNAT{
				Gateway: &disable,
			}
		}

		// ensure VPC gets imported and validated first, if it's already set
		if err := importVPC(ctx, ec2API, spec, spec.VPC.ID); err != nil {
			return err
		}
	}
	if makeSubnetAlias == nil {
		makeSubnetAlias = func(subnet *ec2types.Subnet) string {
			return *subnet.AvailabilityZone
		}
	}

	// as subnetMapping will be populated / altered within ImportSubnet,
	// we want to keep an unchanged copy for local against remote VPC config validation
	localSubnetConfig := api.AZSubnetMapping{}
	for k, v := range subnetMapping {
		localSubnetConfig[k] = v
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

		if err := api.ImportSubnet(subnetMapping, localSubnetConfig, &sn, makeSubnetAlias); err != nil {
			return fmt.Errorf("could not import subnet %s: %w", *sn.SubnetId, err)
		}
		spec.AppendAvailabilityZone(*sn.AvailabilityZone)
	}
	return nil
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

	if spec.IsControlPlaneOnOutposts() {
		var invalidSubnetIDs []string
		for _, subnet := range subnets {
			if subnet.OutpostArn == nil || *subnet.OutpostArn != spec.Outpost.ControlPlaneOutpostARN {
				invalidSubnetIDs = append(invalidSubnetIDs, *subnet.SubnetId)
			}
		}
		if len(invalidSubnetIDs) > 0 {
			return fmt.Errorf("all subnets must be on the control plane Outpost when specifying pre-existing subnets for a cluster on Outposts; found invalid %s subnet(s): %v", strings.ToLower(string(topology)), strings.Join(invalidSubnetIDs, ","))
		}
	}
	return ImportSubnets(ctx, ec2API, spec, subnetMapping, subnets, nil)
}

// ImportSubnetsFromIDList will update cluster config with subnets _only specified by ID_
// then pass resulting subnets to ImportSubnets
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func ImportSubnetsFromIDList(ctx context.Context, ec2API awsapi.EC2, spec *api.ClusterConfig, subnetMapping api.AZSubnetMapping, subnetIDs []string) error {
	return ImportSubnetsByIDsWithAlias(ctx, ec2API, spec, subnetMapping, subnetIDs, nil)
}

// ImportSubnetsByIDsWithAlias is like ImportSubnetsFromIDList but allows passing a function that generates an alias
// for a subnet.
func ImportSubnetsByIDsWithAlias(ctx context.Context, ec2API awsapi.EC2, spec *api.ClusterConfig, subnetMapping api.AZSubnetMapping, subnetIDs []string, makeSubnetAlias MakeSubnetAlias) error {
	subnets, err := describeSubnets(ctx, ec2API, spec.VPC.ID, subnetIDs, nil, nil)
	if err != nil {
		return err
	}

	return ImportSubnets(ctx, ec2API, spec, subnetMapping, subnets, makeSubnetAlias)
}

func ValidateLegacySubnetsForNodeGroups(ctx context.Context, spec *api.ClusterConfig, provider api.ClusterProvider) error {
	subnetsToValidate := sets.New[string]()

	selectSubnets := func(np api.NodePool) error {
		if ng := np.BaseNodeGroup(); ng.PrivateNetworking || ng.OutpostARN != "" {
			return nil
		}
		subnetIDs, err := SelectNodeGroupSubnets(ctx, np, spec, provider.EC2())
		if err != nil {
			return fmt.Errorf("could not find public subnets: %w", err)
		}
		if len(subnetIDs) > 0 {
			subnetsToValidate.Insert(subnetIDs...)
		} else {
			// This ng doesn't have AZs defined, so we need to check all public subnets
			subnetsToValidate.Insert(spec.VPC.Subnets.Public.WithIDs()...)
		}
		return nil
	}

	for _, ng := range spec.NodeGroups {
		if err := selectSubnets(ng); err != nil {
			return err
		}
	}

	for _, ng := range spec.ManagedNodeGroups {
		if err := selectSubnets(ng); err != nil {
			return err
		}
	}
	if err := ValidateExistingPublicSubnets(ctx, provider, spec.VPC.ID, sets.List(subnetsToValidate)); err != nil {
		// If the cluster endpoint is reachable from the VPC, nodes might still be able to join
		if spec.HasPrivateEndpointAccess() {
			logger.Warning("public subnets for one or more nodegroups have %q disabled. This means that nodes won't "+
				"get public IP addresses. If they can't reach the cluster through the private endpoint they won't be "+
				"able to join the cluster", "MapPublicIpOnLaunch")
			return nil
		}

		logger.Critical(err.Error())
		return fmt.Errorf("subnets for one or more new nodegroups don't meet requirements. "+
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
			return fmt.Errorf("unable to set MapPublicIpOnLaunch attribute to true for subnet %q: %w", s, err)
		}
	}
	return nil
}

// ImportSubnetsFromSpec will update spec with subnets, it will call describeSubnets first,
// then pass resulting subnets to ImportSubnets
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func ImportSubnetsFromSpec(ctx context.Context, ec2API awsapi.EC2, spec *api.ClusterConfig) error {
	if spec.VPC.ID != "" {
		// ensure VPC gets imported and validated first, if it's already set
		if err := importVPC(ctx, ec2API, spec, spec.VPC.ID); err != nil {
			return err
		}
	}
	if err := importSubnetsForTopology(ctx, ec2API, spec, api.SubnetTopologyPrivate); err != nil {
		return err
	}
	if err := importSubnetsForTopology(ctx, ec2API, spec, api.SubnetTopologyPublic); err != nil {
		return err
	}
	// to clean up invalid subnets based on AZ after importing both private and public subnets
	cleanupSubnets(spec)
	return nil
}

// UseEndpointAccessFromCluster retrieves the Cluster's endpoint access configuration via the SDK
// as the CloudFormation Stack doesn't support that configuration currently
func UseEndpointAccessFromCluster(ctx context.Context, provider api.ClusterProvider, spec *api.ClusterConfig) error {
	input := &awseks.DescribeClusterInput{
		Name: &spec.Metadata.Name,
	}
	output, err := provider.EKS().DescribeCluster(ctx, input)
	if err != nil {
		return err
	}
	if spec.VPC.ClusterEndpoints == nil {
		spec.VPC.ClusterEndpoints = &api.ClusterEndpoints{}
	}
	spec.VPC.ClusterEndpoints.PublicAccess = &output.Cluster.ResourcesVpcConfig.EndpointPublicAccess
	spec.VPC.ClusterEndpoints.PrivateAccess = &output.Cluster.ResourcesVpcConfig.EndpointPrivateAccess
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
				// since we're removing the subnet with invalid AZ from spec, we want to reference it by ID in any subsequent nodegroup creation task
				for _, node := range nodes.ToNodePools(spec) {
					for i, subnetRef := range node.BaseNodeGroup().Subnets {
						if subnetRef == name {
							node.BaseNodeGroup().Subnets[i] = subnet.ID
						}
					}
				}
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
		return ec2types.Subnet{}, fmt.Errorf("subnet with ID %q not found", id)
	}
	return output.Subnets[0], nil
}

// SelectNodeGroupSubnets returns the subnet IDs to use for a nodegroup from the specified availability zones, local zones,
// and subnets.
func SelectNodeGroupSubnets(ctx context.Context, np api.NodePool, clusterConfig *api.ClusterConfig, ec2API awsapi.EC2) ([]string, error) {

	var (
		subnetMapping        api.AZSubnetMapping
		publicSubnetMapping  api.AZSubnetMapping
		privateSubnetMapping api.AZSubnetMapping
		zones                []string
		supportedZones       *[]string
		zoneTypeMapping      *map[string]ZoneType
	)

	ng := np.BaseNodeGroup()

	if nodeGroup, ok := np.(*api.NodeGroup); ok && len(nodeGroup.LocalZones) > 0 {
		zones = nodeGroup.LocalZones
		privateSubnetMapping = clusterConfig.VPC.LocalZoneSubnets.Private
		publicSubnetMapping = clusterConfig.VPC.LocalZoneSubnets.Public
		if nodeGroup.PrivateNetworking {
			subnetMapping = privateSubnetMapping
		} else {
			subnetMapping = publicSubnetMapping
		}
	} else {
		zones = ng.AvailabilityZones
		privateSubnetMapping = clusterConfig.VPC.Subnets.Private
		publicSubnetMapping = clusterConfig.VPC.Subnets.Public
		if ng.PrivateNetworking {
			subnetMapping = privateSubnetMapping
		} else {
			subnetMapping = publicSubnetMapping
		}
	}

	makeErrorDesc := func() string {
		return fmt.Sprintf("(allSubnets=%#v localZones=%#v subnets=%#v)", subnetMapping, zones, ng.Subnets)
	}

	validateZoneType := func(zone string) error {
		if zoneTypeMapping == nil {
			output, err := DiscoverZoneTypes(ctx, ec2API, clusterConfig.Metadata.Region)
			if err != nil {
				return fmt.Errorf("error discovering zone types: %w", err)
			}
			zoneTypeMapping = &output
		}
		zoneType, ok := (*zoneTypeMapping)[zone]
		if !ok {
			return fmt.Errorf("unexpected error finding zone type for zone %q", zone)
		}
		if nodes.IsManaged(np) {
			if zoneType == ZoneTypeLocalZone {
				return fmt.Errorf("managed nodegroups cannot be launched in local zones: %q", ng.Name)
			}
			return nil
		}
		if zoneType == ZoneTypeAvailabilityZone {
			logger.Warning("subnets contain a mix of both local and availability zones")
		}
		return nil
	}

	validateZoneInstanceSupport := func(zone string) error {
		if supportedZones == nil {
			output, err := az.FilterBasedOnAvailability(ctx, clusterConfig.AvailabilityZones, []api.NodePool{np}, ec2API)
			if err != nil {
				return err
			}
			supportedZones = &output
		}
		if zoneTypeMapping == nil {
			output, err := DiscoverZoneTypes(ctx, ec2API, clusterConfig.Metadata.Region)
			if err != nil {
				return fmt.Errorf("error discovering zone types: %w", err)
			}
			zoneTypeMapping = &output
		}
		zoneType, ok := (*zoneTypeMapping)[zone]
		if !ok {
			return fmt.Errorf("unexpected error finding zone type for zone %q", zone)
		}
		// only validate instance support for availability zones
		if zoneType == ZoneTypeAvailabilityZone &&
			slices.Contains(clusterConfig.AvailabilityZones, zone) && // for now, we won't validate support for user specified new zones
			!slices.Contains(*supportedZones, zone) {
			return fmt.Errorf("cannot create nodegroup %s in availability zone %s as it does not support all required instance types",
				np.BaseNodeGroup().Name, zone)
		}
		return nil
	}

	var subnetIDs []string
	if len(zones) > 0 {
		var err error
		if subnetIDs, err = selectNodeGroupZoneSubnets(zones, subnetMapping, validateZoneInstanceSupport); err != nil {
			return nil, fmt.Errorf("could not find %s subnets for zones %q %s: %w", getNetworkType(ng), zones, makeErrorDesc(), err)
		}
	}

	if len(ng.Subnets) > 0 {
		subnetsFromIDs, err := selectNodeGroupSubnetsFromIDs(ctx, ng, publicSubnetMapping, privateSubnetMapping, clusterConfig, ec2API, validateZoneType, validateZoneInstanceSupport)
		if err != nil {
			return nil, fmt.Errorf("could not select subnets from subnet IDs %s: %w", makeErrorDesc(), err)
		}
		subnetIDs = append(subnetIDs, subnetsFromIDs...)
	} else if ng.OutpostARN != "" {
		subnetIDs = subnetMapping.SelectOutpostSubnetIDs()
		if len(subnetIDs) == 0 {
			return nil, fmt.Errorf("no %s subnets exist in Outpost for nodegroup %s", getNetworkType(ng), ng.Name)
		}
	}

	if api.IsEnabled(ng.EFAEnabled) && len(subnetIDs) > 0 {
		subnetIDs = subnetIDs[:1]
		logger.Info("EFA requires all nodes be in a single subnet, arbitrarily choosing one: %s", subnetIDs)
	}

	return subnetIDs, nil
}

func selectNodeGroupZoneSubnets(nodeGroupZones []string,
	subnetMapping api.AZSubnetMapping,
	validateSubnetZoneInstanceSupport func(zone string) error) ([]string, error) {
	makeErrorDesc := func() string {
		return fmt.Sprintf("(allSubnets=%#v zones=%#v)", subnetMapping, nodeGroupZones)
	}
	if len(subnetMapping) < len(nodeGroupZones) {
		return nil, fmt.Errorf("mapping does not have enough subnets: %s", makeErrorDesc())
	}

	var subnetIDs []string
	for _, zone := range nodeGroupZones {
		found := false
		for _, s := range subnetMapping {
			if s.AZ == zone {
				if err := validateSubnetZoneInstanceSupport(zone); err != nil {
					return nil, fmt.Errorf("failed to select subnet %s: %w", zone, err)
				}
				subnetIDs = append(subnetIDs, s.ID)
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("mapping does not have subnet with zone %q: %s", zone, makeErrorDesc())
		}
	}

	return subnetIDs, nil
}

func selectNodeGroupSubnetsFromIDs(
	ctx context.Context,
	ng *api.NodeGroupBase,
	publicSubnetMapping api.AZSubnetMapping,
	privateSubnetMapping api.AZSubnetMapping,
	clusterConfig *api.ClusterConfig,
	ec2API awsapi.EC2,
	validateSubnetZoneType func(zone string) error,
	validateSubnetZoneInstanceSupport func(zone string) error) ([]string, error) {
	var (
		mappedSubnet      *api.AZSubnetSpec
		ec2Subnet         ec2types.Subnet
		selectedSubnetIDs []string
		subnetID          string
		subnetOutpostARN  string
		subnetZone        string
		foundInPrivate    bool
		foundInPublic     bool
		err               error
	)

	outpostARN := getOutpostARN(clusterConfig, ng)

	for _, subnetName := range ng.Subnets {
		// first try to find the specified subnet as part of the VPC inside ClusterConfig, if existent
		if ng.PrivateNetworking {
			if mappedSubnet, foundInPrivate = findInConfiguredVPC(subnetName, privateSubnetMapping); !foundInPrivate {
				if mappedSubnet, foundInPublic = findInConfiguredVPC(subnetName, publicSubnetMapping); foundInPublic {
					logger.Warning("public subnet %s is being used with `privateNetworking` enabled, please ensure this is the desired behaviour", subnetName)
				}
			}
		} else {
			if mappedSubnet, foundInPublic = findInConfiguredVPC(subnetName, publicSubnetMapping); !foundInPublic {
				if mappedSubnet, foundInPrivate = findInConfiguredVPC(subnetName, privateSubnetMapping); foundInPrivate {
					return nil, fmt.Errorf("subnet %s is specified as private in ClusterConfig, thus must only be used when `privateNetworking` is enabled", subnetName)
				}
			}
		}

		if mappedSubnet != nil {
			subnetID = mappedSubnet.ID
			subnetOutpostARN = mappedSubnet.OutpostARN
			subnetZone = mappedSubnet.AZ
		} else {
			// otherwise try to find the subnet as part of the AWS Account
			ec2Subnet, err = getSubnetByID(ctx, ec2API, subnetName)
			if err != nil {
				return nil, err
			}
			if ec2Subnet.VpcId != nil && *ec2Subnet.VpcId != clusterConfig.VPC.ID {
				return nil, fmt.Errorf("subnet with ID %q is not in the attached VPC with ID %q", *ec2Subnet.SubnetId, clusterConfig.VPC.ID)
			}
			subnetID = *ec2Subnet.SubnetId
			subnetOutpostARN = aws.ToString(ec2Subnet.OutpostArn)
			subnetZone = *ec2Subnet.AvailabilityZone
		}

		if err := validateSubnetZoneType(subnetZone); err != nil {
			return nil, fmt.Errorf("failed to select subnet %s: %w", subnetName, err)
		}
		if err := validateSubnetZoneInstanceSupport(subnetZone); err != nil {
			return nil, fmt.Errorf("failed to select subnet %s: %w", subnetName, err)
		}
		if err := validateSubnetOnOutposts(subnetID, subnetOutpostARN, outpostARN); err != nil {
			return nil, fmt.Errorf("failed to select subnet %s: %w", subnetName, err)
		}

		selectedSubnetIDs = append(selectedSubnetIDs, subnetID)
	}

	return selectedSubnetIDs, nil
}

func findInConfiguredVPC(subnetName string, subnetMapping api.AZSubnetMapping) (*api.AZSubnetSpec, bool) {
	if subnet, ok := subnetMapping[subnetName]; !ok {
		// if not found by name, search by id
		for _, s := range subnetMapping {
			if s.ID != subnetName {
				continue
			}
			return &s, true
		}
	} else {
		// if found by name, return the subnet
		return &subnet, true
	}
	return nil, false
}

func getNetworkType(ng *api.NodeGroupBase) string {
	if ng.PrivateNetworking {
		return "private"
	}
	return "public"
}

func getOutpostARN(clusterConfig *api.ClusterConfig, ng *api.NodeGroupBase) string {
	if clusterConfig.IsControlPlaneOnOutposts() {
		return clusterConfig.Outpost.ControlPlaneOutpostARN
	}
	return ng.OutpostARN
}

func validateSubnetOnOutposts(subnetID, subnetOutpostARN, outpostARN string) error {
	if outpostARN == "" {
		return nil
	}
	if subnetOutpostARN == "" {
		return fmt.Errorf("subnet %q is not on Outposts", subnetID)
	}
	if subnetOutpostARN != outpostARN {
		return fmt.Errorf("subnet %q is in a different Outpost ARN (%q) than the control plane or nodegroup Outpost (%q)", subnetID, subnetOutpostARN, outpostARN)
	}
	return nil
}
