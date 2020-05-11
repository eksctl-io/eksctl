package vpc

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/kops/util/pkg/slice"

	"github.com/kris-nova/logger"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	awseks "github.com/aws/aws-sdk-go/service/eks"

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

// describeSubnets fetches subnet metadata from EC2
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

// UseFromCluster retrieves the VPC configuration from an existing cluster
// based on stack outputs
// NOTE: it doesn't expect any fields in spec.VPC to be set, the remote state
// is treated as the source of truth
func UseFromCluster(provider api.ClusterProvider, stack *cfn.Stack, spec *api.ClusterConfig) error {
	if spec.VPC == nil {
		spec.VPC = api.NewClusterVPC()
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

// importVPC will update spec with VPC ID/CIDR
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func importVPC(provider api.ClusterProvider, spec *api.ClusterConfig, id string) error {
	vpc, err := describeVPC(provider, id)
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
		// ensure managed NAT is disabled
		// if we are importing an existing VPC/subnets, the expectation is that the user has
		// already setup NAT, routing, etc. for these subnets
		disable := api.ClusterDisableNAT
		spec.VPC.NAT = &api.ClusterNAT{
			Gateway: &disable,
		}

		// ensure VPC gets imported and validated first, if it's already set
		if err := importVPC(provider, spec, spec.VPC.ID); err != nil {
			return err
		}
	}

	for _, sn := range subnets {
		if spec.VPC.ID == "" {
			// if VPC wasn't defined, import it based on VPC of the first
			// subnet that we have
			if err := importVPC(provider, spec, *sn.VpcId); err != nil {
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

func ValidateLegacySubnetsForNodeGroups(spec *api.ClusterConfig, provider api.ClusterProvider) error {
	subnetsToValidate := sets.NewString()

	selectSubnets := func(azs []string) error {
		if len(azs) > 0 {
			// Check only the public subnets that this ng has
			subnetIDs, err := SelectPublicNodeGroupSubnets(azs, spec)
			if err != nil {
				return err
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
		err := selectSubnets(ng.AvailabilityZones)
		if err != nil {
			return err
		}
	}

	for _, ng := range spec.ManagedNodeGroups {
		err := selectSubnets(ng.AvailabilityZones)
		if err != nil {
			return err
		}
	}

	if err := ValidateExistingPublicSubnets(provider, subnetsToValidate.List()); err != nil {
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
func ValidateExistingPublicSubnets(provider api.ClusterProvider, subnetIDs []string) error {
	if len(subnetIDs) == 0 {
		return nil
	}
	subnets, err := describeSubnets(provider, subnetIDs...)
	if err != nil {
		return err
	}
	if err := validatePublicSubnet(subnets); err != nil {
		return err
	}
	return nil
}

// EnsureMapPublicIPOnLaunchEnabled will enable MapPublicIpOnLaunch in EC2 for all given subnet IDs
func EnsureMapPublicIPOnLaunchEnabled(provider api.ClusterProvider, subnetIDs []string) error {
	if len(subnetIDs) == 0 {
		logger.Debug("no subnets to update")
		return nil
	}

	for _, s := range subnetIDs {
		input := &ec2.ModifySubnetAttributeInput{
			SubnetId:            aws.String(s),
			MapPublicIpOnLaunch: &ec2.AttributeBooleanValue{Value: aws.Bool(true)},
		}

		logger.Debug("enabling MapPublicIpOnLaunch for subnet %q", s)
		_, err := provider.EC2().ModifySubnetAttribute(input)
		if err != nil {
			return errors.Wrapf(err, "unable to set MapPublicIpOnLaunch attribute to true for subnet %q", s)
		}
	}
	return nil
}

// ImportAllSubnets will update spec with subnets, it will call describeSubnets first,
// then pass resulting subnets to ImportSubnets
// NOTE: it does respect all fields set in spec.VPC, and will error if
// there is a mismatch of local vs remote states
func ImportAllSubnets(provider api.ClusterProvider, spec *api.ClusterConfig) error {
	if spec.VPC.ID != "" {
		// ensure VPC gets imported and validated first, if it's already set
		if err := importVPC(provider, spec, spec.VPC.ID); err != nil {
			return err
		}
	}
	if err := ImportSubnetsFromList(provider, spec, api.SubnetTopologyPrivate, spec.PrivateSubnetIDs()); err != nil {
		return err
	}
	if err := ImportSubnetsFromList(provider, spec, api.SubnetTopologyPublic, spec.PublicSubnetIDs()); err != nil {
		return err
	}
	// to clean up invalid subnets based on AZ after imported both private and public subnets
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
	for id := range spec.VPC.Subnets.Private {
		if !slice.Contains(spec.AvailabilityZones, id) {
			delete(spec.VPC.Subnets.Private, id)
		}
	}

	for id := range spec.VPC.Subnets.Public {
		if !slice.Contains(spec.AvailabilityZones, id) {
			delete(spec.VPC.Subnets.Public, id)
		}
	}
}

func validatePublicSubnet(subnets []*ec2.Subnet) error {
	legacySubnets := make([]string, 0)
	for _, sn := range subnets {
		if sn.MapPublicIpOnLaunch == nil || !*sn.MapPublicIpOnLaunch {
			legacySubnets = append(legacySubnets, *sn.SubnetId)
		}
	}
	if len(legacySubnets) > 0 {
		return fmt.Errorf("found mis-configured subnets %q. Expected public subnets with property "+
			"\"MapPublicIpOnLaunch\" enabled. Without it new nodes won't get an IP assigned", legacySubnets)
	}

	return nil
}

func SelectPublicNodeGroupSubnets(nodegroupAZs []string, clusterSpec *api.ClusterConfig) ([]string, error) {
	numNodeGroupsAZs := len(nodegroupAZs)
	if numNodeGroupsAZs == 0 {
		return nil, nil
	}

	subnets := clusterSpec.VPC.Subnets.Public
	makeErrorDesc := func() string {
		return fmt.Sprintf("(subnets=%#v AZs=%#v)", subnets, nodegroupAZs)
	}
	if len(subnets) < numNodeGroupsAZs {
		return nil, fmt.Errorf("VPC doesn't have enough subnets for nodegroup AZs %s", makeErrorDesc())
	}
	subnetIDs := make([]string, numNodeGroupsAZs)
	for i, az := range nodegroupAZs {
		subnet, ok := subnets[az]
		if !ok {
			return nil, fmt.Errorf("VPC doesn't have subnets in %s %s", az, makeErrorDesc())
		}

		subnetIDs[i] = subnet.ID
	}
	return subnetIDs, nil
}
