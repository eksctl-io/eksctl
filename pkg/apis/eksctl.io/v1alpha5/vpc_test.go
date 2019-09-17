package v1alpha5

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type EndpointAccessCases struct {
	vpc       *ClusterVPC
	endpoints *ClusterEndpoints
	Public    *bool
	Private   *bool
}

var _ = Describe("VPC Configuration", func() {
	DescribeTable("Can determine if VPC config in config file has cluster endpoints",
		func(e EndpointAccessCases) {
			cc := &ClusterConfig{}
			Expect(cc.HasClusterEndpointAccess()).Should(BeTrue())
			cc.VPC = e.vpc
			Expect(cc.HasClusterEndpointAccess()).Should(BeTrue())
			if cc.VPC != nil {
				cc.VPC.ClusterEndpoints = e.endpoints
			}
			if e.Public != nil && cc.VPC.ClusterEndpoints != nil {
				cc.VPC.ClusterEndpoints.PublicAccess = e.Public
			}
			if e.Private != nil && cc.VPC.ClusterEndpoints != nil {
				cc.VPC.ClusterEndpoints.PrivateAccess = e.Private
			}
			if e.Public != nil && e.Private != nil {
				if *e.Public == true || *e.Private == true {
					Expect(cc.HasClusterEndpointAccess()).Should(BeTrue())
				}
				if *e.Public == true && *e.Private == false {
					Expect(cc.HasClusterEndpointAccess()).Should(BeTrue())
				}
				if *e.Public == false && *e.Private == true {
					Expect(cc.HasClusterEndpointAccess()).Should(BeTrue())
				}
				if *e.Public == false && *e.Private == false {
					Expect(cc.HasClusterEndpointAccess()).Should(BeFalse())
				}
			}
		},
		Entry("No VPC", EndpointAccessCases{
			vpc:       nil,
			endpoints: nil,
			Public:    nil,
			Private:   nil,
		}),
		Entry("Has Empty VPC", EndpointAccessCases{
			vpc:       &ClusterVPC{},
			endpoints: nil,
			Public:    nil,
			Private:   nil,
		}),
		Entry("Public=True, Private=true", EndpointAccessCases{
			vpc:       &ClusterVPC{},
			endpoints: &ClusterEndpoints{},
			Public:    &True,
			Private:   &True,
		}),
		Entry("Public=true, Private=false", EndpointAccessCases{
			vpc:       &ClusterVPC{},
			endpoints: &ClusterEndpoints{},
			Public:    &True,
			Private:   &False,
		}),
		Entry("Public=false, Private=true", EndpointAccessCases{
			vpc:       &ClusterVPC{},
			endpoints: &ClusterEndpoints{},
			Public:    nil,
			Private:   nil,
		}),
		Entry("Public=false, Private=false", EndpointAccessCases{
			vpc:       &ClusterVPC{},
			endpoints: &ClusterEndpoints{},
			Public:    &False,
			Private:   &False,
		}),
	)
})
