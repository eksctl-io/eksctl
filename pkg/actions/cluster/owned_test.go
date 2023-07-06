package cluster_test

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/weaveworks/eksctl/pkg/actions/cluster"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type drainerMockOwned struct {
	mock.Mock
}

func (d *drainerMockOwned) Drain(ctx context.Context, input *nodegroup.DrainInput) error {
	args := d.Called(input)
	return args.Error(0)
}

var _ = Describe("Delete", func() {
	var (
		clusterName              string
		ctx                      context.Context
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
		ctx = context.Background()
		p = mockprovider.NewMockProvider()
		cfg = api.NewClusterConfig()
		cfg.Metadata.Name = clusterName
		fakeStackManager = new(fakes.FakeStackManager)
		ranDeleteDeprecatedTasks = false
		ranDeleteClusterTasks = false
		ctl = &eks.ClusterProvider{AWSProvider: p, Status: &eks.ProviderStatus{
			ClusterInfo: &eks.ClusterInfo{
				Cluster: testutils.NewFakeCluster(clusterName, ekstypes.ClusterStatusActive),
			},
		}}
	})

	Context("when the cluster is operable", func() {
		It("deletes the cluster", func() {
			//mocks are in order of being called
			p.MockEKS().On("DescribeCluster", mock.Anything, mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
				Expect(*input.Name).To(Equal(clusterName))
				return true
			})).Return(&awseks.DescribeClusterOutput{
				Cluster: testutils.NewFakeCluster(clusterName, ekstypes.ClusterStatusActive),
			}, nil)

			p.MockEKS().On("ListFargateProfiles", mock.Anything, &awseks.ListFargateProfilesInput{
				ClusterName: strings.Pointer(clusterName),
			}).Once().Return(&awseks.ListFargateProfilesOutput{FargateProfileNames: []string{"fargate-1"}}, nil)

			p.MockEKS().On("DeleteFargateProfile", mock.Anything, &awseks.DeleteFargateProfileInput{
				ClusterName:        aws.String(clusterName),
				FargateProfileName: aws.String("fargate-1"),
			}).Once().Return(&awseks.DeleteFargateProfileOutput{}, nil)

			p.MockEKS().On("ListFargateProfiles", mock.Anything, &awseks.ListFargateProfilesInput{
				ClusterName: strings.Pointer(clusterName),
			}).Once().Return(&awseks.ListFargateProfilesOutput{}, nil)

			fakeStackManager.DeleteTasksForDeprecatedStacksReturns(&tasks.TaskTree{
				Tasks: []tasks.Task{&tasks.GenericTask{Doer: func() error {
					ranDeleteDeprecatedTasks = true
					return nil
				}}},
			}, nil)

			p.MockEC2().On("DescribeKeyPairs", mock.Anything, mock.Anything).Return(&ec2.DescribeKeyPairsOutput{}, nil)

			p.MockEC2().On("DescribeSecurityGroups", mock.Anything, mock.Anything).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

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

			c, err := cluster.NewOwnedCluster(ctx, cfg, ctl, nil, fakeStackManager)
			Expect(err).NotTo(HaveOccurred())
			fakeClientSet = fake.NewSimpleClientset()

			c.SetNewClientSet(func() (kubernetes.Interface, error) {
				return fakeClientSet, nil
			})
			c.SetNewNodeGroupManager(func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) cluster.NodeGroupDrainer {
				mockedDrainer := &drainerMockOwned{}
				mockedDrainer.On("Drain", mock.Anything).Return(nil)
				return mockedDrainer
			})

			err = c.Delete(context.Background(), time.Microsecond, 0, false, false, false, 1)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.DeleteTasksForDeprecatedStacksCallCount()).To(Equal(1))
			Expect(ranDeleteDeprecatedTasks).To(BeTrue())
			Expect(fakeStackManager.NewTasksToDeleteClusterWithNodeGroupsCallCount()).To(Equal(1))
			Expect(ranDeleteClusterTasks).To(BeTrue())
			Expect(fakeStackManager.DeleteStackSyncCallCount()).To(Equal(1))
			_, stack := fakeStackManager.DeleteStackSyncArgsForCall(0)
			Expect(*stack.StackName).To(Equal("karpenter"))
		})

		When("force flag is set to true", func() {
			It("ignoring nodes draining error", func() {
				ctl.Status = &eks.ProviderStatus{
					ClusterInfo: &eks.ClusterInfo{
						Cluster: &ekstypes.Cluster{
							Status:  ekstypes.ClusterStatusActive,
							Version: aws.String("1.22"),
						},
					},
				}
				//mocks are in order of being called
				p.MockEKS().On("DescribeCluster", mock.Anything, mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
					Expect(*input.Name).To(Equal(clusterName))
					return true
				})).Return(&awseks.DescribeClusterOutput{
					Cluster: testutils.NewFakeCluster(clusterName, ekstypes.ClusterStatusActive),
				}, nil)

				p.MockEKS().On("ListFargateProfiles", mock.Anything, &awseks.ListFargateProfilesInput{
					ClusterName: strings.Pointer(clusterName),
				}).Once().Return(&awseks.ListFargateProfilesOutput{FargateProfileNames: []string{}}, nil)

				p.MockEKS().On("DeleteFargateProfile", mock.Anything, &awseks.DeleteFargateProfileInput{
					ClusterName:        aws.String(clusterName),
					FargateProfileName: aws.String("fargate-1"),
				}).Once().Return(&awseks.DeleteFargateProfileOutput{}, nil)

				p.MockEKS().On("ListFargateProfiles", mock.Anything, &awseks.ListFargateProfilesInput{
					ClusterName: strings.Pointer(clusterName),
				}).Once().Return(&awseks.ListFargateProfilesOutput{}, nil)

				fakeStackManager.DeleteTasksForDeprecatedStacksReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{},
				}, nil)

				fakeStackManager.ListNodeGroupStacksWithStatusesReturns([]manager.NodeGroupStack{{NodeGroupName: "ng-1"}}, nil)

				p.MockEC2().On("DescribeKeyPairs", mock.Anything, mock.Anything).Return(&ec2.DescribeKeyPairsOutput{}, nil)

				p.MockEC2().On("DescribeSecurityGroups", mock.Anything, mock.Anything).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

				fakeStackManager.NewTasksToDeleteClusterWithNodeGroupsReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{},
				}, nil)

				c, err := cluster.NewOwnedCluster(ctx, cfg, ctl, nil, fakeStackManager)
				Expect(err).NotTo(HaveOccurred())
				fakeClientSet = fake.NewSimpleClientset()

				c.SetNewClientSet(func() (kubernetes.Interface, error) {
					return fakeClientSet, nil
				})

				mockedDrainInput := &nodegroup.DrainInput{
					NodeGroups:     cmdutils.ToKubeNodeGroups(cfg),
					MaxGracePeriod: ctl.AWSProvider.WaitTimeout(),
					Parallel:       1,
				}

				mockedDrainer := &drainerMockOwned{}
				mockedDrainer.On("Drain", mockedDrainInput).Return(errors.New("Mocked error"))
				c.SetNewNodeGroupManager(func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) cluster.NodeGroupDrainer {
					return mockedDrainer
				})

				err = c.Delete(context.Background(), time.Microsecond, 0, false, true, false, 1)
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
				p.MockEKS().On("DescribeCluster", mock.Anything, mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
					Expect(*input.Name).To(Equal(clusterName))
					return true
				})).Return(&awseks.DescribeClusterOutput{
					Cluster: testutils.NewFakeCluster(clusterName, ekstypes.ClusterStatusActive),
				}, nil)

				p.MockEKS().On("ListFargateProfiles", mock.Anything, &awseks.ListFargateProfilesInput{
					ClusterName: strings.Pointer(clusterName),
				}).Once().Return(&awseks.ListFargateProfilesOutput{FargateProfileNames: []string{}}, nil)

				p.MockEKS().On("DeleteFargateProfile", mock.Anything, &awseks.DeleteFargateProfileInput{
					ClusterName:        aws.String(clusterName),
					FargateProfileName: aws.String("fargate-1"),
				}).Once().Return(&awseks.DeleteFargateProfileOutput{}, nil)

				p.MockEKS().On("ListFargateProfiles", mock.Anything, &awseks.ListFargateProfilesInput{
					ClusterName: strings.Pointer(clusterName),
				}).Once().Return(&awseks.ListFargateProfilesOutput{}, nil)

				fakeStackManager.DeleteTasksForDeprecatedStacksReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{},
				}, nil)
				fakeStackManager.ListNodeGroupStacksWithStatusesReturns([]manager.NodeGroupStack{{NodeGroupName: "ng-1"}}, nil)

				p.MockEC2().On("DescribeKeyPairs", mock.Anything, mock.Anything).Return(&ec2.DescribeKeyPairsOutput{}, nil)

				p.MockEC2().On("DescribeSecurityGroups", mock.Anything, mock.Anything).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

				fakeStackManager.NewTasksToDeleteClusterWithNodeGroupsReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{},
				}, nil)

				c, err := cluster.NewOwnedCluster(ctx, cfg, ctl, nil, fakeStackManager)
				Expect(err).NotTo(HaveOccurred())
				fakeClientSet = fake.NewSimpleClientset()

				c.SetNewClientSet(func() (kubernetes.Interface, error) {
					return fakeClientSet, nil
				})

				mockedDrainInput := &nodegroup.DrainInput{
					NodeGroups:     cmdutils.ToKubeNodeGroups(cfg),
					MaxGracePeriod: ctl.AWSProvider.WaitTimeout(),
					Parallel:       1,
				}
				ctl.Status = &eks.ProviderStatus{
					ClusterInfo: &eks.ClusterInfo{
						Cluster: &ekstypes.Cluster{
							Status:  ekstypes.ClusterStatusActive,
							Version: aws.String("1.22"),
						},
					},
				}
				errorMessage := "Mocked error"
				mockedDrainer := &drainerMockOwned{}
				mockedDrainer.On("Drain", mockedDrainInput).Return(errors.New(errorMessage))
				c.SetNewNodeGroupManager(func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) cluster.NodeGroupDrainer {
					return mockedDrainer
				})

				err = c.Delete(context.Background(), time.Microsecond, time.Second*0, false, false, false, 1)
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
			p.MockEKS().On("DescribeCluster", mock.Anything, mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
				Expect(*input.Name).To(Equal(clusterName))
				return true
			})).Return(&awseks.DescribeClusterOutput{
				Cluster: testutils.NewFakeCluster(clusterName, ekstypes.ClusterStatusFailed),
			}, nil)

			p.MockEKS().On("ListFargateProfiles", mock.Anything, &awseks.ListFargateProfilesInput{
				ClusterName: strings.Pointer(clusterName),
			}).Once().Return(&awseks.ListFargateProfilesOutput{FargateProfileNames: []string{}}, nil)

			fakeStackManager.DeleteTasksForDeprecatedStacksReturns(&tasks.TaskTree{
				Tasks: []tasks.Task{&tasks.GenericTask{Doer: func() error {
					ranDeleteDeprecatedTasks = true
					return nil
				}}},
			}, nil)

			p.MockEC2().On("DescribeKeyPairs", mock.Anything, mock.Anything).Return(&ec2.DescribeKeyPairsOutput{}, nil)

			p.MockEC2().On("DescribeSecurityGroups", mock.Anything, mock.Anything).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

			fakeStackManager.NewTasksToDeleteClusterWithNodeGroupsReturns(&tasks.TaskTree{
				Tasks: []tasks.Task{&tasks.GenericTask{Doer: func() error {
					ranDeleteClusterTasks = true
					return nil
				}}},
			}, nil)

			c, err := cluster.NewOwnedCluster(ctx, cfg, ctl, nil, fakeStackManager)
			Expect(err).NotTo(HaveOccurred())
			c.SetNewNodeGroupManager(func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) cluster.NodeGroupDrainer {
				mockedDrainer := &drainerMockOwned{}
				mockedDrainer.On("Drain", mock.Anything).Return(nil)
				return mockedDrainer
			})
			c.SetNewClientSet(func() (kubernetes.Interface, error) {
				return fake.NewSimpleClientset(), nil
			})

			err = c.Delete(context.Background(), time.Microsecond, time.Second*0, false, false, false, 1)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.DeleteTasksForDeprecatedStacksCallCount()).To(Equal(1))
			Expect(ranDeleteDeprecatedTasks).To(BeTrue())
			Expect(fakeStackManager.NewTasksToDeleteClusterWithNodeGroupsCallCount()).To(Equal(1))
			Expect(ranDeleteClusterTasks).To(BeTrue())
		})
	})
})
