package v1alpha5

import (
	"fmt"
	"net"

	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
)

type (
	// ClusterVPC holds global subnet and all child public/private subnet
	ClusterVPC struct {
		// +optional
		Network `json:",inline"` // global CIDR and VPC ID
		// +optional
		SecurityGroup string `json:"securityGroup,omitempty"` // cluster SG
		// subnets are either public or private for use with separate nodegroups
		// these are keyed by AZ for convenience
		// +optional
		Subnets *ClusterSubnets `json:"subnets,omitempty"`
		// for additional CIDR associations, e.g. to use with separate CIDR for
		// private subnets or any ad-hoc subnets
		// +optional
		ExtraCIDRs []*ipnet.IPNet `json:"extraCIDRs,omitempty"`
		// for pre-defined shared node SG
		SharedNodeSecurityGroup string `json:"sharedNodeSecurityGroup,omitempty"`
		// +optional
		AutoAllocateIPv6 *bool `json:"autoAllocateIPv6,omitempty"`
		// +optional
		NAT *ClusterNAT `json:"nat,omitempty"`
		// +optional
		ClusterEndpoints *ClusterEndpoints `json:"clusterEndpoints,omitempty"`
	}
	// ClusterSubnets holds private and public subnets
	ClusterSubnets struct {
		Private map[string]Network `json:"private,omitempty"`
		Public  map[string]Network `json:"public,omitempty"`
	}
	// SubnetTopology can be SubnetTopologyPrivate or SubnetTopologyPublic
	SubnetTopology string
	// Network holds ID and CIDR
	Network struct {
		// +optional
		ID string `json:"id,omitempty"`
		// +optional
		CIDR *ipnet.IPNet `json:"cidr,omitempty"`
	}
	// ClusterNAT holds NAT gateway configuration options
	ClusterNAT struct {
		Gateway *string `json:"gateway,omitempty"`
	}

	// ClusterEndpoints holds cluster api server endpoint access information
	ClusterEndpoints struct {
		PrivateAccess *bool `json:"privateAccess,omitempty,false"`
		PublicAccess  *bool `json:"publicAccess,omitempty,true"`
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

var (
	// True holds true value and can have it's address taken
	True = true
	// False holds false value and can have it's address taken
	False = false
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

// PrivateSubnetIDs returns list of subnets
func (c *ClusterConfig) PrivateSubnetIDs() []string {
	subnets := []string{}
	if c.VPC.Subnets != nil {
		for _, s := range c.VPC.Subnets.Private {
			subnets = append(subnets, s.ID)
		}
	}
	return subnets
}

// PublicSubnetIDs returns list of subnets
func (c *ClusterConfig) PublicSubnetIDs() []string {
	subnets := []string{}
	if c.VPC.Subnets != nil {
		for _, s := range c.VPC.Subnets.Public {
			subnets = append(subnets, s.ID)
		}
	}
	return subnets
}

// ImportSubnet loads a given subnet into cluster config
func (c *ClusterConfig) ImportSubnet(topology SubnetTopology, az, subnetID, cidr string) error {
	if c.VPC.Subnets == nil {
		c.VPC.Subnets = &ClusterSubnets{}
	}

	switch topology {
	case SubnetTopologyPrivate:
		if c.VPC.Subnets.Private == nil {
			c.VPC.Subnets.Private = make(map[string]Network)
		}
		return doImportSubnet(c.VPC.Subnets.Private, az, subnetID, cidr)
	case SubnetTopologyPublic:
		if c.VPC.Subnets.Public == nil {
			c.VPC.Subnets.Public = make(map[string]Network)
		}
		return doImportSubnet(c.VPC.Subnets.Public, az, subnetID, cidr)
	default:
		return fmt.Errorf("unexpected subnet topology: %s", topology)
	}
}

func doImportSubnet(subnets map[string]Network, az, subnetID, cidr string) error {
	subnetCIDR, _ := ipnet.ParseCIDR(cidr)

	if network, ok := subnets[az]; !ok {
		subnets[az] = Network{ID: subnetID, CIDR: subnetCIDR}
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
		subnets[az] = network
	}
	return nil
}

// HasAnySubnets checks if any subnets were set
func (c *ClusterConfig) HasAnySubnets() bool {
	return c.VPC.Subnets != nil && len(c.VPC.Subnets.Private)+len(c.VPC.Subnets.Public) != 0
}

// HasSufficientPrivateSubnets validates if there is a sufficient
// number of private subnets available to create a cluster
func (c *ClusterConfig) HasSufficientPrivateSubnets() bool {
	return len(c.PrivateSubnetIDs()) >= MinRequiredSubnets
}

// HasSufficientPublicSubnets validates if there is a sufficient
// number of public subnets available to create a cluster
func (c *ClusterConfig) HasSufficientPublicSubnets() bool {
	return len(c.PublicSubnetIDs()) >= MinRequiredSubnets
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
	numPublic := len(c.PublicSubnetIDs())
	if numPublic > 0 && numPublic < MinRequiredSubnets {
		return errInsufficientSubnets
	}

	numPrivate := len(c.PrivateSubnetIDs())
	if numPrivate > 0 && numPrivate < MinRequiredSubnets {
		return errInsufficientSubnets
	}

	if numPublic == 0 && numPrivate == 0 {
		return errInsufficientSubnets
	}

	return nil
}

//HasClusterEndpointAccess determines if endpoint access was configured in config file or not
func (c *ClusterConfig) HasClusterEndpointAccess() bool {
	hasAccess := false
	if !(c.VPC == nil || c.VPC.ClusterEndpoints == nil) {
		hasPublicAccess := &c.VPC.ClusterEndpoints.PublicAccess != nil && *c.VPC.ClusterEndpoints.PublicAccess
		hasPrivateAccess := &c.VPC.ClusterEndpoints.PrivateAccess != nil && *c.VPC.ClusterEndpoints.PrivateAccess

		hasAccess = hasPublicAccess || hasPrivateAccess
	} else {
		// an empty VPC or ClusterEndpoints config defaults to public true/private false
		hasAccess = true
	}
	return hasAccess
}
