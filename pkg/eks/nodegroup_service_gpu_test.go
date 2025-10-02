package eks

import (
	"context"
	"testing"

	"github.com/aws/amazon-ec2-instance-selector/v3/pkg/selector"
	"github.com/stretchr/testify/assert"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// MockInstanceSelector for testing
type MockInstanceSelector struct {
	returnInstances []string
	returnError     error
}

func (m *MockInstanceSelector) Filter(ctx context.Context, filters selector.Filters) ([]string, error) {
	return m.returnInstances, m.returnError
}

func TestExpandInstanceSelector_GPUFiltering(t *testing.T) {
	tests := []struct {
		name                string
		instanceSelector    *api.InstanceSelector
		mockReturnInstances []string
		expectedInstances   []string
		expectError         bool
	}{
		{
			name: "filters out GPU instances when GPUs=0",
			instanceSelector: &api.InstanceSelector{
				VCPUs: 8,
				GPUs:  newIntPtr(0),
			},
			mockReturnInstances: []string{"m5.2xlarge", "g6f.2xlarge", "c5.2xlarge", "g4dn.xlarge"},
			expectedInstances:   []string{"m5.2xlarge", "c5.2xlarge"},
			expectError:         false,
		},
		{
			name: "includes GPU instances when GPUs=1",
			instanceSelector: &api.InstanceSelector{
				VCPUs: 8,
				GPUs:  newIntPtr(1),
			},
			mockReturnInstances: []string{"m5.2xlarge", "g6f.2xlarge", "c5.2xlarge", "g4dn.xlarge"},
			expectedInstances:   []string{"m5.2xlarge", "g6f.2xlarge", "c5.2xlarge", "g4dn.xlarge"},
			expectError:         false,
		},
		{
			name: "no filtering when GPUs is nil",
			instanceSelector: &api.InstanceSelector{
				VCPUs: 8,
			},
			mockReturnInstances: []string{"m5.2xlarge", "g6f.2xlarge", "c5.2xlarge"},
			expectedInstances:   []string{"m5.2xlarge", "g6f.2xlarge", "c5.2xlarge"},
			expectError:         false,
		},
		{
			name: "error when all instances are filtered out",
			instanceSelector: &api.InstanceSelector{
				VCPUs: 8,
				GPUs:  newIntPtr(0),
			},
			mockReturnInstances: []string{"g6f.2xlarge", "g4dn.xlarge"},
			expectedInstances:   nil,
			expectError:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSelector := &MockInstanceSelector{
				returnInstances: tt.mockReturnInstances,
			}

			service := &NodeGroupService{
				instanceSelector: mockSelector,
			}

			result, err := service.expandInstanceSelector(tt.instanceSelector, []string{"us-west-2a"})

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedInstances, result)
			}
		})
	}
}

func newIntPtr(i int) *int {
	return &i
}
