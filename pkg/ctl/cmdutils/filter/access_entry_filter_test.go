package filter_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter/filterfakes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
)

var _ = Describe("Access Entry", func() {
	type accessEntryTest struct {
		clusterName    string
		accessEntries  []api.AccessEntry
		existingStacks []string

		expectedAccessEntries []api.AccessEntry
	}

	DescribeTable("access entry filter", func(aet accessEntryTest) {
		lister := &filterfakes.FakeAccessEntryLister{}
		lister.ListAccessEntryStackNamesReturns(aet.existingStacks, nil)
		f := &filter.AccessEntry{
			Lister:      lister,
			ClusterName: aet.clusterName,
		}
		actual, err := f.FilterOutExistingStacks(context.Background(), aet.accessEntries)
		Expect(err).NotTo(HaveOccurred())
		Expect(actual).To(Equal(aet.expectedAccessEntries))
	},
		Entry("one entry exists", accessEntryTest{
			clusterName: "web",
			accessEntries: []api.AccessEntry{
				{
					PrincipalARN:       api.MustParseARN("arn:aws:iam::012345689:role/AdministratorAccess"),
					KubernetesUsername: "admin",
				},
				{
					PrincipalARN:       api.MustParseARN("arn:aws:iam::012345689:role/Viewer"),
					KubernetesUsername: "viewer",
				},
			},
			existingStacks: []string{"eksctl-web-accessentry-WMWIOJ7RIE7SKMRBEO5HI6YVVUKNBV4I"},
			expectedAccessEntries: []api.AccessEntry{
				{
					PrincipalARN:       api.MustParseARN("arn:aws:iam::012345689:role/Viewer"),
					KubernetesUsername: "viewer",
				},
			},
		}),

		Entry("no entry exists", accessEntryTest{
			clusterName: "production",
			accessEntries: []api.AccessEntry{
				{
					PrincipalARN:       api.MustParseARN("arn:aws:iam::012345689:role/AdministratorAccess"),
					KubernetesUsername: "admin",
				},
				{
					PrincipalARN:       api.MustParseARN("arn:aws:iam::012345689:role/Viewer"),
					KubernetesUsername: "viewer",
				},
			},

			expectedAccessEntries: []api.AccessEntry{
				{
					PrincipalARN:       api.MustParseARN("arn:aws:iam::012345689:role/AdministratorAccess"),
					KubernetesUsername: "admin",
				},
				{
					PrincipalARN:       api.MustParseARN("arn:aws:iam::012345689:role/Viewer"),
					KubernetesUsername: "viewer",
				},
			},
		}),

		Entry("all entries exist", accessEntryTest{
			clusterName: "staging",
			accessEntries: []api.AccessEntry{
				{
					PrincipalARN:       api.MustParseARN("arn:aws:iam::012345689:role/AdministratorAccess"),
					KubernetesUsername: "admin",
				},
				{
					PrincipalARN:       api.MustParseARN("arn:aws:iam::012345689:role/Viewer"),
					KubernetesUsername: "viewer",
				},
			},

			existingStacks: []string{"eksctl-staging-accessentry-WMWIOJ7RIE7SKMRBEO5HI6YVVUKNBV4I", "eksctl-staging-accessentry-ZPCUBSOXPMTW5RRIV4YCDOYWDKIO4CMV"},
		}),
	)
})
