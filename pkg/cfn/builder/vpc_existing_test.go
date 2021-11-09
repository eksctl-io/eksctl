package builder_test

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/builder/fakes"
	"github.com/weaveworks/eksctl/pkg/eks/mocks"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

var _ = Describe("Existing VPC", func() {
	var (
		vpcRs   *builder.ExistingVPCResourceSet
		cfg     *api.ClusterConfig
		mockEC2 *mocks.EC2API
	)

	BeforeEach(func() {
		cfg = api.NewClusterConfig()
		cfg.VPC = vpcConfig()
		cfg.AvailabilityZones = []string{azA, azB}
		cfg.VPC.ID = "custom-vpc"
		mockEC2 = &mocks.EC2API{}
	})

	JustBeforeEach(func() {
		mockEC2.On("DescribeVpcs", &awsec2.DescribeVpcsInput{
			VpcIds: aws.StringSlice([]string{"custom-vpc"}),
		}).Return(&awsec2.DescribeVpcsOutput{
			Vpcs: []*awsec2.Vpc{
				{
					VpcId: aws.String("custom-vpc"),
					Ipv6CidrBlockAssociationSet: []*awsec2.VpcIpv6CidrBlockAssociation{
						{
							Ipv6CidrBlock: aws.String("foo"),
						},
					},
				},
			},
		}, nil)
		vpcRs = builder.NewExistingVPCResourceSet(builder.NewRS(), cfg, mockEC2)
	})

	Describe("CreateTemplate", func() {
		var (
			addErr        error
			vpcID         *gfnt.Value
			subnetDetails *builder.SubnetDetails
			vpcTemplate   *fakes.FakeTemplate
		)

		JustBeforeEach(func() {
			vpcID, subnetDetails, addErr = vpcRs.CreateTemplate()
			vpcTemplate = &fakes.FakeTemplate{}
			templateBody, err := vpcRs.RenderJSON()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(json.Unmarshal(templateBody, vpcTemplate)).To(Succeed())
		})

		It("uses the existing VPC", func() {
			Expect(addErr).NotTo(HaveOccurred())
			By("the correct VPC resource values are loaded into the VPCResource")
			Expect(vpcID).To(Equal(gfnt.NewString("custom-vpc")))

			By("no resources are added to the set")
			Expect(vpcTemplate.Resources).To(HaveLen(0))

			By("the private subnet resource values are loaded into the VPCResource")
			Expect(subnetDetails.Private).To(HaveLen(2))
			Expect(subnetDetails.Private).To(ContainElement(builder.SubnetResource{
				Subnet:           gfnt.NewString(privateSubnet2),
				AvailabilityZone: azB,
			}))
			Expect(subnetDetails.Private).To(ContainElement(builder.SubnetResource{
				Subnet:           gfnt.NewString(privateSubnet1),
				AvailabilityZone: azA,
			}))

			By("the public subnet resource values are loaded into the VPCResource")
			Expect(subnetDetails.Public).To(HaveLen(2))
			Expect(subnetDetails.Public).To(ContainElement(builder.SubnetResource{
				Subnet:           gfnt.NewString(publicSubnet2),
				AvailabilityZone: azB,
			}))
			Expect(subnetDetails.Public).To(ContainElement(builder.SubnetResource{
				Subnet:           gfnt.NewString(publicSubnet1),
				AvailabilityZone: azA,
			}))
		})

		When("and the VPC does not exist", func() {
			BeforeEach(func() {
				mockEC2.On("DescribeVpcs", &awsec2.DescribeVpcsInput{
					VpcIds: aws.StringSlice([]string{"custom-vpc"}),
				}).Return(&awsec2.DescribeVpcsOutput{
					Vpcs: []*awsec2.Vpc{},
				}, nil)
			})

			It("errors", func() {
				Expect(addErr).To(MatchError("VPC \"custom-vpc\" does not exist"))
			})
		})

		When("describing the VPC fails", func() {
			BeforeEach(func() {
				mockEC2.On("DescribeVpcs", &awsec2.DescribeVpcsInput{
					VpcIds: aws.StringSlice([]string{"custom-vpc"}),
				}).Return(nil, fmt.Errorf("foo"))
			})

			It("errors", func() {
				Expect(addErr).To(MatchError("failed to describe VPC \"custom-vpc\": foo"))
			})
		})

		Context("when ipv6 is true", func() {
			BeforeEach(func() {
				cfg.VPC.IPFamily = api.IPV6Family
			})

			It("succeeds", func() {
				Expect(addErr).NotTo(HaveOccurred())
			})

			When("and the VPC does not have ipv6 enabled", func() {
				BeforeEach(func() {
					mockEC2.On("DescribeVpcs", &awsec2.DescribeVpcsInput{
						VpcIds: aws.StringSlice([]string{"custom-vpc"}),
					}).Return(&awsec2.DescribeVpcsOutput{
						Vpcs: []*awsec2.Vpc{
							{
								VpcId: aws.String("custom-vpc"),
							},
						},
					}, nil)
				})

				It("errors", func() {
					Expect(addErr).To(MatchError("VPC \"custom-vpc\" does not have any associated IPv6 cidr blocks"))
				})
			})
		})

		Context("PrivateCluster is enabled", func() {
			var rtOutput *ec2.DescribeRouteTablesOutput

			BeforeEach(func() {
				cfg.PrivateCluster.Enabled = true
				rtOutput = makeRTOutput([]string{privateSubnet2, privateSubnet1}, false)
			})

			Context("ec2 call succeeds", func() {
				BeforeEach(func() {
					mockResultFn := func(_ *ec2.DescribeRouteTablesInput) *ec2.DescribeRouteTablesOutput {
						return rtOutput
					}

					mockEC2.On("DescribeRouteTables", mock.MatchedBy(func(input *ec2.DescribeRouteTablesInput) bool {
						return len(input.Filters) > 0
					})).Return(mockResultFn, nil)
				})

				It("the private subnet resource values are loaded into the VPCResource with route table association", func() {
					Expect(subnetDetails.Private).To(HaveLen(2))
					Expect(subnetDetails.Private).To(ContainElement(builder.SubnetResource{
						Subnet:           gfnt.NewString(privateSubnet2),
						RouteTable:       gfnt.NewString("this-is-a-route-table"),
						AvailabilityZone: azB,
					}))
					Expect(subnetDetails.Private).To(ContainElement(builder.SubnetResource{
						Subnet:           gfnt.NewString(privateSubnet1),
						RouteTable:       gfnt.NewString("this-is-a-route-table"),
						AvailabilityZone: azA,
					}))
				})
			})

			Context("importing route tables fails because the rt association points to main", func() {
				BeforeEach(func() {
					rtOutput.RouteTables[0].Associations[0].Main = aws.Bool(true)
					mockEC2.On("DescribeRouteTables", mock.MatchedBy(func(input *ec2.DescribeRouteTablesInput) bool {
						return len(input.Filters) > 0
					})).Return(rtOutput, nil)
				})

				It("returns an error", func() {
					Expect(addErr).To(MatchError(ContainSubstring("subnets must be associated with a non-main route table; eksctl does not modify the main route table")))
				})
			})

			Context("adding the route table to the subnet resource fails", func() {
				BeforeEach(func() {
					rtOutput.RouteTables[0].Associations[0].SubnetId = aws.String("fake")
					mockEC2.On("DescribeRouteTables", mock.MatchedBy(func(input *ec2.DescribeRouteTablesInput) bool {
						return len(input.Filters) > 0
					})).Return(rtOutput, nil)
				})

				It("returns an error", func() {
					Expect(addErr).To(MatchError(ContainSubstring("failed to find an explicit route table associated with subnet \"subnet-0f98135715dfcf55a\"; eksctl does not modify the main route table if a subnet is not associated with an explicit route table")))
				})
			})
		})
	})
})
