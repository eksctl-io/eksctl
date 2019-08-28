package eks_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/ami"
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

	Context("AMI selection", func() {
		var (
			cfg *api.ClusterConfig
			ctl *ClusterProvider
			ng  *api.NodeGroup
			p   *mockprovider.MockProvider
		)
		BeforeEach(func() {
			ami.DefaultResolvers = []ami.Resolver{&ami.StaticGPUResolver{}, &ami.StaticDefaultResolver{}}

			cfg = &api.ClusterConfig{}
			ng = cfg.NewNodeGroup()
			ng.AMIFamily = api.DefaultNodeImageFamily

			p = mockprovider.NewMockProvider()

			ctl = &ClusterProvider{
				Provider: p,
				Status:   &ProviderStatus{},
			}

			mockDescribeImages(p, "something", "abc123")
		})

		It("should retrieve the AMI from EC2 when AMI is auto", func() {
			ng.AMI = "auto"
			ng.InstanceType = "p2.xlarge"

			err := ctl.EnsureAMI("1.12", ng)

			Expect(err).ToNot(HaveOccurred())
			Expect(ng.AMI).To(Equal("abc123"))
		})
		It("should pick a valid AMI for GPU instances when AMI is static", func() {
			ng.AMI = "static"
			ng.InstanceType = "p2.xlarge"

			err := ctl.EnsureAMI("1.12", ng)

			Expect(err).ToNot(HaveOccurred())
			Expect(ng.AMI).To(Equal("ami-0c9156d7fcd3c2948"))
		})
		It("should pick a valid AMI for normal instances when AMI is static", func() {
			ng.AMI = "static"
			ng.InstanceType = "m5.xlarge"

			err := ctl.EnsureAMI("1.12", ng)

			Expect(err).ToNot(HaveOccurred())
			Expect(ng.AMI).To(Equal("ami-0b520e822d42998c1"))
		})
		It("should pick a valid AMI for mixed normal instances", func() {
			ng.AMI = "static"
			ng.InstanceType = "mixed"
			ng.InstancesDistribution = &api.NodeGroupInstancesDistribution{
				InstanceTypes: []string{"t3.large", "m5.large", "m5a.large"},
			}

			err := ctl.EnsureAMI("1.12", ng)

			Expect(err).ToNot(HaveOccurred())
			Expect(ng.AMI).To(Equal("ami-0b520e822d42998c1"))
		})
		It("should pick a GPU AMI for mixed instances with GPU instance types", func() {
			ng.AMI = "static"
			ng.InstanceType = "mixed"
			ng.InstancesDistribution = &api.NodeGroupInstancesDistribution{
				InstanceTypes: []string{"t3.large", "m5.large", "m5a.large", "p3.2xlarge"},
			}

			err := ctl.EnsureAMI("1.12", ng)

			Expect(err).ToNot(HaveOccurred())
			Expect(ng.AMI).To(Equal("ami-0c9156d7fcd3c2948"))
		})
	})
})

func mockDescribeImages(p *mockprovider.MockProvider, expectedNamePattern string, amiId string) {
	p.MockEC2().On("DescribeImages",
		mock.MatchedBy(func(input *ec2.DescribeImagesInput) bool {
			return true
		})).
		Return(&ec2.DescribeImagesOutput{
			Images: []*ec2.Image{
				&ec2.Image{
					ImageId:        aws.String(amiId),
					State:          aws.String("available"),
					OwnerId:        aws.String("123"),
					RootDeviceType: aws.String("ebs"),
					BlockDeviceMappings: []*ec2.BlockDeviceMapping{
						{
							Ebs: &ec2.EbsBlockDevice{
								Encrypted: aws.Bool(false),
							},
						},
					},
				},
			},
		}, nil)
}
