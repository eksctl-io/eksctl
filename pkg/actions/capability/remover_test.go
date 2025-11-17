package capability_test

import (
	"context"
	"testing"
	"time"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/capability"
	"github.com/weaveworks/eksctl/pkg/actions/capability/mocks"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
)

func TestRemover_Delete(t *testing.T) {
	mockStackRemover := mocks.NewStackRemover(t)

	// Capability stack operations
	mockStackRemover.EXPECT().DescribeStack(mock.Anything, mock.Anything).Return(&cfntypes.Stack{}, nil).Times(2)
	mockStackRemover.EXPECT().DeleteStackBySpecSync(mock.Anything, mock.Anything, mock.Anything).Run(func(ctx context.Context, s *cfntypes.Stack, errs chan error) {
		go func() {
			errs <- nil
			close(errs)
		}()
	}).Return(nil).Times(2)

	remover := capability.NewRemover("test-cluster", mockStackRemover, nil, 5*time.Minute)

	capabilities := []capability.CapabilitySummary{
		{
			Capability: api.Capability{
				Name: "test-capability",
				Type: "ACK",
			},
		},
	}

	err := remover.Delete(context.Background(), capabilities)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRemover_DeleteTasks(t *testing.T) {
	mockStackRemover := mocks.NewStackRemover(t)

	remover := capability.NewRemover("test-cluster", mockStackRemover, nil, 5*time.Second)

	capabilities := []capability.CapabilitySummary{
		{
			Capability: api.Capability{
				Name: "cap1",
				Type: "ACK",
			},
		},
		{
			Capability: api.Capability{
				Name: "cap2",
				Type: "KRO",
			},
		},
	}

	taskTree, err := remover.DeleteTasks(context.Background(), capabilities)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should have 2 sequential task groups: capabilities and IAM roles
	if taskTree.Len() != 2 {
		t.Errorf("Expected 2 task groups, got %d", taskTree.Len())
	}
	// Should not be parallel (sequential execution)
	if taskTree.Parallel {
		t.Error("Expected sequential task execution, got parallel")
	}
}

func TestRemover_DeleteWithWait(t *testing.T) {
	mockEKSClient := mocksv2.NewEKS(t)
	mockStackRemover := mocks.NewStackRemover(t)

	// Mock DeleteCapability call - use MatchedBy to handle any context
	mockEKSClient.EXPECT().DeleteCapability(
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		&awseks.DeleteCapabilityInput{
			ClusterName:    &[]string{"test-cluster"}[0],
			CapabilityName: &[]string{"test-capability"}[0],
		}).Return(&awseks.DeleteCapabilityOutput{}, nil)

	// Mock DescribeCapability calls - first call returns not found (capability deleted)
	mockEKSClient.EXPECT().DescribeCapability(
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		&awseks.DescribeCapabilityInput{
			ClusterName:    &[]string{"test-cluster"}[0],
			CapabilityName: &[]string{"test-capability"}[0],
		}).Return(nil, &ekstypes.ResourceNotFoundException{
		Message: &[]string{"Capability not found"}[0],
	})

	// IAM stack operations
	mockStackRemover.EXPECT().DescribeStack(mock.Anything, mock.Anything).Return(&cfntypes.Stack{}, nil)
	mockStackRemover.EXPECT().DeleteStackBySpecSync(mock.Anything, mock.Anything, mock.Anything).Run(func(ctx context.Context, s *cfntypes.Stack, errs chan error) {
		go func() {
			errs <- nil
			close(errs)
		}()
	}).Return(nil)

	remover := capability.NewRemover("test-cluster", mockStackRemover, mockEKSClient, 5*time.Second)

	capabilities := []capability.CapabilitySummary{
		{
			Capability: api.Capability{
				Name: "test-capability",
				Type: "ACK",
			},
		},
	}

	err := remover.Delete(context.Background(), capabilities)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
