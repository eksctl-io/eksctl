package cluster_test

import (
	"time"

	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/weaveworks/eksctl/pkg/utils/strings"

	"github.com/weaveworks/eksctl/pkg/utils/tasks"

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
		ctl = &eks.ClusterProvider{Provider: p, Status: &eks.ProviderStatus{}}
	})

	Context("when the cluster is operable", func() {
		It("deletes the cluster", func() {
			//mocks are in order of being called
			p.MockEKS().On("DescribeCluster", mock.MatchedBy(func(input *awseks.DescribeClusterInput) bool {
				return *input.Name == clusterName
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

			fargateStackName := aws.String("eksctl-my-cluster-fargate")
			p.MockCloudFormation().On("DescribeStacks", &cloudformation.DescribeStacksInput{
				StackName: fargateStackName,
			}).Return(&cloudformation.DescribeStacksOutput{
				Stacks: []*cloudformation.Stack{
					{
						StackName: fargateStackName,
						Tags: []*cloudformation.Tag{
							{
								Key:   aws.String("alpha.eksctl.io/cluster-name"),
								Value: aws.String(clusterName),
							},
						},
						StackStatus: aws.String(cloudformation.StackStatusCreateComplete),
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

			p.MockEC2().On("DescribeKeyPairs", mock.Anything).Return(&ec2.DescribeKeyPairsOutput{}, nil)

			p.MockEC2().On("DescribeSecurityGroupsWithContext", mock.Anything, mock.Anything).Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

			fakeStackManager.GetFargateStackReturns(&cloudformation.Stack{StackName: aws.String("fargate-role")}, nil)
			fakeStackManager.DeleteStackByNameReturns(nil, nil)

			p.MockEKS().On("ListNodegroupsPages", mock.MatchedBy(func(input *awseks.ListNodegroupsInput) bool {
				Expect(*input.ClusterName).To(Equal(clusterName))
				return true
			}), mock.Anything).Run(func(args mock.Arguments) {
				consume := args[1].(func(ng *awseks.ListNodegroupsOutput, _ bool) bool)
				out := &awseks.ListNodegroupsOutput{
					Nodegroups: aws.StringSlice([]string{"ng-1"}),
				}
				cont := consume(out, true)
				if !cont {
					panic("unexpected return value from the paging function: shouldContinue was false which isn't expected in this test scenario")
				}
			}).Return(nil)

			p.MockEKS().On("ListNodegroups", mock.Anything).Return(&awseks.ListNodegroupsOutput{}, nil)

			p.MockEKS().On("DeleteNodegroup", mock.MatchedBy(func(input *awseks.DeleteNodegroupInput) bool {
				Expect(*input.ClusterName).To(Equal(clusterName))
				Expect(*input.NodegroupName).To(Equal("ng-1"))
				return true
			})).Return(&awseks.DeleteNodegroupOutput{}, nil)

			p.MockEKS().On("DeleteCluster", mock.Anything).Return(&awseks.DeleteClusterOutput{}, nil)
			c := cluster.NewUnownedCluster(cfg, ctl, fakeStackManager)
			fakeClientSet := fake.NewSimpleClientset()

			c.SetNewClientSet(func() (kubernetes.Interface, error) {
				return fakeClientSet, nil
			})

			err := c.Delete(time.Microsecond, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.DeleteTasksForDeprecatedStacksCallCount()).To(Equal(1))
			Expect(ranDeleteDeprecatedTasks).To(BeTrue())
			Expect(fakeStackManager.DeleteStackBySpecCallCount()).To(Equal(1))
			Expect(*fakeStackManager.DeleteStackBySpecArgsForCall(0).StackName).To(Equal("fargate-role"))
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

			p.MockEKS().On("ListNodegroupsPages", mock.MatchedBy(func(input *awseks.ListNodegroupsInput) bool {
				Expect(*input.ClusterName).To(Equal(clusterName))
				return true
			}), mock.Anything).Run(func(args mock.Arguments) {
				consume := args[1].(func(ng *awseks.ListNodegroupsOutput, _ bool) bool)
				out := &awseks.ListNodegroupsOutput{
					Nodegroups: aws.StringSlice([]string{"ng-1"}),
				}
				cont := consume(out, true)
				if !cont {
					panic("unexpected return value from the paging function: shouldContinue was false which isn't expected in this test scenario")
				}
			}).Return(nil)

			p.MockEKS().On("ListNodegroups", mock.Anything).Return(&awseks.ListNodegroupsOutput{}, nil)

			p.MockEKS().On("DeleteNodegroup", mock.MatchedBy(func(input *awseks.DeleteNodegroupInput) bool {
				Expect(*input.ClusterName).To(Equal(clusterName))
				Expect(*input.NodegroupName).To(Equal("ng-1"))
				return true
			})).Return(&awseks.DeleteNodegroupOutput{}, nil)

			p.MockEKS().On("DeleteCluster", mock.Anything).Return(&awseks.DeleteClusterOutput{}, nil)

			c := cluster.NewUnownedCluster(cfg, ctl, fakeStackManager)
			err := c.Delete(time.Microsecond, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeStackManager.DeleteTasksForDeprecatedStacksCallCount()).To(Equal(1))
			Expect(ranDeleteDeprecatedTasks).To(BeTrue())
		})
	})
})
