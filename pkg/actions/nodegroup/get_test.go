package nodegroup_test

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
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
		fakeClientSet                  *fake.Clientset
	)

	BeforeEach(func() {
		t = time.Now()
		ngName = "my-nodegroup"
		clusterName = "my-cluster"
		cfg = api.NewClusterConfig()
		cfg.Metadata.Name = clusterName
		p = mockprovider.NewMockProvider()
		fakeClientSet = fake.NewSimpleClientset()
		m = nodegroup.New(cfg, &eks.ClusterProvider{Provider: p}, fakeClientSet)

		fakeStackManager = new(fakes.FakeStackManager)
		m.SetStackManager(fakeStackManager)
	})

	Describe("GetAll", func() {
		BeforeEach(func() {
			p.MockEKS().On("ListNodegroups", &awseks.ListNodegroupsInput{
				ClusterName: aws.String(clusterName),
			}).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(1))
				Expect(args[0]).To(BeAssignableToTypeOf(&awseks.ListNodegroupsInput{
					ClusterName: aws.String(clusterName),
				}))
			}).Return(&awseks.ListNodegroupsOutput{
				Nodegroups: []*string{
					aws.String(ngName),
				},
			}, nil)
		})

		Context("when getting managed nodegroups", func() {
			When("a nodegroup is associated to a CF Stack", func() {
				BeforeEach(func() {
					p.MockEKS().On("DescribeNodegroup", &awseks.DescribeNodegroupInput{
						ClusterName:   aws.String(clusterName),
						NodegroupName: aws.String(ngName),
					}).Run(func(args mock.Arguments) {
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
							Version: aws.String("1.18"),
						},
					}, nil)
				})

				It("returns a summary of the node group and its StackName", func() {
					fakeStackManager.DescribeNodeGroupStackReturns(&cloudformation.Stack{
						StackName: aws.String(stackName),
					}, nil)

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
					Expect(ngSummary[0].Version).To(Equal("1.18"))
				})
			})

			When("a nodegroup is not associated to a CF Stack", func() {
				BeforeEach(func() {
					p.MockEKS().On("DescribeNodegroup", &awseks.DescribeNodegroupInput{
						ClusterName:   aws.String(clusterName),
						NodegroupName: aws.String(ngName),
					}).Run(func(args mock.Arguments) {
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
							Version: aws.String("1.18"),
						},
					}, nil)
				})

				It("returns a summary of the node group without a StackName", func() {
					fakeStackManager.DescribeNodeGroupStackReturns(nil, fmt.Errorf("error describing cloudformation stack"))

					ngSummary, err := m.GetAll()
					Expect(err).NotTo(HaveOccurred())
					Expect(ngSummary[0].StackName).To(Equal(""))
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
					Expect(ngSummary[0].Version).To(Equal("1.18"))
				})
			})

			When("a nodegroup has multiple ASGs", func() {
				BeforeEach(func() {
					p.MockEKS().On("DescribeNodegroup", &awseks.DescribeNodegroupInput{
						ClusterName:   aws.String(clusterName),
						NodegroupName: aws.String(ngName),
					}).Run(func(args mock.Arguments) {
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
							InstanceTypes: aws.StringSlice([]string{"m5.xlarge"}),
							AmiType:       aws.String("ami-type"),
							CreatedAt:     &t,
							NodeRole:      aws.String("node-role"),
							Resources: &awseks.NodegroupResources{
								AutoScalingGroups: []*awseks.AutoScalingGroup{
									{
										Name: aws.String("asg-1"),
									},
									{
										Name: aws.String("asg-2"),
									},
								},
							},
							Version: aws.String("1.18"),
						},
					}, nil)
				})

				It("returns all ASG names in the summary", func() {
					summary, err := m.Get(ngName)
					Expect(err).NotTo(HaveOccurred())
					Expect(summary.AutoScalingGroupName).To(Equal("asg-1,asg-2"))
					Expect(summary.Name).To(Equal(ngName))
					Expect(summary.Cluster).To(Equal(clusterName))
					Expect(summary.Status).To(Equal("my-status"))
					Expect(summary.MaxSize).To(Equal(4))
					Expect(summary.MinSize).To(Equal(0))
					Expect(summary.DesiredCapacity).To(Equal(2))
					Expect(summary.InstanceType).To(Equal("m5.xlarge"))
					Expect(summary.ImageID).To(Equal("ami-type"))
					Expect(summary.CreationTime).To(Equal(&t))
					Expect(summary.NodeInstanceRoleARN).To(Equal("node-role"))
					Expect(summary.Version).To(Equal("1.18"))
				})
			})

			When("a nodegroup is associated with a launch template", func() {
				BeforeEach(func() {
					p.MockEKS().On("DescribeNodegroup", &awseks.DescribeNodegroupInput{
						ClusterName:   aws.String(clusterName),
						NodegroupName: aws.String(ngName),
					}).Run(func(args mock.Arguments) {
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
							Version: aws.String("1.18"),
							LaunchTemplate: &awseks.LaunchTemplateSpecification{
								Id:      aws.String("4"),
								Version: aws.String("5"),
							},
						},
					}, nil)

					p.MockEC2().On("DescribeLaunchTemplateVersions", &ec2.DescribeLaunchTemplateVersionsInput{
						LaunchTemplateId: aws.String("4"),
					}).Return(&ec2.DescribeLaunchTemplateVersionsOutput{LaunchTemplateVersions: []*ec2.LaunchTemplateVersion{
						{
							LaunchTemplateData: &ec2.ResponseLaunchTemplateData{
								InstanceType: aws.String("big"),
							},
							VersionNumber: aws.Int64(5),
						},
					}}, nil)
				})

				It("returns a summary of the node group with the instance type from the launch template", func() {
					fakeStackManager.DescribeNodeGroupStackReturns(nil, fmt.Errorf("error describing cloudformation stack"))

					ngSummary, err := m.GetAll()
					Expect(err).NotTo(HaveOccurred())
					Expect(ngSummary[0].StackName).To(Equal(""))
					Expect(ngSummary[0].Name).To(Equal(ngName))
					Expect(ngSummary[0].Cluster).To(Equal(clusterName))
					Expect(ngSummary[0].Status).To(Equal("my-status"))
					Expect(ngSummary[0].MaxSize).To(Equal(4))
					Expect(ngSummary[0].MinSize).To(Equal(0))
					Expect(ngSummary[0].DesiredCapacity).To(Equal(2))
					Expect(ngSummary[0].InstanceType).To(Equal("big"))
					Expect(ngSummary[0].ImageID).To(Equal("ami-type"))
					Expect(ngSummary[0].CreationTime).To(Equal(&t))
					Expect(ngSummary[0].NodeInstanceRoleARN).To(Equal("node-role"))
					Expect(ngSummary[0].AutoScalingGroupName).To(Equal("asg-name"))
					Expect(ngSummary[0].Version).To(Equal("1.18"))
				})
			})
		})

		Context("when getting unmanaged nodegroups", func() {
			When("the DesiredCapacity is 0", func() {
				BeforeEach(func() {
					p.MockEKS().On("DescribeNodegroup", &awseks.DescribeNodegroupInput{
						ClusterName:   aws.String(clusterName),
						NodegroupName: aws.String(ngName),
					}).Run(func(args mock.Arguments) {
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
							Version: aws.String("1.18"),
						},
					}, nil)

					fakeStackManager.DescribeNodeGroupStackReturns(&cloudformation.Stack{
						StackName: aws.String(stackName),
					}, nil)
				})

				It("does not return the k8s version of the nodes", func() {
					fakeStackManager.GetUnmanagedNodeGroupSummariesReturns([]*manager.NodeGroupSummary{
						{
							DesiredCapacity: 0,
						},
					}, nil)

					ngSummary, err := m.GetAll()
					Expect(err).NotTo(HaveOccurred())
					Expect(ngSummary[0].Version).To(Equal(""))
				})
			})
		})
	})
})
