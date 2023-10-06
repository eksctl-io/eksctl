package irsa_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/irsa"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Create", func() {

	Describe("CreateIAMServiceAccountsTasks", func() {
		var (
			clusterName = "test-cluster"
			roleName    = "test-role-name"
			roleARN     = "test-role-arn"
			region      = "us-west-2"
			creator     *irsa.Creator
		)

		BeforeEach(func() {
			creator = irsa.NewCreator(clusterName, region, nil, nil, nil)
		})

		When("attachRoleARN is provided and RoleOnly is true", func() {
			It("returns a empty tasktree", func() {
				serviceAccounts := []*api.ClusterIAMServiceAccount{
					{
						RoleName:      roleName,
						AttachRoleARN: roleARN,
						RoleOnly:      aws.Bool(true),
					},
				}
				taskTree := creator.CreateIAMServiceAccountsTasks(context.Background(), serviceAccounts)
				Expect(taskTree.Parallel).To(Equal(true))
				Expect(taskTree.IsSubTask).To(Equal(false))
				Expect(len(taskTree.Tasks)).To(Equal(1))
				Expect(taskTree.Tasks[0].Describe()).To(Equal("no tasks"))
			})
		})

		When("attachRoleARN is provided and RoleOnly is false", func() {
			It("returns a tasktree with all expected tasks", func() {
				serviceAccounts := []*api.ClusterIAMServiceAccount{
					{
						RoleName:      roleName,
						AttachRoleARN: roleARN,
					},
				}
				taskTree := creator.CreateIAMServiceAccountsTasks(context.Background(), serviceAccounts)
				Expect(taskTree.Parallel).To(Equal(true))
				Expect(taskTree.IsSubTask).To(Equal(false))
				Expect(len(taskTree.Tasks)).To(Equal(1))
				Expect(taskTree.Tasks[0].Describe()).To(ContainSubstring("create serviceaccount"))
			})
		})

		When("attachRoleARN is not provided and RoleOnly is true", func() {
			It("returns a tasktree with all expected tasks", func() {
				serviceAccounts := []*api.ClusterIAMServiceAccount{
					{
						RoleName: roleName,
						RoleOnly: aws.Bool(true),
					},
				}
				taskTree := creator.CreateIAMServiceAccountsTasks(context.Background(), serviceAccounts)
				Expect(taskTree.Parallel).To(Equal(true))
				Expect(taskTree.IsSubTask).To(Equal(false))
				Expect(len(taskTree.Tasks)).To(Equal(1))
				Expect(taskTree.Tasks[0].Describe()).To(ContainSubstring("create IAM role for serviceaccount"))
			})
		})

		When("attachRoleARN is not provided and RoleOnly is false", func() {
			It("returns a tasktree with all expected tasks", func() {
				serviceAccounts := []*api.ClusterIAMServiceAccount{
					{
						RoleName: roleName,
						RoleOnly: aws.Bool(false),
					},
				}
				taskTree := creator.CreateIAMServiceAccountsTasks(context.Background(), serviceAccounts)
				Expect(taskTree.Parallel).To(Equal(true))
				Expect(taskTree.IsSubTask).To(Equal(false))
				Expect(len(taskTree.Tasks)).To(Equal(1))
				Expect(taskTree.Tasks[0].Describe()).To(ContainSubstring("2 sequential sub-tasks"))
				Expect(taskTree.Tasks[0].Describe()).To(ContainSubstring("create IAM role for serviceaccount"))
				Expect(taskTree.Tasks[0].Describe()).To(ContainSubstring("create serviceaccount"))
			})
		})
	})
})
