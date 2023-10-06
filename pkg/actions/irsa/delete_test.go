package irsa_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/irsa"
	"github.com/weaveworks/eksctl/pkg/actions/irsa/fakes"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

var _ = Describe("Delete", func() {

	Describe("DeleteIAMServiceAccountsTasks", func() {
		var (
			remover          *irsa.Remover
			fakeStackManager *fakes.FakeStackManager
			stackName1       = "eksctl-test-cluster-addon-iamserviceaccount-default-sa"
			stackName2       = "eksctl-test-cluster-addon-iamserviceaccount-kube-system-sa"
		)

		BeforeEach(func() {
			fakeStackManager = &fakes.FakeStackManager{}
			remover = irsa.NewRemover(&kubernetes.CachedClientSet{}, fakeStackManager)
		})

		When("DescribeIAMServiceAccountStacks fails", func() {
			It("returns an error", func() {
				fakeStackManager.DescribeIAMServiceAccountStacksReturns(nil, fmt.Errorf("foo"))
				_, err := remover.DeleteIAMServiceAccountsTasks(context.Background(), []string{}, false)
				Expect(err).To(MatchError(ContainSubstring("failed to describe IAM Service Account CFN Stacks")))
			})
		})

		When("there is no IAM Service Account stack", func() {
			It("returns an empty tasktree", func() {
				fakeStackManager.DescribeIAMServiceAccountStacksReturns([]*cfntypes.Stack{}, nil)
				taskTree, err := remover.DeleteIAMServiceAccountsTasks(context.Background(), []string{}, false)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(taskTree.Tasks)).To(Equal(0))
			})
		})

		When("there are multiple IAM Service Account stacks", func() {
			When("there is a stack with invalid name string", func() {
				It("returns an error", func() {
					fakeStackManager.DescribeIAMServiceAccountStacksReturns([]*cfntypes.Stack{
						{
							StackName: &stackName1,
						},
					}, nil)
					_, err := remover.DeleteIAMServiceAccountsTasks(context.Background(), []string{"invalid-name"}, false)
					Expect(err).To(MatchError(ContainSubstring("unexpected serviceaccount name format")))
				})
			})

			When("all stacks have valid names", func() {
				It("returns a tasktree with all expected tasks", func() {
					fakeStackManager.DescribeIAMServiceAccountStacksReturns([]*cfntypes.Stack{
						{
							StackName: &stackName1,
							Tags: []cfntypes.Tag{
								{
									Key:   aws.String(api.IAMServiceAccountNameTag),
									Value: aws.String("default/sa"),
								},
							},
						},
						{
							StackName: &stackName2,
							Tags: []cfntypes.Tag{
								{
									Key:   aws.String(api.IAMServiceAccountNameTag),
									Value: aws.String("kube-system/sa"),
								},
							},
						},
					}, nil)
					taskTree, err := remover.DeleteIAMServiceAccountsTasks(context.Background(), []string{
						"default/sa",
						"kube-system/sa",
					}, false)
					Expect(err).ToNot(HaveOccurred())
					Expect(taskTree.Parallel).To(Equal(true))
					Expect(len(taskTree.Tasks)).To(Equal(2))
					Expect(taskTree.Tasks[0].Describe()).To(ContainSubstring("2 sequential sub-tasks"))
					Expect(taskTree.Tasks[0].Describe()).To(ContainSubstring("delete IAM role for serviceaccount"))
					Expect(taskTree.Tasks[0].Describe()).To(ContainSubstring("delete serviceaccount"))
					Expect(taskTree.Tasks[1].Describe()).To(ContainSubstring("2 sequential sub-tasks"))
					Expect(taskTree.Tasks[1].Describe()).To(ContainSubstring("delete IAM role for serviceaccount"))
					Expect(taskTree.Tasks[1].Describe()).To(ContainSubstring("delete serviceaccount"))
				})
			})
		})
	})
})
