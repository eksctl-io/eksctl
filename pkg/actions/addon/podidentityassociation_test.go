package addon_test

import (
	"context"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/addon"
	"github.com/weaveworks/eksctl/pkg/actions/addon/mocks"
	piamocks "github.com/weaveworks/eksctl/pkg/actions/podidentityassociation/mocks"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Update Pod Identity Association", func() {
	type piaMocks struct {
		stackDeleter *piamocks.StackDeleter
		roleCreator  *mocks.IAMRoleCreator
		roleUpdater  *mocks.IAMRoleUpdater
		eks          *mocksv2.EKS
	}
	type updateEntry struct {
		podIdentityAssociations         []api.PodIdentityAssociation
		mockCalls                       func(m piaMocks)
		existingPodIdentityAssociations []addon.PodIdentityAssociationSummary

		expectedAddonPodIdentityAssociations []ekstypes.AddonPodIdentityAssociations

		expectedErr string
	}

	const clusterName = "test"

	makeID := func(i int) string {
		return fmt.Sprintf("a-%d", i+1)
	}

	mockDescribePodIdentityAssociation := func(eksAPI *mocksv2.EKS, roleARNs ...string) {
		for i, roleARN := range roleARNs {
			id := aws.String(makeID(i))
			eksAPI.On("DescribePodIdentityAssociation", mock.Anything, &eks.DescribePodIdentityAssociationInput{
				ClusterName:   aws.String(clusterName),
				AssociationId: id,
			}).Return(&eks.DescribePodIdentityAssociationOutput{
				Association: &ekstypes.PodIdentityAssociation{
					AssociationId: id,
					RoleArn:       aws.String(roleARN),
				},
			}, nil).Once()
		}
	}

	DescribeTable("update pod identity association", func(e updateEntry) {
		provider := mockprovider.NewMockProvider()
		var (
			roleCreator  mocks.IAMRoleCreator
			roleUpdater  mocks.IAMRoleUpdater
			stackDeleter piamocks.StackDeleter
		)

		piaUpdater := &addon.PodIdentityAssociationUpdater{
			ClusterName:             clusterName,
			IAMRoleCreator:          &roleCreator,
			IAMRoleUpdater:          &roleUpdater,
			EKSPodIdentityDescriber: provider.MockEKS(),
			StackDeleter:            &stackDeleter,
		}
		if e.mockCalls != nil {
			e.mockCalls(piaMocks{
				stackDeleter: &stackDeleter,
				roleCreator:  &roleCreator,
				roleUpdater:  &roleUpdater,
				eks:          provider.MockEKS(),
			})
		}
		addonPodIdentityAssociations, err := piaUpdater.UpdateRole(context.Background(), e.podIdentityAssociations, "main", e.existingPodIdentityAssociations)
		if e.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			return
		}
		Expect(err).NotTo(HaveOccurred())
		Expect(addonPodIdentityAssociations).To(Equal(e.expectedAddonPodIdentityAssociations))
		for _, asserter := range []interface {
			AssertExpectations(t mock.TestingT) bool
		}{
			&roleCreator,
			&roleUpdater,
			provider.MockEKS(),
		} {
			asserter.AssertExpectations(GinkgoT())
		}
	},
		Entry("addon contains pod identity that does not exist", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
				},
			},
			mockCalls: func(m piaMocks) {
				m.roleCreator.On("Create", mock.Anything, &api.PodIdentityAssociation{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
				}, "main").Return("role-1", nil)
				m.stackDeleter.On("DescribeStack", mock.Anything, &manager.Stack{
					StackName: aws.String("eksctl-test-addon-main"),
				}).Return(nil, &smithy.OperationError{
					Err: errors.New("ValidationError"),
				}).Once()

			},
			expectedAddonPodIdentityAssociations: []ekstypes.AddonPodIdentityAssociations{
				{
					ServiceAccount: aws.String("vpc-cni"),
					RoleArn:        aws.String("role-1"),
				},
			},
		}),

		Entry("addon contains pod identities, some of which do not exist", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
				},
				{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-ebs-csi-driver",
				},
				{
					Namespace:          "karpenter",
					ServiceAccountName: "karpenter",
				},
			},
			existingPodIdentityAssociations: []addon.PodIdentityAssociationSummary{
				{
					Namespace:      "kube-system",
					ServiceAccount: "vpc-cni",
					AssociationID:  "a-1",
				},
			},
			mockCalls: func(m piaMocks) {
				mockDescribePodIdentityAssociation(m.eks, "cni-role")

				m.roleUpdater.On("Update", mock.Anything, api.PodIdentityAssociation{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
				}, "eksctl-test-addon-main-podidentityrole-vpc-cni", "a-1").Return("cni-role-2", true, nil).Once()
				m.stackDeleter.On("DescribeStack", mock.Anything, &manager.Stack{
					StackName: aws.String("eksctl-test-addon-main-podidentityrole-vpc-cni"),
				}).Return(&manager.Stack{
					StackName: aws.String("eksctl-test-addon-main-podidentityrole-vpc-cni"),
				}, nil).Twice()
				m.stackDeleter.On("DescribeStack", mock.Anything, &manager.Stack{
					StackName: aws.String("eksctl-test-addon-main"),
				}).Return(nil, &smithy.OperationError{
					Err: errors.New("ValidationError"),
				}).Twice()

				m.roleCreator.On("Create", mock.Anything, &api.PodIdentityAssociation{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-ebs-csi-driver",
				}, "main").Return("csi-role", nil).Once()
				m.roleCreator.On("Create", mock.Anything, &api.PodIdentityAssociation{
					Namespace:          "karpenter",
					ServiceAccountName: "karpenter",
				}, "main").Return("karpenter-role", nil).Once()
			},
			expectedAddonPodIdentityAssociations: []ekstypes.AddonPodIdentityAssociations{
				{
					ServiceAccount: aws.String("vpc-cni"),
					RoleArn:        aws.String("cni-role-2"),
				},
				{
					ServiceAccount: aws.String("aws-ebs-csi-driver"),
					RoleArn:        aws.String("csi-role"),
				},
				{
					ServiceAccount: aws.String("karpenter"),
					RoleArn:        aws.String("karpenter-role"),
				},
			},
		}),

		Entry("addon contains pod identities that already exist", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
				},
				{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-ebs-csi-driver",
				},
				{
					Namespace:          "karpenter",
					ServiceAccountName: "karpenter",
				},
			},
			existingPodIdentityAssociations: []addon.PodIdentityAssociationSummary{
				{
					Namespace:      "kube-system",
					ServiceAccount: "vpc-cni",
					AssociationID:  "a-1",
				},
				{
					Namespace:      "kube-system",
					ServiceAccount: "aws-ebs-csi-driver",
					AssociationID:  "a-2",
				},
				{
					Namespace:      "karpenter",
					ServiceAccount: "karpenter",
					AssociationID:  "a-3",
				},
			},
			mockCalls: func(m piaMocks) {
				mockDescribePodIdentityAssociation(m.eks, "cni-role", "csi-role", "karpenter-role")

				for i, updateInput := range []struct {
					namespace            string
					serviceAccount       string
					hasIAMResourcesStack bool
					stackName            string
					returnRole           string
				}{
					{
						namespace:            "kube-system",
						serviceAccount:       "vpc-cni",
						hasIAMResourcesStack: true,
						stackName:            "kube-system-vpc-cni",
						returnRole:           "cni-role-2",
					},
					{
						namespace:            "kube-system",
						serviceAccount:       "aws-ebs-csi-driver",
						hasIAMResourcesStack: true,
						stackName:            "kube-system-aws-ebs-csi-driver",
						returnRole:           "csi-role-2",
					},
					{
						namespace:            "karpenter",
						serviceAccount:       "karpenter",
						hasIAMResourcesStack: true,
						stackName:            "karpenter-karpenter",
						returnRole:           "karpenter-role-2",
					},
				} {
					id := makeID(i)

					stackName := fmt.Sprintf("eksctl-test-addon-main-podidentityrole-%s", updateInput.serviceAccount)
					m.roleUpdater.On("Update", mock.Anything, api.PodIdentityAssociation{
						Namespace:          updateInput.namespace,
						ServiceAccountName: updateInput.serviceAccount,
					}, stackName, id).Return(updateInput.returnRole, true, nil).Once()

					m.stackDeleter.On("DescribeStack", mock.Anything, &manager.Stack{
						StackName: aws.String(stackName),
					}).Return(&manager.Stack{
						StackName: aws.String(stackName),
					}, nil)
				}
			},
			expectedAddonPodIdentityAssociations: []ekstypes.AddonPodIdentityAssociations{
				{
					ServiceAccount: aws.String("vpc-cni"),
					RoleArn:        aws.String("cni-role-2"),
				},
				{
					ServiceAccount: aws.String("aws-ebs-csi-driver"),
					RoleArn:        aws.String("csi-role-2"),
				},
				{
					ServiceAccount: aws.String("karpenter"),
					RoleArn:        aws.String("karpenter-role-2"),
				},
			},
		}),

		Entry("addon contains pod identities that do not exist and have a pre-existing roleARN", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
					RoleARN:            "role-1",
				},
				{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-ebs-csi-driver",
					RoleARN:            "role-2",
				},
				{
					Namespace:          "karpenter",
					ServiceAccountName: "karpenter",
					RoleARN:            "role-3",
				},
			},
			expectedAddonPodIdentityAssociations: []ekstypes.AddonPodIdentityAssociations{
				{
					ServiceAccount: aws.String("vpc-cni"),
					RoleArn:        aws.String("role-1"),
				},
				{
					ServiceAccount: aws.String("aws-ebs-csi-driver"),
					RoleArn:        aws.String("role-2"),
				},
				{
					ServiceAccount: aws.String("karpenter"),
					RoleArn:        aws.String("role-3"),
				},
			},
		}),

		Entry("addon contains pod identities that already exist and have a pre-existing roleARN", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
					RoleARN:            "role-1",
				},
				{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-ebs-csi-driver",
					RoleARN:            "role-2",
				},
				{
					Namespace:          "karpenter",
					ServiceAccountName: "karpenter",
					RoleARN:            "role-3",
				},
			},
			expectedAddonPodIdentityAssociations: []ekstypes.AddonPodIdentityAssociations{
				{
					ServiceAccount: aws.String("vpc-cni"),
					RoleArn:        aws.String("role-1"),
				},
				{
					ServiceAccount: aws.String("aws-ebs-csi-driver"),
					RoleArn:        aws.String("role-2"),
				},
				{
					ServiceAccount: aws.String("karpenter"),
					RoleArn:        aws.String("role-3"),
				},
			},
		}),

		Entry("addon contains pod identity IAM resources created by eksctl but are being updated with a new roleARN", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
					RoleARN:            "role-1",
				},
				{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-ebs-csi-driver",
					RoleARN:            "role-2",
				},
				{
					Namespace:          "karpenter",
					ServiceAccountName: "karpenter",
					RoleARN:            "karpenter-role",
				},
			},
			existingPodIdentityAssociations: []addon.PodIdentityAssociationSummary{
				{
					Namespace:      "kube-system",
					ServiceAccount: "vpc-cni",
					AssociationID:  "a-1",
				},
				{
					Namespace:      "kube-system",
					ServiceAccount: "aws-ebs-csi-driver",
					AssociationID:  "a-2",
				},
				{
					Namespace:      "karpenter",
					ServiceAccount: "karpenter",
					AssociationID:  "a-3",
				},
			},
			mockCalls: func(m piaMocks) {
				mockDescribePodIdentityAssociation(m.eks, "role-1", "role-2", "role-3")
				for _, serviceAccount := range []string{"vpc-cni", "aws-ebs-csi-driver"} {
					m.stackDeleter.On("DescribeStack", mock.Anything, &manager.Stack{
						StackName: aws.String(fmt.Sprintf("eksctl-test-addon-main-podidentityrole-%s", serviceAccount)),
					}).Return(nil, &smithy.OperationError{
						Err: fmt.Errorf("ValidationError"),
					}).Once()

					m.stackDeleter.On("DescribeStack", mock.Anything, &manager.Stack{
						StackName: aws.String("eksctl-test-addon-main"),
					}).Return(nil, &smithy.OperationError{
						Err: fmt.Errorf("ValidationError"),
					}).Once()
				}
				m.stackDeleter.On("DescribeStack", mock.Anything, &manager.Stack{
					StackName: aws.String("eksctl-test-addon-main-podidentityrole-karpenter"),
				}).Return(&manager.Stack{}, nil).Once()
			},
			expectedErr: "cannot change podIdentityAssociation.roleARN since the role was created by eksctl",
		}),

		Entry("addon contains pod identity created with a pre-existing roleARN and is being updated", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
					RoleARN:            "vpc-cni-role-2",
				},
			},
			existingPodIdentityAssociations: []addon.PodIdentityAssociationSummary{
				{
					Namespace:      "kube-system",
					ServiceAccount: "vpc-cni",
					AssociationID:  "a-1",
				},
			},
			mockCalls: func(m piaMocks) {
				mockDescribePodIdentityAssociation(m.eks, "vpc-cni-role")
				m.stackDeleter.On("DescribeStack", mock.Anything, &manager.Stack{
					StackName: aws.String("eksctl-test-addon-main-podidentityrole-vpc-cni"),
				}).Return(&manager.Stack{}, nil).Once()
			},
			expectedAddonPodIdentityAssociations: []ekstypes.AddonPodIdentityAssociations{
				{
					RoleArn:        aws.String("vpc-cni-role-2"),
					ServiceAccount: aws.String("vpc-cni"),
				},
			},
			expectedErr: "cannot change podIdentityAssociation.roleARN since the role was created by eksctl",
		}),

		Entry("addon contains pod identity created with a pre-existing roleARN but it is no longer set", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
				},
			},
			existingPodIdentityAssociations: []addon.PodIdentityAssociationSummary{
				{
					Namespace:      "kube-system",
					ServiceAccount: "vpc-cni",
					AssociationID:  "a-1",
				},
			},
			mockCalls: func(m piaMocks) {
				mockDescribePodIdentityAssociation(m.eks, "vpc-cni-role")
				m.stackDeleter.On("DescribeStack", mock.Anything, &manager.Stack{
					StackName: aws.String("eksctl-test-addon-main-podidentityrole-vpc-cni"),
				}).Return(nil, &smithy.OperationError{
					Err: errors.New("ValidationError"),
				})
				m.stackDeleter.On("DescribeStack", mock.Anything, &manager.Stack{
					StackName: aws.String("eksctl-test-addon-main"),
				}).Return(nil, &smithy.OperationError{
					Err: errors.New("ValidationError"),
				}).Once()
			},
			expectedErr: "podIdentityAssociation.roleARN is required since the role was not created by eksctl",
		}),

		Entry("addon contains pod identities that do not exist and have a pre-existing roleARN", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
					RoleARN:            "role-1",
				},
				{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-ebs-csi-driver",
					RoleARN:            "role-2",
				},
				{
					Namespace:          "karpenter",
					ServiceAccountName: "karpenter",
					RoleARN:            "role-3",
				},
			},
			expectedAddonPodIdentityAssociations: []ekstypes.AddonPodIdentityAssociations{
				{
					ServiceAccount: aws.String("vpc-cni"),
					RoleArn:        aws.String("role-1"),
				},
				{
					ServiceAccount: aws.String("aws-ebs-csi-driver"),
					RoleArn:        aws.String("role-2"),
				},
				{
					ServiceAccount: aws.String("karpenter"),
					RoleArn:        aws.String("role-3"),
				},
			},
		}),

		Entry("addon contains a pod identity with an IRSA stack", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
				},
			},
			existingPodIdentityAssociations: []addon.PodIdentityAssociationSummary{
				{
					Namespace:      "kube-system",
					ServiceAccount: "vpc-cni",
					AssociationID:  "a-1",
				},
			},
			mockCalls: func(m piaMocks) {
				mockDescribePodIdentityAssociation(m.eks, "cni-role")

				m.roleUpdater.On("Update", mock.Anything, api.PodIdentityAssociation{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
				}, "eksctl-test-addon-main", "a-1").Return("cni-role-2", true, nil).Once()
				m.stackDeleter.On("DescribeStack", mock.Anything, &manager.Stack{
					StackName: aws.String("eksctl-test-addon-main-podidentityrole-vpc-cni"),
				}).Return(nil, &smithy.OperationError{
					Err: fmt.Errorf("ValidationError"),
				}).Once()

				m.stackDeleter.On("DescribeStack", mock.Anything, &manager.Stack{
					StackName: aws.String("eksctl-test-addon-main"),
				}).Return(&manager.Stack{
					StackName: aws.String("eksctl-test-addon-main"),
				}, nil).Once()
			},
			expectedAddonPodIdentityAssociations: []ekstypes.AddonPodIdentityAssociations{
				{
					ServiceAccount: aws.String("vpc-cni"),
					RoleArn:        aws.String("cni-role-2"),
				},
			},
		}),

		Entry("addon using IRSA is updated to use pod identity", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
				},
			},
			mockCalls: func(m piaMocks) {
				m.roleCreator.On("Create", mock.Anything, &api.PodIdentityAssociation{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
				}, "main").Return("role-1", nil)
				irsaStack := &manager.Stack{
					StackName: aws.String("eksctl-test-addon-main"),
				}
				m.stackDeleter.On("DescribeStack", mock.Anything, irsaStack).Return(irsaStack, nil).Once()
				m.stackDeleter.EXPECT().DeleteStackBySpecSync(mock.Anything, irsaStack, mock.Anything).RunAndReturn(func(ctx context.Context, stack *cfntypes.Stack, errCh chan error) error {
					close(errCh)
					return nil
				}).Once()

			},
			expectedAddonPodIdentityAssociations: []ekstypes.AddonPodIdentityAssociations{
				{
					ServiceAccount: aws.String("vpc-cni"),
					RoleArn:        aws.String("role-1"),
				},
			},
		}),
	)
})
