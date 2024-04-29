package podidentityassociation_test

import (
	"context"
	"crypto/sha1"
	"fmt"

	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	managerfakes "github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Pod Identity Update", func() {
	type updateEntry struct {
		podIdentityAssociations []api.PodIdentityAssociation
		mockCalls               func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS)

		expectedCalls func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS)
		expectedErr   string
	}

	mockStackManager := func(stackManager *managerfakes.FakeStackManager, stackName string, outputs []cfntypes.Output, capabilities []cfntypes.Capability) {
		stackManager.DescribeStackReturns(&cfntypes.Stack{
			StackName:    aws.String(stackName),
			Outputs:      outputs,
			Capabilities: capabilities,
		}, nil)
	}

	type mockOptions struct {
		podIdentifier             podidentityassociation.Identifier
		updateRoleARN             string
		describeStackOutputs      []cfntypes.Output
		describeStackCapabilities []cfntypes.Capability
		makeStackName             func(podidentityassociation.Identifier) string
	}

	mockCalls := func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS, o mockOptions) {
		stackName := makeIRSAv2StackName(o.podIdentifier)
		if o.makeStackName != nil {
			stackName = o.makeStackName(o.podIdentifier)
		}
		associationID := fmt.Sprintf("%x", sha1.Sum([]byte(stackName)))
		mockListPodIdentityAssociations(eksAPI, o.podIdentifier, []ekstypes.PodIdentityAssociationSummary{
			{
				AssociationId: aws.String(associationID),
			},
		}, nil)
		eksAPI.On("DescribePodIdentityAssociation", mock.Anything, &eks.DescribePodIdentityAssociationInput{
			AssociationId: aws.String(associationID),
			ClusterName:   aws.String(clusterName),
		}).Return(&eks.DescribePodIdentityAssociationOutput{
			Association: &ekstypes.PodIdentityAssociation{
				AssociationId: aws.String(associationID),
				RoleArn:       aws.String("arn:aws:iam::1234567:role/Role"),
			},
		}, nil)
		if o.updateRoleARN != "" {
			eksAPI.On("UpdatePodIdentityAssociation", mock.Anything, &eks.UpdatePodIdentityAssociationInput{
				AssociationId: aws.String(associationID),
				ClusterName:   aws.String(clusterName),
				RoleArn:       aws.String(o.updateRoleARN),
			}).Return(&eks.UpdatePodIdentityAssociationOutput{}, nil)
		}
		mockStackManager(stackManager, stackName, o.describeStackOutputs, o.describeStackCapabilities)
	}

	DescribeTable("update pod identity associations", func(e updateEntry) {
		provider := mockprovider.NewMockProvider()
		var stackManager managerfakes.FakeStackManager

		e.mockCalls(&stackManager, provider.MockEKS())
		updater := podidentityassociation.Updater{
			ClusterName:  clusterName,
			StackUpdater: &stackManager,
			APIUpdater:   provider.EKS(),
		}
		err := updater.Update(context.Background(), e.podIdentityAssociations)
		if e.expectedErr != "" {
			Expect(err).To(MatchError(e.expectedErr))
		} else {
			Expect(err).NotTo(HaveOccurred())
		}
		e.expectedCalls(&stackManager, provider.MockEKS())
	},
		Entry("pod identity association does not exist", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "default",
					ServiceAccountName: "default",
				},
			},
			mockCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				podID := podidentityassociation.Identifier{
					Namespace:          "default",
					ServiceAccountName: "default",
				}
				mockListStackNames(stackManager, nil)
				mockListPodIdentityAssociations(eksAPI, podID, nil, &ekstypes.NotFoundException{
					Message: aws.String("not found"),
				})
			},

			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListPodIdentityStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(0))
				Expect(stackManager.MustUpdateStackCallCount()).To(Equal(0))
				eksAPI.AssertExpectations(GinkgoT())
			},
			expectedErr: `error updating pod identity association "default/default": pod identity association does not exist: NotFoundException: not found`,
		}),

		Entry("attempting to update a pod identity associated with an addon ", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
				},
			},
			mockCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				podID := podidentityassociation.Identifier{
					Namespace:          "kube-system",
					ServiceAccountName: "vpc-cni",
				}
				mockListStackNames(stackManager, nil)
				mockListPodIdentityAssociations(eksAPI, podID, []ekstypes.PodIdentityAssociationSummary{
					{
						AssociationId: aws.String("a-1"),
						OwnerArn:      aws.String("vpc-cni"),
					},
				}, nil)
			},

			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListPodIdentityStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(0))
				Expect(stackManager.MustUpdateStackCallCount()).To(Equal(0))
				eksAPI.AssertExpectations(GinkgoT())
			},
			expectedErr: "error updating pod identity association \"kube-system/vpc-cni\": cannot update podidentityassociation kube-system/vpc-cni as it is in use by addon vpc-cni; " +
				"please use `eksctl update addon` instead",
		}),

		Entry("role ARN specified when the IAM resources were created by eksctl", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "default",
					ServiceAccountName: "default",
					RoleARN:            "arn:aws:iam::00000000:role/new-role",
				},
			},
			mockCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				id := podidentityassociation.Identifier{
					Namespace:          "default",
					ServiceAccountName: "default",
				}
				mockListStackNames(stackManager, []podidentityassociation.Identifier{id})
				mockCalls(stackManager, eksAPI, mockOptions{
					podIdentifier: id,
				})
			},

			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListPodIdentityStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(0))
				Expect(stackManager.MustUpdateStackCallCount()).To(Equal(0))
				eksAPI.AssertExpectations(GinkgoT())
			},

			expectedErr: `error updating pod identity association "default/default": cannot change podIdentityAssociation.roleARN since the role was created by eksctl`,
		}),

		Entry("role ARN specified when the IAM resources were not created by eksctl", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "default",
					ServiceAccountName: "default",
					RoleARN:            "arn:aws:iam::00000000:role/new-role",
				},
			},
			mockCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				mockListStackNames(stackManager, nil)
				mockCalls(stackManager, eksAPI, mockOptions{
					podIdentifier: podidentityassociation.Identifier{
						Namespace:          "default",
						ServiceAccountName: "default",
					},
					updateRoleARN: "arn:aws:iam::00000000:role/new-role",
				})
			},

			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListPodIdentityStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(0))
				Expect(stackManager.MustUpdateStackCallCount()).To(Equal(0))
				eksAPI.AssertExpectations(GinkgoT())
			},
		}),

		Entry("pod identity association has changes", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "default",
					ServiceAccountName: "default",
				},
				{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-node",
				},
			},
			mockCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				podIdentifiers := []podidentityassociation.Identifier{
					{
						Namespace:          "default",
						ServiceAccountName: "default",
					},
					{
						Namespace:          "kube-system",
						ServiceAccountName: "aws-node",
					},
				}
				mockListStackNamesWithIRSAv1(stackManager, podIdentifiers[:1], podIdentifiers[1:])
				describeStackOutputs := []cfntypes.Output{
					{
						OutputKey:   aws.String(outputs.IAMServiceAccountRoleName),
						OutputValue: aws.String("arn:aws:iam::1234567:role/Role"),
					},
				}
				for _, options := range []mockOptions{
					{
						podIdentifier:        podIdentifiers[0],
						updateRoleARN:        "arn:aws:iam::1234567:role/Role",
						describeStackOutputs: describeStackOutputs,
						makeStackName:        makeIRSAv1StackName,
					},
					{
						podIdentifier:        podIdentifiers[1],
						updateRoleARN:        "arn:aws:iam::1234567:role/Role",
						describeStackOutputs: describeStackOutputs,
					},
				} {
					mockCalls(stackManager, eksAPI, options)
				}

				stackManager.MustUpdateStackReturns(nil)
			},

			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListPodIdentityStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(4))
				Expect(stackManager.MustUpdateStackCallCount()).To(Equal(2))
				eksAPI.AssertExpectations(GinkgoT())
			},
		}),

		Entry("pod identity association has no changes", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "default",
					ServiceAccountName: "default",
				},
				{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-node",
				},
			},
			mockCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				podIdentifiers := []podidentityassociation.Identifier{
					{
						Namespace:          "default",
						ServiceAccountName: "default",
					},
					{
						Namespace:          "kube-system",
						ServiceAccountName: "aws-node",
					},
				}
				mockListStackNames(stackManager, podIdentifiers)
				for _, options := range []mockOptions{
					{
						podIdentifier: podIdentifiers[0],
					},
					{
						podIdentifier: podIdentifiers[1],
					},
				} {
					mockCalls(stackManager, eksAPI, options)
				}

				stackManager.MustUpdateStackReturns(&manager.NoChangeError{
					Msg: "no changes found",
				})
			},

			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListPodIdentityStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(2))
				Expect(stackManager.MustUpdateStackCallCount()).To(Equal(2))
				eksAPI.AssertExpectations(GinkgoT())
			},
		}),

		Entry("fields that cannot be updated specified when the IAM resources were not created by eksctl", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-node",
					RoleARN:            "arn:aws:iam::00000000:role/new-role",
					WellKnownPolicies: api.WellKnownPolicies{
						AutoScaler: true,
					},
					PermissionPolicyARNs: []string{"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"},
				},
			},
			mockCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				mockListStackNames(stackManager, nil)
				mockCalls(stackManager, eksAPI, mockOptions{
					podIdentifier: podidentityassociation.Identifier{
						Namespace:          "kube-system",
						ServiceAccountName: "aws-node",
					},
				})
			},

			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListPodIdentityStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(0))
				Expect(stackManager.MustUpdateStackCallCount()).To(Equal(0))
				eksAPI.AssertExpectations(GinkgoT())
			},

			expectedErr: `error updating pod identity association "kube-system/aws-node": only namespace, serviceAccountName and roleARN can be specified if the role was not created by eksctl`,
		}),

		Entry("roleName specified when the pod identity association was not created with a roleName", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-node",
					RoleName:           "default-role",
				},
			},
			mockCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				podIdentifier := podidentityassociation.Identifier{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-node",
				}
				mockListStackNames(stackManager, []podidentityassociation.Identifier{podIdentifier})

				mockCalls(stackManager, eksAPI, mockOptions{
					podIdentifier:             podIdentifier,
					describeStackCapabilities: []cfntypes.Capability{cfntypes.CapabilityCapabilityIam},
				})
			},

			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListPodIdentityStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(1))
				Expect(stackManager.MustUpdateStackCallCount()).To(Equal(0))
				eksAPI.AssertExpectations(GinkgoT())
			},

			expectedErr: `error updating pod identity association "kube-system/aws-node": cannot update role name if the pod identity association was not created with a role name`,
		}),

		Entry("roleName specified when the pod identity association was created with a roleName", updateEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "default",
					ServiceAccountName: "default",
				},
				{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-node",
					RoleName:           "default-role",
				},
			},
			mockCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				podIdentifiers := []podidentityassociation.Identifier{
					{
						Namespace:          "default",
						ServiceAccountName: "default",
					},
					{
						Namespace:          "kube-system",
						ServiceAccountName: "aws-node",
					},
				}
				mockListStackNamesWithIRSAv1(stackManager, podIdentifiers[:1], podIdentifiers[1:])
				describeStackOutputs := []cfntypes.Output{
					{
						OutputKey:   aws.String(outputs.IAMServiceAccountRoleName),
						OutputValue: aws.String("arn:aws:iam::1234567:role/Role"),
					},
				}
				for _, options := range []mockOptions{
					{
						podIdentifier:        podIdentifiers[0],
						describeStackOutputs: describeStackOutputs,
						updateRoleARN:        "arn:aws:iam::1234567:role/Role",
						makeStackName:        makeIRSAv1StackName,
					},
					{
						podIdentifier:             podIdentifiers[1],
						updateRoleARN:             "arn:aws:iam::1234567:role/Role",
						describeStackOutputs:      describeStackOutputs,
						describeStackCapabilities: []cfntypes.Capability{cfntypes.CapabilityCapabilityIam, cfntypes.CapabilityCapabilityNamedIam},
					},
				} {
					mockCalls(stackManager, eksAPI, options)
				}
				stackManager.MustUpdateStackReturns(nil)
			},

			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListPodIdentityStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(4))
				Expect(stackManager.MustUpdateStackCallCount()).To(Equal(2))
				eksAPI.AssertExpectations(GinkgoT())
			},
		}),
	)
})

