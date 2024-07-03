package eks_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awsoutposts "github.com/aws/aws-sdk-go-v2/service/outposts"
	outpoststypes "github.com/aws/aws-sdk-go-v2/service/outposts/types"

	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/outposts"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/nodes"
)

type normalizeEntry struct {
	clusterConfig         *api.ClusterConfig
	expectedInstanceTypes []string

	expectedCallsCount callsCount
	expectedErr        string
}

type callsCount struct {
	getOutpostInstanceTypes int
	describeInstanceTypes   int
}

var _ = Describe("NodeGroupService", func() {
	DescribeTable("Normalize nodegroup", func(ne normalizeEntry) {
		provider := mockprovider.NewMockProvider()
		var outpostsService *outposts.Service
		if ne.clusterConfig.IsControlPlaneOnOutposts() {
			mockOutpostInstanceTypes(provider)
			outpostsService = &outposts.Service{
				EC2API:      provider.EC2(),
				OutpostsAPI: provider.Outposts(),
				OutpostID:   "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
			}
		}
		provider.MockEC2().On("DescribeImages", mock.Anything, &ec2.DescribeImagesInput{
			ImageIds: []string{"ami-test"},
		}).Return(&ec2.DescribeImagesOutput{
			Images: []ec2types.Image{
				{
					ImageId: aws.String("ami-test"),
				},
			},
		}, nil)

		nodeGroupService := eks.NewNodeGroupService(provider, nil, outpostsService)
		nodePools := nodes.ToNodePools(ne.clusterConfig)
		err := nodeGroupService.Normalize(context.Background(), nodePools, ne.clusterConfig)
		provider.MockOutposts().AssertNumberOfCalls(GinkgoT(), "GetOutpostInstanceTypes", ne.expectedCallsCount.getOutpostInstanceTypes)
		provider.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeInstanceTypes", ne.expectedCallsCount.describeInstanceTypes)
		if ne.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(ne.expectedErr)))
			return
		}

		Expect(err).NotTo(HaveOccurred())
		var actualInstanceTypes []string
		for _, ng := range nodePools {
			actualInstanceTypes = append(actualInstanceTypes, ng.BaseNodeGroup().InstanceType)
		}

		Expect(actualInstanceTypes).To(Equal(ne.expectedInstanceTypes))
	},

		Entry("[Outposts] nodeGroup.instanceType should be set to the smallest available instance type", normalizeEntry{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Version: api.DefaultVersion,
				},
				Outpost: &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				},
				NodeGroups: []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							AMI: "ami-test",
							SSH: &api.NodeGroupSSH{
								Allow: api.Disabled(),
							},
						},
					},
					{
						NodeGroupBase: &api.NodeGroupBase{
							AMI:          "ami-test",
							InstanceType: "m5a.16xlarge",
							SSH: &api.NodeGroupSSH{
								Allow: api.Disabled(),
							},
						},
					},
					{
						NodeGroupBase: &api.NodeGroupBase{
							AMI: "ami-test",
							SSH: &api.NodeGroupSSH{
								Allow: api.Disabled(),
							},
						},
					},
				},
			},

			expectedInstanceTypes: []string{"m5a.large", "m5a.16xlarge", "m5a.large"},
			expectedCallsCount: callsCount{
				getOutpostInstanceTypes: 1,
				describeInstanceTypes:   1,
			},
		}),

		Entry("[Outposts] unavailable instance type should return an error", normalizeEntry{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Version: api.DefaultVersion,
				},
				Outpost: &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				},
				NodeGroups: []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							AMI: "ami-test",
							SSH: &api.NodeGroupSSH{
								Allow: api.Disabled(),
							},
							InstanceType: "t2.medium",
						},
					},
				},
			},
			expectedErr: `instance type "t2.medium" does not exist in Outpost`,
			expectedCallsCount: callsCount{
				getOutpostInstanceTypes: 1,
			},
		}),

		Entry("[Outposts] available instance type should not return an error", normalizeEntry{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Version: api.DefaultVersion,
				},
				Outpost: &api.Outpost{
					ControlPlaneOutpostARN: "arn:aws:outposts:us-west-2:1234:outpost/op-1234",
				},
				NodeGroups: []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							AMI: "ami-test",
							SSH: &api.NodeGroupSSH{
								Allow: api.Disabled(),
							},
							InstanceType: "m5a.large",
						},
					},
					{
						NodeGroupBase: &api.NodeGroupBase{
							AMI: "ami-test",
							SSH: &api.NodeGroupSSH{
								Allow: api.Disabled(),
							},
							InstanceType: "m5.xlarge",
						},
					},
				},
			},
			expectedInstanceTypes: []string{"m5a.large", "m5.xlarge"},
			expectedCallsCount: callsCount{
				getOutpostInstanceTypes: 1,
			},
		}),

		Entry("instance type should be set to the default instance type", normalizeEntry{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Version: api.DefaultVersion,
				},
				NodeGroups: []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							AMI: "ami-test",
							SSH: &api.NodeGroupSSH{
								Allow: api.Disabled(),
							},
							InstanceSelector: &api.InstanceSelector{},
						},
					},
				},
				ManagedNodeGroups: []*api.ManagedNodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							AMI: "ami-test",
							SSH: &api.NodeGroupSSH{
								Allow: api.Disabled(),
							},
							InstanceSelector: &api.InstanceSelector{},
						},
					},
				},
			},
			expectedInstanceTypes: []string{api.DefaultNodeType, api.DefaultNodeType},
		}),

		Entry(`instance type should be set to "mixed" when using mixed instance types`, normalizeEntry{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Version: api.DefaultVersion,
				},
				NodeGroups: []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							AMI: "ami-test",
							SSH: &api.NodeGroupSSH{
								Allow: api.Disabled(),
							},
							InstanceSelector: &api.InstanceSelector{},
						},
						InstancesDistribution: &api.NodeGroupInstancesDistribution{
							InstanceTypes: []string{"t2.medium", "t2.large"},
						},
					},
				},
			},
			expectedInstanceTypes: []string{"mixed"},
		}),

		Entry("instance type should be left unset when instanceSelector or launchTemplate is set", normalizeEntry{
			clusterConfig: &api.ClusterConfig{
				Metadata: &api.ClusterMeta{
					Version: api.DefaultVersion,
				},
				ManagedNodeGroups: []*api.ManagedNodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							AMI: "ami-test",
							SSH: &api.NodeGroupSSH{
								Allow: api.Disabled(),
							},
							InstanceSelector: &api.InstanceSelector{
								VCPUs: 2,
							},
						},
					},
					{
						NodeGroupBase: &api.NodeGroupBase{
							AMI: "ami-test",
							SSH: &api.NodeGroupSSH{
								Allow: api.Disabled(),
							},
						},
						LaunchTemplate: &api.LaunchTemplate{
							ID: "lt-123",
						},
					},
				},
			},
			expectedInstanceTypes: []string{"", ""},
		}),
	)
})

