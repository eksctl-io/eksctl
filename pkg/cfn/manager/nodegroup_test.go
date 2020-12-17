package manager

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
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

		return ng
	}

	Describe("ScaleNodeGroup", func() {
		var (
			ng *api.NodeGroup
		)

		JustBeforeEach(func() {
			p = mockprovider.NewMockProvider()
		})

		Context("With an existing NodeGroup", func() {
			JustBeforeEach(func() {
				cc = newClusterConfig("test-cluster")
				ng = newNodeGroup(cc)
				sc = NewStackCollection(p, cc)

				p.MockCloudFormation().
					On("DescribeStacks", mock.MatchedBy(func(input *cfn.DescribeStacksInput) bool {
						return input.StackName != nil && *input.StackName == "eksctl-test-cluster-nodegroup-12345"
					})).Return(&cfn.DescribeStacksOutput{
					Stacks: []*Stack{
						{
							Tags: []*cfn.Tag{
								{
									Key:   aws.String(api.NodeGroupNameTag),
									Value: aws.String("12345"),
								},
							},
						},
					},
				}, nil).
					On("GetTemplate", mock.MatchedBy(func(input *cfn.GetTemplateInput) bool {
						return input.StackName != nil && *input.StackName == "eksctl-test-cluster-nodegroup-12345"
					})).Return(&cfn.GetTemplateOutput{
					TemplateBody: aws.String(`{
						"Resources": {
							"NodeGroup": {
								"Properties": {
									"DesiredCapacity": 2,
									"MinSize": 1,
									"MaxSize": 3
								}
							}
						}
					}`),
				}, nil)
			})

			It("should be a no-op if attempting to scale to the existing desired capacity", func() {
				ng.Name = "12345"
				capacity := 2
				ng.DesiredCapacity = &capacity
				err := sc.ScaleNodeGroup(ng)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should be a no-op if attempting to scale to the existing desired capacity, min size", func() {
				ng.Name = "12345"
				minSize := 1
				capacity := 2
				ng.MinSize = &minSize
				ng.DesiredCapacity = &capacity
				err := sc.ScaleNodeGroup(ng)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should be a no-op if attempting to scale to the existing desired capacity, max size", func() {
				ng.Name = "12345"
				capacity := 2
				maxSize := 3
				ng.DesiredCapacity = &capacity
				ng.MaxSize = &maxSize
				err := sc.ScaleNodeGroup(ng)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should be a no-op if attempting to scale to the existing desired capacity, min size and max size", func() {
				ng.Name = "12345"
				minSize := 1
				capacity := 2
				maxSize := 3
				ng.MinSize = &minSize
				ng.DesiredCapacity = &capacity
				ng.MaxSize = &maxSize
				err := sc.ScaleNodeGroup(ng)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should be a error if the desired capacity is greater than current CF maxSize", func() {
				ng.Name = "12345"
				capacity := 5
				ng.DesiredCapacity = &capacity
				err := sc.ScaleNodeGroup(ng)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("the desired nodes 5 is greater than current nodes-max/maxSize 3"))
			})

			It("should be a error if the desired capacity is less than current CF minSize", func() {
				ng.Name = "12345"
				capacity := 0
				ng.DesiredCapacity = &capacity
				err := sc.ScaleNodeGroup(ng)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("the desired nodes 0 is less than current nodes-min/minSize 1"))
			})
		})
	})

	Describe("GetNodeGroupSummaries", func() {
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
					TemplateBody: aws.String("TEMPLATE_BODY"),
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
						PhysicalResourceId: aws.String("eksctl-test-cluster-nodegroup-123451-NodeGroup-1N68LL8H1EH27"),
					},
				}, nil)

				p.MockCloudFormation().On("DescribeStackResource", mock.Anything).Return(nil, fmt.Errorf("DescribeStackResource failed"))

			})

			Context("With no matching stacks", func() {
				BeforeEach(func() {
					clusterName = "test-cluster-non-existent"
				})

				JustBeforeEach(func() {
					out, err = sc.GetNodeGroupSummaries("")
				})

				It("should not error", func() {
					Expect(err).ToNot(HaveOccurred())
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
					out, err = sc.GetNodeGroupSummaries("")
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
				Expect(err).ToNot(HaveOccurred())
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
				Expect(err).ToNot(HaveOccurred())
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
