package autonomousmode_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"

	autonomousmodeactions "github.com/weaveworks/eksctl/pkg/actions/autonomousmode"
	"github.com/weaveworks/eksctl/pkg/actions/autonomousmode/mocks"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/autonomousmode"
	autonomousmodemocks "github.com/weaveworks/eksctl/pkg/autonomousmode/mocks"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
)

type updaterTest struct {
	autonomousModeConfig *api.AutonomousModeConfig
	vpc                  *api.ClusterVPC
	currentCluster       *ekstypes.Cluster
	updateMocks          func(*updaterMocks)
	drainNodes           bool

	expectedErr string
}

type updaterMocks struct {
	roleManager      mocks.RoleManager
	drainer          mocks.NodeGroupDrainer
	subnetsLoader    mocks.SubnetsLoader
	nodeClassApplier mocks.NodeClassApplier
	eksUpdater       mocksv2.EKS
	rawClient        autonomousmodemocks.RawClient
	clientSet        *kubernetesfake.Clientset
}

var _ = DescribeTable("Autonomous Mode Updater", func(t updaterTest) {
	var um updaterMocks
	um.clientSet = kubernetesfake.NewSimpleClientset()
	clusterConfig := api.NewClusterConfig()
	clusterConfig.Metadata.Name = "cluster"
	if t.updateMocks != nil {
		t.updateMocks(&um)
	}
	clusterConfig.AutonomousModeConfig = t.autonomousModeConfig
	clusterConfig.VPC = t.vpc
	updater := &autonomousmodeactions.Updater{
		RoleManager:      &um.roleManager,
		CoreV1Interface:  um.clientSet.CoreV1(),
		EKSUpdater:       &um.eksUpdater,
		SubnetsLoader:    &um.subnetsLoader,
		NodeClassApplier: &um.nodeClassApplier,
		RBACApplier: &autonomousmode.RBACApplier{
			RawClient: &um.rawClient,
		},
	}
	if t.drainNodes {
		updater.Drainer = &um.drainer
	}
	err := updater.Update(context.Background(), clusterConfig, t.currentCluster)
	if t.expectedErr != "" {
		Expect(err).To(MatchError(t.expectedErr))
	} else {
		Expect(err).NotTo(HaveOccurred())
	}
	for _, asserter := range []interface {
		AssertExpectations(t mock.TestingT) bool
	}{
		&um.roleManager,
		&um.eksUpdater,
		&um.drainer,
		&um.subnetsLoader,
		&um.nodeClassApplier,
		&um.rawClient,
	} {
		asserter.AssertExpectations(GinkgoT())
	}

},
	Entry("attempt to update nodeRoleARN", updaterTest{
		autonomousModeConfig: &api.AutonomousModeConfig{
			Enabled:     api.Enabled(),
			NodeRoleARN: api.MustParseARN("arn:aws:iam::000:role/CustomNodeRole"),
			NodePools:   &[]string{"general-purpose", "system"},
		},
		currentCluster: &ekstypes.Cluster{
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled:     aws.Bool(true),
				NodeRoleArn: aws.String("arn:aws:iam::000:role/NodeRole"),
			},
		},
		expectedErr: "autonomousModeConfig.nodeRoleARN cannot be modified",
	}),
	Entry("Autonomous Mode enabled and up-to-date", updaterTest{
		autonomousModeConfig: &api.AutonomousModeConfig{
			Enabled:   api.Enabled(),
			NodePools: &[]string{api.AutonomousModeNodePoolGeneralPurpose, api.AutonomousModeNodePoolSystem},
		},
		currentCluster: &ekstypes.Cluster{
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled:   aws.Bool(true),
				NodePools: []string{api.AutonomousModeNodePoolSystem, api.AutonomousModeNodePoolGeneralPurpose},
			},
		},
	}),
	Entry("enabling Autonomous Mode with default values", updaterTest{
		autonomousModeConfig: &api.AutonomousModeConfig{
			Enabled:   api.Enabled(),
			NodePools: &[]string{api.AutonomousModeNodePoolGeneralPurpose, api.AutonomousModeNodePoolSystem},
		},
		currentCluster: &ekstypes.Cluster{
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled: aws.Bool(false),
			},
		},
		updateMocks: func(u *updaterMocks) {
			u.roleManager.EXPECT().CreateOrImport(mock.Anything, "cluster").Return("arn:aws:iam::000:role/NodeRole", nil).Once()
			mockEnableAutonomousMode(u, "arn:aws:iam::000:role/NodeRole", []string{api.AutonomousModeNodePoolGeneralPurpose, api.AutonomousModeNodePoolSystem})
		},
	}),
	Entry("enabling Autonomous Mode with a custom nodeRoleARN and node pools", updaterTest{
		autonomousModeConfig: &api.AutonomousModeConfig{
			Enabled:     aws.Bool(true),
			NodeRoleARN: api.MustParseARN("arn:aws:iam::000:role/CustomNodeRole"),
			NodePools:   &[]string{api.AutonomousModeNodePoolGeneralPurpose},
		},
		currentCluster: &ekstypes.Cluster{
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled: aws.Bool(false),
			},
		},
		updateMocks: func(u *updaterMocks) {
			mockEnableAutonomousMode(u, "arn:aws:iam::000:role/CustomNodeRole", []string{api.AutonomousModeNodePoolGeneralPurpose})
		},
	}),
	Entry("enabling Autonomous Mode with a pre-existing VPC", updaterTest{
		autonomousModeConfig: &api.AutonomousModeConfig{
			Enabled:     aws.Bool(true),
			NodeRoleARN: api.MustParseARN("arn:aws:iam::000:role/CustomNodeRole"),
			NodePools:   &[]string{api.AutonomousModeNodePoolGeneralPurpose},
		},
		vpc: &api.ClusterVPC{
			Network: api.Network{
				ID: "vpc-123",
			},
		},
		currentCluster: &ekstypes.Cluster{
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled: aws.Bool(false),
			},
		},
		updateMocks: func(u *updaterMocks) {
			mockUpdateClusterConfig(u, &ekstypes.ComputeConfigRequest{
				Enabled:     aws.Bool(true),
				NodePools:   []string{api.AutonomousModeNodePoolGeneralPurpose},
				NodeRoleArn: aws.String("arn:aws:iam::000:role/CustomNodeRole"),
			})
			u.rawClient.EXPECT().CreateOrReplace(mock.Anything, false).Return(nil)
		},
	}),
	Entry("disabling Autonomous Mode", updaterTest{
		autonomousModeConfig: &api.AutonomousModeConfig{
			Enabled: aws.Bool(false),
		},
		currentCluster: &ekstypes.Cluster{
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled: aws.Bool(true),
			},
		},
		updateMocks: func(u *updaterMocks) {
			mockUpdateClusterConfig(u, &ekstypes.ComputeConfigRequest{
				Enabled: aws.Bool(false),
			})
			u.roleManager.EXPECT().DeleteIfRequired(mock.Anything).Return(nil).Once()
			u.rawClient.EXPECT().Delete(mock.Anything).Return(nil)
		},
	}),
	Entry("Karpenter pods exist in the cluster when enabling Autonomous Mode", updaterTest{
		autonomousModeConfig: &api.AutonomousModeConfig{
			Enabled:   api.Enabled(),
			NodePools: &[]string{api.AutonomousModeNodePoolGeneralPurpose, api.AutonomousModeNodePoolSystem},
		},
		currentCluster: &ekstypes.Cluster{
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled: aws.Bool(false),
			},
		},
		drainNodes: true,
		updateMocks: func(u *updaterMocks) {
			var uu unstructured.Unstructured
			uu.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   "eks.amazonaws.com",
				Version: "v1",
				Kind:    "NodeClass",
			})
			u.clientSet = kubernetesfake.NewSimpleClientset(&uu)
			_, err := u.clientSet.CoreV1().Pods("karpenter").Create(context.Background(), &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "karpenter-123-123",
					Namespace: "karpenter",
					Labels: map[string]string{
						"app.kubernetes.io/instance": "karpenter",
					},
				},
			}, metav1.CreateOptions{})
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
		},
		expectedErr: "enabling Autonomous Mode: found Karpenter pods in namespace karpenter; either delete Karpenter or scale it down to zero and rerun the command",
	}),
	Entry("drain existing nodes after enabling Autonomous Mode", updaterTest{
		autonomousModeConfig: &api.AutonomousModeConfig{
			Enabled:   api.Enabled(),
			NodePools: &[]string{api.AutonomousModeNodePoolGeneralPurpose, api.AutonomousModeNodePoolSystem},
		},
		currentCluster: &ekstypes.Cluster{
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled: aws.Bool(false),
			},
		},
		drainNodes: true,
		updateMocks: func(u *updaterMocks) {
			u.drainer.EXPECT().Drain(mock.Anything).Return(nil).Once()
			u.roleManager.EXPECT().CreateOrImport(mock.Anything, "cluster").Return("arn:aws:iam::000:role/NodeRole", nil).Once()
			mockEnableAutonomousMode(u, "arn:aws:iam::000:role/NodeRole", []string{api.AutonomousModeNodePoolGeneralPurpose, api.AutonomousModeNodePoolSystem})
		},
	}),
)

