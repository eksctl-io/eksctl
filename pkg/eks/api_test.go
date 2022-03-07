package eks_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/stretchr/testify/mock"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	. "github.com/weaveworks/eksctl/pkg/eks"
)

var _ = Describe("eksctl API", func() {

	Context("loading config files", func() {

		BeforeEach(func() {
			err := api.Register()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should load a valid YAML config without error", func() {
			cfg, err := LoadConfigFromFile("../../examples/01-simple-cluster.yaml")
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Metadata.Name).To(Equal("cluster-1"))
			Expect(cfg.NodeGroups).To(HaveLen(1))
		})

		It("should load a valid JSON config without error", func() {
			cfg, err := LoadConfigFromFile("testdata/example.json")
			Expect(err).NotTo(HaveOccurred())
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
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
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

var _ = Describe("Setting Availability Zones", func() {
	var (
		provider *mockprovider.MockProvider
		cfg      *api.ClusterConfig
	)

	BeforeEach(func() {
		cfg = api.NewClusterConfig()
		provider = mockprovider.NewMockProvider()
	})

	When("the AZs were set as CLI params", func() {
		When("the given params contain enough AZs", func() {
			It("sets them as the AZs to be used", func() {
				err := eks.SetAvailabilityZones(cfg, []string{"us-east-2a", "us-east-2b"}, provider.EC2(), "")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("the given params contain too few AZs", func() {
			It("returns an error", func() {
				err := eks.SetAvailabilityZones(cfg, []string{"us-east-2a"}, provider.EC2(), "")
				Expect(err).To(MatchError("only 1 zone(s) specified [us-east-2a], 2 are required (can be non-unique)"))
			})
		})
	})

	When("the AZs were set in the config file", func() {
		When("the config file contains enough AZs", func() {
			It("sets them as the AZs to be used", func() {
				cfg.AvailabilityZones = []string{"us-east-2a", "us-east-2b"}
				err := eks.SetAvailabilityZones(cfg, []string{}, provider.EC2(), "")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("the config file contains too few AZs", func() {
			It("returns an error", func() {
				cfg.AvailabilityZones = []string{"us-east-2a"}
				err := eks.SetAvailabilityZones(cfg, []string{}, provider.EC2(), "")
				Expect(err).To(MatchError("only 1 zone(s) specified [us-east-2a], 2 are required (can be non-unique)"))
			})
		})
	})

	When("no AZs were set", func() {
		When("the call to fetch AZs fails", func() {
			It("returns an error", func() {
				region := "us-east-2"
				provider.MockEC2().On("DescribeAvailabilityZones", &ec2.DescribeAvailabilityZonesInput{
					Filters: []*ec2.Filter{{
						Name:   aws.String("region-name"),
						Values: []*string{aws.String(region)},
					}, {
						Name:   aws.String("state"),
						Values: []*string{aws.String(ec2.AvailabilityZoneStateAvailable)},
					}},
				}).Return(&ec2.DescribeAvailabilityZonesOutput{}, fmt.Errorf("err"))
				err := eks.SetAvailabilityZones(cfg, []string{}, provider.EC2(), region)
				Expect(err).To(MatchError("getting availability zones: error getting availability zones for region us-east-2: err"))
			})
		})

		When("the call to fetch AZs succeeds", func() {
			It("sets random AZs", func() {
				region := "us-east-2"
				provider.MockEC2().On("DescribeAvailabilityZones", &ec2.DescribeAvailabilityZonesInput{
					Filters: []*ec2.Filter{{
						Name:   aws.String("region-name"),
						Values: []*string{aws.String(region)},
					}, {
						Name:   aws.String("state"),
						Values: []*string{aws.String(ec2.AvailabilityZoneStateAvailable)},
					}},
				}).Return(&ec2.DescribeAvailabilityZonesOutput{
					AvailabilityZones: []*ec2.AvailabilityZone{
						{
							GroupName: aws.String("name"),
							ZoneName:  aws.String(region),
							ZoneId:    aws.String("id"),
						},
						{
							GroupName: aws.String("name"),
							ZoneName:  aws.String(region),
							ZoneId:    aws.String("id"),
						}},
				}, nil)
				err := eks.SetAvailabilityZones(cfg, []string{}, provider.EC2(), region)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
