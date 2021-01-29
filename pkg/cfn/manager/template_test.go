package manager

import (
	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
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

	rawJsonTemplate := `{"Outputs":{"ARN":{"Value":{"Fn::GetAtt":["ControlPlane","Arn"]},"Export":{"Name":{"Fn::Sub":"${AWS::StackName}::ARN"}}}}}`
	rawYamlTemplate := `
Outputs:
  ARN:
    Value: !GetAtt
      - ControlPlane
      - Arn
    Export:
      Name: !Sub '${AWS::StackName}::ARN'`
	expectedJSONResponse := `{
  "Outputs": {
    "ARN": {
      "Value": {
        "Fn::GetAtt": [
          "ControlPlane",
          "Arn"
        ]
      },
      "Export": {
        "Name": {
          "Fn::Sub": "${AWS::StackName}::ARN"
        }
      }
    }
  }
}`

	Describe("GetTemplate", func() {
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
		})

		Context("With an existing stack name", func() {
			Context("and a json template response", func() {
				BeforeEach(func() {
					p.MockCloudFormation().On("GetTemplate", mock.MatchedBy(func(input *cfn.GetTemplateInput) bool {
						return input.StackName != nil && *input.StackName == "foobar"
					})).Return(&cfn.GetTemplateOutput{
						TemplateBody: &rawJsonTemplate,
					}, nil)
					out, err = sc.GetStackTemplate("foobar")
				})

				It("returns the template", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "GetTemplate", 1)).To(BeTrue())
					Expect(out).To(Equal(expectedJSONResponse))
				})
			})

			Context("and a yaml template response", func() {
				BeforeEach(func() {
					p.MockCloudFormation().On("GetTemplate", mock.MatchedBy(func(input *cfn.GetTemplateInput) bool {
						return input.StackName != nil && *input.StackName == "foobar"
					})).Return(&cfn.GetTemplateOutput{
						TemplateBody: &rawYamlTemplate,
					}, nil)
					out, err = sc.GetStackTemplate("foobar")
				})

				It("returns the template in json format", func() {
					Expect(p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "GetTemplate", 1)).To(BeTrue())
					Expect(err).NotTo(HaveOccurred())
					Expect(out).To(Equal(expectedJSONResponse))
				})
			})
		})

		Context("With a non-existing stack name", func() {
			BeforeEach(func() {
				p.MockCloudFormation().On("GetTemplate", mock.Anything).Return(nil, errors.New("GetTemplate failed"))
				out, err = sc.GetStackTemplate("non_existing_stack")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "GetTemplate", 1)).To(BeTrue())
			})
		})

		When("the response isn't valid json or yaml", func() {
			BeforeEach(func() {
				p.MockCloudFormation().On("GetTemplate", mock.MatchedBy(func(input *cfn.GetTemplateInput) bool {
					return input.StackName != nil && *input.StackName == "foobar"
				})).Return(&cfn.GetTemplateOutput{
					TemplateBody: aws.String("~123"),
				}, nil)

				out, err = sc.GetStackTemplate("foobar")
			})

			It("returns an error", func() {
				Expect(p.MockCloudFormation().AssertNumberOfCalls(GinkgoT(), "GetTemplate", 1)).To(BeTrue())
				Expect(out).To(Equal(""))
				Expect(err).To(MatchError(ContainSubstring("failed to parse GetStackTemplate response")))
			})
		})
	})
})
