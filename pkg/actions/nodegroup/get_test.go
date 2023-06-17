package nodegroup_test

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	asgtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	cftypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"

	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/smithy-go"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	. "github.com/onsi/ginkgo/v2"
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
		ngName           = "my-nodegroup"
		stackName        = "stack-name"
		clusterName      = "my-cluster"
		t                = time.Now()
		p                *mockprovider.MockProvider
		cfg              *api.ClusterConfig
		m                *nodegroup.Manager
		fakeStackManager *fakes.FakeStackManager
		fakeClientSet    *fake.Clientset
	)

	BeforeEach(func() {
		cfg = api.NewClusterConfig()
		cfg.Metadata.Name = clusterName
		p = mockprovider.NewMockProvider()
		fakeClientSet = fake.NewSimpleClientset()
		m = nodegroup.New(cfg, &eks.ClusterProvider{AWSProvider: p}, fakeClientSet, nil)
		fakeStackManager = new(fakes.FakeStackManager)
		m.SetStackManager(fakeStackManager)
	})

	Describe("GetAll", func() {
		BeforeEach(func() {
			p.MockEKS().On("ListNodegroups", mock.Anything, &awseks.ListNodegroupsInput{
				ClusterName: aws.String(clusterName),
			}).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.ListNodegroupsInput{
					ClusterName: aws.String(clusterName),
				}))
			}).Return(&awseks.ListNodegroupsOutput{
				Nodegroups: []string{ngName},
			}, nil)
		})

		Context("when getting managed nodegroups", func() {
			When("a nodegroup is associated to a CF Stack", func() {
				BeforeEach(func() {
					p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
						ClusterName:   aws.String(clusterName),
						NodegroupName: aws.String(ngName),
					}).Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeNodegroupInput{
							ClusterName:   aws.String(clusterName),
							NodegroupName: aws.String(ngName),
						}))
					}).Return(&awseks.DescribeNodegroupOutput{
						Nodegroup: &ekstypes.Nodegroup{
							NodegroupName: aws.String(ngName),
							ClusterName:   aws.String(clusterName),
							Status:        "my-status",
							ScalingConfig: &ekstypes.NodegroupScalingConfig{
								DesiredSize: aws.Int32(2),
								MaxSize:     aws.Int32(4),
								MinSize:     aws.Int32(0),
							},
							InstanceTypes: []string{},
							AmiType:       "ami-type",
							CreatedAt:     &t,
							NodeRole:      aws.String("node-role"),
							Resources: &ekstypes.NodegroupResources{
								AutoScalingGroups: []ekstypes.AutoScalingGroup{
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
					fakeStackManager.DescribeNodeGroupStackReturns(&cftypes.Stack{
						StackName: aws.String(stackName),
					}, nil)

					summaries, err := m.GetAll(context.Background())
					Expect(err).NotTo(HaveOccurred())
					Expect(summaries).To(HaveLen(1))

					ngSummary := *summaries[0]
					Expect(ngSummary).To(Equal(nodegroup.Summary{
						StackName:            stackName,
						Cluster:              clusterName,
						Name:                 ngName,
						Status:               "my-status",
						MaxSize:              4,
						MinSize:              0,
						DesiredCapacity:      2,
						InstanceType:         "-",
						ImageID:              "ami-type",
						CreationTime:         t,
						NodeInstanceRoleARN:  "node-role",
						AutoScalingGroupName: "asg-name",
						Version:              "1.18",
						NodeGroupType:        api.NodeGroupTypeManaged,
					}))
				})
			})

			When("a nodegroup is not associated to a CF Stack", func() {
				BeforeEach(func() {
					p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
						ClusterName:   aws.String(clusterName),
						NodegroupName: aws.String(ngName),
					}).Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeNodegroupInput{
							ClusterName:   aws.String(clusterName),
							NodegroupName: aws.String(ngName),
						}))
					}).Return(&awseks.DescribeNodegroupOutput{
						Nodegroup: &ekstypes.Nodegroup{
							NodegroupName: aws.String(ngName),
							ClusterName:   aws.String(clusterName),
							Status:        "my-status",
							ScalingConfig: &ekstypes.NodegroupScalingConfig{
								DesiredSize: aws.Int32(2),
								MaxSize:     aws.Int32(4),
								MinSize:     aws.Int32(0),
							},
							InstanceTypes: []string{},
							AmiType:       "ami-type",
							CreatedAt:     &t,
							NodeRole:      aws.String("node-role"),
							Resources: &ekstypes.NodegroupResources{
								AutoScalingGroups: []ekstypes.AutoScalingGroup{
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

					summaries, err := m.GetAll(context.Background())
					Expect(err).NotTo(HaveOccurred())
					Expect(summaries).To(HaveLen(1))

					ngSummary := *summaries[0]
					Expect(ngSummary).To(Equal(nodegroup.Summary{
						StackName:            "",
						Cluster:              clusterName,
						Name:                 ngName,
						Status:               "my-status",
						MaxSize:              4,
						MinSize:              0,
						DesiredCapacity:      2,
						InstanceType:         "-",
						ImageID:              "ami-type",
						CreationTime:         t,
						NodeInstanceRoleARN:  "node-role",
						AutoScalingGroupName: "asg-name",
						Version:              "1.18",
						NodeGroupType:        api.NodeGroupTypeManaged,
					}))
				})
			})

			When("a nodegroup is associated with a launch template", func() {
				BeforeEach(func() {
					p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
						ClusterName:   aws.String(clusterName),
						NodegroupName: aws.String(ngName),
					}).Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeNodegroupInput{
							ClusterName:   aws.String(clusterName),
							NodegroupName: aws.String(ngName),
						}))
					}).Return(&awseks.DescribeNodegroupOutput{
						Nodegroup: &ekstypes.Nodegroup{
							NodegroupName: aws.String(ngName),
							ClusterName:   aws.String(clusterName),
							Status:        "my-status",
							ScalingConfig: &ekstypes.NodegroupScalingConfig{
								DesiredSize: aws.Int32(2),
								MaxSize:     aws.Int32(4),
								MinSize:     aws.Int32(0),
							},
							InstanceTypes: []string{},
							AmiType:       "ami-type",
							CreatedAt:     &t,
							NodeRole:      aws.String("node-role"),
							Resources: &ekstypes.NodegroupResources{
								AutoScalingGroups: []ekstypes.AutoScalingGroup{
									{
										Name: aws.String("asg-name"),
									},
								},
							},
							Version: aws.String("1.18"),
							LaunchTemplate: &ekstypes.LaunchTemplateSpecification{
								Id:      aws.String("4"),
								Version: aws.String("5"),
							},
						},
					}, nil)

					p.MockEC2().On("DescribeLaunchTemplateVersions", mock.Anything, &ec2.DescribeLaunchTemplateVersionsInput{
						LaunchTemplateId: aws.String("4"),
					}).Return(&ec2.DescribeLaunchTemplateVersionsOutput{LaunchTemplateVersions: []ec2types.LaunchTemplateVersion{
						{
							LaunchTemplateData: &ec2types.ResponseLaunchTemplateData{
								InstanceType: "big",
							},
							VersionNumber: aws.Int64(5),
						},
					}}, nil)
				})

				It("returns a summary of the node group with the instance type from the launch template", func() {
					fakeStackManager.DescribeNodeGroupStackReturns(nil, fmt.Errorf("error describing cloudformation stack"))

					summaries, err := m.GetAll(context.Background())
					Expect(err).NotTo(HaveOccurred())
					Expect(summaries).To(HaveLen(1))

					ngSummary := *summaries[0]
					Expect(ngSummary).To(Equal(nodegroup.Summary{
						StackName:            "",
						Cluster:              clusterName,
						Name:                 ngName,
						Status:               "my-status",
						MaxSize:              4,
						MinSize:              0,
						DesiredCapacity:      2,
						InstanceType:         "big",
						ImageID:              "ami-type",
						CreationTime:         t,
						NodeInstanceRoleARN:  "node-role",
						AutoScalingGroupName: "asg-name",
						Version:              "1.18",
						NodeGroupType:        api.NodeGroupTypeManaged,
					}))
				})
			})

			When("a nodegroup has a custom AMI", func() {
				BeforeEach(func() {
					p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
						ClusterName:   aws.String(clusterName),
						NodegroupName: aws.String(ngName),
					}).Run(func(args mock.Arguments) {
						Expect(args).To(HaveLen(2))
						Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeNodegroupInput{
							ClusterName:   aws.String(clusterName),
							NodegroupName: aws.String(ngName),
						}))
					}).Return(&awseks.DescribeNodegroupOutput{
						Nodegroup: &ekstypes.Nodegroup{
							NodegroupName: aws.String(ngName),
							ClusterName:   aws.String(clusterName),
							Status:        "my-status",
							ScalingConfig: &ekstypes.NodegroupScalingConfig{
								DesiredSize: aws.Int32(2),
								MaxSize:     aws.Int32(4),
								MinSize:     aws.Int32(0),
							},
							InstanceTypes:  []string{"m5.xlarge"},
							AmiType:        "CUSTOM",
							CreatedAt:      &t,
							NodeRole:       aws.String("node-role"),
							ReleaseVersion: aws.String("ami-custom"),
							Resources: &ekstypes.NodegroupResources{
								AutoScalingGroups: []ekstypes.AutoScalingGroup{
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

					fakeStackManager.DescribeNodeGroupStackReturns(&cftypes.Stack{
						StackName: aws.String(stackName),
					}, nil)
				})

				It("returns the AMI ID instead of `CUSTOM`", func() {
					summaries, err := m.GetAll(context.Background())
					Expect(err).NotTo(HaveOccurred())
					Expect(summaries).To(HaveLen(1))

					ngSummary := *summaries[0]
					Expect(ngSummary).To(Equal(nodegroup.Summary{
						StackName:            stackName,
						Cluster:              clusterName,
						Name:                 ngName,
						Status:               "my-status",
						MaxSize:              4,
						MinSize:              0,
						DesiredCapacity:      2,
						InstanceType:         "m5.xlarge",
						ImageID:              "ami-custom",
						CreationTime:         t,
						NodeInstanceRoleARN:  "node-role",
						AutoScalingGroupName: "asg-1,asg-2",
						Version:              "1.18",
						NodeGroupType:        api.NodeGroupTypeManaged,
					}))
				})
			})
		})

		Context("when getting unmanaged nodegroups", func() {
			var (
				creationTime           = time.Now()
				unmanagedStackName     = "unmanaged-stack"
				unmanagedNodegroupName = "unmanaged-ng"
				unmanagedTemplate      = `
{
  "Resources": {
    "NodeGroup": {
      "Type": "AWS::AutoScaling::AutoScalingGroup",
      "Properties": {
        "DesiredCapacity": "3",
        "MaxSize": "6",
        "MinSize": "1"
      }
    }
  }
}

`
			)

			BeforeEach(func() {
				//unmanaged nodegroup
				fakeStackManager.ListNodeGroupStacksReturns([]*cftypes.Stack{
					{
						StackName: aws.String(unmanagedStackName),
						Tags: []cftypes.Tag{
							{
								Key:   aws.String(api.NodeGroupNameTag),
								Value: aws.String(unmanagedNodegroupName),
							},
							{
								Key:   aws.String(api.ClusterNameTag),
								Value: aws.String(clusterName),
							},
						},
						StackStatus:  "CREATE_COMPLETE",
						CreationTime: aws.Time(creationTime),
					},
				}, nil)
				fakeStackManager.GetStackTemplateReturns(unmanagedTemplate, nil)
				fakeStackManager.GetUnmanagedNodeGroupAutoScalingGroupNameReturns("asg", nil)
				fakeStackManager.GetAutoScalingGroupDesiredCapacityReturns(asgtypes.AutoScalingGroup{
					DesiredCapacity: aws.Int32(50),
					MinSize:         aws.Int32(1),
					MaxSize:         aws.Int32(100),
				}, nil)

				_, _ = fakeClientSet.CoreV1().Nodes().Create(context.Background(), &corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"alpha.eksctl.io/nodegroup-name": unmanagedNodegroupName,
						},
					},
					Status: corev1.NodeStatus{
						NodeInfo: corev1.NodeSystemInfo{
							KubeletVersion: "1.23.1",
						},
					},
				}, metav1.CreateOptions{})
				fakeStackManager.GetNodeGroupNameReturns(unmanagedNodegroupName)

				//managed nodegroup
				p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
					ClusterName:   aws.String(clusterName),
					NodegroupName: aws.String(ngName),
				}).Run(func(args mock.Arguments) {
					Expect(args).To(HaveLen(2))
					Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeNodegroupInput{
						ClusterName:   aws.String(clusterName),
						NodegroupName: aws.String(ngName),
					}))
				}).Return(&awseks.DescribeNodegroupOutput{
					Nodegroup: &ekstypes.Nodegroup{
						NodegroupName: aws.String(ngName),
						ClusterName:   aws.String(clusterName),
						Status:        "my-status",
						ScalingConfig: &ekstypes.NodegroupScalingConfig{
							DesiredSize: aws.Int32(2),
							MaxSize:     aws.Int32(4),
							MinSize:     aws.Int32(0),
						},
						InstanceTypes: []string{},
						AmiType:       "ami-type",
						CreatedAt:     &t,
						NodeRole:      aws.String("node-role"),
						Resources: &ekstypes.NodegroupResources{
							AutoScalingGroups: []ekstypes.AutoScalingGroup{
								{
									Name: aws.String("asg-name"),
								},
							},
						},
						Version: aws.String("1.18"),
					},
				}, nil)

				fakeStackManager.DescribeNodeGroupStackReturns(&cftypes.Stack{
					StackName: aws.String(stackName),
				}, nil)
			})

			It("returns the nodegroups with the kubernetes version", func() {
				summaries, err := m.GetAll(context.Background())
				Expect(err).NotTo(HaveOccurred())
				Expect(summaries).To(HaveLen(2))

				unmanagedSummary := *summaries[0]
				Expect(unmanagedSummary).To(Equal(nodegroup.Summary{
					StackName:            unmanagedStackName,
					Cluster:              clusterName,
					Name:                 unmanagedNodegroupName,
					Status:               "CREATE_COMPLETE",
					AutoScalingGroupName: "asg",
					MaxSize:              100,
					DesiredCapacity:      50,
					MinSize:              1,
					Version:              "1.23.1",
					CreationTime:         creationTime,
					NodeGroupType:        api.NodeGroupTypeUnmanaged,
				}))

				Expect(*summaries[1]).To(Equal(nodegroup.Summary{
					StackName:            stackName,
					Cluster:              clusterName,
					Name:                 ngName,
					Status:               "my-status",
					MaxSize:              4,
					MinSize:              0,
					DesiredCapacity:      2,
					InstanceType:         "-",
					ImageID:              "ami-type",
					CreationTime:         t,
					NodeInstanceRoleARN:  "node-role",
					AutoScalingGroupName: "asg-name",
					Version:              "1.18",
					NodeGroupType:        api.NodeGroupTypeManaged,
				}))
			})
		})
	})

	Describe("Get", func() {
		BeforeEach(func() {
			fakeStackManager.DescribeNodeGroupStackReturns(&cftypes.Stack{
				StackName: aws.String(stackName),
				Tags: []cftypes.Tag{
					{
						Key:   aws.String(api.NodeGroupNameTag),
						Value: aws.String(ngName),
					},
					{
						Key:   aws.String(api.ClusterNameTag),
						Value: aws.String(clusterName),
					},
					{
						Key:   aws.String(api.NodeGroupTypeTag),
						Value: aws.String(string(api.NodeGroupTypeManaged)),
					},
				},
			}, nil)

			p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
				ClusterName:   aws.String(clusterName),
				NodegroupName: aws.String(ngName),
			}).Run(func(args mock.Arguments) {
				Expect(args).To(HaveLen(2))
				Expect(args[1]).To(BeAssignableToTypeOf(&awseks.DescribeNodegroupInput{
					ClusterName:   aws.String(clusterName),
					NodegroupName: aws.String(ngName),
				}))
			}).Return(&awseks.DescribeNodegroupOutput{
				Nodegroup: &ekstypes.Nodegroup{
					NodegroupName: aws.String(ngName),
					ClusterName:   aws.String(clusterName),
					Status:        "my-status",
					ScalingConfig: &ekstypes.NodegroupScalingConfig{
						DesiredSize: aws.Int32(2),
						MaxSize:     aws.Int32(4),
						MinSize:     aws.Int32(0),
					},
					InstanceTypes: []string{"m5.xlarge"},
					AmiType:       "ami-type",
					CreatedAt:     &t,
					NodeRole:      aws.String("node-role"),
					Resources: &ekstypes.NodegroupResources{
						AutoScalingGroups: []ekstypes.AutoScalingGroup{
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

		It("returns the summary and the stack name", func() {
			summary, err := m.Get(context.Background(), ngName)
			Expect(err).NotTo(HaveOccurred())

			Expect(*summary).To(Equal(nodegroup.Summary{
				StackName:            stackName,
				Cluster:              clusterName,
				Name:                 ngName,
				Status:               "my-status",
				MaxSize:              4,
				MinSize:              0,
				DesiredCapacity:      2,
				InstanceType:         "m5.xlarge",
				ImageID:              "ami-type",
				CreationTime:         t,
				NodeInstanceRoleARN:  "node-role",
				AutoScalingGroupName: "asg-1,asg-2",
				Version:              "1.18",
				NodeGroupType:        api.NodeGroupTypeManaged,
			}))
		})
		When("there is no associated stack to the nodegroup", func() {
			It("returns a summary of the node group without a StackName", func() {
				fakeStackManager.DescribeNodeGroupStackReturns(nil, errors.Wrap(&smithy.OperationError{
					Err: fmt.Errorf("ValidationError"),
				}, "nope"))
				summary, err := m.Get(context.Background(), ngName)
				Expect(err).NotTo(HaveOccurred())

				Expect(*summary).To(Equal(nodegroup.Summary{
					StackName:            "",
					Cluster:              clusterName,
					Name:                 ngName,
					Status:               "my-status",
					MaxSize:              4,
					MinSize:              0,
					DesiredCapacity:      2,
					InstanceType:         "m5.xlarge",
					ImageID:              "ami-type",
					CreationTime:         t,
					NodeInstanceRoleARN:  "node-role",
					AutoScalingGroupName: "asg-1,asg-2",
					Version:              "1.18",
					NodeGroupType:        api.NodeGroupTypeManaged,
				}))
			})
		})
	})
})
