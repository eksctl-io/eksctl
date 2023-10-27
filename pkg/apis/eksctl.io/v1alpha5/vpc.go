package v1alpha5

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"

	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
)

// Values for `ClusterNAT`
const (
	// ClusterHighlyAvailableNAT configures a highly available NAT gateway
	ClusterHighlyAvailableNAT = "HighlyAvailable"

	// ClusterSingleNAT configures a single NAT gateway
	ClusterSingleNAT = "Single"

	// ClusterDisableNAT disables NAT
	ClusterDisableNAT = "Disable"

	// (default)
	ClusterNATDefault = ClusterSingleNAT
)

// AZSubnetMapping holds subnet to AZ mappings.
// If the key is an AZ, that also becomes the name of the subnet
// otherwise use the key to refer to this subnet.
// Schema type is `map[string]AZSubnetSpec`
type AZSubnetMapping map[string]AZSubnetSpec

func NewAZSubnetMapping() AZSubnetMapping {
	return make(map[string]AZSubnetSpec)
}

func AZSubnetMappingFromMap(m map[string]AZSubnetSpec) AZSubnetMapping {
	for k := range m {
		v := m[k]
		if v.AZ == "" {
			v.AZ = k
			m[k] = v
		}
	}
	return m
}

func (m *AZSubnetMapping) Set(name string, spec AZSubnetSpec) {
	if m == nil {
		m = &AZSubnetMapping{}
	}
	(*m)[name] = spec
}

func (m *AZSubnetMapping) SetAZ(az string, spec Network) {
	if m == nil {
		m = &AZSubnetMapping{}
	}
	(*m)[az] = AZSubnetSpec{
		ID:   spec.ID,
		AZ:   az,
		CIDR: spec.CIDR,
	}
}

// WithIDs returns list of subnet ids
func (m *AZSubnetMapping) WithIDs() []string {
	if m == nil {
		return nil
	}
	subnets := []string{}
	for _, s := range *m {
		if s.ID != "" {
			subnets = append(subnets, s.ID)
		}
	}
	return subnets
}

// WithCIDRs returns list of subnet CIDRs
func (m *AZSubnetMapping) WithCIDRs() []string {
	if m == nil {
		return nil
	}
	subnets := []string{}
	for _, s := range *m {
		if s.CIDR != nil && s.ID == "" {
			subnets = append(subnets, s.CIDR.String())
		}
	}
	return subnets
}

// WithAZs returns list of subnet AZs
func (m *AZSubnetMapping) WithAZs() []string {
	if m == nil {
		return nil
	}
	subnets := []string{}
	for _, s := range *m {
		if s.AZ != "" && s.CIDR == nil && s.ID == "" {
			subnets = append(subnets, s.AZ)
		}
	}
	return subnets
}

// UnmarshalJSON parses JSON data into a value
func (m *AZSubnetMapping) UnmarshalJSON(b []byte) error {
	// TODO we need to validate that the AZ property is maintained
	var raw map[string]AZSubnetSpec
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	*m = AZSubnetMappingFromMap(raw)
	return nil
}

