package capability_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/capability"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
)

func TestGetter_Get(t *testing.T) {
	mockEKSClient := mocksv2.NewEKS(t)

	mockEKSClient.EXPECT().ListCapabilities(mock.Anything, mock.MatchedBy(func(input *eks.ListCapabilitiesInput) bool {
		return *input.ClusterName == "test-cluster"
	})).Return(&eks.ListCapabilitiesOutput{
		Capabilities: []ekstypes.CapabilitySummary{
			{
				CapabilityName: stringPtr("test-capability"),
				Type:           ekstypes.CapabilityTypeAck,
			},
		},
	}, nil)

	mockEKSClient.EXPECT().DescribeCapability(mock.Anything, mock.MatchedBy(func(input *eks.DescribeCapabilityInput) bool {
		return *input.ClusterName == "test-cluster" && *input.CapabilityName == "test-capability"
	})).Return(&eks.DescribeCapabilityOutput{
		Capability: &ekstypes.Capability{
			CapabilityName: stringPtr("test-capability"),
			Type:           ekstypes.CapabilityTypeAck,
			Status:         ekstypes.CapabilityStatusActive,
		},
	}, nil)

	getter := capability.NewGetter("test-cluster", mockEKSClient)

	capabilities, err := getter.Get(context.Background(), "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(capabilities) != 1 {
		t.Errorf("Expected 1 capability, got %d", len(capabilities))
	}

	if capabilities[0].Name != "test-capability" {
		t.Errorf("Expected capability name 'test-capability', got %s", capabilities[0].Name)
	}

	// Mock expectations are automatically verified by mocksv2
}

func stringPtr(s string) *string {
	return &s
}
