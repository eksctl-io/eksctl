package manager

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("StackCollection NodeGroup", func() {
	var (
		cc *api.ClusterConfig
		sc *StackCollection

		p *mockprovider.MockProvider
	)

	const nodegroupResource = `
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
	testAZs := []string{"us-west-2b", "us-west-2a", "us-west-2c"}

	newClusterConfig := func(clusterName string) *api.ClusterConfig {
		cfg := api.NewClusterConfig()

		cfg.Metadata.Region = "us-west-2"
		cfg.Metadata.Name = clusterName
		cfg.AvailabilityZones = testAZs

		*cfg.VPC.CIDR = api.DefaultCIDR()

		return cfg
	}

	newNodeGroup := func(cfg *api.ClusterConfig) *api.NodeGroup {
		ng := cfg.NewNodeGroup()
		ng.InstanceType = "t2.medium"
		ng.AMIFamily = "AmazonLinux2"
		ng.Name = "12345"

		return ng
	}

	Describe("GetUnmanagedNodeGroupSummaries", func() {
		Context("With a cluster name", func() {
			var (
				clusterName string
				err         error
				out         []*NodeGroupSummary
			)

			JustBeforeEach(func() {
				p = mockprovider.NewMockProvider()

				cc = newClusterConfig(clusterName)

				newNodeGroup(cc)

				sc = NewStackCollection(p, cc)

				p.MockCloudFormation().On("GetTemplate", mock.MatchedBy(func(input *cfn.GetTemplateInput) bool {
					return input.StackName != nil && *input.StackName == "eksctl-test-cluster-nodegroup-12345"
				})).Return(&cfn.GetTemplateOutput{
					TemplateBody: aws.String(nodegroupResource),
				}, nil)

				p.MockCloudFormation().On("GetTemplate", mock.Anything).Return(nil, fmt.Errorf("GetTemplate failed"))

				p.MockCloudFormation().On("ListStacksPages", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					consume := args[1].(func(p *cfn.ListStacksOutput, last bool) (shouldContinue bool))
					out := &cfn.ListStacksOutput{
						StackSummaries: []*cfn.StackSummary{
							{
								StackName: aws.String("eksctl-test-cluster-nodegroup-12345"),
							},
						},
					}
					cont := consume(out, true)
					if !cont {
						panic("unexpected return value from the paging function: shouldContinue was false. It becomes false only when subsequent DescribeStacks call(s) fail, which isn't expected in this test scenario")
					}
				}).Return(nil)

				p.MockCloudFormation().On("DescribeStacks", mock.MatchedBy(func(input *cfn.DescribeStacksInput) bool {
					return input.StackName != nil && *input.StackName == "eksctl-test-cluster-nodegroup-12345"
				})).Return(&cfn.DescribeStacksOutput{
					Stacks: []*cfn.Stack{
						{
							StackName:   aws.String("eksctl-test-cluster-nodegroup-12345"),
							StackId:     aws.String("eksctl-test-cluster-nodegroup-12345-id"),
							StackStatus: aws.String("CREATE_COMPLETE"),
							Tags: []*cfn.Tag{
								{
									Key:   aws.String(api.NodeGroupNameTag),
									Value: aws.String("12345"),
								},
							},
							Outputs: []*cfn.Output{
								{
									OutputKey:   aws.String("InstanceRoleARN"),
									OutputValue: aws.String("arn:aws:iam::1111:role/eks-nodes-base-role"),
								},
							},
						},
					},
				}, nil)

				p.MockCloudFormation().On("DescribeStacks", mock.Anything).Return(nil, fmt.Errorf("DescribeStacks failed"))

				p.MockCloudFormation().On("DescribeStackResource", mock.MatchedBy(func(input *cfn.DescribeStackResourceInput) bool {
					return input.StackName != nil && *input.StackName == "eksctl-test-cluster-nodegroup-12345" && input.LogicalResourceId != nil && *input.LogicalResourceId == "NodeGroup"
				})).Return(&cfn.DescribeStackResourceOutput{
					StackResourceDetail: &cfn.StackResourceDetail{
						PhysicalResourceId: aws.String("eksctl-test-cluster-nodegroup-12345-NodeGroup-1N68LL8H1EH27"),
					},
				}, nil)

				p.MockCloudFormation().On("DescribeStackResource", mock.Anything).Return(nil, fmt.Errorf("DescribeStackResource failed"))

				p.MockASG().On("DescribeAutoScalingGroups", mock.MatchedBy(func(input *autoscaling.DescribeAutoScalingGroupsInput) bool {
					return len(input.AutoScalingGroupNames) == 1 && *input.AutoScalingGroupNames[0] == "eksctl-test-cluster-nodegroup-12345-NodeGroup-1N68LL8H1EH27"
				})).Return(&autoscaling.DescribeAutoScalingGroupsOutput{
					AutoScalingGroups: []*autoscaling.Group{
						{
							DesiredCapacity: aws.Int64(7),
							MinSize:         aws.Int64(1),
							MaxSize:         aws.Int64(10),
						},
					},
				}, nil)

				p.MockASG().On("DescribeAutoScalingGroups", mock.Anything).Return(nil, fmt.Errorf("DescribeAutoScalingGroups failed"))
			})

			Context("With no matching stacks", func() {
				BeforeEach(func() {
					clusterName = "test-cluster-non-existent"
				})

				JustBeforeEach(func() {
					out, err = sc.GetUnmanagedNodeGroupSummaries("")
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should not have called AWS CloudFormation GetTemplate", func() {
					Expect(p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "GetTemplate", 0)).To(BeTrue())
				})

				It("the output should equal the expectation", func() {
					Expect(out).To(HaveLen(0))
				})
			})

			Context("With matching stacks", func() {
				BeforeEach(func() {
					clusterName = "test-cluster"
				})

				JustBeforeEach(func() {
					out, err = sc.GetUnmanagedNodeGroupSummaries("")
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should not have called AWS CloudFormation GetTemplate", func() {
					Expect(p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "GetTemplate", 1)).To(BeTrue())
				})

				It("should have called AWS CloudFormation DescribeStacks once", func() {
					Expect(p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "DescribeStacks", 1)).To(BeTrue())
				})

				It("the output should equal the expectation", func() {
					Expect(out).To(HaveLen(1))
					Expect(out[0].StackName).To(Equal("eksctl-test-cluster-nodegroup-12345"))
					Expect(out[0].NodeInstanceRoleARN).To(Equal("arn:aws:iam::1111:role/eks-nodes-base-role"))
					Expect(out[0].DesiredCapacity).To(Equal(7))
					Expect(out[0].MinSize).To(Equal(1))
					Expect(out[0].MaxSize).To(Equal(10))
				})
			})
		})
	})

	Describe("GetNodeGroupType", func() {

		createTags := func(tags map[string]string) []*cfn.Tag {
			cfnTags := make([]*cfn.Tag, 0)
			for k, v := range tags {
				cfnTags = append(cfnTags, &cfn.Tag{
					Key:   aws.String(k),
					Value: aws.String(v),
				})
			}
			return cfnTags
		}

		DescribeTable("with tag for the nodegroup type", func(inputTags map[string]string, expectedType api.NodeGroupType) {
			ngType, err := GetNodeGroupType(createTags(inputTags))

			if expectedType == "" {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(ngType).To(Equal(expectedType))
			}
		},

			Entry("finds the type of a managed nodegroup",
				map[string]string{
					api.NodeGroupNameTag: "mng-1",
					api.NodeGroupTypeTag: "managed",
				},
				api.NodeGroupTypeManaged),

			Entry("finds the type of an un-managed nodegroup",
				map[string]string{
					api.NodeGroupNameTag: "ng-1",
					api.NodeGroupTypeTag: "unmanaged",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("finds the type of an legacy un-managed nodegroup",
				map[string]string{
					api.OldNodeGroupNameTag: "ng-1",
					api.NodeGroupTypeTag:    "unmanaged",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("finds the type of a legacy un-managed nodegroup",
				map[string]string{
					api.OldNodeGroupIDTag: "ng-1",
					api.NodeGroupTypeTag:  "unmanaged",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("doesn't return the type if the stack tags don't contain any ng name tag",
				map[string]string{
					"some-other-tag": "ng-1",
					"name":           "ng-1",
				},
				api.NodeGroupType("")),
		)
		DescribeTable("for legacy ngs without tag for the type", func(inputTags map[string]string, expectedType api.NodeGroupType) {
			ngType, err := GetNodeGroupType(createTags(inputTags))

			if expectedType == "" {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
				Expect(ngType).To(Equal(expectedType))
			}
		},

			Entry("legacy ngs with old name tags are un-managed by default",
				map[string]string{
					api.NodeGroupNameTag: "ng-1",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("legacy ngs with old name tags are un-managed by default",
				map[string]string{
					api.OldNodeGroupNameTag: "ng-1",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("legacy ngs with old name tag group-id are un-managed by default",
				map[string]string{
					api.OldNodeGroupIDTag: "ng-1",
				},
				api.NodeGroupTypeUnmanaged),

			Entry("doesn't return the type if the stack tags don't contain any ng name tag",
				map[string]string{
					"some-other-tag": "ng-1",
					"name":           "ng-1",
				},
				api.NodeGroupType("")),
		)
	})
})
