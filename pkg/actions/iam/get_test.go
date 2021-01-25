package iam_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/iam"
	"github.com/weaveworks/eksctl/pkg/actions/iam/fakes"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Get", func() {

	var (
		iamManager       *iam.Manager
		fakeStackManager *fakes.FakeStackManager
	)

	BeforeEach(func() {
		fakeStackManager = new(fakes.FakeStackManager)

		iamManager = iam.New("my-cluster", fakeStackManager, nil, nil)
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

			serviceAccounts, err := iamManager.Get("", "")
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

			serviceAccounts, err := iamManager.Get("", "test-sa")
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

			serviceAccounts, err := iamManager.Get("not-default", "")
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

			serviceAccounts, err := iamManager.Get("default", "test-sa")
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