type (
	// ClusterVPC holds global subnet and all child subnets
	ClusterVPC struct {
		// global CIDR and VPC ID
		// +optional
		Network
		// SecurityGroup (aka the ControlPlaneSecurityGroup) for communication between control plane and nodes
		// +optional
		SecurityGroup string `json:"securityGroup,omitempty"`
		// Subnets are keyed by AZ for convenience.
		// See [this example](/examples/reusing-iam-and-vpc/)
		// as well as [using existing
		// VPCs](/usage/vpc-networking/#use-existing-vpc-other-custom-configuration).
		// +optional
		Subnets *ClusterSubnets `json:"subnets,omitempty"`

		// LocalZoneSubnets represents subnets in local zones.
		// This field is used internally and is not part of the ClusterConfig schema.
		LocalZoneSubnets *ClusterSubnets `json:"-"`

		// HostnameType is the type of hostname to use for EC2 instances.
		HostnameType string `json:"hostnameType,omitempty"`

		// for additional CIDR associations, e.g. a CIDR for
		// private subnets or any ad-hoc subnets
		// +optional
		ExtraCIDRs []string `json:"extraCIDRs,omitempty"`
		// for additional IPv6 CIDR associations, e.g. a CIDR for
		// private subnets or any ad-hoc subnets
		// +optional
		ExtraIPv6CIDRs []string `json:"extraIPv6CIDRs,omitempty"`
		// for pre-defined shared node SG
		SharedNodeSecurityGroup string `json:"sharedNodeSecurityGroup,omitempty"`
		// Automatically add security group rules to and from the default
		// cluster security group and the shared node security group.
		// This allows unmanaged nodes to communicate with the control plane
		// and managed nodes.
		// This option cannot be disabled when using eksctl created security groups.
		// Defaults to `true`
		// +optional
		ManageSharedNodeSecurityGroupRules *bool `json:"manageSharedNodeSecurityGroupRules,omitempty"`
		// AutoAllocateIPV6 requests an IPv6 CIDR block with /56 prefix for the VPC
		// +optional
		AutoAllocateIPv6 *bool `json:"autoAllocateIPv6,omitempty"`
		// +optional
		NAT *ClusterNAT `json:"nat,omitempty"`
		// See [managing access to API](/usage/vpc-networking/#managing-access-to-the-kubernetes-api-server-endpoints)
		// +optional
		ClusterEndpoints *ClusterEndpoints `json:"clusterEndpoints,omitempty"`
		// PublicAccessCIDRs are which CIDR blocks to allow access to public
		// k8s API endpoint
		// +optional
		PublicAccessCIDRs []string `json:"publicAccessCIDRs,omitempty"`
		// ControlPlaneSubnetIDs configures the subnets for the control plane.
		// +optional
		ControlPlaneSubnetIDs []string `json:"controlPlaneSubnetIDs,omitempty"`
		// ControlPlaneSecurityGroupIDs configures the security groups for the control plane.
		// +optional
		ControlPlaneSecurityGroupIDs []string `json:"controlPlaneSecurityGroupIDs,omitempty"`
	}
	// ClusterSubnets holds private and public subnets
	ClusterSubnets struct {
		Private AZSubnetMapping `json:"private,omitempty"`
		Public  AZSubnetMapping `json:"public,omitempty"`
	}

	// SubnetTopology can be SubnetTopologyPrivate or SubnetTopologyPublic
	SubnetTopology string
	AZSubnetSpec   struct {
		// +optional
		ID string `json:"id,omitempty"`
		// AZ is the zone name for this subnet, it can either be an availability zone name
		// or a local zone name.
		// AZ can be omitted if the key is an AZ.
		// +optional
		AZ string `json:"az,omitempty"`
		// +optional
		CIDR *ipnet.IPNet `json:"cidr,omitempty"`

		CIDRIndex int `json:"-"`

		OutpostARN string `json:"-"`
	}
	// Network holds ID and CIDR
	Network struct {
		// +optional
		ID string `json:"id,omitempty"`
		// +optional
		CIDR *ipnet.IPNet `json:"cidr,omitempty"`
		// +optional
		IPv6Cidr string `json:"ipv6Cidr,omitempty"`
		// +optional
		IPv6Pool string `json:"ipv6Pool,omitempty"`
	}
	// ClusterNAT NAT config
	ClusterNAT struct {
		// Valid variants are `ClusterNAT` constants
		Gateway *string `json:"gateway,omitempty"`
	}

	// ClusterEndpoints holds cluster api server endpoint access information
	ClusterEndpoints struct {
		PrivateAccess *bool `json:"privateAccess,omitempty"`
		PublicAccess  *bool `json:"publicAccess,omitempty"`
	}
)

