package cluster_test

import (
	"time"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"

	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/weaveworks/eksctl/pkg/utils/strings"

	"github.com/weaveworks/eksctl/pkg/utils/tasks"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/actions/cluster"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

type drainerMockOwned struct {
	mock.Mock
}

func (drainer *drainerMockOwned) Drain(input *nodegroup.DrainInput) error {
	args := drainer.Called(input)
	return args.Error(0)
}

var _ = Describe("Delete", func() {
	var (
		clusterName              string
		p                        *mockprovider.MockProvider
		cfg                      *api.ClusterConfig
		fakeStackManager         *fakes.FakeStackManager
		ranDeleteDeprecatedTasks bool
		ranDeleteClusterTasks    bool
		ctl                      *eks.ClusterProvider
		fakeClientSet            *fake.Clientset
	)

	BeforeEach(func() {
		clusterName = "my-cluster"
		p = mockprovider.NewMockProvider()
		cfg = api.NewClusterConfig()
		cfg.Metadata.Name = clusterName
		fakeStackManager = new(fakes.FakeStackManager)
		ranDeleteDeprecatedTasks = false
		ranDeleteClusterTasks = false
		ctl = &eks.ClusterProvider{Provider: p, Status: &eks.ProviderStatus{}}
	})

	Context("when the cluster is operable", func() {
		It("deletes the cluster", func() {
			//mocks are in order of being called
			p.MockEKS().On("DescribeCluster", mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
				Expect(*input.Name).To(Equal(clusterName))
				return true
			})).Return(&awseks.DescribeClusterOutput{
				Cluster: testutils.NewFakeCluster(clusterName, awseks.ClusterStatusActive),
			}, nil)

			p.MockEKS().On("ListFargateProfiles", &awseks.ListFargateProfilesInput{
				ClusterName: strings.Pointer(clusterName),
			}).Once().Return(&awseks.ListFargateProfilesOutput{FargateProfileNames: aws.StringSlice([]string{"fargate-1"})}, nil)

			p.MockEKS().On("DeleteFargateProfile", &awseks.DeleteFargateProfileInput{
				ClusterName:        aws.String(clusterName),
				FargateProfileName: aws.String("fargate-1"),
			}).Once().Return(&awseks.DeleteFargateProfileOutput{}, nil)

			p.MockEKS().On("ListFargateProfiles", &awseks.ListFargateProfilesInput{
				ClusterName: strings.Pointer(clusterName),
			}).Once().Return(&awseks.ListFargateProfilesOutput{}, nil)

			fakeStackManager.DeleteTasksForDeprecatedStacksReturns(&tasks.TaskTree{
				Tasks: []tasks.Task{&tasks.GenericTask{Doer: func() error {
					ranDeleteDeprecatedTasks = true
					return nil
				}}},
			}, nil)

			p.MockEC2().On("DescribeKeyPairs", mock.Anything).Return(&ec2.DescribeKeyPairsOutput{}, nil)

			p.MockEC2().On("DescribeSecurityGroupsWithContext", mock.Anything, mock.Anything).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

			fakeStackManager.NewTasksToDeleteClusterWithNodeGroupsReturns(&tasks.TaskTree{
				Tasks: []tasks.Task{&tasks.GenericTask{Doer: func() error {
					ranDeleteClusterTasks = true
					return nil
				}}},
			}, nil)

			karpenterStack := &manager.Stack{
				StackName: aws.String("karpenter"),
			}

			fakeStackManager.GetKarpenterStackReturns(karpenterStack, nil)

			c := cluster.NewOwnedCluster(cfg, ctl, nil, fakeStackManager)
			fakeClientSet = fake.NewSimpleClientset()

			c.SetNewClientSet(func() (kubernetes.Interface, error) {
				return fakeClientSet, nil
			})

			err := c.Delete(time.Microsecond, false, false, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.DeleteTasksForDeprecatedStacksCallCount()).To(Equal(1))
			Expect(ranDeleteDeprecatedTasks).To(BeTrue())
			Expect(fakeStackManager.NewTasksToDeleteClusterWithNodeGroupsCallCount()).To(Equal(1))
			Expect(ranDeleteClusterTasks).To(BeTrue())
			Expect(fakeStackManager.DeleteStackByNameSyncCallCount()).To(Equal(1))
			Expect(fakeStackManager.DeleteStackByNameSyncArgsForCall(0)).To(Equal("karpenter"))
		})

		When("force flag is set to true", func() {
			It("ignoring nodes draining error", func() {
				//mocks are in order of being called
				p.MockEKS().On("DescribeCluster", mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
					Expect(*input.Name).To(Equal(clusterName))
					return true
				})).Return(&awseks.DescribeClusterOutput{
					Cluster: testutils.NewFakeCluster(clusterName, awseks.ClusterStatusActive),
				}, nil)

				p.MockEKS().On("ListFargateProfiles", &awseks.ListFargateProfilesInput{
					ClusterName: strings.Pointer(clusterName),
				}).Once().Return(&awseks.ListFargateProfilesOutput{FargateProfileNames: aws.StringSlice([]string{})}, nil)

				p.MockEKS().On("DeleteFargateProfile", &awseks.DeleteFargateProfileInput{
					ClusterName:        aws.String(clusterName),
					FargateProfileName: aws.String("fargate-1"),
				}).Once().Return(&awseks.DeleteFargateProfileOutput{}, nil)

				p.MockEKS().On("ListFargateProfiles", &awseks.ListFargateProfilesInput{
					ClusterName: strings.Pointer(clusterName),
				}).Once().Return(&awseks.ListFargateProfilesOutput{}, nil)

				fakeStackManager.DeleteTasksForDeprecatedStacksReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{},
				}, nil)

				fakeStackManager.ListNodeGroupStacksReturns([]manager.NodeGroupStack{{NodeGroupName: "ng-1"}}, nil)

				p.MockEC2().On("DescribeKeyPairs", mock.Anything).Return(&ec2.DescribeKeyPairsOutput{}, nil)

				p.MockEC2().On("DescribeSecurityGroupsWithContext", mock.Anything, mock.Anything).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

				fakeStackManager.NewTasksToDeleteClusterWithNodeGroupsReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{},
				}, nil)

				c := cluster.NewOwnedCluster(cfg, ctl, nil, fakeStackManager)
				fakeClientSet = fake.NewSimpleClientset()

				c.SetNewClientSet(func() (kubernetes.Interface, error) {
					return fakeClientSet, nil
				})

				mockedDrainInput := &nodegroup.DrainInput{
					NodeGroups:     cmdutils.ToKubeNodeGroups(cfg),
					MaxGracePeriod: ctl.Provider.WaitTimeout(),
				}

				mockedDrainer := &drainerMockOwned{}
				mockedDrainer.On("Drain", mockedDrainInput).Return(errors.New("Mocked error"))
				c.SetNewNodeGroupManager(func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) cluster.NodeGroupDrainer {
					return mockedDrainer
				})

				err := c.Delete(time.Microsecond, false, true, false)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeStackManager.DeleteTasksForDeprecatedStacksCallCount()).To(Equal(1))
				Expect(ranDeleteDeprecatedTasks).To(BeFalse())
				Expect(fakeStackManager.NewTasksToDeleteClusterWithNodeGroupsCallCount()).To(Equal(1))
				Expect(ranDeleteClusterTasks).To(BeFalse())
				mockedDrainer.AssertNumberOfCalls(GinkgoT(), "Drain", 1)
			})
		})

		When("force flag is set to false", func() {
			It("nodes draining error thrown", func() {
				//mocks are in order of being called
				p.MockEKS().On("DescribeCluster", mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
					Expect(*input.Name).To(Equal(clusterName))
					return true
				})).Return(&awseks.DescribeClusterOutput{
					Cluster: testutils.NewFakeCluster(clusterName, awseks.ClusterStatusActive),
				}, nil)

				p.MockEKS().On("ListFargateProfiles", &awseks.ListFargateProfilesInput{
					ClusterName: strings.Pointer(clusterName),
				}).Once().Return(&awseks.ListFargateProfilesOutput{FargateProfileNames: aws.StringSlice([]string{})}, nil)

				p.MockEKS().On("DeleteFargateProfile", &awseks.DeleteFargateProfileInput{
					ClusterName:        aws.String(clusterName),
					FargateProfileName: aws.String("fargate-1"),
				}).Once().Return(&awseks.DeleteFargateProfileOutput{}, nil)

				p.MockEKS().On("ListFargateProfiles", &awseks.ListFargateProfilesInput{
					ClusterName: strings.Pointer(clusterName),
				}).Once().Return(&awseks.ListFargateProfilesOutput{}, nil)

				fakeStackManager.DeleteTasksForDeprecatedStacksReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{},
				}, nil)
				fakeStackManager.ListNodeGroupStacksReturns([]manager.NodeGroupStack{{NodeGroupName: "ng-1"}}, nil)

				p.MockEC2().On("DescribeKeyPairs", mock.Anything).Return(&ec2.DescribeKeyPairsOutput{}, nil)

				p.MockEC2().On("DescribeSecurityGroupsWithContext", mock.Anything, mock.Anything).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

				fakeStackManager.NewTasksToDeleteClusterWithNodeGroupsReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{},
				}, nil)

				c := cluster.NewOwnedCluster(cfg, ctl, nil, fakeStackManager)
				fakeClientSet = fake.NewSimpleClientset()

				c.SetNewClientSet(func() (kubernetes.Interface, error) {
					return fakeClientSet, nil
				})

				mockedDrainInput := &nodegroup.DrainInput{
					NodeGroups:     cmdutils.ToKubeNodeGroups(cfg),
					MaxGracePeriod: ctl.Provider.WaitTimeout(),
				}

				errorMessage := "Mocked error"
				mockedDrainer := &drainerMockOwned{}
				mockedDrainer.On("Drain", mockedDrainInput).Return(errors.New(errorMessage))
				c.SetNewNodeGroupManager(func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) cluster.NodeGroupDrainer {
					return mockedDrainer
				})

				err := c.Delete(time.Microsecond, false, false, false)
				Expect(err).To(MatchError(errorMessage))
				Expect(fakeStackManager.DeleteTasksForDeprecatedStacksCallCount()).To(Equal(0))
				Expect(ranDeleteDeprecatedTasks).To(BeFalse())
				Expect(fakeStackManager.NewTasksToDeleteClusterWithNodeGroupsCallCount()).To(Equal(0))
				Expect(ranDeleteClusterTasks).To(BeFalse())
				mockedDrainer.AssertNumberOfCalls(GinkgoT(), "Drain", 1)
			})
		})
	})

	Context("when the cluster is inoperable", func() {
		It("deletes the cluster without trying to query kubernetes", func() {
			//mocks are in order of being called
			p.MockEKS().On("DescribeCluster", mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
				Expect(*input.Name).To(Equal(clusterName))
				return true
			})).Return(&awseks.DescribeClusterOutput{
				Cluster: testutils.NewFakeCluster(clusterName, awseks.ClusterStatusFailed),
			}, nil)

			fakeStackManager.DeleteTasksForDeprecatedStacksReturns(&tasks.TaskTree{
				Tasks: []tasks.Task{&tasks.GenericTask{Doer: func() error {
					ranDeleteDeprecatedTasks = true
					return nil
				}}},
			}, nil)

			p.MockEC2().On("DescribeKeyPairs", mock.Anything).Return(&ec2.DescribeKeyPairsOutput{}, nil)

			p.MockEC2().On("DescribeSecurityGroupsWithContext", mock.Anything, mock.Anything).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

			fakeStackManager.NewTasksToDeleteClusterWithNodeGroupsReturns(&tasks.TaskTree{
				Tasks: []tasks.Task{&tasks.GenericTask{Doer: func() error {
					ranDeleteClusterTasks = true
					return nil
				}}},
			}, nil)

			c := cluster.NewOwnedCluster(cfg, ctl, nil, fakeStackManager)

			err := c.Delete(time.Microsecond, false, false, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.DeleteTasksForDeprecatedStacksCallCount()).To(Equal(1))
			Expect(ranDeleteDeprecatedTasks).To(BeTrue())
			Expect(fakeStackManager.NewTasksToDeleteClusterWithNodeGroupsCallCount()).To(Equal(1))
			Expect(ranDeleteClusterTasks).To(BeTrue())
		})
	})
})
