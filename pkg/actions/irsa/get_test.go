package irsa_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
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

	When("no error occurs", func() {
		It("returns service accounts from GetIAMServiceAccounts", func() {
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

			serviceAccounts, err := irsaManager.Get(context.Background(), irsa.GetOptions{})
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
})