const (
	// MinRequiredSubnets is the minimum required number of subnets
	MinRequiredSubnets = 2
	// OutpostsMinRequiredSubnets is the minimum required number of subnets for Outposts.
	OutpostsMinRequiredSubnets = 1
	// MinRequiredAvailabilityZones defines the minimum number of required availability zones
	MinRequiredAvailabilityZones = MinRequiredSubnets
	// RecommendedSubnets is the recommended number of subnets
	RecommendedSubnets = 3
	// RecommendedAvailabilityZones defines the default number of required availability zones
	RecommendedAvailabilityZones = RecommendedSubnets
	// SubnetTopologyPrivate represents privately-routed subnets
	SubnetTopologyPrivate SubnetTopology = "Private"
	// SubnetTopologyPublic represents publicly-routed subnets
	SubnetTopologyPublic SubnetTopology = "Public"
)

// SubnetTopologies returns a list of topologies
func SubnetTopologies() []SubnetTopology {
	return []SubnetTopology{
		SubnetTopologyPrivate,
		SubnetTopologyPublic,
	}
}

// DefaultCIDR returns default global CIDR for VPC
func DefaultCIDR() ipnet.IPNet {
	return ipnet.IPNet{
		IPNet: net.IPNet{
			IP:   []byte{192, 168, 0, 0},
			Mask: []byte{255, 255, 0, 0},
		},
	}
}

// ImportSubnet loads a given subnet into ClusterConfig.
// Note that the user must use
// either AZs as keys
// OR names as keys and specify
//
//	the ID (optionally with AZ and CIDR)
//	OR AZ, optionally with CIDR.
//
// If a user specifies a subnet by AZ without CIDR and ID but multiple subnets
// exist in this VPC, one will be arbitrarily chosen.
func ImportSubnet(subnets AZSubnetMapping, localSubnetsConfig AZSubnetMapping, subnet *ec2types.Subnet, makeSubnetAlias func(*ec2types.Subnet) string) error {
	if localSubnetsConfig == nil {
		return nil
	}

	remoteSubnet, err := remoteSubnetToAZSubnetSpec(subnet)
	if err != nil {
		return err
	}

	// if a VPC config was provided as part of the config file,
	// we need to validate it against the remote config
	// and return an error in case of mismatch
	subnetKey, err := validateLocalConfigAgainstRemote(localSubnetsConfig, remoteSubnet, makeSubnetAlias(subnet))
	if err != nil {
		return fmt.Errorf("mismatch found between local and remote VPC config: %w", err)
	}

	subnets[subnetKey] = remoteSubnet

	return nil
}

func validateLocalConfigAgainstRemote(localSubnetsConfig AZSubnetMapping, remoteSubnet AZSubnetSpec, subnetAlias string) (string, error) {
	if len(localSubnetsConfig) == 0 {
		return subnetAlias, nil
	}

	// if the subnet is found by alias in config file, validate ID and CIDR
	if localSubnet, ok := localSubnetsConfig[subnetAlias]; ok {
		if localSubnet.ID != "" && localSubnet.ID != remoteSubnet.ID {
			return "", fmt.Errorf("subnet ID %q, found in config file, is not the same as subnet ID %q, found in remote VPC config", localSubnet.ID, remoteSubnet.ID)
		}
		if localSubnet.CIDR.String() != "" && localSubnet.CIDR.String() != remoteSubnet.CIDR.String() {
			return "", fmt.Errorf("subnet CIDR %q, found in config file, is not the same as subnet CIDR %q, found in remote VPC config", localSubnet.CIDR.String(), remoteSubnet.CIDR.String())
		}
		return subnetAlias, nil
	}

	// otherwise look up the remote subnet by ID or <AZ,CIDR> pair
	var foundByIDKey string
	var foundByAZCIDRKey string
	for k, s := range localSubnetsConfig {
		if s.ID == remoteSubnet.ID {
			if s.CIDR.String() != "" && s.CIDR.String() != remoteSubnet.CIDR.String() {
				return "", fmt.Errorf("subnet CIDR %q, found in config file, is not the same as subnet CIDR %q, found in remote VPC config", s.CIDR.String(), remoteSubnet.CIDR.String())
			}
			if foundByIDKey != "" {
				return "", fmt.Errorf("unable to unambiguously import subnet by ID, found both %s and %s", foundByIDKey, k)
			}
			foundByIDKey = k
		} else if s.ID == "" {
			if s.AZ != remoteSubnet.AZ || (s.CIDR.String() != "" && s.CIDR.String() != remoteSubnet.CIDR.String()) {
				continue
			}
			if foundByAZCIDRKey != "" {
				return "", fmt.Errorf("unable to unambiguously import subnet by <AZ,CIDR> pair, found both %s and %s", foundByAZCIDRKey, k)
			}
			foundByAZCIDRKey = k
		}
	}

	if foundByIDKey != "" {
		return foundByIDKey, nil
	}
	if foundByAZCIDRKey != "" {
		return foundByAZCIDRKey, nil
	}

	return subnetAlias, nil
}

