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
		instanceTypes: []string{"t2.medium", "t4.large", "g4dn.xlarge"},

		expectedInstanceType: "g4dn.xlarge",
	}),

	Entry("all GPU instances", instanceTypeCase{
		instanceTypes: []string{"g5.8xlarge", "p3.8xlarge", "g4dn.xlarge"},

		expectedInstanceType: "g5.8xlarge",
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
