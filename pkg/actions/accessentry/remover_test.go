package accessentry_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	"github.com/weaveworks/eksctl/pkg/actions/accessentry"
	"github.com/weaveworks/eksctl/pkg/actions/accessentry/fakes"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type deleteAccessEntryTest struct {
	toBeDeleted []api.AccessEntry
	mockCFN     func(stackManager *fakes.FakeStackRemover)
	mockEKS     func(provider *mockprovider.MockProvider)
	expectedErr string
}

var _ = Describe("Delete", func() {
	var (
		manager          *accessentry.Remover
		FakeStackRemover *fakes.FakeStackRemover
		mockProvider     *mockprovider.MockProvider
	)
	genericErr := fmt.Errorf("ERR")

	DescribeTable("Delete", func(t deleteAccessEntryTest) {
		FakeStackRemover = new(fakes.FakeStackRemover)
		if t.mockCFN != nil {
			t.mockCFN(FakeStackRemover)
		}

		mockProvider = mockprovider.NewMockProvider()
		if t.mockEKS != nil {
			t.mockEKS(mockProvider)
		}

		manager = accessentry.NewRemover(&api.ClusterConfig{Metadata: &api.ClusterMeta{
			Name: clusterName,
		}}, FakeStackRemover, mockProvider.MockEKS())

		err := manager.Delete(context.Background(), t.toBeDeleted)
		if t.expectedErr != "" {
			Expect(err).To(MatchError(ContainSubstring(t.expectedErr)))
			return
		}
		Expect(err).ToNot(HaveOccurred())
	},
		Entry("returns an error if calling ListAccessEntry fails", deleteAccessEntryTest{
			mockEKS: func(provider *mockprovider.MockProvider) {
				provider.MockEKS().
					On("ListAccessEntries", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&eks.ListAccessEntriesInput{}))
					}).
					Return(nil, genericErr)
			},
			expectedErr: "calling EKS API to list access entries",
		}),

		Entry("returns an error if listing stack names fails", deleteAccessEntryTest{
			toBeDeleted: []api.AccessEntry{{PrincipalARN: api.MustParseARN(mockPrincipalArn1)}},
			mockEKS: func(provider *mockprovider.MockProvider) {
				provider.MockEKS().
					On("ListAccessEntries", mock.Anything, mock.Anything).
					Return(nil, nil)
			},
			mockCFN: func(stackManager *fakes.FakeStackRemover) {
				stackManager.ListAccessEntryStackNamesReturns(nil, genericErr)
			},
			expectedErr: "listing access entry stacks",
		}),

		Entry("returns an error if deleting an owned access entry fails", deleteAccessEntryTest{
			toBeDeleted: []api.AccessEntry{{PrincipalARN: api.MustParseARN(mockPrincipalArn1)}},
			mockEKS: func(provider *mockprovider.MockProvider) {
				provider.MockEKS().
					On("ListAccessEntries", mock.Anything, mock.Anything).
					Return(nil, nil)
			},
			mockCFN: func(stackManager *fakes.FakeStackRemover) {
				stackName := accessentry.MakeStackName(clusterName, api.AccessEntry{
					PrincipalARN: api.MustParseARN(mockPrincipalArn1),
				})
				stackManager.ListAccessEntryStackNamesReturns([]string{stackName}, nil)

				stackManager.DescribeStackReturns(&types.Stack{
					StackName: &stackName,
				}, nil)

				stackManager.DeleteStackBySpecReturns(nil, genericErr)
			},
			expectedErr: "failed to delete accessentry(ies)",
		}),
		Entry("returns an error if deleting an unowned access entry fails", deleteAccessEntryTest{
			toBeDeleted: []api.AccessEntry{{PrincipalARN: api.MustParseARN(mockPrincipalArn1)}},
			mockCFN: func(stackManager *fakes.FakeStackRemover) {
				stackManager.ListAccessEntryStackNamesReturns([]string{}, nil)
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				provider.MockEKS().
					On("ListAccessEntries", mock.Anything, mock.Anything).
					Return(nil, nil)
				provider.MockEKS().
					On("DeleteAccessEntry", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&eks.DeleteAccessEntryInput{}))
					}).
					Return(nil, genericErr)
			},
			expectedErr: "failed to delete accessentry(ies)",
		}),
		Entry("deletes all access entries successfully", deleteAccessEntryTest{
			toBeDeleted: []api.AccessEntry{{PrincipalARN: api.MustParseARN(mockPrincipalArn1)}, {PrincipalARN: api.MustParseARN(mockPrincipalArn2)}},
			mockCFN: func(stackManager *fakes.FakeStackRemover) {
				stackName := accessentry.MakeStackName(clusterName, api.AccessEntry{
					PrincipalARN: api.MustParseARN(mockPrincipalArn1),
				})
				stackManager.ListAccessEntryStackNamesReturns([]string{stackName}, nil)
				stackManager.DescribeStackReturns(&types.Stack{
					StackName: &stackName,
				}, nil)
				stackManager.DeleteStackBySpecReturns(&types.Stack{
					StackName: &stackName,
				}, nil)
			},
			mockEKS: func(provider *mockprovider.MockProvider) {
				provider.MockEKS().
					On("ListAccessEntries", mock.Anything, mock.Anything).
					Return(nil, nil)
				provider.MockEKS().
					On("DeleteAccessEntry", mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&eks.DeleteAccessEntryInput{}))
					}).
					Return(&eks.DeleteAccessEntryOutput{}, nil)
			},
		}),
	)
})
