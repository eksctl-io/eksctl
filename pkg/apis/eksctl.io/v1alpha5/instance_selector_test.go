package v1alpha5

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type instanceSelectorCase struct {
	ng     *NodeGroup
	errMsg string
}

var _ = Describe("Instance Selector Validation", func() {
	DescribeTable("Supported and unsupported field combinations", func(n *instanceSelectorCase) {
		SetNodeGroupDefaults(n.ng, &ClusterMeta{Name: "cluster"})
		err := ValidateNodeGroup(0, n.ng)
		if n.errMsg == "" {
			Expect(err).NotTo(HaveOccurred())
			return
		}
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(n.errMsg))

	},
		Entry("valid instanceSelector options", &instanceSelectorCase{
			ng: &NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					InstanceSelector: &InstanceSelector{
						VCPUs:  2,
						Memory: "4",
					},
				},
			},
		}),
		Entry("invalid use of instanceSelector and instance type", &instanceSelectorCase{
			ng: &NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					InstanceSelector: &InstanceSelector{
						VCPUs:  2,
						Memory: "4",
					},
					InstanceType: "m5.large",
				},
			},
			errMsg: `instanceType should be "mixed" or unset`,
		}),
		Entry("invalid use of instanceSelector and instancesDistribution", &instanceSelectorCase{
			ng: &NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					InstanceType:     "m5.large",
					InstanceSelector: &InstanceSelector{},
				},
				InstancesDistribution: &NodeGroupInstancesDistribution{
					InstanceTypes: []string{"m5.large"},
				},
			},
			errMsg: `instanceType should be "mixed" or unset`,
		}),
		Entry("instancesDistribution without instanceTypes and instanceSelector", &instanceSelectorCase{
			ng: &NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					InstanceSelector: &InstanceSelector{},
				},
				InstancesDistribution: &NodeGroupInstancesDistribution{},
			},
			errMsg: `instanceType should be "mixed" or unset`,
		}),
		Entry("valid use of instanceSelector and instancesDistribution", &instanceSelectorCase{
			ng: &NodeGroup{
				NodeGroupBase: &NodeGroupBase{
					InstanceSelector: &InstanceSelector{
						VCPUs:  2,
						Memory: "2",
					},
				},
				InstancesDistribution: &NodeGroupInstancesDistribution{
					InstanceTypes: []string{"m5.large"},
				},
			},
		}),
	)

})
