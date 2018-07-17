package eks_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	. "github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/eks/mocks"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type MockProvider struct {
	cfn *mocks.CloudFormationAPI
	eks *mocks.EKSAPI
	ec2 *mocks.EC2API
	sts *mocks.STSAPI
}

func (m MockProvider) CloudFormation() cloudformationiface.CloudFormationAPI { return m.cfn }
func (m MockProvider) mockCloudFormation() *mocks.CloudFormationAPI {
	return m.CloudFormation().(*mocks.CloudFormationAPI)
}

func (m MockProvider) EKS() eksiface.EKSAPI   { return m.eks }
func (m MockProvider) mockEKS() *mocks.EKSAPI { return m.EKS().(*mocks.EKSAPI) }
func (m MockProvider) EC2() ec2iface.EC2API   { return m.ec2 }
func (m MockProvider) mockEC2() *mocks.EC2API { return m.EC2().(*mocks.EC2API) }
func (m MockProvider) STS() stsiface.STSAPI   { return m.sts }
func (m MockProvider) mockSTS() *mocks.STSAPI { return m.STS().(*mocks.STSAPI) }

var _ = Describe("Eks", func() {

	Describe("When calling SelectAvailabilityZones", func() {
		var (
			zonesToAvoid []*ec2.AvailabilityZone
			c            *ClusterProvider
			p            *MockProvider
			err          error
		)

		BeforeEach(func() {
			zonesToAvoid = avoidedZones(ec2.AvailabilityZoneStateAvailable)
		})

		Context("with 2 zones to avoid and 3 additional zones", func() {
			var (
				zones []*ec2.AvailabilityZone
			)

			Context("and all zones available", func() {
				var (
					selectedZones []string
				)
				BeforeEach(func() {
					zones = append(zonesToAvoid, usWest2Zones(ec2.AvailabilityZoneStateAvailable)...)
					c, p = createProviders()

					p.mockEC2().On("DescribeAvailabilityZones",
						mock.MatchedBy(func(input *ec2.DescribeAvailabilityZonesInput) bool {
							// This will match an valid DescribeAvailabilityZonesInput
							return true
						}),
					).Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: zones,
					}, nil)
				})

				JustBeforeEach(func() {
					selectedZones, err = c.SelectAvailabilityZones()
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should have called AWS EC2 DescribeAvailabilityZones", func() {
					Expect(p.mockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeAvailabilityZones", 1)).To(BeTrue())
				})

				It("should have returned 3 availability zones", func() {
					Expect(len(selectedZones)).To(Equal(3))
				})
			})

			Context("and only 1 of the 3 additional zones is available", func() {
				var (
					selectedZones    []string
					expectedZoneName *string
				)
				BeforeEach(func() {
					westZones := usWest2Zones(ec2.AvailabilityZoneStateUnavailable)
					westZones[0].State = aws.String(ec2.AvailabilityZoneStateAvailable)
					expectedZoneName = westZones[0].ZoneName
					zones = append(zonesToAvoid, westZones...)

					c, p = createProviders()

					p.mockEC2().On("DescribeAvailabilityZones",
						mock.MatchedBy(func(input *ec2.DescribeAvailabilityZonesInput) bool {
							// This will match an valid DescribeAvailabilityZonesInput
							return true
						}),
					).Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: zones,
					}, nil)
				})

				JustBeforeEach(func() {
					selectedZones, err = c.SelectAvailabilityZones()
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should have called AWS EC2 DescribeAvailabilityZones", func() {
					Expect(p.mockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeAvailabilityZones", 1)).To(BeTrue())
				})

				It("should have returned 3 identical availability zones", func() {
					Expect(len(selectedZones)).To(Equal(3))

					for _, actualZoneName := range selectedZones {
						Expect(actualZoneName).To(Equal(*expectedZoneName))
					}
				})
			})
		})

		Context("with an error from AWS", func() {
			var (
				selectedZones []string
			)
			BeforeEach(func() {
				c, p = createProviders()

				p.mockEC2().On("DescribeAvailabilityZones",
					mock.MatchedBy(func(input *ec2.DescribeAvailabilityZonesInput) bool {
						// This will match an valid DescribeAvailabilityZonesInput
						return true
					}),
				).Return(nil, fmt.Errorf("Some random error from AWS"))
			})

			JustBeforeEach(func() {
				selectedZones, err = c.SelectAvailabilityZones()
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should not have returned selected zones", func() {
				Expect(selectedZones).Should(BeNil())
			})

			It("should have called AWS EC2 DescribeAvailabilityZones", func() {
				Expect(p.mockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeAvailabilityZones", 1)).To(BeTrue())
			})
		})
	})
})

func createProviders() (*ClusterProvider, *MockProvider) {
	p := &MockProvider{
		cfn: &mocks.CloudFormationAPI{},
		eks: &mocks.EKSAPI{},
		ec2: &mocks.EC2API{},
		sts: &mocks.STSAPI{},
	}

	c := &ClusterProvider{
		Provider: p,
		Spec: &ClusterConfig{
			Region: "us-west-1",
		},
	}

	return c, p
}

func createAvailabilityZone(region string, state string, zone string) *ec2.AvailabilityZone {
	return &ec2.AvailabilityZone{
		RegionName: aws.String(region),
		State:      aws.String(state),
		ZoneName:   aws.String(zone),
	}
}

func avoidedZones(initialStatus string) []*ec2.AvailabilityZone {
	return []*ec2.AvailabilityZone{
		createAvailabilityZone("US East (N. Virginia)", initialStatus, "us-east1-a"),
		createAvailabilityZone("US East (N. Virginia)", initialStatus, "us-east1-b"),
	}
}

func usEast1Zones(initialStatus string) []*ec2.AvailabilityZone {
	return append(avoidedZones(initialStatus), createAvailabilityZone("US East (N. Virginia)", initialStatus, "us-east1-c"))
}

func usWest2Zones(initialStatus string) []*ec2.AvailabilityZone {
	return []*ec2.AvailabilityZone{
		createAvailabilityZone("US West (N. California)", initialStatus, "us-west2-a"),
		createAvailabilityZone("US West (N. California)", initialStatus, "us-west2-b"),
		createAvailabilityZone("US West (N. California)", initialStatus, "us-west2-c"),
	}
}
