package nodegroup_test

import (
	"context"
	"fmt"
	"time"

	autoscalingtypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	. "github.com/onsi/ginkgo/v2"
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
		ng                  *api.NodeGroupBase
		m                   *nodegroup.Manager
		fakeStackManager    *fakes.FakeStackManager
	)
	BeforeEach(func() {
		clusterName = "my-cluster"
		ngName = "my-ng"
		p = mockprovider.NewMockProvider()
		cfg = api.NewClusterConfig()
		cfg.Metadata.Name = clusterName

		ng = &api.NodeGroupBase{
			Name: ngName,
			ScalingConfig: &api.ScalingConfig{
				MinSize:         aws.Int(1),
				DesiredCapacity: aws.Int(3),
			},
		}

		m = nodegroup.New(cfg, &eks.ClusterProvider{AWSProvider: p}, fake.NewSimpleClientset(), nil)
		fakeStackManager = new(fakes.FakeStackManager)
		m.SetStackManager(fakeStackManager)
		p.MockCloudFormation().On("ListStacksPages", mock.Anything, mock.Anything).Return(nil, nil)
	})

	Describe("Managed NodeGroup", func() {
		BeforeEach(func() {
			nodegroups := make(map[string]manager.StackInfo)
			nodegroups["my-ng"] = manager.StackInfo{
				Stack: &manager.Stack{
					Tags: []types.Tag{
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
			p.MockEKS().On("UpdateNodegroupConfig", mock.Anything, &awseks.UpdateNodegroupConfigInput{
				ScalingConfig: &ekstypes.NodegroupScalingConfig{
					MinSize:     aws.Int32(1),
					DesiredSize: aws.Int32(3),
				},
				ClusterName:   &clusterName,
				NodegroupName: &ngName,
			}).Return(nil, nil)

			err := m.Scale(context.Background(), ng, false)
			Expect(err).NotTo(HaveOccurred())
		})

		It("waits for scaling and times out", func() {
			p.MockEKS().On("UpdateNodegroupConfig", mock.Anything, mock.Anything).Return(nil, nil)

			p.MockEKS().On("DescribeNodegroup", mock.Anything, mock.Anything, mock.Anything).Return(
				&awseks.DescribeNodegroupOutput{
					Nodegroup: &ekstypes.Nodegroup{
						Status: ekstypes.NodegroupStatusUpdating,
					},
				}, nil)

			p.SetWaitTimeout(1 * time.Millisecond)

			err := m.Scale(context.Background(), ng, true)
			Expect(err).To(MatchError(fmt.Sprintf("failed to scale nodegroup %q for cluster %q, error: exceeded max wait time for NodegroupActive waiter", ngName, clusterName)))
		})

		When("update fails", func() {
			It("returns an error", func() {
				p.MockEKS().On("UpdateNodegroupConfig", mock.Anything, &awseks.UpdateNodegroupConfigInput{
					ScalingConfig: &ekstypes.NodegroupScalingConfig{
						MinSize:     aws.Int32(1),
						DesiredSize: aws.Int32(3),
					},
					ClusterName:   &clusterName,
					NodegroupName: &ngName,
				}).Return(nil, fmt.Errorf("foo"))

				err := m.Scale(context.Background(), ng, false)
				Expect(err).To(MatchError(fmt.Sprintf("failed to scale nodegroup %q for cluster %q, error: foo", ngName, clusterName)))
			})
		})
	})

	Describe("Unmanaged Nodegroup", func() {
		mockNodeGroupStack := func(ngName, asgName string) {
			nodegroups := map[string]manager.StackInfo{
				ngName: {
					Stack: &manager.Stack{
						Tags: []types.Tag{
							{
								Key:   aws.String(api.NodeGroupNameTag),
								Value: aws.String(ngName),
							},
							{
								Key:   aws.String(api.NodeGroupTypeTag),
								Value: aws.String(string(api.NodeGroupTypeUnmanaged)),
							},
						},
					},
					Resources: []types.StackResource{
						{
							PhysicalResourceId: aws.String(asgName),
							LogicalResourceId:  aws.String("NodeGroup"),
						},
					},
				},
			}
			fakeStackManager.DescribeNodeGroupStacksAndResourcesReturns(nodegroups, nil)
		}

		mockNodeGroupAMI := func(amiDeprecated bool, asgName string) {
			p.MockASG().On("DescribeAutoScalingGroups", mock.Anything, mock.MatchedBy(func(input *autoscaling.DescribeAutoScalingGroupsInput) bool {
				return len(input.AutoScalingGroupNames) == 1 && input.AutoScalingGroupNames[0] == asgName
			})).Return(&autoscaling.DescribeAutoScalingGroupsOutput{
				AutoScalingGroups: []autoscalingtypes.AutoScalingGroup{
					{
						LaunchTemplate: &autoscalingtypes.LaunchTemplateSpecification{
							LaunchTemplateId:   aws.String("lt-1234"),
							LaunchTemplateName: aws.String("eksctl-test-ng"),
							Version:            aws.String("1"),
						},
					},
				},
			}, nil)
			p.MockEC2().On("DescribeLaunchTemplateVersions", mock.Anything, mock.MatchedBy(func(input *ec2.DescribeLaunchTemplateVersionsInput) bool {
				return len(input.Versions) == 1 && input.LaunchTemplateId != nil && *input.LaunchTemplateId == "lt-1234"
			})).Return(&ec2.DescribeLaunchTemplateVersionsOutput{
				LaunchTemplateVersions: []ec2types.LaunchTemplateVersion{
					{
						LaunchTemplateData: &ec2types.ResponseLaunchTemplateData{
							ImageId: aws.String("ami-1234"),
						},
					},
				},
			}, nil)
			describeImagesOutput := &ec2.DescribeImagesOutput{}
			if !amiDeprecated {
				describeImagesOutput.Images = []ec2types.Image{
					{
						ImageId: aws.String("ami-1234"),
					},
				}
			}
			p.MockEC2().On("DescribeImages", mock.Anything, mock.MatchedBy(func(input *ec2.DescribeImagesInput) bool {
				return len(input.ImageIds) == 1 && input.ImageIds[0] == "ami-1234"
			})).Return(describeImagesOutput, nil)

			if !amiDeprecated {
				p.MockASG().On("UpdateAutoScalingGroup", mock.Anything, &autoscaling.UpdateAutoScalingGroupInput{
					AutoScalingGroupName: aws.String(asgName),
					MinSize:              aws.Int32(1),
					DesiredCapacity:      aws.Int32(3),
				}).Return(nil, nil)
			}
		}

		mockMixedInstanceNodeGroupAMI := func(missingLaunchTemplate bool, asgName string) {
			asgOutputWithMixedInstancesPolicy := &autoscaling.DescribeAutoScalingGroupsOutput{
				AutoScalingGroups: []autoscalingtypes.AutoScalingGroup{
					{
						MixedInstancesPolicy: &autoscalingtypes.MixedInstancesPolicy{},
					},
				},
			}
			if !missingLaunchTemplate {
				asgOutputWithMixedInstancesPolicy.AutoScalingGroups[0].MixedInstancesPolicy.LaunchTemplate = &autoscalingtypes.LaunchTemplate{
					LaunchTemplateSpecification: &autoscalingtypes.LaunchTemplateSpecification{
						LaunchTemplateId:   aws.String("lt-1234"),
						LaunchTemplateName: aws.String("eksctl-test-ng"),
						Version:            aws.String("1"),
					},
				}
			}
			p.MockASG().On("DescribeAutoScalingGroups", mock.Anything, mock.MatchedBy(func(input *autoscaling.DescribeAutoScalingGroupsInput) bool {
				return len(input.AutoScalingGroupNames) == 1 && input.AutoScalingGroupNames[0] == asgName
			})).Return(asgOutputWithMixedInstancesPolicy, nil)
			p.MockEC2().On("DescribeLaunchTemplateVersions", mock.Anything, mock.MatchedBy(func(input *ec2.DescribeLaunchTemplateVersionsInput) bool {
				return len(input.Versions) == 1 && input.LaunchTemplateId != nil && *input.LaunchTemplateId == "lt-1234"
			})).Return(&ec2.DescribeLaunchTemplateVersionsOutput{
				LaunchTemplateVersions: []ec2types.LaunchTemplateVersion{
					{
						LaunchTemplateData: &ec2types.ResponseLaunchTemplateData{
							ImageId: aws.String("ami-1234"),
						},
					},
				},
			}, nil)
			describeImagesOutput := &ec2.DescribeImagesOutput{
				Images: []ec2types.Image{
					{
						ImageId: aws.String("ami-1234"),
					},
				},
			}
			p.MockEC2().On("DescribeImages", mock.Anything, mock.MatchedBy(func(input *ec2.DescribeImagesInput) bool {
				return len(input.ImageIds) == 1 && input.ImageIds[0] == "ami-1234"
			})).Return(describeImagesOutput, nil)

			p.MockASG().On("UpdateAutoScalingGroup", mock.Anything, &autoscaling.UpdateAutoScalingGroupInput{
				AutoScalingGroupName: aws.String(asgName),
				MinSize:              aws.Int32(1),
				DesiredCapacity:      aws.Int32(3),
			}).Return(nil, nil)
		}

		When("the ASG exists with a MixedInstancesPolicy", func() {
			asgName := "asg-name"

			BeforeEach(func() {
				mockNodeGroupStack(ngName, asgName)
			})

			It("fails for missing launch template", func() {
				mockMixedInstanceNodeGroupAMI(true, asgName)
				err := m.Scale(context.Background(), ng, false)
				Expect(err).To(MatchError(fmt.Sprintf("failed to scale nodegroup %q for cluster %q, error: expected the MixedInstancesPolicy in Auto Scaling group %q to include a launch template", ng.Name, clusterName, asgName)))
			})

			It("scales the nodegroup", func() {
				mockMixedInstanceNodeGroupAMI(false, asgName)
				err := m.Scale(context.Background(), ng, false)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("the ASG exists", func() {
			BeforeEach(func() {
				mockNodeGroupStack(ngName, "asg-name")
				mockNodeGroupAMI(false, "asg-name")
				p.MockASG().On("UpdateAutoScalingGroup", mock.Anything, &autoscaling.UpdateAutoScalingGroupInput{
					AutoScalingGroupName: aws.String("asg-name"),
					MinSize:              aws.Int32(1),
					DesiredCapacity:      aws.Int32(3),
				}).Return(nil, nil)

			})

			It("scales the nodegroup", func() {
				err := m.Scale(context.Background(), ng, false)
				Expect(err).NotTo(HaveOccurred())
			})

			It("waits for scaling and times out", func() {
				p.SetWaitTimeout(1 * time.Millisecond)
				err := m.Scale(context.Background(), ng, true)
				Expect(err).To(MatchError(fmt.Sprintf("failed to scale nodegroup %q for cluster %q, error: timed out waiting for at least %d nodes to join the cluster and become ready in %q: context deadline exceeded", ng.Name, clusterName, *ng.ScalingConfig.MinSize, ng.Name)))
			})
		})

		When("the asg resource doesn't exist", func() {
			BeforeEach(func() {
				nodegroups := make(map[string]manager.StackInfo)
				nodegroups["my-ng"] = manager.StackInfo{
					Stack: &manager.Stack{
						Tags: []types.Tag{
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
					Resources: []types.StackResource{},
				}
				fakeStackManager.DescribeNodeGroupStacksAndResourcesReturns(nodegroups, nil)

			})

			It("returns an error", func() {
				err := m.Scale(context.Background(), ng, false)
				Expect(err).To(MatchError(ContainSubstring("failed to find NodeGroup auto scaling group")))
			})
		})

		type amiDeprecationEntry struct {
			deprecated  bool
			expectedErr string
		}

		DescribeTable("AMI deprecation", func(e amiDeprecationEntry) {
			const asgName = "asg-1234"
			mockNodeGroupStack(ngName, asgName)
			mockNodeGroupAMI(e.deprecated, asgName)
			err := m.Scale(context.Background(), ng, false)
			if e.expectedErr != "" {
				Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
		},
			Entry("AMI associated with the nodegroup is either deprecated or removed", amiDeprecationEntry{
				deprecated:  true,
				expectedErr: "AMI associated with the nodegroup is either deprecated or removed; please upgrade the nodegroup before scaling it: https://eksctl.io/usage/nodegroup-upgrade/",
			}),
			Entry("AMI associated with the nodegroup still exists", amiDeprecationEntry{
				deprecated: false,
			}),
		)

	})
})
