package v1alpha5_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("ClusterHasInstanceType", func() {
	isARM := func(instanceType string) bool {
		return instanceType == "m6g.large"
	}

	DescribeTable("checks every nodegroup", func(cfg *api.ClusterConfig, expected bool) {
		Expect(api.ClusterHasInstanceType(cfg, isARM)).To(Equal(expected))
	},
		Entry("no nodegroups", &api.ClusterConfig{}, false),
		Entry("no matching instance type", &api.ClusterConfig{
			NodeGroups: []*api.NodeGroup{{
				NodeGroupBase: &api.NodeGroupBase{InstanceType: "m5.large"},
			}},
			ManagedNodeGroups: []*api.ManagedNodeGroup{{
				NodeGroupBase: &api.NodeGroupBase{},
				InstanceTypes: []string{"m5.large", "m5a.large"},
			}},
		}, false),
		Entry("instanceType of an unmanaged nodegroup", &api.ClusterConfig{
			NodeGroups: []*api.NodeGroup{{
				NodeGroupBase: &api.NodeGroupBase{InstanceType: "m6g.large"},
			}},
		}, true),
		Entry("instancesDistribution of an unmanaged nodegroup", &api.ClusterConfig{
			NodeGroups: []*api.NodeGroup{{
				NodeGroupBase: &api.NodeGroupBase{InstanceType: "mixed"},
				InstancesDistribution: &api.NodeGroupInstancesDistribution{
					InstanceTypes: []string{"m5.large", "m6g.large"},
				},
			}},
		}, true),
		Entry("instanceType of a managed nodegroup", &api.ClusterConfig{
			ManagedNodeGroups: []*api.ManagedNodeGroup{{
				NodeGroupBase: &api.NodeGroupBase{InstanceType: "m6g.large"},
			}},
		}, true),
		Entry("instanceTypes of a managed nodegroup", &api.ClusterConfig{
			ManagedNodeGroups: []*api.ManagedNodeGroup{{
				NodeGroupBase: &api.NodeGroupBase{},
				InstanceTypes: []string{"m5.large", "m6g.large"},
			}},
		}, true),
	)
})
