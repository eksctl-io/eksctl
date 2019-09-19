package v1alpha5

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type endpointAccessCase struct {
	vpc       *ClusterVPC
	endpoints *ClusterEndpoints
	Public    *bool
	Private   *bool
	Expected  bool
}

var _ = Describe("VPC Configuration", func() {
	DescribeTable("Can determine if VPC config in config file has cluster endpoints",
		func(e endpointAccessCase) {
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
				if e.Expected {
					Expect(cc.HasClusterEndpointAccess()).Should(BeTrue())
				}
				if e.Expected {
					Expect(cc.HasClusterEndpointAccess()).Should(BeTrue())
				}
				if e.Expected {
					Expect(cc.HasClusterEndpointAccess()).Should(BeTrue())
				}
				if !e.Expected {
					Expect(cc.HasClusterEndpointAccess()).Should(BeFalse())
				}
			}
		},
		Entry("No VPC", endpointAccessCase{
			vpc:       nil,
			endpoints: nil,
			Public:    nil,
			Private:   nil,
			Expected:  true,
		}),
		Entry("Has Empty VPC", endpointAccessCase{
			vpc:       &ClusterVPC{},
			endpoints: nil,
			Public:    nil,
			Private:   nil,
			Expected:  true,
		}),
		Entry("Public=true, Private=true", endpointAccessCase{
			vpc:       &ClusterVPC{},
			endpoints: &ClusterEndpoints{},
			Public:    Enabled(),
			Private:   Enabled(),
			Expected:  true,
		}),
		Entry("Public=true, Private=false", endpointAccessCase{
			vpc:       &ClusterVPC{},
			endpoints: &ClusterEndpoints{},
			Public:    Enabled(),
			Private:   Disabled(),
			Expected:  true,
		}),
		Entry("Public=false, Private=true", endpointAccessCase{
			vpc:       &ClusterVPC{},
			endpoints: &ClusterEndpoints{},
			Public:    Disabled(),
			Private:   Enabled(),
			Expected:  true,
		}),
		Entry("Public=false, Private=false", endpointAccessCase{
			vpc:       &ClusterVPC{},
			endpoints: &ClusterEndpoints{},
			Public:    Disabled(),
			Private:   Disabled(),
			Expected:  false,
		}),
	)
})
