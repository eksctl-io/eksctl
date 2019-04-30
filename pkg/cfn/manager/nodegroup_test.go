package manager

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	. "github.com/onsi/ginkgo"
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

				p.MockCloudFormation().On("GetTemplate", mock.MatchedBy(func(input *cfn.GetTemplateInput) bool {
					return input.StackName != nil && *input.StackName == "eksctl-test-cluster-nodegroup-12345"
				})).Return(&cfn.GetTemplateOutput{
					TemplateBody: aws.String(`{
						"Resources": {
							"NodeGroup": {
								"Properties": {
									"DesiredCapacity": 2,
									"MinSize": 1,
									"MaxSize: 3
								}
							}
						}
					}`),
				}, nil)
			})

			It("should be a no-op if attempting to scale to the existing desired capacity", func() {
				ng.Name = "12345"
				cap := 2
				ng.DesiredCapacity = &cap

				err := sc.ScaleNodeGroup(ng)

				Expect(err).NotTo(HaveOccurred())
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
								&cfn.Tag{
									Key:   aws.String(api.NodeGroupNameTag),
									Value: aws.String("12345"),
								},
							},
						},
					},
				}, nil)

				p.MockCloudFormation().On("DescribeStacks", mock.Anything).Return(nil, fmt.Errorf("DescribeStacks failed"))
			})

			Context("With no matching stacks", func() {
				BeforeEach(func() {
					clusterName = "test-cluster-non-existent"
				})

				JustBeforeEach(func() {
					out, err = sc.GetNodeGroupSummaries("")
				})

				It("should error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("getting nodegroup stacks: no eksctl-managed CloudFormation stacks found for \"test-cluster-non-existent\""))
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
				})
			})
		})
	})
})
