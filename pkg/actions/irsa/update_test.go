package irsa_test

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/irsa"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
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
			stacks := []*types.Stack{
				{
					StackName: aws.String("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"),
				},
			}
			err := irsaManager.UpdateIAMServiceAccounts(context.Background(), serviceAccount, stacks, false)
			Expect(err).NotTo(HaveOccurred())

			Expect(fakeStackManager.UpdateStackCallCount()).To(Equal(1))
			fakeStackManager.UpdateStackArgsForCall(0)
			_, options := fakeStackManager.UpdateStackArgsForCall(0)
			Expect(options.StackName).To(Equal("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"))
			Expect(options.ChangeSetName).To(ContainSubstring("updating-policy"))
			Expect(options.Description).To(Equal("updating policies for IAMServiceAccount default/test-sa"))
			Expect(options.Wait).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
			Expect(string(options.TemplateData.(manager.TemplateBody))).To(ContainSubstring("arn-123"))
			Expect(string(options.TemplateData.(manager.TemplateBody))).To(ContainSubstring(":sub\":\"system:serviceaccount:default:test-sa"))
		})

		When("in plan mode", func() {
			It("does not trigger an update", func() {
				stacks := []*types.Stack{
					{
						StackName: aws.String("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"),
					},
				}
				err := irsaManager.UpdateIAMServiceAccounts(context.Background(), serviceAccount, stacks, true)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeStackManager.UpdateStackCallCount()).To(Equal(0))
			})
		})

		When("the service account doesn't exist", func() {
			It("errors", func() {
				err := irsaManager.UpdateIAMServiceAccounts(context.Background(), serviceAccount, []*types.Stack{}, false)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeStackManager.UpdateStackCallCount()).To(BeZero())
			})
		})

		When("a custom role name was used during creation", func() {
			It("uses that role name", func() {
				stacks := []*types.Stack{
					{
						StackName: aws.String("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"),
					},
				}
				fakeStackManager.GetStackTemplateReturns(stackTemplateWithRoles, nil)

				err := irsaManager.UpdateIAMServiceAccounts(context.Background(), serviceAccount, stacks, false)
				Expect(err).NotTo(HaveOccurred())

				Expect(fakeStackManager.UpdateStackCallCount()).To(Equal(1))
				fakeStackManager.UpdateStackArgsForCall(0)
				_, options := fakeStackManager.UpdateStackArgsForCall(0)
				Expect(options.StackName).To(Equal("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"))
				Expect(options.ChangeSetName).To(ContainSubstring("updating-policy"))
				Expect(options.Description).To(Equal("updating policies for IAMServiceAccount default/test-sa"))
				Expect(options.Wait).To(BeTrue())
				Expect(err).NotTo(HaveOccurred())
				Expect(string(options.TemplateData.(manager.TemplateBody))).To(ContainSubstring("arn-123"))
				Expect(string(options.TemplateData.(manager.TemplateBody))).To(ContainSubstring(":sub\":\"system:serviceaccount:default:test-sa"))
				Expect(string(options.TemplateData.(manager.TemplateBody))).To(ContainSubstring("\"RoleName\":\"test-role\""))
			})
		})
		When("GetStackTemplate errors", func() {
			It("errors", func() {
				stacks := []*types.Stack{
					{
						StackName: aws.String("eksctl-my-cluster-addon-iamserviceaccount-default-test-sa"),
					},
				}
				fakeStackManager.GetStackTemplateReturns("", errors.New("nope"))

				err := irsaManager.UpdateIAMServiceAccounts(context.Background(), serviceAccount, stacks, false)
				Expect(err).To(MatchError(ContainSubstring("failed to get stack template: nope")))
			})
		})
	})
})

var stackTemplateWithRoles = `{
  "AWSTemplateFormatVersion": "2010-09-09",
  "Description": "IAM role for serviceaccount \"default/test-overwrite\" [created and managed by eksctl]",
  "Resources": {
    "Role1": {
      "Type": "AWS::IAM::Role",
      "Properties": {
        "AssumeRolePolicyDocument": {
          "Statement": [
            {
              "Action": [
                "sts:AssumeRoleWithWebIdentity"
              ],
              "Condition": {
                "StringEquals": {
                  "oidc.eks.us-west-2.amazonaws.com/id/761B7DDCE9618E9AE44143A538F85E8F:aud": "sts.amazonaws.com",
                  "oidc.eks.us-west-2.amazonaws.com/id/761B7DDCE9618E9AE44143A538F85E8F:sub": "system:serviceaccount:default:test-overwrite"
                }
              },
              "Effect": "Allow",
              "Principal": {
                "Federated": "arn:aws:iam::083751696308:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/761B7DDCE9618E9AE44143A538F85E8F"
              }
            }
          ],
          "Version": "2012-10-17"
        },
        "ManagedPolicyArns": [
          "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
        ],
        "RoleName": "test-role"
      }
    }
  },
  "Outputs": {
    "Role1": {
      "Value": {
        "Fn::GetAtt": "Role1.Arn"
      }
    }
  }
}
`
