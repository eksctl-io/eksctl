package podidentityassociation_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var (
	clusterName = "test-cluster"
	roleARN     = "arn:aws:iam::111122223333:role/pod-identity-role"
)

type getPodIdentityAssociationEntry struct {
	namespace            string
	serviceAccountName   string
	mockEKS              func(provider *mockprovider.MockProvider)
	expectedAssociations []podidentityassociation.Summary
	expectedErr          string
}

var _ = Describe("Get", func() {
	var (
		getter       *podidentityassociation.Getter
		mockProvider *mockprovider.MockProvider

		associationID1 = "test-ID1"
		associationID2 = "test-ID2"
		associationID3 = "test-ID3"

		associationARN1 = "arn:aws:eks:us-west-2:111122223333:podidentityassociation/test-cluster/a-1"
		associationARN2 = "arn:aws:eks:us-west-2:111122223333:podidentityassociation/test-cluster/a-2"
		associationARN3 = "arn:aws:eks:us-west-2:111122223333:podidentityassociation/test-cluster/a-3"

		namespace1 = "test-namespace-1"
		namespace2 = "test-namespace-2"

		serviceAccountName1 = "test-sa-name-1"
		serviceAccountName2 = "test-sa-name-2"

		genericErr = fmt.Errorf("ERR")
	)

	DescribeTable("Get", func(e getPodIdentityAssociationEntry) {
		mockProvider = mockprovider.NewMockProvider()
		if e.mockEKS != nil {
			e.mockEKS(mockProvider)
		}

		getter = podidentityassociation.NewGetter(clusterName, mockProvider.MockEKS())

		summaries, err := getter.GetPodIdentityAssociations(context.Background(), e.namespace, e.serviceAccountName)
		if e.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			return
		}
		Expect(err).ToNot(HaveOccurred())
		Expect(summaries).To(ConsistOf(e.expectedAssociations))
	},
		Entry("returns an error if listing associations fails", getPodIdentityAssociationEntry{
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockProvider.MockEKS().
					On("ListPodIdentityAssociations", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.ListPodIdentityAssociationsInput{}))
					}).
					Return(nil, genericErr).
					Once()
			},
			expectedErr: "failed to list pod identity associations",
		}),

		Entry("returns an error if describing associations fails", getPodIdentityAssociationEntry{
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockProvider.MockEKS().
					On("ListPodIdentityAssociations", mock.Anything, mock.Anything).
					Return(&awseks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{
							{AssociationId: aws.String(associationID1)},
						},
					}, nil).
					Once()

				mockProvider.MockEKS().
					On("DescribePodIdentityAssociation", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribePodIdentityAssociationInput{}))
						input := args[1].(*awseks.DescribePodIdentityAssociationInput)
						Expect(aws.ToString(input.AssociationId)).To(Equal(associationID1))
					}).
					Return(nil, genericErr).
					Once()
			},
			expectedErr: "failed to describe pod identity association",
		}),

		Entry("successfully fetches all associations for a cluster", getPodIdentityAssociationEntry{
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockProvider.MockEKS().
					On("ListPodIdentityAssociations", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.ListPodIdentityAssociationsInput{}))
						input := args[1].(*awseks.ListPodIdentityAssociationsInput)
						Expect(aws.ToString(input.Namespace)).To(Equal(""))
						Expect(aws.ToString(input.ServiceAccount)).To(Equal(""))
					}).
					Return(&awseks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{
							{AssociationId: aws.String(associationID1)},
							{AssociationId: aws.String(associationID2)},
							{AssociationId: aws.String(associationID3)},
						},
					}, nil).
					Once()

				mockDescribePodIdentityAssociation(mockProvider, associationID1, associationARN1, namespace1, serviceAccountName1)
				mockDescribePodIdentityAssociationWithOwnerARN(mockProvider, associationID2, associationARN2, namespace1, serviceAccountName2,
					"arn:aws:eks:us-west-2:00:addon/cluster/vpc-cni/14c7a7ae-78a2-2c58-609e-d80af6f7bb3e")
				mockDescribePodIdentityAssociation(mockProvider, associationID3, associationARN3, namespace2, serviceAccountName2)
			},
			expectedAssociations: []podidentityassociation.Summary{
				{
					AssociationARN:     associationARN1,
					Namespace:          namespace1,
					ServiceAccountName: serviceAccountName1,
					RoleARN:            roleARN,
				},
				{
					AssociationARN:     associationARN2,
					Namespace:          namespace1,
					ServiceAccountName: serviceAccountName2,
					RoleARN:            roleARN,
					OwnerARN:           "arn:aws:eks:us-west-2:00:addon/cluster/vpc-cni/14c7a7ae-78a2-2c58-609e-d80af6f7bb3e",
				},
				{
					AssociationARN:     associationARN3,
					Namespace:          namespace2,
					ServiceAccountName: serviceAccountName2,
					RoleARN:            roleARN,
				},
			},
		}),

		Entry("successfully fetches all associations for a cluster within a namespace", getPodIdentityAssociationEntry{
			namespace: namespace1,
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockProvider.MockEKS().
					On("ListPodIdentityAssociations", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.ListPodIdentityAssociationsInput{}))
						input := args[1].(*awseks.ListPodIdentityAssociationsInput)
						Expect(aws.ToString(input.Namespace)).To(Equal(namespace1))
						Expect(aws.ToString(input.ServiceAccount)).To(Equal(""))
					}).
					Return(&awseks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{
							{AssociationId: aws.String(associationID1)},
							{AssociationId: aws.String(associationID2)},
						},
					}, nil).
					Once()

				mockDescribePodIdentityAssociation(mockProvider, associationID1, associationARN1, namespace1, serviceAccountName1)
				mockDescribePodIdentityAssociation(mockProvider, associationID2, associationARN2, namespace1, serviceAccountName2)
			},
			expectedAssociations: []podidentityassociation.Summary{
				{
					AssociationARN:     associationARN1,
					Namespace:          namespace1,
					ServiceAccountName: serviceAccountName1,
					RoleARN:            roleARN,
				},
				{
					AssociationARN:     associationARN2,
					Namespace:          namespace1,
					ServiceAccountName: serviceAccountName2,
					RoleARN:            roleARN,
				},
			},
		}),
		Entry("successfully fetches a single association", getPodIdentityAssociationEntry{
			namespace:          namespace1,
			serviceAccountName: serviceAccountName1,
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockProvider.MockEKS().
					On("ListPodIdentityAssociations", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.ListPodIdentityAssociationsInput{}))
						input := args[1].(*awseks.ListPodIdentityAssociationsInput)
						Expect(aws.ToString(input.Namespace)).To(Equal(namespace1))
						Expect(aws.ToString(input.ServiceAccount)).To(Equal(serviceAccountName1))
					}).
					Return(&awseks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{
							{AssociationId: aws.String(associationID1)},
						},
					}, nil).
					Once()

				mockDescribePodIdentityAssociation(mockProvider, associationID1, associationARN1, namespace1, serviceAccountName1)
			},
			expectedAssociations: []podidentityassociation.Summary{
				{
					AssociationARN:     associationARN1,
					Namespace:          namespace1,
					ServiceAccountName: serviceAccountName1,
					RoleARN:            roleARN,
				},
			},
		}),

		Entry("returns no association", getPodIdentityAssociationEntry{
			namespace:          namespace2,
			serviceAccountName: serviceAccountName1,
			mockEKS: func(provider *mockprovider.MockProvider) {
				mockProvider.MockEKS().
					On("ListPodIdentityAssociations", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.ListPodIdentityAssociationsInput{}))
						input := args[1].(*awseks.ListPodIdentityAssociationsInput)
						Expect(aws.ToString(input.Namespace)).To(Equal(namespace2))
						Expect(aws.ToString(input.ServiceAccount)).To(Equal(serviceAccountName1))
					}).
					Return(&awseks.ListPodIdentityAssociationsOutput{
						Associations: []ekstypes.PodIdentityAssociationSummary{},
					}, nil).
					Once()
			},
			expectedAssociations: []podidentityassociation.Summary{},
		}),
	)
})

func mockDescribePodIdentityAssociation(
	mp *mockprovider.MockProvider,
	associationID, associationARN, namespace, serviceAccount string,
) {
	mockDescribePodIdentityAssociationWithOwnerARN(mp, associationID, associationARN, namespace, serviceAccount, "")
}

func mockDescribePodIdentityAssociationWithOwnerARN(
	mp *mockprovider.MockProvider,
	associationID, associationARN, namespace, serviceAccount, ownerARN string,
) {
	mp.MockEKS().
		On("DescribePodIdentityAssociation", mock.Anything, &awseks.DescribePodIdentityAssociationInput{
			ClusterName:   aws.String(clusterName),
			AssociationId: aws.String(associationID),
		}).
		Return(&awseks.DescribePodIdentityAssociationOutput{
			Association: &ekstypes.PodIdentityAssociation{
				AssociationArn: &associationARN,
				Namespace:      &namespace,
				ServiceAccount: &serviceAccount,
				RoleArn:        &roleARN,
				OwnerArn:       &ownerARN,
			},
		}, nil).
		Once()
}
