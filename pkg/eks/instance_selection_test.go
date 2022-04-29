package eks_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type instanceTypeCase struct {
	instanceTypes []string
	instanceType  string

	expectedInstanceType string
}

var _ = DescribeTable("Instance type selection", func(t instanceTypeCase) {
	var (
		ng  *api.NodeGroup
		mng *api.ManagedNodeGroup
	)

	if len(t.instanceTypes) > 0 {
		ng = &api.NodeGroup{
			InstancesDistribution: &api.NodeGroupInstancesDistribution{
				InstanceTypes: t.instanceTypes,
			},
		}
		mng = &api.ManagedNodeGroup{
			InstanceTypes: t.instanceTypes,
		}
	} else {
		ngBase := &api.NodeGroupBase{
			InstanceType: t.instanceType,
		}
		ng = &api.NodeGroup{
			NodeGroupBase: ngBase,
		}
		mng = &api.ManagedNodeGroup{
			NodeGroupBase: ngBase,
		}
	}

	for _, np := range []api.NodePool{ng, mng} {
		instanceType := api.SelectInstanceType(np)
		Expect(instanceType).To(Equal(t.expectedInstanceType))
	}
},
	Entry("all ARM instances", instanceTypeCase{
		instanceTypes: []string{"t4g.xlarge", "m6g.xlarge", "r6g.xlarge"},

		expectedInstanceType: "t4g.xlarge",
	}),

	Entry("one GPU instance", instanceTypeCase{
		instanceTypes: []string{"t2.medium", "t4.large", "g3.xlarge"},

		expectedInstanceType: "g3.xlarge",
	}),

	Entry("all GPU instances", instanceTypeCase{
		instanceTypes: []string{"p2.large", "p3.large", "g3.large"},

		expectedInstanceType: "p2.large",
	}),

	Entry("single instance type", instanceTypeCase{
		instanceType: "t4.large",

		expectedInstanceType: "t4.large",
	}),

	Entry("single GPU instance type", instanceTypeCase{
		instanceType: "t4.large",

		expectedInstanceType: "t4.large",
	}),
)
