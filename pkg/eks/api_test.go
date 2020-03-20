package eks_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	. "github.com/weaveworks/eksctl/pkg/eks"
)

var _ = Describe("eksctl API", func() {

	Context("loading config files", func() {

		BeforeEach(func() {
			err := api.Register()
			Expect(err).ToNot(HaveOccurred())
		})

		It("should load a valid YAML config without error", func() {
			cfg, err := LoadConfigFromFile("../../examples/01-simple-cluster.yaml")
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg.Metadata.Name).To(Equal("cluster-1"))
			Expect(cfg.NodeGroups).To(HaveLen(1))
		})

		It("should load a valid JSON config without error", func() {
			cfg, err := LoadConfigFromFile("testdata/example.json")
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg.Metadata.Name).To(Equal("cluster-1"))
			Expect(cfg.NodeGroups).To(HaveLen(1))
		})

		It("should error when version is a float, not a string", func() {
			_, err := LoadConfigFromFile("testdata/bad-type-1.yaml")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(`loading config file "testdata/bad-type-1.yaml": v1alpha5.ClusterConfig.Metadata: v1alpha5.ClusterMeta.Version: ReadString: expects " or n, but found 1`))
		})

		It("should reject unknown field in a YAML config", func() {
			_, err := LoadConfigFromFile("testdata/bad-field-1.yaml")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(`loading config file "testdata/bad-field-1.yaml": error unmarshaling JSON: while decoding JSON: json: unknown field "zone"`))
		})

		It("should reject unknown field in a YAML config", func() {
			_, err := LoadConfigFromFile("testdata/bad-field-2.yaml")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(`loading config file "testdata/bad-field-2.yaml": error unmarshaling JSON: while decoding JSON: json: unknown field "bar"`))
		})

		It("should reject unknown field in a JSON config", func() {
			_, err := LoadConfigFromFile("testdata/bad-field-1.json")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(`loading config file "testdata/bad-field-1.json": error unmarshaling JSON: while decoding JSON: json: unknown field "nodes"`))
		})

		It("should reject old API version", func() {
			_, err := LoadConfigFromFile("testdata/old-version.json")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(HavePrefix(`loading config file "testdata/old-version.json": no kind "ClusterConfig" is registered for version "eksctl.io/v1alpha3" in scheme`))
		})

		It("should error when cannot read a file", func() {
			_, err := LoadConfigFromFile("../../examples/nothing.xml")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`reading config file "../../examples/nothing.xml": open ../../examples/nothing.xml: no such file or directory`))
		})
	})

	Context("Static AMI selection", func() {
		var (
			ng       *api.NodeGroup
			provider *mockprovider.MockProvider
		)
		BeforeEach(func() {
			ng = api.NewNodeGroup()
			ng.AMIFamily = api.DefaultNodeImageFamily

			provider = mockprovider.NewMockProvider()

			mockDescribeImages(provider, "ami-123", func(input *ec2.DescribeImagesInput) bool {
				return len(input.ImageIds) == 1
			})
		})

		It("should pick a valid AMI for GPU instances when AMI is static", func() {
			ng.AMI = "static"
			ng.InstanceType = "p2.xlarge"

			err := EnsureAMI(provider, "1.12", ng)

			Expect(err).ToNot(HaveOccurred())
			Expect(ng.AMI).To(Equal("ami-02551cb499388bebb"))
		})
		It("should pick a valid AMI for normal instances when AMI is static", func() {
			ng.AMI = "static"
			ng.InstanceType = "m5.xlarge"

			err := EnsureAMI(provider, "1.12", ng)

			Expect(err).ToNot(HaveOccurred())
			Expect(ng.AMI).To(Equal("ami-0267968f4310157f1"))
		})
		It("should pick a valid AMI for mixed normal instances", func() {
			ng.AMI = "static"
			ng.InstanceType = "mixed"
			ng.InstancesDistribution = &api.NodeGroupInstancesDistribution{
				InstanceTypes: []string{"t3.large", "m5.large", "m5a.large"},
			}

			err := EnsureAMI(provider, "1.12", ng)

			Expect(err).ToNot(HaveOccurred())
			Expect(ng.AMI).To(Equal("ami-0267968f4310157f1"))
		})
		It("should pick a GPU AMI for mixed instances with GPU instance types", func() {
			ng.AMI = "static"
			ng.InstanceType = "mixed"
			ng.InstancesDistribution = &api.NodeGroupInstancesDistribution{
				InstanceTypes: []string{"t3.large", "m5.large", "m5a.large", "p3.2xlarge"},
			}

			err := EnsureAMI(provider, "1.12", ng)

			Expect(err).ToNot(HaveOccurred())
			Expect(ng.AMI).To(Equal("ami-02551cb499388bebb"))
		})
	})

	Context("Dynamic AMI Resolution", func() {
		var (
			ng       *api.NodeGroup
			provider *mockprovider.MockProvider
		)

		BeforeEach(func() {
			ng = api.NewNodeGroup()
			ng.AMIFamily = api.DefaultNodeImageFamily

			provider = mockprovider.NewMockProvider()
			mockDescribeImages(provider, "ami-123", func(input *ec2.DescribeImagesInput) bool {
				return len(input.ImageIds) == 1
			})

		})

		testEnsureAMI := func(expectedAMI string) {
			err := EnsureAMI(provider, "1.14", ng)
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, ng.AMI).To(Equal(expectedAMI))
		}

		It("should resolve AMI using SSM Parameter Store by default", func() {
			provider.MockSSM().On("GetParameter", &ssm.GetParameterInput{
				Name: aws.String("/aws/service/eks/optimized-ami/1.14/amazon-linux-2/recommended/image_id"),
			}).Return(&ssm.GetParameterOutput{
				Parameter: &ssm.Parameter{
					Value: aws.String("ami-ssm"),
				},
			}, nil)

			testEnsureAMI("ami-ssm")
		})

		It("should use static resolution when specified", func() {
			ng.AMI = "static"
			testEnsureAMI("ami-0c13bb9cbfd007e56")
		})

		It("should fall back to auto resolution for Ubuntu", func() {
			ng.AMIFamily = api.NodeImageFamilyUbuntu1804
			mockDescribeImages(provider, "ami-ubuntu", func(input *ec2.DescribeImagesInput) bool {
				return *input.Owners[0] == "099720109477"
			})
			testEnsureAMI("ami-ubuntu")
		})

		It("should retrieve the AMI from EC2 when AMI is auto", func() {
			ng.AMI = "auto"
			ng.InstanceType = "p2.xlarge"
			mockDescribeImages(provider, "ami-auto", func(input *ec2.DescribeImagesInput) bool {
				return len(input.ImageIds) == 0
			})

			testEnsureAMI("ami-auto")
		})
	})

})

func mockDescribeImages(p *mockprovider.MockProvider, amiId string, matcher func(*ec2.DescribeImagesInput) bool) {
	p.MockEC2().On("DescribeImages", mock.MatchedBy(matcher)).
		Return(&ec2.DescribeImagesOutput{
			Images: []*ec2.Image{
				{
					ImageId:        aws.String(amiId),
					State:          aws.String("available"),
					OwnerId:        aws.String("123"),
					RootDeviceType: aws.String("ebs"),
					RootDeviceName: aws.String("/dev/sda1"),
					BlockDeviceMappings: []*ec2.BlockDeviceMapping{
						{
							DeviceName: aws.String("/dev/sda1"),
							Ebs: &ec2.EbsBlockDevice{
								Encrypted: aws.Bool(false),
							},
						},
					},
				},
			},
		}, nil)
}
