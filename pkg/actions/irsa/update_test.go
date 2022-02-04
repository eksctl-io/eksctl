package irsa_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/irsa"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
)

var _ = Describe("Update", func() {

	var (
		irsaManager      *irsa.Manager
		oidc             *iamoidc.OpenIDConnectManager
		fakeStackManager *fakes.FakeStackManager
		serviceAccount   []*api.ClusterIAMServiceAccount
	)

	BeforeEach(func() {
		serviceAccount = []*api.ClusterIAMServiceAccount{
			{
				ClusterIAMMeta: api.ClusterIAMMeta{
					Name:      "test-sa",
					Namespace: "default",
				},
				RoleName:         "test-role",
				AttachPolicyARNs: []string{"arn-123"},
			},
		}
		var err error

		fakeStackManager = new(fakes.FakeStackManager)

		oidc, err = iamoidc.NewOpenIDConnectManager(nil, "456123987123", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws", nil)
		Expect(err).NotTo(HaveOccurred())
		oidc.ProviderARN = "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"
		irsaManager = irsa.New("my-cluster", fakeStackManager, oidc, nil)
	})

	When("the IAMServiceAccount exists", func() {
		It("updates the role", func() {
			fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{
				{
					StackName: aws.String("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"),
					Outputs: []*cloudformation.Output{
						{
							OutputKey:   aws.String(outputs.IAMServiceAccountRoleName),
							OutputValue: aws.String("arn:aws:iam::123456789111:role/test-role"),
						},
					},
				},
			}, nil)

			err := irsaManager.UpdateIAMServiceAccounts(serviceAccount, false)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.ListStacksMatchingCallCount()).To(Equal(1))
			Expect(fakeStackManager.ListStacksMatchingArgsForCall(0)).To(Equal("eksctl-.*-addon-iamserviceaccount"))
			Expect(fakeStackManager.UpdateStackCallCount()).To(Equal(1))
			fakeStackManager.UpdateStackArgsForCall(0)
			options := fakeStackManager.UpdateStackArgsForCall(0)
			Expect(options.StackName).To(Equal("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"))
			Expect(options.ChangeSetName).To(ContainSubstring("updating-policy"))
			Expect(options.Description).To(Equal("updating policies for IAMServiceAccount default/test-sa"))
			Expect(options.Wait).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
			Expect(string(options.TemplateData.(manager.TemplateBody))).To(ContainSubstring("arn-123"))
			Expect(string(options.TemplateData.(manager.TemplateBody))).To(ContainSubstring(":sub\":\"system:serviceaccount:default:test-sa"))
			Expect(string(options.TemplateData.(manager.TemplateBody))).To(ContainSubstring("\"RoleName\":\"test-role\""))
		})

		When("in plan mode", func() {
			It("does not trigger an update", func() {
				fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{
					{
						StackName: aws.String("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"),
						Outputs: []*cloudformation.Output{
							{
								OutputKey:   aws.String(outputs.IAMServiceAccountRoleName),
								OutputValue: aws.String("arn:aws:iam::123456789111:role/test-role"),
							},
						},
					},
				}, nil)

				err := irsaManager.UpdateIAMServiceAccounts(serviceAccount, true)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeStackManager.ListStacksMatchingCallCount()).To(Equal(1))
				Expect(fakeStackManager.ListStacksMatchingArgsForCall(0)).To(Equal("eksctl-.*-addon-iamserviceaccount"))
				Expect(fakeStackManager.UpdateStackCallCount()).To(Equal(0))
			})
		})

		When("the role does not exist for a service account", func() {
			It("errors", func() {
				fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{
					{
						StackName: aws.String("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"),
					},
				}, nil)

				err := irsaManager.UpdateIAMServiceAccounts(serviceAccount, false)
				Expect(err).To(MatchError(ContainSubstring("failed to find role name service account")))
			})
		})

		When("the service account doesn't exist", func() {
			It("errors", func() {
				fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{}, nil)

				err := irsaManager.UpdateIAMServiceAccounts(serviceAccount, false)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeStackManager.UpdateStackCallCount()).To(BeZero())
			})
		})

		When("the role arn has an invalid format", func() {
			It("errors", func() {
				fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{
					{
						StackName: aws.String("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"),
						Outputs: []*cloudformation.Output{
							{
								OutputKey:   aws.String(outputs.IAMServiceAccountRoleName),
								OutputValue: aws.String("invalid-arn"),
							},
						},
					},
				}, nil)

				err := irsaManager.UpdateIAMServiceAccounts(serviceAccount, false)
				Expect(err).To(MatchError(ContainSubstring("failed to parse role arn \"invalid-arn\"")))
			})
		})

		When("the role is missing from the arn", func() {
			It("errors", func() {
				fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{
					{
						StackName: aws.String("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"),
						Outputs: []*cloudformation.Output{
							{
								OutputKey:   aws.String(outputs.IAMServiceAccountRoleName),
								OutputValue: aws.String("arn:aws:iam::123456789111:asdf"),
							},
						},
					},
				}, nil)

				err := irsaManager.UpdateIAMServiceAccounts(serviceAccount, false)
				Expect(err).To(MatchError(ContainSubstring("failed to parse resource: asdf")))
			})
		})
	})
})
