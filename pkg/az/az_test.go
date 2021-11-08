package az_test

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	. "github.com/weaveworks/eksctl/pkg/az"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("AZ", func() {

	Describe("When calling SelectZones", func() {
		var (
			p   *mockprovider.MockProvider
			err error
		)

		BeforeEach(func() {
			_ = avoidedZones(ec2.AvailabilityZoneStateAvailable)
		})

		Context("with a region that has no zones to avoid", func() {
			var (
				zones  []*ec2.AvailabilityZone
				region string
			)

			Context("and all zones available", func() {
				var (
					selectedZones []string
					azSelector    *AvailabilityZoneSelector
				)
				BeforeEach(func() {
					region = "us-west-2"

					zones = usWest2Zones(ec2.AvailabilityZoneStateAvailable)
					p = mockprovider.NewMockProvider()

					p.MockEC2().On("DescribeAvailabilityZones",
						mock.MatchedBy(func(input *ec2.DescribeAvailabilityZonesInput) bool {
							filter := input.Filters[0]
							return *filter.Name == "region-name" && *filter.Values[0] == region
						}),
					).Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: zones,
					}, nil)

					azSelector = NewSelectorWithDefaults(p.MockEC2(), region)
				})

				JustBeforeEach(func() {
					selectedZones, err = azSelector.SelectZones()
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should have called AWS EC2 DescribeAvailabilityZones", func() {
					Expect(p.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeAvailabilityZones", 1)).To(BeTrue())
				})

				It("should have returned 3 availability zones", func() {
					Expect(selectedZones).To(HaveLen(3))
				})
			})

			Context("and only 1 zone is available", func() {
				var (
					selectedZones    []string
					expectedZoneName *string
					azSelector       *AvailabilityZoneSelector
				)
				BeforeEach(func() {
					westZone := usWest2Zones(ec2.AvailabilityZoneStateAvailable)[0]
					expectedZoneName = westZone.ZoneName
					zones = []*ec2.AvailabilityZone{westZone}

					p = mockprovider.NewMockProvider()

					p.MockEC2().On("DescribeAvailabilityZones",
						mock.MatchedBy(func(input *ec2.DescribeAvailabilityZonesInput) bool {
							filter := input.Filters[0]
							return *filter.Name == "region-name" && *filter.Values[0] == region
						}),
					).Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: zones,
					}, nil)

					azSelector = NewSelectorWithDefaults(p.MockEC2(), region)
				})

				JustBeforeEach(func() {
					selectedZones, err = azSelector.SelectZones()
				})

				It("should not error", func() {
					Expect(err).NotTo(HaveOccurred())
				})

				It("should have called AWS EC2 DescribeAvailabilityZones", func() {
					Expect(p.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeAvailabilityZones", 1)).To(BeTrue())
				})

				It("should have returned 3 identical availability zones", func() {
					Expect(selectedZones).To(HaveLen(3))

					for _, actualZoneName := range selectedZones {
						Expect(actualZoneName).To(Equal(*expectedZoneName))
					}
				})
			})
		})

		Context("with a region that has zones to avoid", func() {
			var (
				zones            []*ec2.AvailabilityZone
				selectedZones    []string
				region           string
				azSelector       *AvailabilityZoneSelector
				expectedZoneName *string
			)
			BeforeEach(func() {
				region = "us-east-1"
				expectedZoneName = aws.String("us-east-1c")

				zones = usEast1Zones(ec2.AvailabilityZoneStateAvailable)
				p = mockprovider.NewMockProvider()

				p.MockEC2().On("DescribeAvailabilityZones",
					mock.MatchedBy(func(input *ec2.DescribeAvailabilityZonesInput) bool {
						filter := input.Filters[0]
						return *filter.Name == "region-name" && *filter.Values[0] == region
					}),
				).Return(&ec2.DescribeAvailabilityZonesOutput{
					AvailabilityZones: zones,
				}, nil)

				azSelector = NewSelectorWithDefaults(p.MockEC2(), region)
			})

			JustBeforeEach(func() {
				selectedZones, err = azSelector.SelectZones()
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should have called AWS EC2 DescribeAvailabilityZones", func() {
				Expect(p.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeAvailabilityZones", 1)).To(BeTrue())
			})

			It("should have returned 3 availability zones", func() {
				Expect(selectedZones).To(HaveLen(3))
			})

			It("should have returned none of the zones to avoid", func() {
				for _, actualZoneName := range selectedZones {
					Expect(actualZoneName).To(Equal(*expectedZoneName))
				}
			})
		})

		Context("with an error from AWS", func() {
			var (
				selectedZones []string
				azSelector    *AvailabilityZoneSelector
			)
			BeforeEach(func() {
				p = mockprovider.NewMockProvider()

				p.MockEC2().On("DescribeAvailabilityZones",
					mock.Anything,
				).Return(nil, errors.New("some random error from AWS"))

				azSelector = NewSelectorWithDefaults(p.MockEC2(), "us-west-2")
			})

			JustBeforeEach(func() {
				selectedZones, err = azSelector.SelectZones()
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("should not have returned selected zones", func() {
				Expect(selectedZones).Should(BeNil())
			})

			It("should have called AWS EC2 DescribeAvailabilityZones", func() {
				Expect(p.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeAvailabilityZones", 1)).To(BeTrue())
			})
		})

		Context("with min required zones selector", func() {
			var (
				region        string
				selectedZones []string
				azSelector    *AvailabilityZoneSelector
				zones         []*ec2.AvailabilityZone
			)

			BeforeEach(func() {
				region = "us-east-1"
				zones = usEast1Zones(ec2.AvailabilityZoneStateAvailable)
				p = mockprovider.NewMockProvider()

				p.MockEC2().On("DescribeAvailabilityZones",
					mock.MatchedBy(func(input *ec2.DescribeAvailabilityZonesInput) bool {
						filter := input.Filters[0]
						return *filter.Name == "region-name" && *filter.Values[0] == region
					}),
				).Return(&ec2.DescribeAvailabilityZonesOutput{
					AvailabilityZones: zones,
				}, nil)

				azSelector = NewSelectorWithMinRequired(p.MockEC2(), region)
			})

			JustBeforeEach(func() {
				selectedZones, err = azSelector.SelectZones()
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should have called AWS EC2 DescribeAvailabilityZones", func() {
				Expect(p.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeAvailabilityZones", 1)).To(BeTrue())
			})

			It("should have returned 2 availability zones", func() {
				Expect(selectedZones).To(HaveLen(2))
			})
		})

		Context("with Beijing region that has an unsupported zone", func() {
			const region = "cn-north-1"
			var p *mockprovider.MockProvider

			BeforeEach(func() {
				p = mockprovider.NewMockProvider()
			})

			It("should avoid unsupported zones", func() {
				p.MockEC2().On("DescribeAvailabilityZones", mock.MatchedBy(func(input *ec2.DescribeAvailabilityZonesInput) bool {
					filter := input.Filters[0]
					return *filter.Name == "region-name" && *filter.Values[0] == region
				})).Return(&ec2.DescribeAvailabilityZonesOutput{
					AvailabilityZones: []*ec2.AvailabilityZone{
						{
							ZoneName: aws.String("cn-north-1a"),
							ZoneId:   aws.String("cnn1-az2"),
						},
						{
							ZoneName: aws.String("cn-north-1d"),
							ZoneId:   aws.String("cnn1-az4"),
						},
						{
							ZoneName: aws.String("cn-north-1b"),
							ZoneId:   aws.String("cnn1-az3"),
						},
						{
							ZoneName: aws.String("cn-north-1e"),
							ZoneId:   aws.String("cnn1-az1"),
						},
					},
				}, nil)

				azSelector := NewSelectorWithDefaults(p.EC2(), region)
				selectedZones, err := azSelector.SelectZones()
				Expect(err).NotTo(HaveOccurred())
				Expect(p.MockEC2().AssertNumberOfCalls(GinkgoT(), "DescribeAvailabilityZones", 1)).To(BeTrue())
				Expect(selectedZones).To(ConsistOf("cn-north-1a", "cn-north-1b", "cn-north-1e"))
			})
		})
	})
})

func createAvailabilityZone(region string, state string, zone string) *ec2.AvailabilityZone {
	return &ec2.AvailabilityZone{
		RegionName: aws.String(region),
		State:      aws.String(state),
		ZoneName:   aws.String(zone),
	}
}

func avoidedZones(initialStatus string) []*ec2.AvailabilityZone {
	return []*ec2.AvailabilityZone{
		// createAvailabilityZone("US East (N. Virginia)", initialStatus, "us-east-1a"),
		// createAvailabilityZone("US East (N. Virginia)", initialStatus, "us-east-1b"),
	}
}

func usEast1Zones(initialStatus string) []*ec2.AvailabilityZone {
	return append(avoidedZones(initialStatus), createAvailabilityZone("US East (N. Virginia)", initialStatus, "us-east-1c"))
}

func usWest2Zones(initialStatus string) []*ec2.AvailabilityZone {
	return []*ec2.AvailabilityZone{
		createAvailabilityZone("US West (N. California)", initialStatus, "us-west-2a"),
		createAvailabilityZone("US West (N. California)", initialStatus, "us-west-2b"),
		createAvailabilityZone("US West (N. California)", initialStatus, "us-west-2c"),
	}
}
