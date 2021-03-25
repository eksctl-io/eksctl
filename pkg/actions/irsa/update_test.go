package irsa_test

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	"github.com/weaveworks/eksctl/pkg/actions/irsa"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var _ = Describe("Update", func() {

	var (
		irsaManager      *irsa.Manager
		oidc             *iamoidc.OpenIDConnectManager
		fakeStackManager *fakes.FakeStackManager
		serviceAccounts  []*api.ClusterIAMServiceAccount
		serviceAccount   api.ClusterIAMServiceAccount
		clientSet        *fake.Clientset
	)

	BeforeEach(func() {
		serviceAccount = api.ClusterIAMServiceAccount{
			ClusterIAMMeta: api.ClusterIAMMeta{
				Name:      "test-sa",
				Namespace: "default",
				Labels:    map[string]string{"foo": "bar"},
			},
			AttachPolicyARNs: []string{"arn-123"},
		}
		serviceAccounts = []*api.ClusterIAMServiceAccount{
			&serviceAccount,
		}
		var err error

		fakeStackManager = new(fakes.FakeStackManager)
		clientSet = fake.NewSimpleClientset()

		oidc, err = iamoidc.NewOpenIDConnectManager(nil, "456123987123", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws")
		Expect(err).ToNot(HaveOccurred())
		oidc.ProviderARN = "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"
		irsaManager = irsa.New("my-cluster", fakeStackManager, oidc, clientSet)
	})

	Context("UpdateIAMServiceAccounts", func() {
		When("the IAMServiceAccount exists", func() {
			It("updates the role", func() {
				fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{
					{
						StackName: aws.String("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"),
					},
				}, nil)

				err := irsaManager.UpdateIAMServiceAccounts(serviceAccounts, false)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeStackManager.ListStacksMatchingCallCount()).To(Equal(1))
				Expect(fakeStackManager.ListStacksMatchingArgsForCall(0)).To(Equal("eksctl-.*-addon-iamserviceaccount"))
				Expect(fakeStackManager.UpdateStackCallCount()).To(Equal(1))
				fakeStackManager.UpdateStackArgsForCall(0)
				stackName, changeSetName, description, templateData, _ := fakeStackManager.UpdateStackArgsForCall(0)
				Expect(stackName).To(Equal("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"))
				Expect(changeSetName).To(Equal("updating-policy"))
				Expect(description).To(Equal("updating policies for IAMServiceAccount default/test-sa"))
				Expect(err).NotTo(HaveOccurred())
				Expect(string(templateData.(manager.TemplateBody))).To(ContainSubstring("arn-123"))
				Expect(string(templateData.(manager.TemplateBody))).To(ContainSubstring(":sub\":\"system:serviceaccount:default:test-sa"))
			})

			When("in plan mode", func() {
				It("does not trigger an update", func() {
					fakeStackManager.ListStacksMatchingReturns([]*cloudformation.Stack{
						{
							StackName: aws.String("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"),
						},
					}, nil)

					err := irsaManager.UpdateIAMServiceAccounts(serviceAccounts, true)
					Expect(err).NotTo(HaveOccurred())

					Expect(fakeStackManager.ListStacksMatchingCallCount()).To(Equal(1))
					Expect(fakeStackManager.ListStacksMatchingArgsForCall(0)).To(Equal("eksctl-.*-addon-iamserviceaccount"))
					Expect(fakeStackManager.UpdateStackCallCount()).To(Equal(0))
				})
			})
		})
	})

	Context("IsUpToDate", func() {
		ctx := context.Background()
		const template = `{"AWSTemplateFormatVersion":"2010-09-09","Description":"IAM role for serviceaccount \"default/test-sa\" [created and managed by eksctl]","Resources":{"Role1":{"Type":"AWS::IAM::Role","Properties":{"AssumeRolePolicyDocument":{"Statement":[{"Action":["sts:AssumeRoleWithWebIdentity"],"Condition":{"StringEquals":{"oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E:aud":"sts.amazonaws.com","oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E:sub":"system:serviceaccount:default:test-sa"}},"Effect":"Allow","Principal":{"Federated":"arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"}}],"Version":"2012-10-17"},"ManagedPolicyArns":["arn-123"]}}},"Outputs":{"Role1":{"Value":{"Fn::GetAtt":"Role1.Arn"}}}}`
		var stack *manager.Stack
		BeforeEach(func() {
			_, err := clientSet.CoreV1().ServiceAccounts("default").Create(ctx, &corev1.ServiceAccount{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ServiceAccount",
					APIVersion: corev1.SchemeGroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace:   "default",
					Name:        "test-sa",
					Labels:      map[string]string{"foo": "bar"},
					Annotations: map[string]string{api.AnnotationEKSRoleARN: "arn123"},
				},
			}, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
			stack = &cloudformation.Stack{
				StackName: aws.String("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"),
				Outputs: []*cloudformation.Output{
					{
						OutputKey:   aws.String("Role1"),
						OutputValue: aws.String("arn123"),
					},
				},
			}
		})

		It("returns true when its up to date", func() {
			fakeStackManager.GetStackTemplateReturns(template, nil)
			updateNeeded, err := irsaManager.IsUpToDate(serviceAccount, stack)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.GetStackTemplateCallCount()).To(Equal(1))
			Expect(fakeStackManager.GetStackTemplateArgsForCall(0)).To(Equal("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"))
			Expect(updateNeeded).To(BeTrue())
		})

		It("returns false when the stack template has changed", func() {
			fakeStackManager.GetStackTemplateReturns(`{"new":"template"}`, nil)
			updateNeeded, err := irsaManager.IsUpToDate(serviceAccount, stack)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.GetStackTemplateCallCount()).To(Equal(1))
			Expect(fakeStackManager.GetStackTemplateArgsForCall(0)).To(Equal("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"))
			Expect(updateNeeded).To(BeFalse())
		})

		It("returns false when the service account doesn't exist", func() {
			err := clientSet.CoreV1().ServiceAccounts("default").Delete(ctx, "test-sa", metav1.DeleteOptions{})
			Expect(err).NotTo(HaveOccurred())
			stack := &cloudformation.Stack{
				StackName: aws.String("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"),
			}

			fakeStackManager.GetStackTemplateReturns(template, nil)
			updateNeeded, err := irsaManager.IsUpToDate(serviceAccount, stack)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.GetStackTemplateCallCount()).To(Equal(1))
			Expect(fakeStackManager.GetStackTemplateArgsForCall(0)).To(Equal("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"))
			Expect(updateNeeded).To(BeFalse())
		})

		It("returns false when the service account labels are out of date", func() {
			serviceAccount.Labels = map[string]string{"new": "label"}
			fakeStackManager.GetStackTemplateReturns(template, nil)
			updateNeeded, err := irsaManager.IsUpToDate(serviceAccount, stack)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.GetStackTemplateCallCount()).To(Equal(1))
			Expect(fakeStackManager.GetStackTemplateArgsForCall(0)).To(Equal("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"))
			Expect(updateNeeded).To(BeFalse())
		})

		It("returns false when the service account annotations are out of date", func() {
			stack.Outputs = []*cloudformation.Output{}
			fakeStackManager.GetStackTemplateReturns(template, nil)
			updateNeeded, err := irsaManager.IsUpToDate(serviceAccount, stack)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.GetStackTemplateCallCount()).To(Equal(1))
			Expect(fakeStackManager.GetStackTemplateArgsForCall(0)).To(Equal("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"))
			Expect(updateNeeded).To(BeFalse())
		})
	})
})