func remoteSubnetToAZSubnetSpec(subnet *ec2types.Subnet) (AZSubnetSpec, error) {
	subnetCIDR, err := ipnet.ParseCIDR(*subnet.CidrBlock)
	if err != nil {
		return AZSubnetSpec{}, fmt.Errorf("unexpected error parsing subnet CIDR %q: %w", *subnet.CidrBlock, err)
	}

	subnetSpec := AZSubnetSpec{
		ID:   *subnet.SubnetId,
		AZ:   *subnet.AvailabilityZone,
		CIDR: subnetCIDR,
	}

	if subnet.OutpostArn != nil {
		subnetSpec.OutpostARN = *subnet.OutpostArn
	}

	return subnetSpec, nil
}

// SelectOutpostSubnetIDs returns all subnets that are on Outposts.
func (m AZSubnetMapping) SelectOutpostSubnetIDs() []string {
	var subnetIDs []string
	for _, s := range m {
		if s.OutpostARN != "" {
			subnetIDs = append(subnetIDs, s.ID)
		}
	}
	return subnetIDs
}

func (m AZSubnetMapping) getOutpostARN() (outpostARN string, found bool) {
	for _, s := range m {
		if s.OutpostARN != "" {
			return s.OutpostARN, true
		}
	}
	return "", false
}

// FindOutpostSubnetsARN finds all subnets that are on Outposts and returns the Outpost ARN.
func (v *ClusterVPC) FindOutpostSubnetsARN() (outpostARN string, found bool) {
	outpostARN, found = v.Subnets.Private.getOutpostARN()
	if found {
		return outpostARN, true
	}
	return v.Subnets.Public.getOutpostARN()

}

// SubnetInfo returns a string containing VPC subnet information
// Useful for error messages and logs
func (c *ClusterConfig) SubnetInfo() string {
	return fmt.Sprintf("VPC (%s) and subnets (private:%v public:%v)",
		c.VPC.ID, c.VPC.Subnets.Private, c.VPC.Subnets.Public)
}

// HasAnySubnets checks if any subnets were set
func (c *ClusterConfig) HasAnySubnets() bool {
	return c.VPC.Subnets != nil && (len(c.VPC.Subnets.Private) > 0 || len(c.VPC.Subnets.Public) > 0)
}

// HasSufficientPrivateSubnets validates if there is a sufficient
// number of private subnets available to create a cluster
func (c *ClusterConfig) HasSufficientPrivateSubnets() bool {
	subnetsCount := len(c.VPC.Subnets.Private)
	if c.IsControlPlaneOnOutposts() {
		return subnetsCount >= OutpostsMinRequiredSubnets
	}
	return subnetsCount >= MinRequiredSubnets
}

