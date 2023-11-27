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
	managerfakes "github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks/mocksv2"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Pod Identity Deleter", func() {
	type deleteEntry struct {
		podIdentityAssociations []api.PodIdentityAssociation
		mockCalls               func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS)

		expectedCalls func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS)
		expectedErr   string
	}

	mockStackManager := func(stackManager *managerfakes.FakeStackManager, stackName string) {
		stackManager.DescribeStackReturns(&cfntypes.Stack{
			StackName: aws.String(stackName),
		}, nil)
		stackManager.DeleteStackBySpecSyncStub = func(_ context.Context, _ *cfntypes.Stack, errCh chan error) error {
			close(errCh)
			return nil
		}
	}
	mockCalls := func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS, podID podidentityassociation.Identifier) {
		stackName := makeStackName(podID)
		associationID := fmt.Sprintf("%x", sha1.Sum([]byte(stackName)))
		mockListPodIdentityAssociations(eksAPI, podID, []ekstypes.PodIdentityAssociationSummary{
			{
				AssociationId: aws.String(associationID),
			},
		}, nil)
		eksAPI.On("DeletePodIdentityAssociation", mock.Anything, &eks.DeletePodIdentityAssociationInput{
			ClusterName:   aws.String(clusterName),
			AssociationId: aws.String(associationID),
		}).Return(&eks.DeletePodIdentityAssociationOutput{}, nil)
		mockStackManager(stackManager, stackName)
	}

	DescribeTable("delete pod identity association", func(e deleteEntry) {
		provider := mockprovider.NewMockProvider()
		var stackManager managerfakes.FakeStackManager
		e.mockCalls(&stackManager, provider.MockEKS())
		deleter := podidentityassociation.Deleter{
			ClusterName:  clusterName,
			StackDeleter: &stackManager,
			APIDeleter:   provider.EKS(),
		}
		err := deleter.Delete(context.Background(), podidentityassociation.ToIdentifiers(e.podIdentityAssociations))

		if e.expectedErr != "" {
			Expect(err).To(MatchError(e.expectedErr))
		} else {
			Expect(err).NotTo(HaveOccurred())
		}
		e.expectedCalls(&stackManager, provider.MockEKS())
	},
		Entry("one pod identity association exists", deleteEntry{
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
				mockListStackNames(stackManager, []podidentityassociation.Identifier{podID})
				mockCalls(stackManager, eksAPI, podID)
			},

			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(1))
				eksAPI.AssertExpectations(GinkgoT())
			},
		}),

		Entry("multiple pod identity associations exist", deleteEntry{
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
				podIDs := []podidentityassociation.Identifier{
					{
						Namespace:          "default",
						ServiceAccountName: "default",
					},
					{
						Namespace:          "kube-system",
						ServiceAccountName: "aws-node",
					},
				}
				mockListStackNames(stackManager, podIDs)
				for _, podID := range podIDs {
					mockCalls(stackManager, eksAPI, podID)
				}
			},

			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(2))
				eksAPI.AssertExpectations(GinkgoT())
			},
		}),

		Entry("some pod identity associations do not exist", deleteEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "default",
					ServiceAccountName: "default",
				},
				{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-node",
				},
				{
					Namespace:          "kube-system",
					ServiceAccountName: "kube-proxy",
				},
				{
					Namespace:          "kube-system",
					ServiceAccountName: "coredns",
				},
			},
			mockCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				podIDs := []podidentityassociation.Identifier{
					{
						Namespace:          "default",
						ServiceAccountName: "default",
					},
					{
						Namespace:          "kube-system",
						ServiceAccountName: "aws-node",
					},
					{
						Namespace:          "kube-system",
						ServiceAccountName: "coredns",
					},
				}
				mockListStackNames(stackManager, podIDs)
				for _, podID := range podIDs {
					mockCalls(stackManager, eksAPI, podID)
				}
				mockListPodIdentityAssociations(eksAPI, podidentityassociation.Identifier{
					Namespace:          "kube-system",
					ServiceAccountName: "kube-proxy",
				}, nil, nil)
			},

			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(3))
				Expect(stackManager.DeleteStackBySpecSyncCallCount()).To(Equal(3))
				eksAPI.AssertExpectations(GinkgoT())
			},
		}),

		Entry("pod identity association resource does not exist but IAM resources exist", deleteEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{
				{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-node",
				},
			},
			mockCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				podID := podidentityassociation.Identifier{
					Namespace:          "kube-system",
					ServiceAccountName: "aws-node",
				}
				mockListStackNames(stackManager, []podidentityassociation.Identifier{podID})
				mockListPodIdentityAssociations(eksAPI, podID, nil, nil)
				mockStackManager(stackManager, makeStackName(podID))
			},

			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(1))
				Expect(stackManager.DeleteStackBySpecSyncCallCount()).To(Equal(1))
				eksAPI.AssertExpectations(GinkgoT())
			},
		}),

		Entry("no pod identity associations exist", deleteEntry{
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
				podIDs := []podidentityassociation.Identifier{
					{
						Namespace:          "default",
						ServiceAccountName: "default",
					},
					{
						Namespace:          "kube-system",
						ServiceAccountName: "aws-node",
					},
				}
				mockListStackNames(stackManager, nil)
				for _, podID := range podIDs {
					mockListPodIdentityAssociations(eksAPI, podID, nil, nil)
				}
			},

			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(0))
				Expect(stackManager.DeleteStackBySpecSyncCallCount()).To(Equal(0))
				eksAPI.AssertExpectations(GinkgoT())
			},
		}),

		Entry("delete IAM resources on cluster deletion", deleteEntry{
			podIdentityAssociations: []api.PodIdentityAssociation{},
			mockCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				podIDs := []podidentityassociation.Identifier{
					{
						Namespace:          "default",
						ServiceAccountName: "default",
					},
					{
						Namespace:          "kube-system",
						ServiceAccountName: "aws-node",
					},
					{
						Namespace:          "kube-system",
						ServiceAccountName: "default",
					},
				}
				mockListStackNames(stackManager, podIDs)
				mockStackManager(stackManager, "")
			},
			expectedCalls: func(stackManager *managerfakes.FakeStackManager, eksAPI *mocksv2.EKS) {
				Expect(stackManager.ListStackNamesCallCount()).To(Equal(1))
				Expect(stackManager.DescribeStackCallCount()).To(Equal(3))
				Expect(stackManager.DeleteStackBySpecSyncCallCount()).To(Equal(3))

				var names []string
				for i := 0; i < stackManager.DescribeStackCallCount(); i++ {
					_, stack := stackManager.DescribeStackArgsForCall(i)
					names = append(names, *stack.StackName)
				}
				Expect(names).To(ConsistOf(
					makeStackName(podidentityassociation.Identifier{
						Namespace:          "default",
						ServiceAccountName: "default",
					}),
					makeStackName(podidentityassociation.Identifier{
						Namespace:          "kube-system",
						ServiceAccountName: "default",
					}),
					makeStackName(podidentityassociation.Identifier{
						Namespace:          "kube-system",
						ServiceAccountName: "aws-node",
					}),
				))
			},
		}),
	)

})