func mockListPodIdentityAssociations(eksAPI *mocksv2.EKS, podID podidentityassociation.Identifier, output []ekstypes.PodIdentityAssociationSummary, err error) {
	eksAPI.On("ListPodIdentityAssociations", mock.Anything, &eks.ListPodIdentityAssociationsInput{
		ClusterName:    aws.String(clusterName),
		Namespace:      aws.String(podID.Namespace),
		ServiceAccount: aws.String(podID.ServiceAccountName),
	}).Return(&eks.ListPodIdentityAssociationsOutput{
		Associations: output,
	}, err)
}

func makeIRSAv1StackName(podID podidentityassociation.Identifier) string {
	return fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-%s-%s", clusterName, podID.Namespace, podID.ServiceAccountName)
}

func makeIRSAv2StackName(podID podidentityassociation.Identifier) string {
	return podidentityassociation.MakeStackName(clusterName, podID.Namespace, podID.ServiceAccountName)
}

func mockListStackNames(stackManager *managerfakes.FakeStackManager, podIDs []podidentityassociation.Identifier) {
	mockListStackNamesWithIRSAv1(stackManager, []podidentityassociation.Identifier{}, podIDs)
}

func mockListStackNamesWithIRSAv1(
	stackManager *managerfakes.FakeStackManager,
	irsaV1podIdentifiers []podidentityassociation.Identifier,
	irsaV2podIdentifiers []podidentityassociation.Identifier,
) {
	var stackNames []string
	for _, id := range irsaV1podIdentifiers {
		stackNames = append(stackNames, makeIRSAv1StackName(id))
	}
	for _, id := range irsaV2podIdentifiers {
		stackNames = append(stackNames, makeIRSAv2StackName(id))
	}
	stackManager.ListPodIdentityStackNamesReturns(stackNames, nil)
}
