package v1alpha5

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"

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
	return AZSubnetMapping(make(map[string]AZSubnetSpec))
}

func AZSubnetMappingFromMap(m map[string]AZSubnetSpec) AZSubnetMapping {
	for k := range m {
		v := m[k]
		if v.AZ == "" {
			v.AZ = k
			m[k] = v
		}
	}
	return AZSubnetMapping(m)
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
		// for additional CIDR associations, e.g. a CIDR for
		// private subnets or any ad-hoc subnets
		// +optional
		ExtraCIDRs []string `json:"extraCIDRs,omitempty"`
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
		// AZ can be omitted if the key is an AZ
		// +optional
		AZ string `json:"az,omitempty"`
		// +optional
		CIDR *ipnet.IPNet `json:"cidr,omitempty"`
	}
	// Network holds ID and CIDR
	Network struct {
		// +optional
		ID string `json:"id,omitempty"`
		// +optional
		CIDR *ipnet.IPNet `json:"cidr,omitempty"`
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
	// RecommendedSubnets is the recommended number of subnets
	RecommendedSubnets = 3
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

// ImportSubnet loads a given subnet into cluster config
func (c *ClusterConfig) ImportSubnet(topology SubnetTopology, az, subnetID, cidr string) error {
	if c.VPC.Subnets == nil {
		c.VPC.Subnets = &ClusterSubnets{
			Private: NewAZSubnetMapping(),
			Public:  NewAZSubnetMapping(),
		}
	}

	var subnetMapping AZSubnetMapping
	switch topology {
	case SubnetTopologyPrivate:
		subnetMapping = c.VPC.Subnets.Private
	case SubnetTopologyPublic:
		subnetMapping = c.VPC.Subnets.Public
	default:
		panic(fmt.Sprintf("unexpected subnet topology: %s", topology))
	}

	if err := doImportSubnet(subnetMapping, az, subnetID, cidr); err != nil {
		return errors.Wrapf(err, "couldn't import subnet %s", subnetID)
	}
	return nil
}

// Note that the user must use
// EITHER AZs as keys
// OR names as keys and specify
//    the ID (optionally with AZ and CIDR)
//    OR AZ, optionally with CIDR
// If a user specifies a subnet by AZ without CIDR and ID but multiple subnets
// exist in this VPC, one will be arbitrarily chosen
func doImportSubnet(subnets AZSubnetMapping, az, subnetID, cidr string) error {
	subnetCIDR, _ := ipnet.ParseCIDR(cidr)

	if subnets == nil {
		return nil
	}

	if network, ok := subnets[az]; !ok {
		newS := AZSubnetSpec{ID: subnetID, AZ: az, CIDR: subnetCIDR}
		// Used if we find an exact ID match
		var idKey string
		// Used if we match to AZ/CIDR
		var guessKey string
		for k, s := range subnets {
			if s.ID == "" {
				if s.AZ != az || (s.CIDR.String() != "" && s.CIDR.String() != subnetCIDR.String()) {
					continue
				}
				if guessKey != "" {
					return fmt.Errorf("unable to unambiguously import subnet, found both %s and %s", guessKey, k)
				}
				guessKey = k
			} else if s.ID == subnetID {
				if s.CIDR.String() != "" && s.CIDR.String() != subnetCIDR.String() {
					return fmt.Errorf("subnet CIDR %q is not the same as %q", s.CIDR.String(), subnetCIDR.String())
				}
				if idKey != "" {
					return fmt.Errorf("unable to unambiguously import subnet, found both %s and %s", idKey, k)
				}
				idKey = k
			}
		}
		if idKey != "" {
			subnets[idKey] = newS
		} else if guessKey != "" {
			subnets[guessKey] = newS
		} else {
			subnets[az] = newS
		}
	} else {
		if network.ID == "" {
			network.ID = subnetID
		} else if network.ID != subnetID {
			return fmt.Errorf("subnet ID %q is not the same as %q", network.ID, subnetID)
		}
		if network.CIDR == nil {
			network.CIDR = subnetCIDR
		} else if network.CIDR.String() != subnetCIDR.String() {
			return fmt.Errorf("subnet CIDR %q is not the same as %q", network.CIDR.String(), subnetCIDR.String())
		}
		network.AZ = az
		subnets[az] = network
	}
	return nil
}

// SubnetInfo returns a string containing VPC subnet information
// Useful for error messages and logs
func (c *ClusterConfig) SubnetInfo() string {
	return fmt.Sprintf("VPC (%s) and subnets (private:%v public:%v)",
		c.VPC.ID, c.VPC.Subnets.Private, c.VPC.Subnets.Public)
}

// HasAnySubnets checks if any subnets were set
func (c *ClusterConfig) HasAnySubnets() bool {
	return c.VPC.Subnets != nil && len(c.VPC.Subnets.Private)+len(c.VPC.Subnets.Public) != 0
}

// HasSufficientPrivateSubnets validates if there is a sufficient
// number of private subnets available to create a cluster
func (c *ClusterConfig) HasSufficientPrivateSubnets() bool {
	return len(c.VPC.Subnets.Private) >= MinRequiredSubnets
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

var errInsufficientSubnets = fmt.Errorf(
	"insufficient number of subnets, at least %dx public and/or %dx private subnets are required",
	MinRequiredSubnets, MinRequiredSubnets)

// HasSufficientSubnets validates if there is a sufficient number
// of either private and/or public subnets available to create
// a cluster, i.e. either non-zero of public or private, and not
// less then MinRequiredSubnets of each, but allowing to have
// public-only or private-only
func (c *ClusterConfig) HasSufficientSubnets() error {
	if !c.HasAnySubnets() {
		return errInsufficientSubnets
	}

	if numPublic := len(c.VPC.Subnets.Public); numPublic > 0 && numPublic < MinRequiredSubnets {
		return errInsufficientSubnets
	}

	if numPrivate := len(c.VPC.Subnets.Private); numPrivate > 0 && numPrivate < MinRequiredSubnets {
		return errInsufficientSubnets
	}

	return nil
}

//DefaultEndpointsMsg returns a message that the EndpointAccess is the same as the default
func (c *ClusterConfig) DefaultEndpointsMsg() string {
	return fmt.Sprintf(
		"Kubernetes API endpoint access will use default of {publicAccess=true, privateAccess=false} for cluster %q in %q", c.Metadata.Name, c.Metadata.Region)
}

//CustomEndpointsMsg returns a message indicating the EndpointAccess given by the user
func (c *ClusterConfig) CustomEndpointsMsg() string {
	return fmt.Sprintf(
		"Kubernetes API endpoint access will use provided values {publicAccess=%v, privateAccess=%v} for cluster %q in %q", *c.VPC.ClusterEndpoints.PublicAccess, *c.VPC.ClusterEndpoints.PrivateAccess, c.Metadata.Name, c.Metadata.Region)
}

//UpdateEndpointsMsg gives message indicating that they need to use eksctl utils to make this config
func (c *ClusterConfig) UpdateEndpointsMsg() string {
	return fmt.Sprintf(
		"you can update Kubernetes API endpoint access with `eksctl utils update-cluster-endpoints --region=%s --name=%s --private-access=bool --public-access=bool`", c.Metadata.Region, c.Metadata.Name)
}

// EndpointsEqual returns true of two endpoints have same values after dereferencing any pointers
func EndpointsEqual(a, b ClusterEndpoints) bool {
	return reflect.DeepEqual(a, b)
}

//HasClusterEndpointAccess determines if endpoint access was configured in config file or not
func (c *ClusterConfig) HasClusterEndpointAccess() bool {
	if c.VPC != nil && c.VPC.ClusterEndpoints != nil {
		pubAccess := c.VPC.ClusterEndpoints.PublicAccess
		privAccess := c.VPC.ClusterEndpoints.PrivateAccess
		hasPublicAccess := pubAccess != nil && *pubAccess
		hasPrivateAccess := privAccess != nil && *privAccess
		return hasPublicAccess || hasPrivateAccess
	}
	return true
}

func (c *ClusterConfig) HasPrivateEndpointAccess() bool {
	return c.VPC != nil &&
		c.VPC.ClusterEndpoints != nil &&
		c.VPC.ClusterEndpoints.PrivateAccess != nil &&
		*c.VPC.ClusterEndpoints.PrivateAccess
}
