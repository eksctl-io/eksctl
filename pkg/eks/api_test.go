package eks_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"
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

		testEnsureAMI := func(matcher gomegatypes.GomegaMatcher) {
			err := ResolveAMI(provider, "1.14", ng)
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, ng.AMI).To(matcher)
		}

		It("should resolve AMI using SSM Parameter Store by default", func() {
			provider.MockSSM().On("GetParameter", &ssm.GetParameterInput{
				Name: aws.String("/aws/service/eks/optimized-ami/1.14/amazon-linux-2/recommended/image_id"),
			}).Return(&ssm.GetParameterOutput{
				Parameter: &ssm.Parameter{
					Value: aws.String("ami-ssm"),
				},
			}, nil)

			testEnsureAMI(Equal("ami-ssm"))
		})

		It("should fall back to auto resolution for Ubuntu", func() {
			ng.AMIFamily = api.NodeImageFamilyUbuntu1804
			mockDescribeImages(provider, "ami-ubuntu", func(input *ec2.DescribeImagesInput) bool {
				return *input.Owners[0] == "099720109477"
			})
			testEnsureAMI(Equal("ami-ubuntu"))
		})

		It("should retrieve the AMI from EC2 when AMI is auto", func() {
			ng.AMI = "auto"
			ng.InstanceType = "p2.xlarge"
			mockDescribeImages(provider, "ami-auto", func(input *ec2.DescribeImagesInput) bool {
				return len(input.ImageIds) == 0
			})

			testEnsureAMI(Equal("ami-auto"))
		})
	})

})

func mockDescribeImages(p *mockprovider.MockProvider, amiID string, matcher func(*ec2.DescribeImagesInput) bool) {
	p.MockEC2().On("DescribeImages", mock.MatchedBy(matcher)).
		Return(&ec2.DescribeImagesOutput{
			Images: []*ec2.Image{
				{
					ImageId:        aws.String(amiID),
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
