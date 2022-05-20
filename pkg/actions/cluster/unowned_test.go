package cluster_test

import (
	"context"
	"time"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
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

type drainerMockUnowned struct {
	mock.Mock
}

func (d *drainerMockUnowned) Drain(ctx context.Context, input *nodegroup.DrainInput) error {
	args := d.Called(input)
	return args.Error(0)
}

var _ = Describe("Delete", func() {
	var (
		clusterName              string
		p                        *mockprovider.MockProvider
		cfg                      *api.ClusterConfig
		fakeStackManager         *fakes.FakeStackManager
		ranDeleteDeprecatedTasks bool
		ctl                      *eks.ClusterProvider
	)

	BeforeEach(func() {
		clusterName = "my-cluster"
		p = mockprovider.NewMockProvider()
		cfg = api.NewClusterConfig()
		cfg.Metadata.Name = clusterName
		fakeStackManager = new(fakes.FakeStackManager)
		ranDeleteDeprecatedTasks = false
		ctl = &eks.ClusterProvider{AWSProvider: p, Status: &eks.ProviderStatus{}}
	})

	Context("when the cluster is operable", func() {
		It("deletes the cluster", func() {
			//mocks are in order of being called
			p.MockEKS().On("DescribeCluster", mock.Anything, mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
				return *input.Name == clusterName
			})).Return(&awseks.DescribeClusterOutput{
				Cluster: testutils.NewFakeCluster(clusterName, ekstypes.ClusterStatusActive),
			}, nil)

			p.MockEKS().On("DeleteAddon", mock.Anything, &awseks.DeleteAddonInput{
				ClusterName: strings.Pointer(clusterName),
				AddonName:   strings.Pointer("vpc-cni"),
			}).Return(&awseks.DeleteAddonOutput{}, nil)

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

			fargateStackName := aws.String("eksctl-my-cluster-fargate")
			p.MockCloudFormation().On("DescribeStacks", mock.Anything, &cloudformation.DescribeStacksInput{
				StackName: fargateStackName,
			}).Return(&cloudformation.DescribeStacksOutput{
				Stacks: []types.Stack{
					{
						StackName: fargateStackName,
						Tags: []types.Tag{
							{
								Key:   aws.String("alpha.eksctl.io/cluster-name"),
								Value: aws.String(clusterName),
							},
						},
						StackStatus: types.StackStatusCreateComplete,
					},
				},
			}, nil)

			p.MockCloudFormation().On("DeleteStack", mock.Anything).Return(nil, nil)

			fakeStackManager.DeleteTasksForDeprecatedStacksReturns(&tasks.TaskTree{
				Tasks: []tasks.Task{&tasks.GenericTask{Doer: func() error {
					ranDeleteDeprecatedTasks = true
					return nil
				}}},
			}, nil)

			p.MockEC2().On("DescribeKeyPairs", mock.Anything, mock.Anything).Return(&ec2.DescribeKeyPairsOutput{}, nil)

			p.MockEC2().On("DescribeSecurityGroups", mock.Anything, mock.Anything).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

			fakeStackManager.GetFargateStackReturns(&types.Stack{StackName: aws.String("fargate-role")}, nil)
			fakeStackManager.DeleteStackBySpecReturns(nil, nil)

			p.MockEKS().On("ListNodegroups", mock.Anything, mock.Anything).Return(&awseks.ListNodegroupsOutput{
				Nodegroups: []string{"ng-1", "ng-2"},
			}, nil)

			fakeStackManager.ListNodeGroupStacksReturns([]manager.NodeGroupStack{{NodeGroupName: "ng-1"}}, nil)

			var deleteCallCount int
			fakeStackManager.NewTasksToDeleteNodeGroupsReturns(&tasks.TaskTree{
				Tasks: []tasks.Task{&tasks.GenericTask{Doer: func() error {
					deleteCallCount++
					return nil
				}}},
			}, nil)

			var unownedDeleteCallCount int
			fakeStackManager.NewTaskToDeleteUnownedNodeGroupReturns(&tasks.TaskTree{
				Tasks: []tasks.Task{&tasks.GenericTask{Doer: func() error {
					unownedDeleteCallCount++
					return nil
				}}},
			})

			p.MockEKS().On("DeleteNodegroup", mock.Anything, &awseks.DeleteNodegroupInput{ClusterName: &clusterName, NodegroupName: aws.String("ng-1")}).Return(&awseks.DeleteNodegroupOutput{}, nil)
			p.MockEKS().On("DeleteNodegroup", mock.Anything, &awseks.DeleteNodegroupInput{ClusterName: &clusterName, NodegroupName: aws.String("ng-2")}).Return(&awseks.DeleteNodegroupOutput{}, nil)

			p.MockEKS().On("DeleteCluster", mock.Anything, mock.Anything).Return(&awseks.DeleteClusterOutput{}, nil)
			c := cluster.NewUnownedCluster(cfg, ctl, fakeStackManager)
			fakeClientSet := fake.NewSimpleClientset()

			c.SetNewClientSet(func() (kubernetes.Interface, error) {
				return fakeClientSet, nil
			})

			err := c.Delete(context.Background(), time.Microsecond, time.Second*0, false, false, false, 1)
			Expect(err).NotTo(HaveOccurred())
			Expect(deleteCallCount).To(Equal(1))
			Expect(unownedDeleteCallCount).To(Equal(1))
			Expect(fakeStackManager.DeleteTasksForDeprecatedStacksCallCount()).To(Equal(1))
			Expect(ranDeleteDeprecatedTasks).To(BeTrue())
			Expect(fakeStackManager.DeleteStackBySpecCallCount()).To(Equal(1))
			_, stack := fakeStackManager.DeleteStackBySpecArgsForCall(0)
			Expect(*stack.StackName).To(Equal("fargate-role"))
		})

		When("force flag is set to true", func() {
			It("ignoring nodes draining error", func() {
				//mocks are in order of being called
				p.MockEKS().On("DescribeCluster", mock.Anything, mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
					return *input.Name == clusterName
				})).Return(&awseks.DescribeClusterOutput{
					Cluster: testutils.NewFakeCluster(clusterName, ekstypes.ClusterStatusActive),
				}, nil)

				p.MockEKS().On("DeleteAddon", mock.Anything, &awseks.DeleteAddonInput{
					ClusterName: strings.Pointer(clusterName),
					AddonName:   strings.Pointer("vpc-cni"),
				}).Return(&awseks.DeleteAddonOutput{}, nil)

				p.MockEKS().On("ListFargateProfiles", mock.Anything, &awseks.ListFargateProfilesInput{
					ClusterName: strings.Pointer(clusterName),
				}).Once().Return(&awseks.ListFargateProfilesOutput{}, nil)

				p.MockEKS().On("DeleteFargateProfile", mock.Anything, &awseks.DeleteFargateProfileInput{
					ClusterName:        aws.String(clusterName),
					FargateProfileName: aws.String("fargate-1"),
				}).Once().Return(&awseks.DeleteFargateProfileOutput{}, nil)

				p.MockEKS().On("ListFargateProfiles", mock.Anything, &awseks.ListFargateProfilesInput{
					ClusterName: strings.Pointer(clusterName),
				}).Once().Return(&awseks.ListFargateProfilesOutput{}, nil)

				fargateStackName := aws.String("eksctl-my-cluster-fargate")
				p.MockCloudFormation().On("DescribeStacks", &cloudformation.DescribeStacksInput{
					StackName: fargateStackName,
				}).Return(&cloudformation.DescribeStacksOutput{}, nil)

				p.MockCloudFormation().On("DeleteStack", mock.Anything).Return(nil, nil)

				fakeStackManager.DeleteTasksForDeprecatedStacksReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{},
				}, nil)

				p.MockEC2().On("DescribeKeyPairs", mock.Anything, mock.Anything).Return(&ec2.DescribeKeyPairsOutput{}, nil)

				p.MockEC2().On("DescribeSecurityGroups", mock.Anything, mock.Anything).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

				fakeStackManager.GetFargateStackReturns(nil, nil)
				fakeStackManager.DeleteStackBySpecReturns(nil, nil)

				p.MockEKS().On("ListNodegroups", mock.Anything, mock.Anything).Return(&awseks.ListNodegroupsOutput{
					Nodegroups: []string{"ng-1", "ng-2"},
				}, nil)

				fakeStackManager.ListNodeGroupStacksReturns([]manager.NodeGroupStack{{NodeGroupName: "ng-1"}}, nil)

				var deleteCallCount int
				fakeStackManager.NewTasksToDeleteNodeGroupsReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{},
				}, nil)

				var unownedDeleteCallCount int
				fakeStackManager.NewTaskToDeleteUnownedNodeGroupReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{},
				})

				p.MockEKS().On("DeleteNodegroup", mock.Anything, nil).Return(&awseks.DeleteNodegroupOutput{}, nil)
				p.MockEKS().On("DeleteNodegroup", mock.Anything, nil).Return(&awseks.DeleteNodegroupOutput{}, nil)

				p.MockEKS().On("DeleteCluster", mock.Anything, mock.Anything).Return(&awseks.DeleteClusterOutput{}, nil)
				ctl.Status = &eks.ProviderStatus{
					ClusterInfo: &eks.ClusterInfo{
						Cluster: &ekstypes.Cluster{
							Status:  ekstypes.ClusterStatusActive,
							Version: aws.String("1.21"),
						},
					},
				}
				c := cluster.NewUnownedCluster(cfg, ctl, fakeStackManager)
				fakeClientSet := fake.NewSimpleClientset()

				c.SetNewClientSet(func() (kubernetes.Interface, error) {
					return fakeClientSet, nil
				})

				mockedDrainInput := &nodegroup.DrainInput{
					NodeGroups:     cmdutils.ToKubeNodeGroups(cfg),
					MaxGracePeriod: ctl.AWSProvider.WaitTimeout(),
					Parallel:       1,
				}

				mockedDrainer := &drainerMockUnowned{}
				mockedDrainer.On("Drain", mockedDrainInput).Return(errors.New("Mocked error"))
				c.SetNewNodeGroupManager(func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) cluster.NodeGroupDrainer {
					return mockedDrainer
				})

				err := c.Delete(context.Background(), time.Microsecond, time.Second*0, false, true, false, 1)
				Expect(err).NotTo(HaveOccurred())
				Expect(deleteCallCount).To(Equal(0))
				Expect(unownedDeleteCallCount).To(Equal(0))
				Expect(fakeStackManager.DeleteTasksForDeprecatedStacksCallCount()).To(Equal(1))
				Expect(ranDeleteDeprecatedTasks).To(BeFalse())
				Expect(fakeStackManager.DeleteStackBySpecCallCount()).To(Equal(0))
				mockedDrainer.AssertNumberOfCalls(GinkgoT(), "Drain", 1)
			})
		})

		When("force flag is set to false", func() {
			It("nodes draining error thrown", func() {
				//mocks are in order of being called
				p.MockEKS().On("DescribeCluster", mock.Anything, mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
					return *input.Name == clusterName
				})).Return(&awseks.DescribeClusterOutput{
					Cluster: testutils.NewFakeCluster(clusterName, ekstypes.ClusterStatusActive),
				}, nil)

				p.MockEKS().On("DeleteAddon", mock.Anything, &awseks.DeleteAddonInput{
					ClusterName: strings.Pointer(clusterName),
					AddonName:   strings.Pointer("vpc-cni"),
				}).Return(&awseks.DeleteAddonOutput{}, nil)

				p.MockEKS().On("ListFargateProfiles", mock.Anything, &awseks.ListFargateProfilesInput{
					ClusterName: strings.Pointer(clusterName),
				}).Once().Return(&awseks.ListFargateProfilesOutput{}, nil)

				p.MockEKS().On("DeleteFargateProfile", mock.Anything, &awseks.DeleteFargateProfileInput{
					ClusterName:        aws.String(clusterName),
					FargateProfileName: aws.String("fargate-1"),
				}).Once().Return(&awseks.DeleteFargateProfileOutput{}, nil)

				p.MockEKS().On("ListFargateProfiles", mock.Anything, &awseks.ListFargateProfilesInput{
					ClusterName: strings.Pointer(clusterName),
				}).Once().Return(&awseks.ListFargateProfilesOutput{}, nil)

				fargateStackName := aws.String("eksctl-my-cluster-fargate")
				p.MockCloudFormation().On("DescribeStacks", &cloudformation.DescribeStacksInput{
					StackName: fargateStackName,
				}).Return(&cloudformation.DescribeStacksOutput{}, nil)

				p.MockCloudFormation().On("DeleteStack", mock.Anything).Return(nil, nil)

				fakeStackManager.DeleteTasksForDeprecatedStacksReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{},
				}, nil)

				p.MockEC2().On("DescribeKeyPairs", mock.Anything, mock.Anything).Return(&ec2.DescribeKeyPairsOutput{}, nil)

				p.MockEC2().On("DescribeSecurityGroups", mock.Anything, mock.Anything).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

				fakeStackManager.GetFargateStackReturns(nil, nil)
				fakeStackManager.DeleteStackBySpecReturns(nil, nil)

				p.MockEKS().On("ListNodegroups", mock.Anything, mock.Anything).Return(&awseks.ListNodegroupsOutput{
					Nodegroups: []string{"ng-1", "ng-2"},
				}, nil)

				fakeStackManager.ListNodeGroupStacksReturns([]manager.NodeGroupStack{{NodeGroupName: "ng-1"}}, nil)

				var deleteCallCount int
				fakeStackManager.NewTasksToDeleteNodeGroupsReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{},
				}, nil)

				var unownedDeleteCallCount int
				fakeStackManager.NewTaskToDeleteUnownedNodeGroupReturns(&tasks.TaskTree{
					Tasks: []tasks.Task{},
				})

				p.MockEKS().On("DeleteNodegroup", mock.Anything, nil).Return(&awseks.DeleteNodegroupOutput{}, nil)
				p.MockEKS().On("DeleteNodegroup", mock.Anything, nil).Return(&awseks.DeleteNodegroupOutput{}, nil)

				p.MockEKS().On("DeleteCluster", mock.Anything, mock.Anything).Return(&awseks.DeleteClusterOutput{}, nil)
				c := cluster.NewUnownedCluster(cfg, ctl, fakeStackManager)
				fakeClientSet := fake.NewSimpleClientset()

				c.SetNewClientSet(func() (kubernetes.Interface, error) {
					return fakeClientSet, nil
				})
				ctl.Status = &eks.ProviderStatus{
					ClusterInfo: &eks.ClusterInfo{
						Cluster: &ekstypes.Cluster{
							Status:  ekstypes.ClusterStatusActive,
							Version: aws.String("1.21"),
						},
					},
				}
				mockedDrainInput := &nodegroup.DrainInput{
					NodeGroups:     cmdutils.ToKubeNodeGroups(cfg),
					MaxGracePeriod: ctl.AWSProvider.WaitTimeout(),
					Parallel:       1,
				}

				errorMessage := "Mocked error"
				mockedDrainer := &drainerMockUnowned{}
				mockedDrainer.On("Drain", mockedDrainInput).Return(errors.New(errorMessage))
				c.SetNewNodeGroupManager(func(cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) cluster.NodeGroupDrainer {
					return mockedDrainer
				})

				err := c.Delete(context.Background(), time.Microsecond, time.Second*0, false, false, false, 1)
				Expect(err).To(MatchError(errorMessage))
				Expect(deleteCallCount).To(Equal(0))
				Expect(unownedDeleteCallCount).To(Equal(0))
				Expect(fakeStackManager.DeleteTasksForDeprecatedStacksCallCount()).To(Equal(0))
				Expect(ranDeleteDeprecatedTasks).To(BeFalse())
				Expect(fakeStackManager.DeleteStackBySpecCallCount()).To(Equal(0))
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

			fakeStackManager.DeleteTasksForDeprecatedStacksReturns(&tasks.TaskTree{
				Tasks: []tasks.Task{&tasks.GenericTask{Doer: func() error {
					ranDeleteDeprecatedTasks = true
					return nil
				}}},
			}, nil)

			p.MockEC2().On("DescribeKeyPairs", mock.Anything, mock.Anything).Return(&ec2.DescribeKeyPairsOutput{}, nil)

			p.MockEC2().On("DescribeSecurityGroups", mock.Anything, mock.Anything).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

			p.MockEKS().On("ListNodegroups", mock.Anything, mock.Anything).Return(&awseks.ListNodegroupsOutput{
				Nodegroups: []string{"ng-1", "ng-2"},
			}, nil)

			fakeStackManager.ListNodeGroupStacksReturns([]manager.NodeGroupStack{{NodeGroupName: "ng-1"}}, nil)

			var deleteCallCount int
			fakeStackManager.NewTasksToDeleteNodeGroupsReturns(&tasks.TaskTree{
				Tasks: []tasks.Task{&tasks.GenericTask{Doer: func() error {
					deleteCallCount++
					return nil
				}}},
			}, nil)

			var unownedDeleteCallCount int
			fakeStackManager.NewTaskToDeleteUnownedNodeGroupReturns(&tasks.TaskTree{
				Tasks: []tasks.Task{&tasks.GenericTask{Doer: func() error {
					unownedDeleteCallCount++
					return nil
				}}},
			})

			p.MockEKS().On("DeleteNodegroup", mock.Anything, mock.MatchedBy(func(input *awseks.DeleteNodegroupInput) bool {
				Expect(*input.ClusterName).To(Equal(clusterName))
				Expect(*input.NodegroupName).To(Equal("ng-1"))
				return true
			})).Return(&awseks.DeleteNodegroupOutput{}, nil)

			p.MockEKS().On("DeleteCluster", mock.Anything, mock.Anything).Return(&awseks.DeleteClusterOutput{}, nil)

			c := cluster.NewUnownedCluster(cfg, ctl, fakeStackManager)
			err := c.Delete(context.Background(), time.Microsecond, time.Second*0, false, false, false, 1)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.DeleteTasksForDeprecatedStacksCallCount()).To(Equal(1))
			Expect(deleteCallCount).To(Equal(1))
			Expect(unownedDeleteCallCount).To(Equal(1))
			Expect(ranDeleteDeprecatedTasks).To(BeTrue())
		})
	})
})
