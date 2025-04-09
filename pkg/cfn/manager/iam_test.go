package manager_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/stretchr/testify/assert"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

// mockStackCollection is a wrapper around StackCollection that allows us to override methods for testing
type mockStackCollection struct {
	*manager.StackCollection
	stacks []*manager.Stack
}

// ListStacks overrides the ListStacks method to return our predefined stacks
func (m *mockStackCollection) ListStacks(ctx context.Context) ([]*manager.Stack, error) {
	return m.stacks, nil
}

func TestGetIAMServiceAccounts(t *testing.T) {
	testCases := []struct {
		name            string
		nameFilter      string
		namespaceFilter string
		expectedCount   int
		expectedNames   []string
	}{
		{
			name:            "No filters - should return all service accounts",
			nameFilter:      "",
			namespaceFilter: "",
			expectedCount:   3,
			expectedNames:   []string{"test-sa1", "test-sa2", "test-sa3"},
		},
		{
			name:            "Filter by name only",
			nameFilter:      "test-sa1",
			namespaceFilter: "",
			expectedCount:   1,
			expectedNames:   []string{"test-sa1"},
		},
		{
			name:            "Filter by namespace only",
			nameFilter:      "",
			namespaceFilter: "kube-system",
			expectedCount:   1,
			expectedNames:   []string{"test-sa2"},
		},
		{
			name:            "Filter by both name and namespace",
			nameFilter:      "test-sa3",
			namespaceFilter: "default",
			expectedCount:   1,
			expectedNames:   []string{"test-sa3"},
		},
		{
			name:            "Filter by name that doesn't exist",
			nameFilter:      "non-existent",
			namespaceFilter: "",
			expectedCount:   0,
			expectedNames:   []string{},
		},
		{
			name:            "Filter by namespace that doesn't exist",
			nameFilter:      "",
			namespaceFilter: "non-existent",
			expectedCount:   0,
			expectedNames:   []string{},
		},
		{
			name:            "Filter by name and namespace that don't match",
			nameFilter:      "test-sa1",
			namespaceFilter: "kube-system",
			expectedCount:   0,
			expectedNames:   []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock provider
			p := mockprovider.NewMockProvider()

			// Create a test cluster config
			cfg := api.NewClusterConfig()
			cfg.Metadata.Name = "test-cluster"
			cfg.Metadata.Region = "us-west-2"

			// Create a stack collection with the mock provider
			realStackCollection := manager.NewStackCollection(p, cfg)

			// Create our mock stacks
			stacks := []*manager.Stack{
				{
					StackName:   "eksctl-test-cluster-addon-iamserviceaccount-default-test-sa1",
					StackStatus: types.StackStatusCreateComplete,
					Tags: []types.Tag{
						{
							Key:   aws.String(api.IAMServiceAccountNameTag),
							Value: aws.String("default/test-sa1"),
						},
					},
					Outputs: []types.Output{
						{
							OutputKey:   aws.String(outputs.IAMServiceAccountRoleName),
							OutputValue: aws.String("arn:aws:iam::123456789012:role/test-sa1-role"),
						},
					},
				},
				{
					StackName:   "eksctl-test-cluster-addon-iamserviceaccount-kube-system-test-sa2",
					StackStatus: types.StackStatusCreateComplete,
					Tags: []types.Tag{
						{
							Key:   aws.String(api.IAMServiceAccountNameTag),
							Value: aws.String("kube-system/test-sa2"),
						},
					},
					Outputs: []types.Output{
						{
							OutputKey:   aws.String(outputs.IAMServiceAccountRoleName),
							OutputValue: aws.String("arn:aws:iam::123456789012:role/test-sa2-role"),
						},
					},
				},
				{
					StackName:   "eksctl-test-cluster-addon-iamserviceaccount-default-test-sa3",
					StackStatus: types.StackStatusCreateComplete,
					Tags: []types.Tag{
						{
							Key:   aws.String(api.IAMServiceAccountNameTag),
							Value: aws.String("default/test-sa3"),
						},
					},
					Outputs: []types.Output{
						{
							OutputKey:   aws.String(outputs.IAMServiceAccountRoleName),
							OutputValue: aws.String("arn:aws:iam::123456789012:role/test-sa3-role"),
						},
					},
				},
				{
					StackName:   "eksctl-test-cluster-nodegroup-ng-1",
					StackStatus: types.StackStatusCreateComplete,
				},
			}

			// Create our mock stack collection
			mockSC := &mockStackCollection{
				StackCollection: realStackCollection,
				stacks:          stacks,
			}

			// Call the function being tested
			serviceAccounts, err := mockSC.GetIAMServiceAccounts(context.Background(), tc.nameFilter, tc.namespaceFilter)

			// Verify results
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCount, len(serviceAccounts))

			// Verify the names of the returned service accounts
			actualNames := make([]string, 0, len(serviceAccounts))
			for _, sa := range serviceAccounts {
				actualNames = append(actualNames, sa.Name)
			}

			// Check that all expected names are in the actual names
			for _, expectedName := range tc.expectedNames {
				assert.Contains(t, actualNames, expectedName)
			}
		})
	}
}
