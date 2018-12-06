package manager

import (
	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"errors"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("StackCollection Template", func() {
	var (
		cc *api.ClusterConfig
		sc *StackCollection

		p *mockprovider.MockProvider
	)

	testAZs := []string{"us-west-2b", "us-west-2a", "us-west-2c"}

	newClusterConfig := func(clusterName string) *api.ClusterConfig {
		cfg := api.NewClusterConfig()
		ng := cfg.NewNodeGroup()

		cfg.Metadata.Region = "us-west-2"
		cfg.Metadata.Name = clusterName
		cfg.AvailabilityZones = testAZs
		ng.InstanceType = "t2.medium"
		ng.AMIFamily = "AmazonLinux2"

		*cfg.VPC.CIDR = api.DefaultCIDR()

		return cfg
	}

	Describe("GetTemplate", func() {
		Context("With a cluster name", func() {
			var (
				clusterName string
				err         error
				out         string
			)

			BeforeEach(func() {
				clusterName = "test-cluster"

				p = mockprovider.NewMockProvider()

				cc = newClusterConfig(clusterName)

				sc = NewStackCollection(p, cc)

				p.MockCloudFormation().On("GetTemplate", mock.MatchedBy(func(input *cfn.GetTemplateInput) bool {
					return input.StackName != nil && *input.StackName == "foobar"
				})).Return(&cfn.GetTemplateOutput{
					TemplateBody: aws.String("TEMPLATE_BODY"),
				}, nil)

				p.MockCloudFormation().On("GetTemplate", mock.Anything).Return(nil, errors.New("GetTemplate failed"))
			})

			Context("With a non-existing stack name", func() {
				JustBeforeEach(func() {
					out, err = sc.GetStackTemplate("non_existing_stack")
				})

				It("should error", func() {
					Expect(err).To(HaveOccurred())
				})

				It("should have called AWS CloudFormation service once", func() {
					Expect(p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "GetTemplate", 1)).To(BeTrue())
				})
			})

			Context("With an existing stack name", func() {
				JustBeforeEach(func() {
					out, err = sc.GetStackTemplate("foobar")
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should have called AWS CloudFormation service once", func() {
					Expect(p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "GetTemplate", 1)).To(BeTrue())
				})

				It("the output should equal the expectation", func() {
					Expect(out).To(Equal("TEMPLATE_BODY"))
				})
			})
		})
	})
})
