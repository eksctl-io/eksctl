package filter

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/weaveworks/eksctl/pkg/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("iamserviceaccount filter", func() {

	Context("Match", func() {
		var (
			filter    *IAMServiceAccountFilter
			cfg       *api.ClusterConfig
			clientSet *fake.Clientset
		)

		BeforeEach(func() {
			filter = NewIAMServiceAccountFilter()
			cfg = api.NewClusterConfig()
			cfg.IAM.WithOIDC = api.Enabled()
			cfg.IAM.ServiceAccounts = serviceAccounts(
				"dev1",
				"dev2",
				"dev3",
				"test1",
				"test2",
				"test3",
			)
			clientSet = fake.NewSimpleClientset()
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
				"kube-system/aws-node",
			)
			err := filter.SetDeleteFilter(context.Background(), mockLister, true, cfg)
			Expect(err).NotTo(HaveOccurred())

			included, excluded := filter.MatchAll(cfg.IAM.ServiceAccounts)
			Expect(included).To(HaveLen(2))
			Expect(sets.List(included)).To(ConsistOf(
				"sa/only-remote-1",
				"sa/only-remote-2"),
			)
			Expect(sets.List(excluded)).To(ConsistOf(
				"sa/dev1",
				"sa/dev2",
				"sa/dev3",
				"sa/test1",
				"sa/test2",
				"sa/test3",
				"kube-system/aws-node",
			))
		})

		It("exclude existing stacks works correctly", func() {
			mockLister := newMockServiceAccountLister(
				"sa/dev1",
				"sa/dev2",
				"sa/dev3",
				"sa/only-remote-1",
				"sa/only-remote-2",
			)
			err := filter.SetExcludeExistingFilter(context.Background(), mockLister, clientSet, cfg.IAM.ServiceAccounts, true)
			Expect(err).NotTo(HaveOccurred())

			included, excluded := filter.MatchAll(cfg.IAM.ServiceAccounts)
			Expect(included).To(HaveLen(3))
			Expect(included.HasAll(
				"sa/test1",
				"sa/test2",
				"sa/test3",
			)).To(BeTrue())
			Expect(excluded).To(HaveLen(3))
			Expect(excluded.HasAll(
				"sa/dev1",
				"sa/dev2",
				"sa/dev3",
			)).To(BeTrue())
		})

		It("exclude existing service accounts works correctly", func() {
			mockLister := newMockServiceAccountLister(
				"sa/dev1",
				"sa/dev2",
				"sa/dev3",
			)
			cfg.IAM.ServiceAccounts = append(cfg.IAM.ServiceAccounts, &api.ClusterIAMServiceAccount{
				ClusterIAMMeta: api.ClusterIAMMeta{
					Namespace: "sa",
					Name:      "role-only",
				},
				RoleOnly: api.Enabled(),
			})
			sa1 := metav1.ObjectMeta{Name: "test1", Namespace: "sa"}
			sa2 := metav1.ObjectMeta{Name: "role-only", Namespace: "sa"}

			err := kubernetes.MaybeCreateServiceAccountOrUpdateMetadata(clientSet, sa1)
			Expect(err).NotTo(HaveOccurred())
			err = kubernetes.MaybeCreateServiceAccountOrUpdateMetadata(clientSet, sa2)
			Expect(err).NotTo(HaveOccurred())

			err = filter.SetExcludeExistingFilter(context.Background(), mockLister, clientSet, cfg.IAM.ServiceAccounts, false)
			Expect(err).NotTo(HaveOccurred())

			included, excluded := filter.MatchAll(cfg.IAM.ServiceAccounts)
			Expect(included).To(HaveLen(3))
			Expect(included.HasAll(
				"sa/test2",
				"sa/test3",
				"sa/role-only",
			)).To(BeTrue())
			Expect(excluded).To(HaveLen(4))
			Expect(excluded.HasAll(
				"sa/test1",
				"sa/dev1",
				"sa/dev2",
				"sa/dev3",
			)).To(BeTrue())
		})
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

func (l *mockSALister) ListIAMServiceAccountStacks(_ context.Context) ([]string, error) {
	return l.result, nil
}

func newMockServiceAccountLister(serviceAccounts ...string) *mockSALister {
	return &mockSALister{
		result: serviceAccounts,
	}
}