// CanUseForPrivateNodeGroups checks whether specified NodeGroups have enough
// private subnets when private networking is enabled
func (c *ClusterConfig) CanUseForPrivateNodeGroups() error {
	for _, ng := range c.NodeGroups {
		if ng.PrivateNetworking && !c.HasSufficientPrivateSubnets() {
			return errors.New("none or too few private subnets to use with --node-private-networking")
		}
	}
	return nil
}

// insufficientSubnetsError represents an error for when the minimum required subnets are not provided.
type insufficientSubnetsError struct {
	controlPlaneOnOutposts bool
}

// Error implements the error interface.
func (e *insufficientSubnetsError) Error() string {
	msg := "insufficient number of subnets, at least %[1]dx public and/or %[1]dx private subnets are required"
	minSubnets := MinRequiredSubnets
	if e.controlPlaneOnOutposts {
		msg += " for Outposts"
		minSubnets = OutpostsMinRequiredSubnets
	}
	return fmt.Sprintf(msg, minSubnets)
}

// HasSufficientSubnets validates if there is a sufficient number
// of either private and/or public subnets available to create
// a cluster, i.e. either non-zero of public or private, and not
// less then MinRequiredSubnets of each, but allowing to have
// public-only or private-only
func (c *ClusterConfig) HasSufficientSubnets() error {
	if !c.HasAnySubnets() {
		return &insufficientSubnetsError{
			controlPlaneOnOutposts: c.IsControlPlaneOnOutposts(),
		}
	}

	if c.IsControlPlaneOnOutposts() {
		return nil
	}

	if numPublic := len(c.VPC.Subnets.Public); numPublic > 0 && numPublic < MinRequiredSubnets {
		return &insufficientSubnetsError{}
	}

	if numPrivate := len(c.VPC.Subnets.Private); numPrivate > 0 && numPrivate < MinRequiredSubnets {
		return &insufficientSubnetsError{}
	}

	return nil
}

// DefaultEndpointsMsg returns a message that the EndpointAccess is the same as the default.
func (c *ClusterConfig) DefaultEndpointsMsg() string {
	return fmt.Sprintf(
		"Kubernetes API endpoint access will use default of {publicAccess=true, privateAccess=false} for cluster %q in %q", c.Metadata.Name, c.Metadata.Region)
}

// CustomEndpointsMsg returns a message indicating the EndpointAccess given by the user.
func (c *ClusterConfig) CustomEndpointsMsg() string {
	return fmt.Sprintf(
		"Kubernetes API endpoint access will use provided values {publicAccess=%v, privateAccess=%v} for cluster %q in %q", *c.VPC.ClusterEndpoints.PublicAccess, *c.VPC.ClusterEndpoints.PrivateAccess, c.Metadata.Name, c.Metadata.Region)
}

// UpdateEndpointsMsg returns a message indicating that they need to use `eksctl utils` to make this config.
func (c *ClusterConfig) UpdateEndpointsMsg() string {
	return fmt.Sprintf(
		"you can update Kubernetes API endpoint access with `eksctl utils update-cluster-endpoints --region=%s --name=%s --private-access=bool --public-access=bool`", c.Metadata.Region, c.Metadata.Name)
}

// EndpointsEqual returns true of two endpoints have same values after dereferencing any pointers
func EndpointsEqual(a, b ClusterEndpoints) bool {
	return reflect.DeepEqual(a, b)
}

// HasClusterEndpointAccess determines if endpoint access was configured in config file or not.
func (c *ClusterConfig) HasClusterEndpointAccess() bool {
	if c.VPC != nil && c.VPC.ClusterEndpoints != nil {
		hasPublicAccess := aws.ToBool(c.VPC.ClusterEndpoints.PublicAccess)
		hasPrivateAccess := aws.ToBool(c.VPC.ClusterEndpoints.PrivateAccess)
		return hasPublicAccess || hasPrivateAccess
	}
	return true
}

func (c *ClusterConfig) HasPrivateEndpointAccess() bool {
	return c.VPC != nil && c.VPC.ClusterEndpoints != nil && IsEnabled(c.VPC.ClusterEndpoints.PrivateAccess)
}
