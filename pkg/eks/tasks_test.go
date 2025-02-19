package eks

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("ClusterTasksForNodeGroups", func() {
	type testCase struct {
		nodeGroups        []*v1alpha5.NodeGroup
		managedNodeGroups []*v1alpha5.ManagedNodeGroup
		expectedTaskCount int
		expectedTaskKind  []string
		installNeuron     bool
		installNvidia     bool
	}

	testCases := []testCase{
		{
			// tasks related to both neuron and nvidia device plugin should be created
			nodeGroups: []*v1alpha5.NodeGroup{
				{
					NodeGroupBase: &v1alpha5.NodeGroupBase{
						InstanceType: "p2.xlarge",
						AMIFamily:    v1alpha5.NodeImageFamilyAmazonLinux2,
					},
				},
			},
			managedNodeGroups: []*v1alpha5.ManagedNodeGroup{
				{
					NodeGroupBase: &v1alpha5.NodeGroupBase{
						InstanceType: "inf1.xlarge",
						AMIFamily:    v1alpha5.NodeImageFamilyAmazonLinux2,
					},
				},
			},
			expectedTaskCount: 2,
			expectedTaskKind:  []string{"Neuron", "Nvidia"},
			installNeuron:     true,
			installNvidia:     true,
		},
		{
			// tasks related to both neuron and nvidia device plugin should NOT be created
			nodeGroups: []*v1alpha5.NodeGroup{
				{
					NodeGroupBase: &v1alpha5.NodeGroupBase{
						InstanceType: "p2.xlarge",
						AMIFamily:    v1alpha5.NodeImageFamilyAmazonLinux2,
					},
				},
			},
			managedNodeGroups: []*v1alpha5.ManagedNodeGroup{
				{
					NodeGroupBase: &v1alpha5.NodeGroupBase{
						InstanceType: "inf1.xlarge",
						AMIFamily:    v1alpha5.NodeImageFamilyAmazonLinux2,
					},
				},
			},
			expectedTaskCount: 0,
			expectedTaskKind:  []string{},
			installNeuron:     false,
			installNvidia:     false,
		},
		{
			// tasks related to either neuron and nvidia device plugin should NOT be created as the
			// instance type does not have supporte either GPU type
			nodeGroups: []*v1alpha5.NodeGroup{
				{
					NodeGroupBase: &v1alpha5.NodeGroupBase{
						InstanceType: "r7a.large",
						AMIFamily:    v1alpha5.NodeImageFamilyAmazonLinux2,
					},
				},
			},
			managedNodeGroups: []*v1alpha5.ManagedNodeGroup{
				{
					NodeGroupBase: &v1alpha5.NodeGroupBase{
						InstanceType: "r7a.large",
						AMIFamily:    v1alpha5.NodeImageFamilyAmazonLinux2023,
					},
				},
			},
			expectedTaskCount: 0,
			expectedTaskKind:  []string{},
			installNeuron:     true,
			installNvidia:     true,
		},
		{
			// tasks related to either neuron and nvidia device plugin should NOT be created as the
			// instance type does not have supporte either GPU type
			nodeGroups: []*v1alpha5.NodeGroup{
				{
					NodeGroupBase: &v1alpha5.NodeGroupBase{
						InstanceType: "g6e.48xlarge",
						AMIFamily:    v1alpha5.NodeImageFamilyAmazonLinux2,
					},
				},
			},
			managedNodeGroups: []*v1alpha5.ManagedNodeGroup{
				{
					NodeGroupBase: &v1alpha5.NodeGroupBase{
						InstanceType: "g6e.48xlarge",
						AMIFamily:    v1alpha5.NodeImageFamilyAmazonLinux2023,
					},
				},
			},
			expectedTaskCount: 1,
			expectedTaskKind:  []string{"Nvidia"},
			installNeuron:     true,
			installNvidia:     true,
		},
		{
			// tasks related to both neuron and nvidia device plugin should be created
			nodeGroups: []*v1alpha5.NodeGroup{
				{
					NodeGroupBase: &v1alpha5.NodeGroupBase{
						InstanceType: "p2.xlarge",
						AMIFamily:    v1alpha5.NodeImageFamilyAmazonLinux2023,
					},
				},
			},
			managedNodeGroups: []*v1alpha5.ManagedNodeGroup{
				{
					NodeGroupBase: &v1alpha5.NodeGroupBase{
						InstanceType: "inf1.xlarge",
						AMIFamily:    v1alpha5.NodeImageFamilyAmazonLinux2023,
					},
				},
			},
			expectedTaskCount: 2,
			expectedTaskKind:  []string{"Neuron", "Nvidia"},
			installNeuron:     true,
			installNvidia:     true,
		},
		{
			// tasks related to both neuron and nvidia device plugin should NOT be created
			nodeGroups: []*v1alpha5.NodeGroup{
				{
					NodeGroupBase: &v1alpha5.NodeGroupBase{
						InstanceType: "p2.xlarge",
						AMIFamily:    v1alpha5.NodeImageFamilyAmazonLinux2023,
					},
				},
			},
			managedNodeGroups: []*v1alpha5.ManagedNodeGroup{
				{
					NodeGroupBase: &v1alpha5.NodeGroupBase{
						InstanceType: "inf1.xlarge",
						AMIFamily:    v1alpha5.NodeImageFamilyAmazonLinux2023,
					},
				},
			},
			expectedTaskCount: 0,
			expectedTaskKind:  []string{},
			installNeuron:     false,
			installNvidia:     false,
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		It("returns the expected tasks", func() {
			provider := mockprovider.NewMockProvider()
			cfg := &v1alpha5.ClusterConfig{
				NodeGroups:        tc.nodeGroups,
				ManagedNodeGroups: tc.managedNodeGroups,
			}
			clusterProvider := &ClusterProvider{
				AWSProvider: provider,
			}
			clusterTasks := clusterProvider.ClusterTasksForNodeGroups(cfg, tc.installNeuron, tc.installNvidia)
			Expect(clusterTasks).NotTo(BeNil())
			Expect(clusterTasks.Tasks).To(HaveLen(tc.expectedTaskCount))
			for i, task := range clusterTasks.Tasks {
				devicePluginTask, _ := task.(*devicePluginTask)
				Expect(devicePluginTask.kind).To(Equal(tc.expectedTaskKind[i]))
			}
		})
	}
})
