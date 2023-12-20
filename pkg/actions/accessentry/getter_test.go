package accessentry_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/actions/accessentry"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
)

type getAccessEntryTest struct {
	principalARN   api.ARN
	mockEKS        func(provider *mockprovider.MockProvider)
	expectedOutput []accessentry.Summary
	expectedErr    string
}

const (
	clusterName       = "my-cluster"
	mockPrincipalArn1 = "arn:aws:iam::111122223333:user/admin"
	mockPrincipalArn2 = "arn:aws:iam::111122223333:user/user-1"
	mockPolicyArn1    = "arn:aws:iam::111122223333:policy/policy-1"
	mockPolicyArn2    = "arn:aws:iam::111122223333:policy/policy-2"
	namespace1        = "default"
	namespace2        = "dev"
	kGroup1           = "group1"
	kGroup2           = "group2"
)

var _ = Describe("Get", func() {
	var (
		manager      *accessentry.Getter
		mockProvider *mockprovider.MockProvider
	)
	eksErr := fmt.Errorf("EKS ERR")

	DescribeTable("Get", func(t getAccessEntryTest) {
		mockProvider = mockprovider.NewMockProvider()
		t.mockEKS(mockProvider)

		manager = accessentry.NewGetter(clusterName, mockProvider.EKS())

		summaries, err := manager.Get(context.Background(), t.principalARN)
		if t.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(t.expectedErr)))
			return
		}
		Expect(err).ToNot(HaveOccurred())
		Expect(summaries).To(Equal(t.expectedOutput))
	},
		Entry("returns an error if calling ListAccessEntries fails", getAccessEntryTest{
			mockEKS: func(provider *mockprovider.MockProvider) {
				provider.MockEKS().
					On("ListAccessEntries", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&eks.ListAccessEntriesInput{}))
					}).
					Return(nil, eksErr)
			},
			expectedErr: "calling EKS API to list access entries",
		}),

		Entry("returns an error if calling DescribeAccessEntry fails", getAccessEntryTest{
			principalARN: api.MustParseARN(mockPrincipalArn1),
			mockEKS: func(provider *mockprovider.MockProvider) {
				provider.MockEKS().
					On("DescribeAccessEntry", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&eks.DescribeAccessEntryInput{}))
					}).
					Return(nil, eksErr)
			},
			expectedErr: fmt.Sprintf("calling EKS API to describe access entry with principal ARN %s", mockPrincipalArn1),
		}),

		Entry("returns an error if calling ListAssociatedAccessPolicies fails", getAccessEntryTest{
			principalARN: api.MustParseARN(mockPrincipalArn1),
			mockEKS: func(provider *mockprovider.MockProvider) {
				provider.MockEKS().
					On("DescribeAccessEntry", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&eks.DescribeAccessEntryInput{}))
					}).
					Return(&eks.DescribeAccessEntryOutput{
						AccessEntry: &ekstypes.AccessEntry{},
					}, nil)

				provider.MockEKS().
					On("ListAssociatedAccessPolicies", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&eks.ListAssociatedAccessPoliciesInput{}))
					}).
					Return(nil, eksErr)
			},
			expectedErr: fmt.Sprintf("calling EKS API to list associated access policies for entry with principal ARN %s", mockPrincipalArn1),
		}),

		Entry("returns access entry matching principal arn", getAccessEntryTest{
			principalARN: api.MustParseARN(mockPrincipalArn1),
			mockEKS: func(provider *mockprovider.MockProvider) {
				provider.MockEKS().
					On("DescribeAccessEntry", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&eks.DescribeAccessEntryInput{}))
					}).
					Return(&eks.DescribeAccessEntryOutput{
						AccessEntry: &ekstypes.AccessEntry{
							KubernetesGroups: []string{kGroup1, kGroup2},
						},
					}, nil)

				provider.MockEKS().
					On("ListAssociatedAccessPolicies", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&eks.ListAssociatedAccessPoliciesInput{}))
					}).
					Return(&eks.ListAssociatedAccessPoliciesOutput{
						AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{
							{
								PolicyArn: aws.String(mockPolicyArn1),
								AccessScope: &ekstypes.AccessScope{
									Type:       ekstypes.AccessScopeTypeNamespace,
									Namespaces: []string{namespace1, namespace2},
								},
							},
						},
					}, nil)
			},
			expectedOutput: []accessentry.Summary{
				{
					PrincipalARN:     mockPrincipalArn1,
					KubernetesGroups: []string{kGroup1, kGroup2},
					AccessPolicies: []api.AccessPolicy{
						{
							PolicyARN: api.MustParseARN(mockPolicyArn1),
							AccessScope: api.AccessScope{
								Type:       "namespace",
								Namespaces: []string{namespace1, namespace2},
							},
						},
					},
				},
			},
		}),

		Entry("returns all access entries for the cluster", getAccessEntryTest{
			mockEKS: func(provider *mockprovider.MockProvider) {
				provider.MockEKS().
					On("ListAccessEntries", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&eks.ListAccessEntriesInput{}))
					}).
					Return(&eks.ListAccessEntriesOutput{
						AccessEntries: []string{mockPrincipalArn1, mockPrincipalArn2},
					}, nil)

				provider.MockEKS().
					On("DescribeAccessEntry", mock.Anything, &eks.DescribeAccessEntryInput{
						ClusterName:  aws.String(clusterName),
						PrincipalArn: aws.String(mockPrincipalArn1),
					}).
					Return(&eks.DescribeAccessEntryOutput{
						AccessEntry: &ekstypes.AccessEntry{
							KubernetesGroups: []string{kGroup1},
						},
					}, nil)

				provider.MockEKS().
					On("DescribeAccessEntry", mock.Anything, &eks.DescribeAccessEntryInput{
						ClusterName:  aws.String(clusterName),
						PrincipalArn: aws.String(mockPrincipalArn2),
					}).
					Return(&eks.DescribeAccessEntryOutput{
						AccessEntry: &ekstypes.AccessEntry{
							KubernetesGroups: []string{kGroup2},
						},
					}, nil)

				provider.MockEKS().
					On("ListAssociatedAccessPolicies", mock.Anything, &eks.ListAssociatedAccessPoliciesInput{
						ClusterName:  aws.String(clusterName),
						PrincipalArn: aws.String(mockPrincipalArn1),
					}).
					Return(&eks.ListAssociatedAccessPoliciesOutput{
						AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{
							{
								PolicyArn: aws.String(mockPolicyArn1),
								AccessScope: &ekstypes.AccessScope{
									Type:       ekstypes.AccessScopeTypeNamespace,
									Namespaces: []string{namespace1, namespace2},
								},
							},
						},
					}, nil)

				provider.MockEKS().
					On("ListAssociatedAccessPolicies", mock.Anything, &eks.ListAssociatedAccessPoliciesInput{
						ClusterName:  aws.String(clusterName),
						PrincipalArn: aws.String(mockPrincipalArn2),
					}).
					Return(&eks.ListAssociatedAccessPoliciesOutput{
						AssociatedAccessPolicies: []ekstypes.AssociatedAccessPolicy{
							{
								PolicyArn: aws.String(mockPolicyArn2),
								AccessScope: &ekstypes.AccessScope{
									Type: ekstypes.AccessScopeTypeCluster,
								},
							},
						},
					}, nil)
			},
			expectedOutput: []accessentry.Summary{
				{
					PrincipalARN:     mockPrincipalArn1,
					KubernetesGroups: []string{kGroup1},
					AccessPolicies: []api.AccessPolicy{
						{
							PolicyARN: api.MustParseARN(mockPolicyArn1),
							AccessScope: api.AccessScope{
								Type:       "namespace",
								Namespaces: []string{namespace1, namespace2},
							},
						},
					},
				},
				{
					PrincipalARN:     mockPrincipalArn2,
					KubernetesGroups: []string{kGroup2},
					AccessPolicies: []api.AccessPolicy{
						{
							PolicyARN: api.MustParseARN(mockPolicyArn2),
							AccessScope: api.AccessScope{
								Type: "cluster",
							},
						},
					},
				},
			},
		}),
	)
})
