package update

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

// mockStackManager is a mock implementation of the StackUpdater interface
type mockStackManager struct {
	mock.Mock
}

func (m *mockStackManager) ListPodIdentityStackNames(ctx context.Context) ([]string, error) {
	return []string{}, nil
}

func (m *mockStackManager) MustUpdateStack(ctx context.Context, options manager.UpdateStackOptions) error {
	return nil
}

func (m *mockStackManager) DescribeStack(ctx context.Context, stack *manager.Stack) (*manager.Stack, error) {
	return stack, nil
}

func (m *mockStackManager) GetIAMServiceAccounts(ctx context.Context, namespace, name string) ([]*api.ClusterIAMServiceAccount, error) {
	return []*api.ClusterIAMServiceAccount{}, nil
}

func (m *mockStackManager) GetStackTemplate(ctx context.Context, stackName string) (string, error) {
	return "", nil
}

func TestUpdatePodIdentityAssociationWithCrossAccountAccess(t *testing.T) {
	// Create a mock provider
	p := mockprovider.NewMockProvider()
	mockEKS := p.MockEKS()

	// Set up the expected API calls
	mockEKS.On("ListPodIdentityAssociations", mock.Anything, &eks.ListPodIdentityAssociationsInput{
		ClusterName:    aws.String("test-cluster"),
		Namespace:      aws.String("default"),
		ServiceAccount: aws.String("test-sa"),
	}).Return(&eks.ListPodIdentityAssociationsOutput{
		Associations: []ekstypes.PodIdentityAssociationSummary{
			{
				AssociationId: aws.String("test-association-id"),
			},
		},
	}, nil)

	mockEKS.On("DescribePodIdentityAssociation", mock.Anything, &eks.DescribePodIdentityAssociationInput{
		ClusterName:   aws.String("test-cluster"),
		AssociationId: aws.String("test-association-id"),
	}).Return(&eks.DescribePodIdentityAssociationOutput{
		Association: &ekstypes.PodIdentityAssociation{
			AssociationId: aws.String("test-association-id"),
			RoleArn:       aws.String("arn:aws:iam::111122223333:role/old-role"),
		},
	}, nil)

	// This is the key part of the test - we're capturing the input to verify the fields
	var capturedInput *eks.UpdatePodIdentityAssociationInput
	mockEKS.On("UpdatePodIdentityAssociation", mock.Anything, mock.MatchedBy(func(input *eks.UpdatePodIdentityAssociationInput) bool {
		capturedInput = input
		return true
	})).Return(&eks.UpdatePodIdentityAssociationOutput{}, nil)

	// Create a command with a mock implementation of doUpdatePodIdentityAssociation
	// that doesn't call NewProviderForExistingCluster
	cmd := &cmdutils.Cmd{
		CobraCommand:   &cobra.Command{},
		ClusterConfig:  api.NewClusterConfig(),
		ProviderConfig: api.ProviderConfig{},
	}
	cmd.ClusterConfig.Metadata.Name = "test-cluster"
	cmd.ProviderConfig.Region = "us-west-2"

	// Set up the options with cross-account access fields
	options := cmdutils.UpdatePodIdentityAssociationOptions{
		PodIdentityAssociationOptions: cmdutils.PodIdentityAssociationOptions{
			Namespace:          "default",
			ServiceAccountName: "test-sa",
		},
		RoleARN:            "arn:aws:iam::111122223333:role/source-role",
		TargetRoleARN:      "arn:aws:iam::444455556666:role/target-role",
		DisableSessionTags: true,
	}

	// Create the pod identity association in the cluster config
	cmd.ClusterConfig.IAM.PodIdentityAssociations = []api.PodIdentityAssociation{
		{
			Namespace:          options.Namespace,
			ServiceAccountName: options.ServiceAccountName,
			RoleARN:            options.RoleARN,
			TargetRoleARN:      options.TargetRoleARN,
			DisableSessionTags: options.DisableSessionTags,
		},
	}

	// Instead of calling doUpdatePodIdentityAssociation, we'll directly call the updater
	// with our mock provider to test the API calls
	stackManager := &mockStackManager{}
	updater := &podidentityassociation.Updater{
		ClusterName:  cmd.ClusterConfig.Metadata.Name,
		APIUpdater:   p.EKS(),
		StackUpdater: stackManager,
	}
	err := updater.Update(context.Background(), cmd.ClusterConfig.IAM.PodIdentityAssociations)
	require.NoError(t, err)

	// Verify that the API was called with the correct parameters
	require.NotNil(t, capturedInput)
	require.Equal(t, "test-association-id", *capturedInput.AssociationId)
	require.Equal(t, "test-cluster", *capturedInput.ClusterName)
	require.Equal(t, "arn:aws:iam::111122223333:role/source-role", *capturedInput.RoleArn)

	// Verify the cross-account access fields
	require.NotNil(t, capturedInput.TargetRoleArn)
	require.Equal(t, "arn:aws:iam::444455556666:role/target-role", *capturedInput.TargetRoleArn)

	require.NotNil(t, capturedInput.DisableSessionTags)
	require.True(t, *capturedInput.DisableSessionTags)

	// Verify all expectations were met
	mockEKS.AssertExpectations(t)
}
