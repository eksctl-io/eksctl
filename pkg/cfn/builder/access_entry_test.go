package builder_test

import (
	"os"
	"path"

	. "github.com/benjamintf1/unmarshalledmatchers"
	. "github.com/onsi/ginkgo/v2"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"

	. "github.com/onsi/gomega"
)

var _ = Describe("Access Entry", func() {
	type accessEntryCase struct {
		clusterName string
		accessEntry api.AccessEntry

		resourceFilename string
	}

	DescribeTable("access entry resources", func(a accessEntryCase) {
		resourceSet := builder.NewAccessEntryResourceSet(a.clusterName, a.accessEntry)
		Expect(resourceSet.AddAllResources()).To(Succeed())
		actual, err := resourceSet.RenderJSON()
		Expect(err).NotTo(HaveOccurred())
		expected, err := os.ReadFile(path.Join("testdata", "access_entry", a.resourceFilename))
		Expect(err).NotTo(HaveOccurred())
		Expect(actual).To(MatchOrderedJSON(expected, WithUnorderedListKeys("Tags")))
	},
		Entry("only principalARN set", accessEntryCase{
			clusterName: "cluster",
			accessEntry: api.AccessEntry{
				PrincipalARN: api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
			},
			resourceFilename: "1.json",
		}),

		Entry("principalARN, groups and username set", accessEntryCase{
			clusterName: "cluster",
			accessEntry: api.AccessEntry{
				PrincipalARN:       api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				KubernetesGroups:   []string{"authenticated", "viewers"},
				KubernetesUsername: "user1",
			},
			resourceFilename: "2.json",
		}),

		Entry("policies set", accessEntryCase{
			clusterName: "cluster",
			accessEntry: api.AccessEntry{
				PrincipalARN:       api.MustParseARN("arn:aws:iam::111122223333:role/role-1"),
				KubernetesGroups:   []string{"viewers"},
				KubernetesUsername: "user1",
				AccessPolicies: []api.AccessPolicy{
					{
						PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSViewPolicy"),
						AccessScope: api.AccessScope{
							Type:       "namespace",
							Namespaces: []string{"kube-system", "default"},
						},
					},
					{
						PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSAdminPolicy"),
						AccessScope: api.AccessScope{
							Type: "cluster",
						},
					},
				},
			},
			resourceFilename: "3.json",
		}),
	)
})
