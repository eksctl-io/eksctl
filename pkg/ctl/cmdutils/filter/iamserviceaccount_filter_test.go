package filter

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("iamserviceaccount filter", func() {

	Context("Match", func() {
		var (
			filter             *IAMServiceAccountFilter
			allServiceAccounts []*api.ClusterIAMServiceAccount
		)

		BeforeEach(func() {
			filter = NewIAMServiceAccountFilter()
			allServiceAccounts = serviceAccounts(
				"dev1",
				"dev2",
				"dev3",
				"test1",
				"test2",
				"test3",
			)
		})

		It("only-missing (only-remote) works correctly", func() {
			mockLister := newMockServiceAccountLister(
				"sa/dev1",
				"sa/dev2",
				"sa/dev3",
				"sa/test1",
				"sa/test2",
				"sa/test3",
				"sa/only-remote-1",
				"sa/only-remote-2",
			)
			err := filter.SetIncludeOrExcludeMissingFilter(mockLister, true, &allServiceAccounts)
			Expect(err).ToNot(HaveOccurred())

			included, excluded := filter.MatchAll(allServiceAccounts)
			Expect(included).To(HaveLen(2))
			Expect(included.HasAll("sa/only-remote-1", "sa/only-remote-2")).To(BeTrue())
			Expect(excluded).To(HaveLen(6))
			Expect(excluded.HasAll(
				"sa/dev1",
				"sa/dev2",
				"sa/dev3",
				"sa/test1",
				"sa/test2",
				"sa/test3",
			)).To(BeTrue())
		})
		// TODO test SetExcludeExistingFilter() which requires mocking different kubernetes functions
	})
})

func serviceAccounts(names ...string) []*api.ClusterIAMServiceAccount {
	serviceAccounts := make([]*api.ClusterIAMServiceAccount, 0)
	for _, name := range names {
		serviceAccounts = append(serviceAccounts, &api.ClusterIAMServiceAccount{
			ClusterIAMMeta: api.ClusterIAMMeta{
				Namespace: "sa",
				Name:      name,
			},
		})
	}
	return serviceAccounts
}

type mockSALister struct {
	result []string
}

func (l *mockSALister) ListIAMServiceAccountStacks() ([]string, error) {
	return l.result, nil
}

func newMockServiceAccountLister(serviceAccounts ...string) *mockSALister {
	return &mockSALister{
		result: serviceAccounts,
	}
}
