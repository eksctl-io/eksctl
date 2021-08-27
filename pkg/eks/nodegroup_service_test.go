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
	nodeGroups                 []api.NodePool
	instanceSelectorValue      *api.InstanceSelector
	createFakeInstanceSelector func() *fakes.FakeInstanceSelector
	expectedInstanceTypes      []string
	clusterAZs                 []string
	expectedErr                string
	expectedAZs                []string
}

var _ = Describe("Instance Selector", func() {
	makeInstanceSelector := func(instanceTypes ...string) func() *fakes.FakeInstanceSelector {
		return func() *fakes.FakeInstanceSelector {
			s := &fakes.FakeInstanceSelector{}
			s.FilterReturns(instanceTypes, nil)
			return s
		}
	}

	DescribeTable("Expand instance selector options", func(isc instanceSelectorCase) {
		for _, np := range isc.nodeGroups {
			np.BaseNodeGroup().InstanceSelector = isc.instanceSelectorValue
		}
		instanceSelectorFake := isc.createFakeInstanceSelector()
		nodeGroupService := eks.NewNodeGroupService(nil, instanceSelectorFake)
		err := nodeGroupService.ExpandInstanceSelectorOptions(isc.nodeGroups, isc.clusterAZs)
		if isc.expectedErr != "" {
			Expect(err.Error()).To(ContainSubstring(isc.expectedErr))
			return
		}
		Expect(instanceSelectorFake.FilterCallCount()).To(Equal(len(isc.nodeGroups)))

		for i := range isc.nodeGroups {
			Expect(*instanceSelectorFake.FilterArgsForCall(i).AvailabilityZones).To(Equal(isc.expectedAZs))
		}

		Expect(err).ToNot(HaveOccurred())
		for _, np := range isc.nodeGroups {
			switch ng := np.(type) {
			case *api.NodeGroup:
				if len(isc.expectedInstanceTypes) > 0 {
					Expect(ng.InstancesDistribution.InstanceTypes).To(Equal(isc.expectedInstanceTypes))
				} else {
					Expect(ng.InstancesDistribution == nil || len(ng.InstancesDistribution.InstanceTypes) == 0).To(BeTrue())
				}

			case *api.ManagedNodeGroup:
				Expect(ng.InstanceTypes).To(Equal(isc.expectedInstanceTypes))

			default:
				panic("unreachable code")
			}
		}
	},

		Entry("valid instance selector criteria", instanceSelectorCase{
			nodeGroups: []api.NodePool{
				&api.NodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
				},
				&api.ManagedNodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
				},
			},
			instanceSelectorValue: &api.InstanceSelector{
				VCPUs:           1,
				CPUArchitecture: "x86_64",
			},
			createFakeInstanceSelector: makeInstanceSelector("t2.medium"),
			expectedInstanceTypes:      []string{"t2.medium"},
			clusterAZs:                 []string{"az1", "az2"},
			expectedAZs:                []string{"az1", "az2"},
		}),

		Entry("valid instance selector criteria with ng-specific AZs", instanceSelectorCase{
			nodeGroups: []api.NodePool{
				&api.NodeGroup{
					NodeGroupBase: &api.NodeGroupBase{
						AvailabilityZones: []string{"az3", "az4"},
					},
				},
				&api.ManagedNodeGroup{
					NodeGroupBase: &api.NodeGroupBase{
						AvailabilityZones: []string{"az3", "az4"},
					},
				},
			},
			instanceSelectorValue: &api.InstanceSelector{
				VCPUs:           1,
				CPUArchitecture: "x86_64",
			},
			createFakeInstanceSelector: makeInstanceSelector("t2.medium"),
			expectedInstanceTypes:      []string{"t2.medium"},
			clusterAZs:                 []string{"az1", "az2"},
			expectedAZs:                []string{"az3", "az4"},
		}),

		Entry("no matching instances", instanceSelectorCase{
			nodeGroups: []api.NodePool{
				&api.NodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
				},
				&api.ManagedNodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
				},
			},
			instanceSelectorValue: &api.InstanceSelector{
				VCPUs:  1000,
				Memory: "400GiB",
			},
			createFakeInstanceSelector: makeInstanceSelector(),
			expectedErr:                "instance selector criteria matched no instances",
		}),

		Entry("too many matching instances", instanceSelectorCase{
			nodeGroups: []api.NodePool{
				&api.NodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
				},
				&api.ManagedNodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
				},
			},
			instanceSelectorValue: &api.InstanceSelector{
				CPUArchitecture: "arm64",
			},
			createFakeInstanceSelector: makeInstanceSelector(tooManyTypes()...),
			expectedErr:                "instance selector filters resulted in 41 instance types, which is greater than the maximum of 40, please set more selector options",
		}),

		Entry("nodeGroup with instancesDistribution set", instanceSelectorCase{
			nodeGroups: []api.NodePool{
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
			instanceSelectorValue: &api.InstanceSelector{
				VCPUs:  2,
				Memory: "4",
			},
			createFakeInstanceSelector: makeInstanceSelector("m5.large", "m5.xlarge"),
			expectedInstanceTypes:      []string{"m5.large", "m5.xlarge"},
		}),

		Entry("mismatching instanceTypes and instance selector criteria for unmanaged nodegroup", instanceSelectorCase{
			nodeGroups: []api.NodePool{
				&api.NodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
					InstancesDistribution: &api.NodeGroupInstancesDistribution{
						SpotInstancePools: aws.Int(4),
						InstanceTypes:     []string{"c3.large", "c4.large"},
					},
				},
			},
			instanceSelectorValue: &api.InstanceSelector{
				VCPUs:  2,
				Memory: "4",
			},
			createFakeInstanceSelector: makeInstanceSelector("c3.large", "c4.xlarge", "c5.large"),
			expectedErr:                "instance types matched by instance selector criteria do not match",
		}),

		Entry("mismatching instanceTypes and instance selector criteria for managed nodegroup", instanceSelectorCase{
			nodeGroups: []api.NodePool{
				&api.ManagedNodeGroup{
					NodeGroupBase: &api.NodeGroupBase{},
					InstanceTypes: []string{"c3.large", "c4.large"},
				},
			},
			instanceSelectorValue: &api.InstanceSelector{
				VCPUs:  2,
				Memory: "4",
			},
			createFakeInstanceSelector: makeInstanceSelector("c3.large", "c4.xlarge", "c5.large"),
			expectedErr:                "instance types matched by instance selector criteria do not match",
		}),

		Entry("matching instanceTypes and instance selector criteria", instanceSelectorCase{
			nodeGroups: []api.NodePool{
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
			instanceSelectorValue: &api.InstanceSelector{
				VCPUs:  2,
				Memory: "4",
			},
			createFakeInstanceSelector: makeInstanceSelector("c3.large", "c4.large", "c5.large"),
			expectedInstanceTypes:      []string{"c3.large", "c4.large", "c5.large"},
		}),
	)
})

func tooManyTypes() []string {
	instances := make([]string, 41)
	for i := range instances {
		instances[i] = "c3.large"
	}
	return instances
}
