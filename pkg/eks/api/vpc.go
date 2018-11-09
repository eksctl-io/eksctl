package api

import (
	"fmt"
	"net"
)

type (
	// ClusterVPC holds global subnet and all child public/private subnet
	ClusterVPC struct {
		Network              // global CIDR and VPC ID
		SecurityGroup string // cluster SG
		// subnets are either public or private for use with separate nodegroups
		// these are keyed by AZ for convenience
		Subnets map[SubnetTopology]map[string]Network
		// for additional CIDR associations, e.g. to use with separate CIDR for
		// private subnets or any ad-hoc subnets
		ExtraCIDRs []*net.IPNet
	}
	// SubnetTopology can be SubnetTopologyPrivate or SubnetTopologyPublic
	SubnetTopology string
	// Network holds ID and CIDR
	Network struct {
		ID   string
		CIDR *net.IPNet
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

// DefaultCIDR returns default global CIDR for VPC
func DefaultCIDR() net.IPNet {
	return net.IPNet{
		IP:   []byte{192, 168, 0, 0},
		Mask: []byte{255, 255, 0, 0},
	}
}

// SubnetIDs returns list of subnets
func (c *ClusterConfig) SubnetIDs(topology SubnetTopology) []string {
	subnets := []string{}
	if t, ok := c.VPC.Subnets[topology]; ok {
		for _, s := range t {
			subnets = append(subnets, s.ID)
		}
	}
	return subnets
}

// ImportSubnet loads a given subnet into cluster config
func (c *ClusterConfig) ImportSubnet(topology SubnetTopology, az, subnetID, cidr string) {
	_, subnetCIDR, _ := net.ParseCIDR(cidr)

	if c.VPC.Subnets == nil {
		c.VPC.Subnets = make(map[SubnetTopology]map[string]Network)
	}
	if _, ok := c.VPC.Subnets[topology]; !ok {
		c.VPC.Subnets[topology] = map[string]Network{}
	}
	if network, ok := c.VPC.Subnets[topology][az]; !ok {
		c.VPC.Subnets[topology][az] = Network{ID: subnetID, CIDR: subnetCIDR}
	} else {
		network.ID = subnetID
		network.CIDR = subnetCIDR
		c.VPC.Subnets[topology][az] = network
	}
}

// HasSufficientPublicSubnets validates if there is a sufficient
// number of public subnets available to create a cluster
func (c *ClusterConfig) HasSufficientPublicSubnets() bool {
	return len(c.SubnetIDs(SubnetTopologyPublic)) >= MinRequiredSubnets
}

// HasSufficientPrivateSubnets validates if there is a sufficient
// number of private subnets available to create a cluster
func (c *ClusterConfig) HasSufficientPrivateSubnets() bool {
	return len(c.SubnetIDs(SubnetTopologyPrivate)) >= MinRequiredSubnets
}

var errInsufficientSubnets = fmt.Errorf(
	"inssuficient number of subnets, at least %dx public and/or %dx private subnets are required",
	MinRequiredSubnets, MinRequiredSubnets)

// HasSufficientSubnets validates if there is a sufficient number
// of either private and/or public subnets available to create
// a cluster, i.e. either non-zero of public or private, and not
// less then MinRequiredSubnets of each, but allowing to have
// public-only or private-only
func (c *ClusterConfig) HasSufficientSubnets() error {
	numPublic := len(c.SubnetIDs(SubnetTopologyPublic))
	if numPublic > 0 && numPublic < MinRequiredSubnets {
		return errInsufficientSubnets
	}

	numPrivate := len(c.SubnetIDs(SubnetTopologyPrivate))
	if numPrivate > 0 && numPrivate < MinRequiredSubnets {
		return errInsufficientSubnets
	}

	if numPublic == 0 && numPrivate == 0 {
		return errInsufficientSubnets
	}

	return nil
}
