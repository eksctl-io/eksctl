package eks_test

import (
	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/fakes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type instanceSelectorCase struct {
	nodePools              []api.NodePool
	selector               *api.InstanceSelector
	createInstanceSelector func() eks.InstanceSelector
	instanceTypes          []string
	err                    string
}

var _ = Describe("Instance Selector", func() {

	makeInstanceSelector := func(instanceTypes ...string) func() eks.InstanceSelector {
		return func() eks.InstanceSelector {
			s := &fakes.FakeInstanceSelector{}
			s.FilterReturns(instanceTypes, nil)
			return s
		}
	}

	DescribeTable("Expand instance selector options", func(isc instanceSelectorCase) {
		for _, np := range isc.nodePools {
			np.BaseNodeGroup().InstanceSelector = isc.selector
		}
		nodeGroupService := eks.NewNodeGroupService(nil, isc.createInstanceSelector())
		err := nodeGroupService.ExpandInstanceSelectorOptions(isc.nodePools)
		if isc.err != "" {
			Expect(err.Error()).To(ContainSubstring(isc.err))
			return
		}

		Expect(err).ToNot(HaveOccurred())
		for _, np := range isc.nodePools {
			switch ng := np.(type) {
			case *api.NodeGroup:
				if len(isc.instanceTypes) > 0 {
					Expect(ng.InstancesDistribution.InstanceTypes).To(Equal(isc.instanceTypes))
				} else {
					Expect(ng.InstancesDistribution == nil || len(ng.InstancesDistribution.InstanceTypes) == 0).To(BeTrue())
				}

			case *api.ManagedNodeGroup:
				Expect(ng.InstanceTypes).To(Equal(isc.instanceTypes))

			default:
				panic("unreachable code")
			}
		}
	},

		Entry("valid instance selector criteria", instanceSelectorCase{
			nodePools: []api.NodePool{
				&api.NodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
				},
				&api.ManagedNodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
				},
			},
			selector: &api.InstanceSelector{
				VCPUs:           1,
				CPUArchitecture: "x86_64",
			},
			createInstanceSelector: makeInstanceSelector("t2.medium"),
			instanceTypes:          []string{"t2.medium"},
		}),

		Entry("no instance selector criteria", instanceSelectorCase{
			nodePools: []api.NodePool{
				&api.NodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
				},
				&api.ManagedNodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
				},
			},
			selector:               &api.InstanceSelector{},
			createInstanceSelector: makeInstanceSelector("c5.large", "c4.large"),
		}),

		Entry("no matching instances", instanceSelectorCase{
			nodePools: []api.NodePool{
				&api.NodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
				},
				&api.ManagedNodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
				},
			},
			selector: &api.InstanceSelector{
				VCPUs:  1000,
				Memory: "400GiB",
			},
			createInstanceSelector: makeInstanceSelector(),
			err:                    "instance selector criteria matched no instances",
		}),

		Entry("nodeGroup with instancesDistribution set", instanceSelectorCase{
			nodePools: []api.NodePool{
				&api.NodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
					InstancesDistribution: &api.NodeGroupInstancesDistribution{
						SpotInstancePools: aws.Int(4),
					},
				},
				&api.ManagedNodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
				},
			},
			selector: &api.InstanceSelector{
				VCPUs:  2,
				Memory: "4",
			},
			createInstanceSelector: makeInstanceSelector("m5.large", "m5.xlarge"),
			instanceTypes:          []string{"m5.large", "m5.xlarge"},
		}),

		Entry("mismatching instanceTypes and instance selector criteria for unmanaged nodegroup", instanceSelectorCase{
			nodePools: []api.NodePool{
				&api.NodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
					InstancesDistribution: &api.NodeGroupInstancesDistribution{
						SpotInstancePools: aws.Int(4),
						InstanceTypes:     []string{"c3.large", "c4.large"},
					},
				},
			},
			selector: &api.InstanceSelector{
				VCPUs:  2,
				Memory: "4",
			},
			createInstanceSelector: makeInstanceSelector("c3.large", "c4.xlarge", "c5.large"),
			err:                    "instance types matched by instance selector criteria do not match",
		}),

		Entry("mismatching instanceTypes and instance selector criteria for managed nodegroup", instanceSelectorCase{
			nodePools: []api.NodePool{
				&api.ManagedNodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
					InstanceTypes: []string{"c3.large", "c4.large"},
				},
			},
			selector: &api.InstanceSelector{
				VCPUs:  2,
				Memory: "4",
			},
			createInstanceSelector: makeInstanceSelector("c3.large", "c4.xlarge", "c5.large"),
			err:                    "instance types matched by instance selector criteria do not match",
		}),

		Entry("matching instanceTypes and instance selector criteria", instanceSelectorCase{
			nodePools: []api.NodePool{
				&api.ManagedNodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
					InstanceTypes: []string{"c3.large", "c4.large", "c5.large"},
				},
				&api.NodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
					InstancesDistribution: &api.NodeGroupInstancesDistribution{
						SpotInstancePools: aws.Int(4),
						InstanceTypes:     []string{"c3.large", "c4.large", "c5.large"},
					},
				},
			},
			selector: &api.InstanceSelector{
				VCPUs:  2,
				Memory: "4",
			},
			createInstanceSelector: makeInstanceSelector("c3.large", "c4.large", "c5.large"),
			instanceTypes:          []string{"c3.large", "c4.large", "c5.large"},
		}),
	)
})
