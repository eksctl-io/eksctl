package capability_test

import (
	"context"
	"strings"
	"testing"

	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/capability"
	"github.com/weaveworks/eksctl/pkg/actions/capability/mocks"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
)

func TestCreator_Create(t *testing.T) {
	mockStackCreator := mocks.NewStackCreator(t)
	mockEKSClient := mocksv2.NewEKS(t)

	mockEKSClient.EXPECT().DescribeCluster(context.Background(), &awseks.DescribeClusterInput{
		Name: &[]string{"test-cluster"}[0],
	}).Return(&awseks.DescribeClusterOutput{
		Cluster: &ekstypes.Cluster{
			Status: ekstypes.ClusterStatusActive,
		},
	}, nil)

	mockEKSClient.EXPECT().CreateCapability(mock.Anything, mock.Anything).Return(&awseks.CreateCapabilityOutput{}, nil)
	mockEKSClient.EXPECT().DescribeCapability(mock.Anything, mock.Anything).Return(&awseks.DescribeCapabilityOutput{
		Capability: &ekstypes.Capability{
			Status: ekstypes.CapabilityStatusActive,
		},
	}, nil)

	creator := capability.NewCreator("test-cluster", mockStackCreator, mockEKSClient, nil)

	capabilities := []api.Capability{
		{
			Name:    "test-capability",
			Type:    "ACK",
			RoleARN: "arn:aws:iam::123456789012:role/test-role",
		},
	}

	err := creator.Create(context.Background(), capabilities)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreator_Create_ClusterNotReady(t *testing.T) {
	mockStackCreator := mocks.NewStackCreator(t)
	mockEKSClient := mocksv2.NewEKS(t)

	mockEKSClient.EXPECT().DescribeCluster(context.Background(), &awseks.DescribeClusterInput{
		Name: &[]string{"test-cluster"}[0],
	}).Return(&awseks.DescribeClusterOutput{
		Cluster: &ekstypes.Cluster{
			Status: ekstypes.ClusterStatusCreating,
		},
	}, nil)

	creator := capability.NewCreator("test-cluster", mockStackCreator, mockEKSClient, nil)

	capabilities := []api.Capability{
		{
			Name:    "test-capability",
			Type:    "ACK",
			RoleARN: "arn:aws:iam::123456789012:role/test-role",
		},
	}

	err := creator.Create(context.Background(), capabilities)
	if err == nil {
		t.Error("Expected error when cluster is not ready")
	}
	if !strings.Contains(err.Error(), "cluster not ready") {
		t.Errorf("Expected cluster not ready error, got %v", err)
	}
}

func TestCreator_CreateTasks(t *testing.T) {
	mockStackCreator := mocks.NewStackCreator(t)
	mockEKSClient := mocksv2.NewEKS(t)
	creator := capability.NewCreator("test-cluster", mockStackCreator, mockEKSClient, nil)

	capabilities := []api.Capability{
		{Name: "cap1", Type: "ACK"},
		{Name: "cap2", Type: "KRO"},
	}

	taskTree := creator.CreateTasks(context.Background(), capabilities)

	if taskTree.Len() != 2 {
		t.Errorf("Expected 2 tasks, got %d", taskTree.Len())
	}
}

func TestCreator_MakeStackName(t *testing.T) {
	cap := api.Capability{
		Name: "test-capability",
		Type: "ACK",
	}

	stackName := capability.MakeIAMRoleStackName("test-cluster", &cap)

	if stackName == "" {
		t.Error("Expected non-empty stack name")
	}
	if !strings.Contains(stackName, "test-cluster") {
		t.Error("Expected stack name to contain cluster name")
	}
	if !strings.Contains(stackName, "capability") {
		t.Error("Expected stack name to contain 'capability'")
	}
}
