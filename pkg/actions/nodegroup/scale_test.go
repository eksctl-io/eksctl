package nodegroup_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Scale", func() {
	var (
		clusterName, ngName string
		p                   *mockprovider.MockProvider
		cfg                 *api.ClusterConfig
		ng                  *api.NodeGroup
		m                   *nodegroup.Manager
		fakeStackManager    *fakes.FakeStackManager
	)
	BeforeEach(func() {
		clusterName = "my-cluster"
		ngName = "my-ng"
		p = mockprovider.NewMockProvider()
		cfg = api.NewClusterConfig()
		cfg.Metadata.Name = clusterName

		ng = &api.NodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name: ngName,
				ScalingConfig: &api.ScalingConfig{
					MinSize:         aws.Int(1),
					DesiredCapacity: aws.Int(3),
				},
			},
		}
		m = nodegroup.New(cfg, &eks.ClusterProvider{Provider: p}, nil)
		fakeStackManager = new(fakes.FakeStackManager)
		m.SetStackManager(fakeStackManager)
		p.MockCloudFormation().On("ListStacksPages", mock.Anything, mock.Anything).Return(nil, nil)
	})

	Describe("Managed NodeGroup", func() {
		BeforeEach(func() {
			nodegroups := make(map[string]manager.StackInfo)
			nodegroups["my-ng"] = manager.StackInfo{
				Stack: &manager.Stack{
					Tags: []*cloudformation.Tag{
						{
							Key:   aws.String(api.NodeGroupNameTag),
							Value: aws.String("my-ng"),
						},
						{
							Key:   aws.String(api.NodeGroupTypeTag),
							Value: aws.String(string(api.NodeGroupTypeManaged)),
						},
					},
				},
			}
			fakeStackManager.DescribeNodeGroupStacksAndResourcesReturns(nodegroups, nil)
		})

		It("scales the nodegroup using the values provided", func() {
			p.MockEKS().On("UpdateNodegroupConfig", &awseks.UpdateNodegroupConfigInput{
				ScalingConfig: &awseks.NodegroupScalingConfig{
					MinSize:     aws.Int64(1),
					DesiredSize: aws.Int64(3),
				},
				ClusterName:   &clusterName,
				NodegroupName: &ngName,
			}).Return(nil, nil)

			p.MockEKS().On("DescribeNodegroupRequest", &awseks.DescribeNodegroupInput{
				ClusterName:   &clusterName,
				NodegroupName: &ngName,
			}).Return(&request.Request{}, nil)

			waitCallCount := 0
			m.SetWaiter(func(name, msg string, acceptors []request.WaiterAcceptor, newRequest func() *request.Request, waitTimeout time.Duration, troubleshoot func(string) error) error {
				waitCallCount++
				return nil
			})

			err := m.Scale(ng)

			Expect(err).NotTo(HaveOccurred())
			Expect(waitCallCount).To(Equal(1))
		})

		When("update fails", func() {
			It("returns an error", func() {
				p.MockEKS().On("UpdateNodegroupConfig", &awseks.UpdateNodegroupConfigInput{
					ScalingConfig: &awseks.NodegroupScalingConfig{
						MinSize:     aws.Int64(1),
						DesiredSize: aws.Int64(3),
					},
					ClusterName:   &clusterName,
					NodegroupName: &ngName,
				}).Return(nil, fmt.Errorf("foo"))

				err := m.Scale(ng)

				Expect(err).To(MatchError(fmt.Sprintf("failed to scale nodegroup for cluster %q, error: foo", clusterName)))
			})
		})
	})

	Describe("Unmanaged Nodegroup", func() {
		When("the ASG exists", func() {
			BeforeEach(func() {
				nodegroups := make(map[string]manager.StackInfo)
				nodegroups["my-ng"] = manager.StackInfo{
					Stack: &manager.Stack{
						Tags: []*cloudformation.Tag{
							{
								Key:   aws.String(api.NodeGroupNameTag),
								Value: aws.String("my-ng"),
							},
							{
								Key:   aws.String(api.NodeGroupTypeTag),
								Value: aws.String(string(api.NodeGroupTypeUnmanaged)),
							},
						},
					},
					Resources: []*cloudformation.StackResource{
						{
							PhysicalResourceId: aws.String("asg-name"),
							LogicalResourceId:  aws.String("NodeGroup"),
						},
					},
				}
				fakeStackManager.DescribeNodeGroupStacksAndResourcesReturns(nodegroups, nil)

				p.MockASG().On("UpdateAutoScalingGroup", &autoscaling.UpdateAutoScalingGroupInput{
					AutoScalingGroupName: aws.String("asg-name"),
					MinSize:              aws.Int64(1),
					DesiredCapacity:      aws.Int64(3),
				}).Return(nil, nil)

			})

			It("scales the nodegroup", func() {
				err := m.Scale(ng)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("the asg resource doesn't exist", func() {
			BeforeEach(func() {
				nodegroups := make(map[string]manager.StackInfo)
				nodegroups["my-ng"] = manager.StackInfo{
					Stack: &manager.Stack{
						Tags: []*cloudformation.Tag{
							{
								Key:   aws.String(api.NodeGroupNameTag),
								Value: aws.String("my-ng"),
							},
							{
								Key:   aws.String(api.NodeGroupTypeTag),
								Value: aws.String(string(api.NodeGroupTypeUnmanaged)),
							},
						},
					},
					Resources: []*cloudformation.StackResource{},
				}
				fakeStackManager.DescribeNodeGroupStacksAndResourcesReturns(nodegroups, nil)

			})

			It("returns an error", func() {
				err := m.Scale(ng)
				Expect(err).To(MatchError(ContainSubstring("failed to find NodeGroup auto scaling group")))
			})
		})
	})
})
