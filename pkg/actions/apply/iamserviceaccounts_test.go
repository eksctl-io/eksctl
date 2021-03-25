package apply_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/actions/apply"
	"github.com/weaveworks/eksctl/pkg/actions/apply/fakes"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	fakestack "github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

var _ = Describe("Iamserviceaccounts", func() {
	var (
		reconciler             *apply.Reconciler
		cfg                    *api.ClusterConfig
		mockProvider           *mockprovider.MockProvider
		fakeManager            *fakes.FakeIRSAManager
		fakeStackManager       *fakestack.FakeStackManager
		createTasksReturnValue *tasks.TaskTree
		updateTasksReturnValue *tasks.TaskTree
		deleteTasksReturnValue *tasks.TaskTree
	)
	BeforeEach(func() {
		cfg = api.NewClusterConfig()
		cfg.Metadata.Name = "mycluster"
		mockProvider = mockprovider.NewMockProvider()
		fakeManager = new(fakes.FakeIRSAManager)
		fakeStackManager = new(fakestack.FakeStackManager)

		reconciler = apply.New(cfg, &eks.ClusterProvider{Provider: mockProvider}, nil, nil, nil, false)
		reconciler.SetIRSAManager(fakeManager)
		reconciler.SetStackManager(fakeStackManager)

		createTasksReturnValue = &tasks.TaskTree{Tasks: []tasks.Task{&tasks.GenericTask{Description: "create"}}}
		updateTasksReturnValue = &tasks.TaskTree{Tasks: []tasks.Task{&tasks.GenericTask{Description: "update"}}}
		deleteTasksReturnValue = &tasks.TaskTree{Tasks: []tasks.Task{&tasks.GenericTask{Description: "delete"}}}

		fakeManager.CreateTasksReturns(createTasksReturnValue)
		fakeManager.UpdateTaskReturns(updateTasksReturnValue, nil)
		fakeManager.DeleteTasksReturns(deleteTasksReturnValue, nil)
	})

	When("an IAM account each needs creating, updating, no-op and deleting", func() {
		BeforeEach(func() {
			cfg.IAM.ServiceAccounts = []*api.ClusterIAMServiceAccount{
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Namespace: "create",
						Name:      "me",
					},
				},
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Namespace: "update",
						Name:      "me",
					},
				},
				{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Namespace: "up-to-date",
						Name:      "me",
					},
				},
			}

			fakeStackManager.DescribeIAMServiceAccountStacksReturns([]*manager.Stack{
				{
					StackName: aws.String(manager.MakeIAMServiceAccountStackName("mycluster", "delete", "me")),
					Tags: []*cloudformation.Tag{
						{
							Key:   aws.String(api.IAMServiceAccountNameTag),
							Value: aws.String("delete/me"),
						},
					},
				},
				{
					StackName: aws.String(manager.MakeIAMServiceAccountStackName("mycluster", "update", "me")),
					Tags: []*cloudformation.Tag{
						{
							Key:   aws.String(api.IAMServiceAccountNameTag),
							Value: aws.String("update/me"),
						},
					},
				},
				{
					StackName: aws.String(manager.MakeIAMServiceAccountStackName("mycluster", "up-to-date", "me")),
					Tags: []*cloudformation.Tag{
						{
							Key:   aws.String(api.IAMServiceAccountNameTag),
							Value: aws.String("up-to-date/me"),
						},
					},
				},
			}, nil)

			fakeManager.IsUpToDateStub = func(sa api.ClusterIAMServiceAccount, _ *cloudformation.Stack) (bool, error) {
				if sa.NameString() == "update/me" {
					return false, nil
				}

				if sa.NameString() == "up-to-date/me" {
					return true, nil
				}
				panic(fmt.Sprintf("isUpToDate called with incorrect SA %s", sa.NameString()))
			}
		})

		It("returns the correct tasks", func() {
			createTasks, updateTasks, deleteTasks, err := reconciler.ReconcileIAMServiceAccounts()
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.DescribeIAMServiceAccountStacksCallCount()).To(Equal(1))

			By("creating the missing SA")
			Expect(fakeManager.CreateTasksCallCount()).To(Equal(1))
			Expect(fakeManager.CreateTasksArgsForCall(0)).To(ConsistOf(
				&api.ClusterIAMServiceAccount{
					ClusterIAMMeta: api.ClusterIAMMeta{
						Namespace: "create",
						Name:      "me",
					},
				},
			))
			Expect(createTasks.Tasks).To(HaveLen(1))
			Expect(createTasks).To(Equal(createTasksReturnValue))

			By("updating only the out of date SA")
			Expect(fakeManager.IsUpToDateCallCount()).To(Equal(2))
			Expect(fakeManager.UpdateTaskCallCount()).To(Equal(1))
			sa, stack := fakeManager.UpdateTaskArgsForCall(0)
			Expect(sa).To(Equal(&api.ClusterIAMServiceAccount{
				ClusterIAMMeta: api.ClusterIAMMeta{
					Namespace: "update",
					Name:      "me",
				},
			}))
			Expect(stack).To(Equal(&manager.Stack{
				StackName: aws.String(manager.MakeIAMServiceAccountStackName("mycluster", "update", "me")),
				Tags: []*cloudformation.Tag{
					{
						Key:   aws.String(api.IAMServiceAccountNameTag),
						Value: aws.String("update/me"),
					},
				},
			}))
			Expect(updateTasks).To(Equal(&tasks.TaskTree{Tasks: []tasks.Task{updateTasksReturnValue}}))

			By("deleting the undesired SA")
			Expect(fakeManager.DeleteTasksCallCount()).To(Equal(1))
			Expect(fakeManager.DeleteTasksArgsForCall(0)).To(HaveKey("delete/me"))
			Expect(*fakeManager.DeleteTasksArgsForCall(0)["delete/me"]).To(Equal(manager.Stack{
				StackName: aws.String(manager.MakeIAMServiceAccountStackName("mycluster", "delete", "me")),
				Tags: []*cloudformation.Tag{
					{
						Key:   aws.String(api.IAMServiceAccountNameTag),
						Value: aws.String("delete/me"),
					},
				},
			}))
			Expect(deleteTasks.Tasks).To(HaveLen(1))
			Expect(deleteTasks).To(Equal(deleteTasksReturnValue))
		})

		When("the describe stack call fails", func() {
			BeforeEach(func() {
				fakeStackManager.DescribeIAMServiceAccountStacksReturns(nil, fmt.Errorf("foo"))
			})

			It("returns an error", func() {
				_, _, _, err := reconciler.ReconcileIAMServiceAccounts()
				Expect(err).To(MatchError("failed to discover existing service accounts: foo"))
			})
		})

		When("the isUpToDate call fails", func() {
			BeforeEach(func() {
				fakeManager.IsUpToDateReturns(false, fmt.Errorf("foo"))
			})

			It("returns an error", func() {
				_, _, _, err := reconciler.ReconcileIAMServiceAccounts()
				Expect(err).To(MatchError("failed to check if service account is up to date: foo"))
			})
		})

		When("the generate update tasks call fails", func() {
			BeforeEach(func() {
				fakeManager.UpdateTaskReturns(nil, fmt.Errorf("foo"))
			})

			It("returns an error", func() {
				_, _, _, err := reconciler.ReconcileIAMServiceAccounts()
				Expect(err).To(MatchError("failed to generate update tasks: foo"))
			})
		})

		When("the generate delete tasks call fails", func() {
			BeforeEach(func() {
				fakeManager.DeleteTasksReturns(nil, fmt.Errorf("foo"))
			})

			It("returns an error", func() {
				_, _, _, err := reconciler.ReconcileIAMServiceAccounts()
				Expect(err).To(MatchError("failed to generate delete tasks: foo"))
			})
		})
	})
})
