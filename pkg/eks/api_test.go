package eks_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	. "github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
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
			err := ResolveAMI(context.Background(), provider, "1.14", ng)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
			ExpectWithOffset(1, ng.AMI).To(matcher)
		}

		It("should resolve AMI using SSM Parameter Store by default", func() {
			provider.MockSSM().On("GetParameter", mock.Anything, &ssm.GetParameterInput{
				Name: aws.String("/aws/service/eks/optimized-ami/1.14/amazon-linux-2/recommended/image_id"),
			}).Return(&ssm.GetParameterOutput{
				Parameter: &ssmtypes.Parameter{
					Value: aws.String("ami-ssm"),
				},
			}, nil)

			testEnsureAMI(Equal("ami-ssm"))
		})

		It("should fall back to auto resolution for Ubuntu", func() {
			ng.AMIFamily = api.NodeImageFamilyUbuntu1804
			mockDescribeImages(provider, "ami-ubuntu", func(input *ec2.DescribeImagesInput) bool {
				return input.Owners[0] == "099720109477"
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
	p.MockEC2().On("DescribeImages", mock.Anything, mock.MatchedBy(matcher)).
		Return(&ec2.DescribeImagesOutput{
			Images: []ec2types.Image{
				{
					ImageId:        aws.String(amiID),
					State:          ec2types.ImageStateAvailable,
					OwnerId:        aws.String("123"),
					RootDeviceType: ec2types.DeviceTypeEbs,
					RootDeviceName: aws.String("/dev/sda1"),
					BlockDeviceMappings: []ec2types.BlockDeviceMapping{
						{
							DeviceName: aws.String("/dev/sda1"),
							Ebs: &ec2types.EbsBlockDevice{
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
				userProvider, err := eks.SetAvailabilityZones(context.Background(), cfg, []string{"us-east-2a", "us-east-2b"}, provider.EC2(), "")
				Expect(err).NotTo(HaveOccurred())
				Expect(userProvider).To(BeTrue())
			})
		})

		When("the given params contain too few AZs", func() {
			It("returns an error", func() {
				userProvider, err := eks.SetAvailabilityZones(context.Background(), cfg, []string{"us-east-2a"}, provider.EC2(), "")
				Expect(err).To(MatchError("only 1 zone(s) specified [us-east-2a], 2 are required (can be non-unique)"))
				Expect(userProvider).To(BeFalse())
			})
		})
	})

	When("the AZs were set in the config file", func() {
		When("the config file contains enough AZs", func() {
			It("sets them as the AZs to be used", func() {
				cfg.AvailabilityZones = []string{"us-east-2a", "us-east-2b"}
				userProvider, err := eks.SetAvailabilityZones(context.Background(), cfg, []string{}, provider.EC2(), "")
				Expect(err).NotTo(HaveOccurred())
				Expect(userProvider).To(BeTrue())
			})
		})

		When("the config file contains too few AZs", func() {
			It("returns an error", func() {
				cfg.AvailabilityZones = []string{"us-east-2a"}
				userProvider, err := eks.SetAvailabilityZones(context.Background(), cfg, []string{}, provider.EC2(), "")
				Expect(err).To(MatchError("only 1 zone(s) specified [us-east-2a], 2 are required (can be non-unique)"))
				Expect(userProvider).To(BeFalse())
			})
		})
	})

	When("no AZs were set", func() {
		When("the call to fetch AZs fails", func() {
			It("returns an error", func() {
				region := "us-east-2"
				provider.MockEC2().On("DescribeAvailabilityZones", mock.Anything, &ec2.DescribeAvailabilityZonesInput{
					Filters: []ec2types.Filter{{
						Name:   aws.String("region-name"),
						Values: []string{region},
					}, {
						Name:   aws.String("state"),
						Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
					}, {
						Name:   aws.String("zone-type"),
						Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
					}},
				}).Return(&ec2.DescribeAvailabilityZonesOutput{}, fmt.Errorf("err"))
				userProvider, err := eks.SetAvailabilityZones(context.Background(), cfg, []string{}, provider.EC2(), region)
				Expect(err).To(MatchError("getting availability zones: error getting availability zones for region us-east-2: err"))
				Expect(userProvider).To(BeFalse())
			})
		})

		When("the call to fetch AZs succeeds", func() {
			It("sets random AZs", func() {
				region := "us-east-2"
				provider.MockEC2().On("DescribeAvailabilityZones", mock.Anything, &ec2.DescribeAvailabilityZonesInput{
					Filters: []ec2types.Filter{{
						Name:   aws.String("region-name"),
						Values: []string{region},
					}, {
						Name:   aws.String("state"),
						Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
					}, {
						Name:   aws.String("zone-type"),
						Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
					}},
				}).Return(&ec2.DescribeAvailabilityZonesOutput{
					AvailabilityZones: []ec2types.AvailabilityZone{
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
				userProvider, err := eks.SetAvailabilityZones(context.Background(), cfg, []string{}, provider.EC2(), region)
				Expect(err).NotTo(HaveOccurred())
				Expect(userProvider).To(BeFalse())
			})
		})
	})
})

var _ = Describe("CheckInstanceAvailability", func() {
	var (
		provider *mockprovider.MockProvider
		cfg      *api.ClusterConfig
	)

	BeforeEach(func() {
		cfg = api.NewClusterConfig()
		provider = mockprovider.NewMockProvider()
		provider.MockEC2().On("DescribeInstanceTypeOfferings", mock.Anything, &ec2.DescribeInstanceTypeOfferingsInput{
			Filters: []ec2types.Filter{
				{
					Name:   aws.String("instance-type"),
					Values: []string{"t2.nano"},
				},
			},
			LocationType: ec2types.LocationTypeAvailabilityZone,
			MaxResults:   aws.Int32(100),
		}).Return(&ec2.DescribeInstanceTypeOfferingsOutput{
			InstanceTypeOfferings: []ec2types.InstanceTypeOffering{
				{
					InstanceType: "t2.nano",
					Location:     aws.String("dummy-zone-1a"),
					LocationType: "availability-zone",
				},
			},
		}, nil)
	})

	When("instance not available in nodegroup AZ", func() {
		It("errors", func() {
			cfg.NodeGroups = []*api.NodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:              "ng-1",
						InstanceType:      "t2.nano",
						AvailabilityZones: []string{"dummy-zone-1b"},
					},
				},
			}
			err := eks.CheckInstanceAvailability(context.Background(), cfg, provider.EC2())
			Expect(err).To(MatchError(`none of the provided AZs "dummy-zone-1b" support instance type t2.nano in nodegroup ng-1`))
		})
	})
	When("uses instance distribution", func() {
		When("azs aren't supported", func() {
			It("errors", func() {
				cfg.NodeGroups = []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:              "ng-1",
							AvailabilityZones: []string{"dummy-zone-1b"},
						},
						InstancesDistribution: &api.NodeGroupInstancesDistribution{
							InstanceTypes: []string{"t2.nano"},
						},
					},
				}
				err := eks.CheckInstanceAvailability(context.Background(), cfg, provider.EC2())
				Expect(err).To(MatchError(`none of the provided AZs "dummy-zone-1b" support instance type t2.nano in nodegroup ng-1`))
			})
		})
	})
	When("instance available in nodegroup AZ", func() {
		It("allows the usage of the instance", func() {
			cfg.NodeGroups = []*api.NodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:              "ng-1",
						AvailabilityZones: []string{"dummy-zone-1a"},
					},
					InstancesDistribution: &api.NodeGroupInstancesDistribution{
						InstanceTypes: []string{"t2.nano"},
					},
				},
			}
			Expect(eks.CheckInstanceAvailability(context.Background(), cfg, provider.EC2())).To(Succeed())
		})
	})
	When("mixed instances are used", func() {
		It("allows the usage of the instance", func() {
			cfg.NodeGroups = []*api.NodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:              "ng-1",
						AvailabilityZones: []string{"dummy-zone-1a"},
						InstanceType:      "mixed",
					},
					InstancesDistribution: &api.NodeGroupInstancesDistribution{
						InstanceTypes: []string{"t2.nano"},
					},
				},
			}
			Expect(eks.CheckInstanceAvailability(context.Background(), cfg, provider.EC2())).To(Succeed())
		})
	})
	When("instance available in nodegroup AZ", func() {
		It("list is deduplicated", func() {
			cfg.NodeGroups = []*api.NodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:              "ng-1",
						AvailabilityZones: []string{"dummy-zone-1a"},
					},
					InstancesDistribution: &api.NodeGroupInstancesDistribution{
						InstanceTypes: []string{"t2.nano", "t2.nano"},
					},
				},
			}
			Expect(eks.CheckInstanceAvailability(context.Background(), cfg, provider.EC2())).To(Succeed())
		})
	})
	When("instance not available in any of the global AZs", func() {
		It("errors", func() {
			cfg.AvailabilityZones = []string{"dummy-zone-1b"}
			cfg.NodeGroups = []*api.NodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name: "ng-1",
					},
					InstancesDistribution: &api.NodeGroupInstancesDistribution{
						InstanceTypes: []string{"t2.nano"},
					},
				},
			}
			err := eks.CheckInstanceAvailability(context.Background(), cfg, provider.EC2())
			Expect(err).To(MatchError(`none of the provided AZs "dummy-zone-1b" support instance type t2.nano in nodegroup ng-1`))
		})
	})
	When("az is overridden by local nodegroup's AZ", func() {
		It("uses the az defined in the nodegroup", func() {
			cfg.AvailabilityZones = []string{"dummy-zone-1b"}
			cfg.NodeGroups = []*api.NodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:              "ng-1",
						AvailabilityZones: []string{"dummy-zone-1a"},
					},
					InstancesDistribution: &api.NodeGroupInstancesDistribution{
						InstanceTypes: []string{"t2.nano"},
					},
				},
			}
			Expect(eks.CheckInstanceAvailability(context.Background(), cfg, provider.EC2())).To(Succeed())
		})
	})
	When("more than one AZ is available and more than one AZ is returned", func() {
		When("one of the instances doesn't support any AZs", func() {
			It("errors", func() {
				provider.MockEC2().On("DescribeInstanceTypeOfferings", mock.Anything, &ec2.DescribeInstanceTypeOfferingsInput{
					Filters: []ec2types.Filter{
						{
							Name:   aws.String("instance-type"),
							Values: []string{"t2.large", "t2.micro", "t2.nano"},
						},
					},
					LocationType: ec2types.LocationTypeAvailabilityZone,
					MaxResults:   aws.Int32(100),
				}).Return(&ec2.DescribeInstanceTypeOfferingsOutput{
					InstanceTypeOfferings: []ec2types.InstanceTypeOffering{
						{
							InstanceType: "t2.nano",
							Location:     aws.String("dummy-zone-1a"),
							LocationType: "availability-zone",
						},
						{
							InstanceType: "t2.micro",
							Location:     aws.String("dummy-zone-1b"),
							LocationType: "availability-zone",
						},
						{
							InstanceType: "t2.large",
							Location:     aws.String("dummy-zone-1c"),
							LocationType: "availability-zone",
						},
					},
				}, nil)
				cfg.NodeGroups = []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:              "ng-1",
							AvailabilityZones: []string{"dummy-zone-1a", "dummy-zone-1b"},
							InstanceType:      "t2.large",
						},
						InstancesDistribution: &api.NodeGroupInstancesDistribution{
							InstanceTypes: []string{"t2.nano", "t2.micro"},
						},
					},
				}
				cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:              "mng-1",
							AvailabilityZones: []string{"dummy-zone-1a", "dummy-zone-1b"},
							InstanceType:      "t2.large",
						},
					},
				}
				err := eks.CheckInstanceAvailability(context.Background(), cfg, provider.EC2())
				Expect(err).To(MatchError(`none of the provided AZs "dummy-zone-1a,dummy-zone-1b" support instance type t2.large in nodegroup mng-1`))
			})
		})
		When("all instances are available in at least one of the provided AZs", func() {
			It("allows the selection", func() {
				provider.MockEC2().On("DescribeInstanceTypeOfferings", mock.Anything, &ec2.DescribeInstanceTypeOfferingsInput{
					Filters: []ec2types.Filter{
						{
							Name:   aws.String("instance-type"),
							Values: []string{"t2.large", "t2.micro", "t2.nano"},
						},
					},
					LocationType: ec2types.LocationTypeAvailabilityZone,
					MaxResults:   aws.Int32(100),
				}).Return(&ec2.DescribeInstanceTypeOfferingsOutput{
					InstanceTypeOfferings: []ec2types.InstanceTypeOffering{
						{
							InstanceType: "t2.nano",
							Location:     aws.String("dummy-zone-1a"),
							LocationType: "availability-zone",
						},
						{
							InstanceType: "t2.nano",
							Location:     aws.String("dummy-zone-1b"),
							LocationType: "availability-zone",
						},
						{
							InstanceType: "t2.micro",
							Location:     aws.String("dummy-zone-1b"),
							LocationType: "availability-zone",
						},
						{
							InstanceType: "t2.large",
							Location:     aws.String("dummy-zone-1a"),
							LocationType: "availability-zone",
						},
						{
							InstanceType: "t2.large",
							Location:     aws.String("dummy-zone-1c"),
							LocationType: "availability-zone",
						},
					},
				}, nil)
				cfg.NodeGroups = []*api.NodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:              "ng-1",
							AvailabilityZones: []string{"dummy-zone-1a", "dummy-zone-1b"},
							InstanceType:      "t2.large",
						},
						InstancesDistribution: &api.NodeGroupInstancesDistribution{
							InstanceTypes: []string{"t2.nano", "t2.micro"},
						},
					},
				}
				cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name:              "mng-1",
							AvailabilityZones: []string{"dummy-zone-1a", "dummy-zone-1b"},
							InstanceType:      "t2.large",
						},
					},
				}
				Expect(eks.CheckInstanceAvailability(context.Background(), cfg, provider.EC2())).To(Succeed())
			})
		})
	})
})