func mockOutpostInstanceTypes(provider *mockprovider.MockProvider) {
	instanceTypeInfoList := []ec2types.InstanceTypeInfo{
		{
			InstanceType: "m5a.12xlarge",
			VCpuInfo: &ec2types.VCpuInfo{
				DefaultVCpus:          aws.Int32(48),
				DefaultCores:          aws.Int32(24),
				DefaultThreadsPerCore: aws.Int32(2),
			},
			MemoryInfo: &ec2types.MemoryInfo{
				SizeInMiB: aws.Int64(196608),
			},
		},
		{
			InstanceType: "m5a.large",
			VCpuInfo: &ec2types.VCpuInfo{
				DefaultVCpus:          aws.Int32(2),
				DefaultCores:          aws.Int32(1),
				DefaultThreadsPerCore: aws.Int32(2),
			},
			MemoryInfo: &ec2types.MemoryInfo{
				SizeInMiB: aws.Int64(196608),
			},
		},
		{
			InstanceType: "m5.xlarge",
			VCpuInfo: &ec2types.VCpuInfo{
				DefaultVCpus:          aws.Int32(4),
				DefaultCores:          aws.Int32(2),
				DefaultThreadsPerCore: aws.Int32(2),
			},
			MemoryInfo: &ec2types.MemoryInfo{
				SizeInMiB: aws.Int64(16384),
			},
		},
		{
			InstanceType: "m5a.16xlarge",
			VCpuInfo: &ec2types.VCpuInfo{
				DefaultVCpus:          aws.Int32(64),
				DefaultCores:          aws.Int32(32),
				DefaultThreadsPerCore: aws.Int32(2),
			},
			MemoryInfo: &ec2types.MemoryInfo{
				SizeInMiB: aws.Int64(262144),
			},
		},
	}

	instanceTypeItems := make([]outpoststypes.InstanceTypeItem, len(instanceTypeInfoList))
	instanceTypes := make([]ec2types.InstanceType, len(instanceTypeInfoList))
	for i, it := range instanceTypeInfoList {
		instanceTypeItems[i] = outpoststypes.InstanceTypeItem{
			InstanceType: aws.String(string(it.InstanceType)),
		}
		instanceTypes[i] = it.InstanceType
	}

	provider.MockOutposts().On("GetOutpostInstanceTypes", mock.Anything, mock.Anything, mock.Anything).Return(&awsoutposts.GetOutpostInstanceTypesOutput{
		InstanceTypes: instanceTypeItems,
	}, nil)

	provider.MockEC2().On("DescribeInstanceTypes", mock.Anything, &ec2.DescribeInstanceTypesInput{
		InstanceTypes: instanceTypes,
	}, mock.Anything).Return(&ec2.DescribeInstanceTypesOutput{
		InstanceTypes: instanceTypeInfoList,
	}, nil)
}