func mockEnableAutonomousMode(u *updaterMocks, nodeRoleARN string, nodePools []string) {
	mockUpdateClusterConfig(u, &ekstypes.ComputeConfigRequest{
		Enabled:     aws.Bool(true),
		NodePools:   nodePools,
		NodeRoleArn: aws.String(nodeRoleARN),
	})
	subnetIDs := []string{"subnet-1", "subnet-2"}
	u.subnetsLoader.EXPECT().LoadSubnets(mock.Anything, mock.Anything).Return(subnetIDs, true, nil).Once()
	u.nodeClassApplier.EXPECT().PatchSubnets(mock.Anything, subnetIDs).Return(nil).Once()
	u.rawClient.EXPECT().CreateOrReplace(mock.Anything, false).Return(nil)
}

func mockUpdateClusterConfig(u *updaterMocks, computeConfig *ekstypes.ComputeConfigRequest) {
	u.eksUpdater.EXPECT().UpdateClusterConfig(mock.Anything, &awseks.UpdateClusterConfigInput{
		Name:          aws.String("cluster"),
		ComputeConfig: computeConfig,
		KubernetesNetworkConfig: &ekstypes.KubernetesNetworkConfigRequest{
			ElasticLoadBalancing: &ekstypes.ElasticLoadBalancingRequest{
				Enabled: computeConfig.Enabled,
			},
		},
		StorageConfig: &ekstypes.StorageConfigRequest{
			BlockStorage: &ekstypes.BlockStorageRequest{
				Enabled: computeConfig.Enabled,
			},
		},
	}).Return(&awseks.UpdateClusterConfigOutput{
		Update: &ekstypes.Update{
			Id: aws.String("update-123"),
		},
	}, nil).Once()
	u.eksUpdater.EXPECT().DescribeUpdate(mock.Anything, &awseks.DescribeUpdateInput{
		Name:     aws.String("cluster"),
		UpdateId: aws.String("update-123"),
	}, mock.Anything).Return(&awseks.DescribeUpdateOutput{
		Update: &ekstypes.Update{
			Status: ekstypes.UpdateStatusSuccessful,
		},
	}, nil).Once()
}
