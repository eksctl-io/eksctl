package automode_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	automodeactions "github.com/weaveworks/eksctl/pkg/actions/automode"
	"github.com/weaveworks/eksctl/pkg/actions/automode/mocks"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
)

type updaterTest struct {
	autoModeConfig *api.AutoModeConfig
	vpc            *api.ClusterVPC
	currentCluster *ekstypes.Cluster
	updateMocks    func(*updaterMocks)
	drainNodes     bool

	expectedErr string
}

type updaterMocks struct {
	roleManager        mocks.RoleManager
	clusterRoleManager mocks.ClusterRoleManager
	drainer            mocks.NodeGroupDrainer
	eksUpdater         mocksv2.EKS
	clientSet          *kubernetesfake.Clientset
}

var _ = DescribeTable("Auto Mode Updater", func(t updaterTest) {
	var um updaterMocks
	um.clientSet = kubernetesfake.NewSimpleClientset()
	clusterConfig := api.NewClusterConfig()
	clusterConfig.Metadata.Name = "cluster"
	if t.updateMocks != nil {
		t.updateMocks(&um)
	}
	clusterConfig.AutoModeConfig = t.autoModeConfig
	clusterConfig.VPC = t.vpc
	updater := &automodeactions.Updater{
		RoleManager:        &um.roleManager,
		ClusterRoleManager: &um.clusterRoleManager,
		PodsGetter:         um.clientSet.CoreV1(),
		EKSUpdater:         &um.eksUpdater,
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
		&um.clusterRoleManager,
		&um.eksUpdater,
		&um.drainer,
	} {
		asserter.AssertExpectations(GinkgoT())
	}

},
	Entry("attempt to update nodeRoleARN", updaterTest{
		autoModeConfig: &api.AutoModeConfig{
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
		expectedErr: "autoModeConfig.nodeRoleARN cannot be modified",
	}),
	Entry("Auto Mode enabled and up-to-date", updaterTest{
		autoModeConfig: &api.AutoModeConfig{
			Enabled:   api.Enabled(),
			NodePools: &[]string{api.AutoModeNodePoolGeneralPurpose, api.AutoModeNodePoolSystem},
		},
		currentCluster: &ekstypes.Cluster{
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled:   aws.Bool(true),
				NodePools: []string{api.AutoModeNodePoolSystem, api.AutoModeNodePoolGeneralPurpose},
			},
		},
	}),
	Entry("enabling Auto Mode with default values", updaterTest{
		autoModeConfig: &api.AutoModeConfig{
			Enabled:   api.Enabled(),
			NodePools: &[]string{api.AutoModeNodePoolGeneralPurpose, api.AutoModeNodePoolSystem},
		},
		currentCluster: &ekstypes.Cluster{
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled: aws.Bool(false),
			},
		},
		updateMocks: func(u *updaterMocks) {
			u.roleManager.EXPECT().CreateOrImport(mock.Anything, "cluster").Return("arn:aws:iam::000:role/NodeRole", nil).Once()
			mockEnableAutoMode(u, "arn:aws:iam::000:role/NodeRole", []string{api.AutoModeNodePoolGeneralPurpose, api.AutoModeNodePoolSystem})
		},
	}),
	Entry("enabling Auto Mode with a custom nodeRoleARN and node pools", updaterTest{
		autoModeConfig: &api.AutoModeConfig{
			Enabled:     aws.Bool(true),
			NodeRoleARN: api.MustParseARN("arn:aws:iam::000:role/CustomNodeRole"),
			NodePools:   &[]string{api.AutoModeNodePoolGeneralPurpose},
		},
		currentCluster: &ekstypes.Cluster{
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled: aws.Bool(false),
			},
		},
		updateMocks: func(u *updaterMocks) {
			mockEnableAutoMode(u, "arn:aws:iam::000:role/CustomNodeRole", []string{api.AutoModeNodePoolGeneralPurpose})
		},
	}),
	Entry("enabling Auto Mode with a pre-existing VPC", updaterTest{
		autoModeConfig: &api.AutoModeConfig{
			Enabled:     aws.Bool(true),
			NodeRoleARN: api.MustParseARN("arn:aws:iam::000:role/CustomNodeRole"),
			NodePools:   &[]string{api.AutoModeNodePoolGeneralPurpose},
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
			u.clusterRoleManager.EXPECT().UpdateRoleForAutoMode(mock.Anything).Return(nil).Once()
			mockUpdateClusterConfig(u, &ekstypes.ComputeConfigRequest{
				Enabled:     aws.Bool(true),
				NodePools:   []string{api.AutoModeNodePoolGeneralPurpose},
				NodeRoleArn: aws.String("arn:aws:iam::000:role/CustomNodeRole"),
			})
		},
	}),
	Entry("disabling Auto Mode", updaterTest{
		autoModeConfig: &api.AutoModeConfig{
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
			u.clusterRoleManager.EXPECT().DeleteAutoModePolicies(mock.Anything).Return(nil).Once()
		},
	}),
	Entry("Karpenter pods exist in the cluster when enabling Auto Mode", updaterTest{
		autoModeConfig: &api.AutoModeConfig{
			Enabled:   api.Enabled(),
			NodePools: &[]string{api.AutoModeNodePoolGeneralPurpose, api.AutoModeNodePoolSystem},
		},
		currentCluster: &ekstypes.Cluster{
			ComputeConfig: &ekstypes.ComputeConfigResponse{
				Enabled: aws.Bool(false),
			},
		},
		drainNodes: true,
		updateMocks: func(u *updaterMocks) {
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
		expectedErr: "enabling Auto Mode: found Karpenter pods in namespace karpenter; either delete Karpenter or scale it down to zero and rerun the command",
	}),
	Entry("drain existing nodes after enabling Auto Mode", updaterTest{
		autoModeConfig: &api.AutoModeConfig{
			Enabled:   api.Enabled(),
			NodePools: &[]string{api.AutoModeNodePoolGeneralPurpose, api.AutoModeNodePoolSystem},
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
			mockEnableAutoMode(u, "arn:aws:iam::000:role/NodeRole", []string{api.AutoModeNodePoolGeneralPurpose, api.AutoModeNodePoolSystem})
		},
	}),
)

func mockEnableAutoMode(u *updaterMocks, nodeRoleARN string, nodePools []string) {
	u.clusterRoleManager.EXPECT().UpdateRoleForAutoMode(mock.Anything).Return(nil).Once()
	mockUpdateClusterConfig(u, &ekstypes.ComputeConfigRequest{
		Enabled:     aws.Bool(true),
		NodePools:   nodePools,
		NodeRoleArn: aws.String(nodeRoleARN),
	})
}

func mockUpdateClusterConfig(u *updaterMocks, computeConfig *ekstypes.ComputeConfigRequest) {
	u.eksUpdater.EXPECT().UpdateClusterConfig(mock.Anything, &awseks.UpdateClusterConfigInput{
		Name:          aws.String("cluster"),
		ComputeConfig: computeConfig,
		KubernetesNetworkConfig: &ekstypes.KubernetesNetworkConfigRequest{
			ElasticLoadBalancing: &ekstypes.ElasticLoadBalancing{
				Enabled: computeConfig.Enabled,
			},
		},
		StorageConfig: &ekstypes.StorageConfigRequest{
			BlockStorage: &ekstypes.BlockStorage{
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
