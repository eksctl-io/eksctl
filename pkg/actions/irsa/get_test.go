package irsa_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/irsa"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
)

var _ = Describe("Get", func() {

	var (
		irsaManager      *irsa.Manager
		fakeStackManager *fakes.FakeStackManager
	)

	BeforeEach(func() {
		fakeStackManager = new(fakes.FakeStackManager)

		irsaManager = irsa.New("my-cluster", fakeStackManager, nil, nil)
	})

	When("no options are specified", func() {
		It("returns all service accounts", func() {
			fakeStackManager.GetIAMServiceAccountsReturns([]*api.ClusterIAMServiceAccount{
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Name:      "test-sa",
						Namespace: "default",
					},
					AttachPolicyARNs: []string{"arn-123"},
				},
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Name:      "test-sa-2",
						Namespace: "not-default",
					},
					AttachPolicyARNs: []string{"arn-123"},
				},
			}, nil)

			serviceAccounts, err := irsaManager.Get(irsa.GetOptions{})
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.GetIAMServiceAccountsCallCount()).To(Equal(1))
			Expect(serviceAccounts).To(Equal([]*api.ClusterIAMServiceAccount{
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Name:      "test-sa",
						Namespace: "default",
					},
					AttachPolicyARNs: []string{"arn-123"},
				},
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Name:      "test-sa-2",
						Namespace: "not-default",
					},
					AttachPolicyARNs: []string{"arn-123"},
				},
			}))
		})
	})

	When("name option is specified", func() {
		It("returns only the service account matching the name", func() {
			fakeStackManager.GetIAMServiceAccountsReturns([]*api.ClusterIAMServiceAccount{
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Name:      "test-sa",
						Namespace: "default",
					},
					AttachPolicyARNs: []string{"arn-123"},
				},
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Name:      "test-sa-2",
						Namespace: "not-default",
					},
					AttachPolicyARNs: []string{"arn-123"},
				},
			}, nil)

			serviceAccounts, err := irsaManager.Get(irsa.GetOptions{Name: "test-sa"})
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.GetIAMServiceAccountsCallCount()).To(Equal(1))
			Expect(serviceAccounts).To(Equal([]*api.ClusterIAMServiceAccount{
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Name:      "test-sa",
						Namespace: "default",
					},
					AttachPolicyARNs: []string{"arn-123"},
				},
			}))
		})
	})

	When("namespace option is specified", func() {
		It("returns only the service account matching the name", func() {
			fakeStackManager.GetIAMServiceAccountsReturns([]*api.ClusterIAMServiceAccount{
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Name:      "test-sa",
						Namespace: "default",
					},
					AttachPolicyARNs: []string{"arn-123"},
				},
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Name:      "test-sa-2",
						Namespace: "not-default",
					},
					AttachPolicyARNs: []string{"arn-123"},
				},
			}, nil)

			serviceAccounts, err := irsaManager.Get(irsa.GetOptions{Namespace: "not-default"})
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.GetIAMServiceAccountsCallCount()).To(Equal(1))
			Expect(serviceAccounts).To(Equal([]*api.ClusterIAMServiceAccount{
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Name:      "test-sa-2",
						Namespace: "not-default",
					},
					AttachPolicyARNs: []string{"arn-123"},
				},
			}))
		})
	})

	When("name and namespace option is specified", func() {
		It("returns only the service account matching the name", func() {
			fakeStackManager.GetIAMServiceAccountsReturns([]*api.ClusterIAMServiceAccount{
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Name:      "test-sa",
						Namespace: "default",
					},
					AttachPolicyARNs: []string{"arn-123"},
				},
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Name:      "some-other-sa",
						Namespace: "default",
					},
					AttachPolicyARNs: []string{"arn-123"},
				},
			}, nil)

			serviceAccounts, err := irsaManager.Get(irsa.GetOptions{Namespace: "default", Name: "test-sa"})
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.GetIAMServiceAccountsCallCount()).To(Equal(1))
			Expect(serviceAccounts).To(Equal([]*api.ClusterIAMServiceAccount{
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Name:      "test-sa",
						Namespace: "default",
					},
					AttachPolicyARNs: []string{"arn-123"},
				},
			}))
		})
	})
})
