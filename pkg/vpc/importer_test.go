package vpc

import (
	gfnt "goformation/v4/cloudformation/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
)

var _ = Describe("SpecConfigImporter", func() {
	clusterSecurityGroup := "sg-222222222"

	existingVpc := api.ClusterVPC{
		Network: api.Network{
			ID:       "vpc-111111111",
			CIDR:     &ipnet.IPNet{},
			IPv6Cidr: "",
			IPv6Pool: "",
		},
		SecurityGroup: "sg-111111111",
		Subnets: &api.ClusterSubnets{
			Private: map[string]api.AZSubnetSpec{
				"us-west-2a": {
					ID:   "subnet-1111111",
					AZ:   "us-west-2a",
					CIDR: &ipnet.IPNet{},
				},
				"us-west-2b": {
					ID:   "subnet-2222222",
					AZ:   "us-west-2b",
					CIDR: &ipnet.IPNet{},
				},
			},
			Public: map[string]api.AZSubnetSpec{
				"us-west-2a": {
					ID:   "subnet-3333333",
					AZ:   "us-west-2a",
					CIDR: &ipnet.IPNet{},
				},
				"us-west-2b": {
					ID:   "subnet-4444444",
					AZ:   "us-west-2b",
					CIDR: &ipnet.IPNet{},
				}},
		},
		ExtraCIDRs:                         []string{},
		ExtraIPv6CIDRs:                     []string{},
		SharedNodeSecurityGroup:            "sg-333333333",
		ManageSharedNodeSecurityGroupRules: new(bool),
		AutoAllocateIPv6:                   new(bool),
		NAT:                                &api.ClusterNAT{},
		ClusterEndpoints:                   &api.ClusterEndpoints{},
		PublicAccessCIDRs:                  []string{},
	}

	importer := NewSpecConfigImporter(clusterSecurityGroup, &existingVpc)

	Describe("VPC", func() {
		It("returns the gfnt value of the cluster config VPC ID", func() {
			Expect(importer.VPC()).To(Equal(gfnt.NewString(existingVpc.ID)))
		})
	})

	Describe("ClusterSecurityGroup", func() {
		It("returns the gfnt value of the default cluser security group", func() {
			Expect(importer.ClusterSecurityGroup()).To(Equal(gfnt.NewString(clusterSecurityGroup)))
		})
	})

	Describe("ControlPlaneSecurityGroup", func() {
		It("returns the gfnt value of the cluster config VPC securityGroup", func() {
			Expect(importer.ControlPlaneSecurityGroup()).To(Equal(gfnt.NewString(existingVpc.SecurityGroup)))
		})
	})

	Describe("SharedNodeSecurityGroup", func() {
		It("returns the gfnt value of the cluster config VPC sharedNodeSecurityGroup", func() {
			Expect(importer.SharedNodeSecurityGroup()).To(Equal(gfnt.NewString(existingVpc.SharedNodeSecurityGroup)))
		})

		Context("when the VPC has no shared node security group", func() {
			noSharedSgImporter := NewSpecConfigImporter(clusterSecurityGroup, &api.ClusterVPC{})

			It("returns the default cluster security group", func() {
				Expect(noSharedSgImporter.SharedNodeSecurityGroup()).To(Equal(gfnt.NewString(clusterSecurityGroup)))
			})
		})
	})

	Describe("SecurityGroups", func() {
		It("returns a gfnt slice of the ClusterSecurityGroup", func() {
			Expect(importer.SecurityGroups()).To(HaveLen(1))
			Expect(importer.SecurityGroups()).To(ContainElement(gfnt.NewString(clusterSecurityGroup)))
		})
	})

	Describe("SubnetsPublic", func() {
		It("returns a gfnt string slice of the Public subnets from the cluster config VPC subnets spec", func() {
			Expect(importer.SubnetsPublic().Raw()).To(HaveLen(2))
			Expect(importer.SubnetsPublic().Raw()).To(ContainElement(gfnt.NewString("subnet-3333333")))
			Expect(importer.SubnetsPublic().Raw()).To(ContainElement(gfnt.NewString("subnet-4444444")))
		})
	})

	Describe("SubnetsPrivate", func() {
		It("returns a gfnt string slice of the Private subnets from the cluster config VPC subnets spec", func() {
			Expect(importer.SubnetsPrivate().Raw()).To(HaveLen(2))
			Expect(importer.SubnetsPrivate().Raw()).To(ContainElement(gfnt.NewString("subnet-1111111")))
			Expect(importer.SubnetsPrivate().Raw()).To(ContainElement(gfnt.NewString("subnet-2222222")))
		})
	})
})
