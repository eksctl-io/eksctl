package nodegroup_test

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Get", func() {
	var (
		clusterName, stackName, ngName string
		t                              time.Time
		p                              *mockprovider.MockProvider
		cfg                            *api.ClusterConfig
		m                              *nodegroup.Manager
		fakeStackManager               *fakes.FakeStackManager
	)
	BeforeEach(func() {
		t = time.Now()
		ngName = "my-nodegroup"
		clusterName = "my-cluster"
		cfg = api.NewClusterConfig()
		cfg.Metadata.Name = clusterName
		p = mockprovider.NewMockProvider()
		m = nodegroup.New(cfg, &eks.ClusterProvider{Provider: p}, nil)

		fakeStackManager = new(fakes.FakeStackManager)
		m.SetStackManager(fakeStackManager)
	})

	Describe("Managed Nodegroup", func() {
		BeforeEach(func() {
			p.MockEKS().On("ListNodegroups", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.ListNodegroupsInput{
					ClusterName: aws.String(clusterName),
				}))
			}).Return(&awseks.ListNodegroupsOutput{
				Nodegroups: []*string{
					aws.String(ngName),
				},
			}, nil)

			p.MockEKS().On("DescribeNodegroup", mock.Anything).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.DescribeNodegroupInput{
					ClusterName:   aws.String(clusterName),
					NodegroupName: aws.String(ngName),
				}))
			}).Return(&awseks.DescribeNodegroupOutput{
				Nodegroup: &awseks.Nodegroup{
					NodegroupName: aws.String(ngName),
					ClusterName:   aws.String(clusterName),
					Status:        aws.String("my-status"),
					ScalingConfig: &awseks.NodegroupScalingConfig{
						DesiredSize: aws.Int64(2),
						MaxSize:     aws.Int64(4),
						MinSize:     aws.Int64(0),
					},
					InstanceTypes: []*string{},
					AmiType:       aws.String("ami-type"),
					CreatedAt:     &t,
					NodeRole:      aws.String("node-role"),
					Resources: &awseks.NodegroupResources{
						AutoScalingGroups: []*awseks.AutoScalingGroup{
							{
								Name: aws.String("asg-name"),
							},
						},
					},
				},
			}, nil)

			p.MockCloudFormation().On("DescribeNodeGroupStack", mock.Anything, mock.Anything).Return(nil, nil)
			cfStack := cloudformation.Stack{
				StackName: aws.String(stackName),
			}
			fakeStackManager.DescribeNodeGroupStackReturns(&cfStack, nil)
		})

		It("returns a summary of the node group", func() {
			ngSummary, err := m.GetAll()
			Expect(err).NotTo(HaveOccurred())
			Expect(ngSummary[0].StackName).To(Equal(stackName))
			Expect(ngSummary[0].Name).To(Equal(ngName))
			Expect(ngSummary[0].Cluster).To(Equal(clusterName))
			Expect(ngSummary[0].Status).To(Equal("my-status"))
			Expect(ngSummary[0].MaxSize).To(Equal(4))
			Expect(ngSummary[0].MinSize).To(Equal(0))
			Expect(ngSummary[0].DesiredCapacity).To(Equal(2))
			Expect(ngSummary[0].InstanceType).To(Equal("-"))
			Expect(ngSummary[0].ImageID).To(Equal("ami-type"))
			Expect(ngSummary[0].CreationTime).To(Equal(&t))
			Expect(ngSummary[0].NodeInstanceRoleARN).To(Equal("node-role"))
			Expect(ngSummary[0].AutoScalingGroupName).To(Equal("asg-name"))
		})
	})
})
