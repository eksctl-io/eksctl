package builder

import (
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	vpcfakes "github.com/weaveworks/eksctl/pkg/vpc/fakes"
	"github.com/weaveworks/goformation/v4"
	gfneks "github.com/weaveworks/goformation/v4/cloudformation/eks"
)

type amiTypeEntry struct {
	nodeGroup *api.ManagedNodeGroup

	expectedAMIType string
}

var _ = DescribeTable("Managed Nodegroup AMI type", func(e amiTypeEntry) {
	clusterConfig := api.NewClusterConfig()
	clusterConfig.Status = &api.ClusterStatus{
		Endpoint: "https://test.com",
	}
	api.SetManagedNodeGroupDefaults(e.nodeGroup, clusterConfig.Metadata)
	p := mockprovider.NewMockProvider()
	fakeVPCImporter := new(vpcfakes.FakeImporter)
	bootstrapper := nodebootstrap.NewManagedBootstrapper(clusterConfig, e.nodeGroup)
	stack := NewManagedNodeGroup(p.EC2(), clusterConfig, e.nodeGroup, nil, bootstrapper, false, fakeVPCImporter)

	Expect(stack.AddAllResources()).To(Succeed())
	bytes, err := stack.RenderJSON()
	Expect(err).NotTo(HaveOccurred())

	template, err := goformation.ParseJSON(bytes)
	Expect(err).NotTo(HaveOccurred())
	ngResource, ok := template.Resources["ManagedNodeGroup"]
	Expect(ok).To(BeTrue())
	ng, ok := ngResource.(*gfneks.Nodegroup)
	Expect(ok).To(BeTrue())
	Expect(ng.AmiType.String()).To(Equal(e.expectedAMIType))
},
	Entry("default AMI type", amiTypeEntry{
		nodeGroup: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name: "test",
			},
		},
		expectedAMIType: "AL2_x86_64",
	}),

	Entry("AL2 AMI type", amiTypeEntry{
		nodeGroup: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:      "test",
				AMIFamily: api.NodeImageFamilyAmazonLinux2,
			},
		},
		expectedAMIType: "AL2_x86_64",
	}),

	Entry("AMI type", amiTypeEntry{
		nodeGroup: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name: "test",
			},
		},
		expectedAMIType: "AL2_x86_64",
	}),

	Entry("default GPU instance type", amiTypeEntry{
		nodeGroup: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:         "test",
				InstanceType: "p2.xlarge",
			},
		},
		expectedAMIType: "AL2_x86_64_GPU",
	}),

	Entry("AL2 GPU instance type", amiTypeEntry{
		nodeGroup: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:         "test",
				AMIFamily:    api.NodeImageFamilyAmazonLinux2,
				InstanceType: "p2.xlarge",
			},
		},
		expectedAMIType: "AL2_x86_64_GPU",
	}),

	Entry("AL2 ARM instance type", amiTypeEntry{
		nodeGroup: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:         "test",
				AMIFamily:    api.NodeImageFamilyAmazonLinux2,
				InstanceType: "a1.2xlarge",
			},
		},
		expectedAMIType: "AL2_ARM_64",
	}),

	Entry("Bottlerocket AMI type", amiTypeEntry{
		nodeGroup: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:      "test",
				AMIFamily: api.NodeImageFamilyBottlerocket,
			},
		},
		expectedAMIType: "BOTTLEROCKET_x86_64",
	}),

	Entry("Bottlerocket on ARM", amiTypeEntry{
		nodeGroup: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:         "test",
				AMIFamily:    api.NodeImageFamilyBottlerocket,
				InstanceType: "a1.2xlarge",
			},
		},
		expectedAMIType: "BOTTLEROCKET_ARM_64",
	}),

	Entry("Bottlerocket on ARM", amiTypeEntry{
		nodeGroup: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:         "test",
				AMIFamily:    api.NodeImageFamilyBottlerocket,
				InstanceType: "a1.2xlarge",
			},
		},
		expectedAMIType: "BOTTLEROCKET_ARM_64",
	}),

	Entry("non-native Ubuntu", amiTypeEntry{
		nodeGroup: &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name:      "test",
				AMIFamily: api.NodeImageFamilyUbuntu2004,
			},
		},
		expectedAMIType: "CUSTOM",
	}),
)
