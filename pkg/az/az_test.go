package az_test

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/az"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("AZ", func() {
	var (
		region string
		p      *mockprovider.MockProvider
		spec   *api.ClusterConfig
	)

	BeforeEach(func() {
		region = "us-west-1"
		p = mockprovider.NewMockProvider()
		spec = api.NewClusterConfig()
	})

	When("1 AZ is available", func() {
		BeforeEach(func() {
			p.MockEC2().On("DescribeAvailabilityZones", mock.Anything, &ec2.DescribeAvailabilityZonesInput{
				Filters: []ec2types.Filter{
					{
						Name:   aws.String("region-name"),
						Values: []string{region},
					},
					{
						Name:   aws.String("state"),
						Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
					},
					{
						Name:   aws.String("zone-type"),
						Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
					},
				},
			}).Return(&ec2.DescribeAvailabilityZonesOutput{
				AvailabilityZones: []ec2types.AvailabilityZone{
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone1"),
				},
			}, nil)
		})

		It("errors", func() {
			_, err := az.GetAvailabilityZones(context.Background(), p.MockEC2(), region, spec)
			Expect(err).To(MatchError("only 1 zones discovered [zone1], at least 2 are required"))
		})
	})

	When("2 AZs are available", func() {
		BeforeEach(func() {
			p.MockEC2().On("DescribeAvailabilityZones", mock.Anything, &ec2.DescribeAvailabilityZonesInput{
				Filters: []ec2types.Filter{
					{
						Name:   aws.String("region-name"),
						Values: []string{region},
					},
					{
						Name:   aws.String("state"),
						Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
					},
					{
						Name:   aws.String("zone-type"),
						Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
					},
				},
			}).Return(&ec2.DescribeAvailabilityZonesOutput{
				AvailabilityZones: []ec2types.AvailabilityZone{
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone1"),
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone2"),
				},
			}, nil)
		})

		It("should return the 2 available AZs", func() {
			zones, err := az.GetAvailabilityZones(context.Background(), p.MockEC2(), region, spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(zones).To(HaveLen(2))
			Expect(zones).To(ConsistOf("zone1", "zone2"))
		})
	})

	When("3 AZs are available", func() {
		BeforeEach(func() {
			p.MockEC2().On("DescribeAvailabilityZones", mock.Anything, &ec2.DescribeAvailabilityZonesInput{
				Filters: []ec2types.Filter{
					{
						Name:   aws.String("region-name"),
						Values: []string{region},
					},
					{
						Name:   aws.String("state"),
						Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
					},
					{
						Name:   aws.String("zone-type"),
						Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
					},
				},
			}).Return(&ec2.DescribeAvailabilityZonesOutput{
				AvailabilityZones: []ec2types.AvailabilityZone{
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone1"),
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone2"),
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone3"),
				},
			}, nil)
		})

		It("should return the 3 available AZs", func() {
			zones, err := az.GetAvailabilityZones(context.Background(), p.MockEC2(), region, spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(zones).To(HaveLen(3))
			Expect(zones).To(ConsistOf("zone1", "zone2", "zone3"))
		})
	})

	When("more than 3 AZs are available", func() {
		BeforeEach(func() {
			p.MockEC2().On("DescribeAvailabilityZones", mock.Anything, &ec2.DescribeAvailabilityZonesInput{
				Filters: []ec2types.Filter{
					{
						Name:   aws.String("region-name"),
						Values: []string{region},
					},
					{
						Name:   aws.String("state"),
						Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
					},
					{
						Name:   aws.String("zone-type"),
						Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
					},
				},
			}).Return(&ec2.DescribeAvailabilityZonesOutput{
				AvailabilityZones: []ec2types.AvailabilityZone{
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone1"),
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone2"),
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone3"),
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone4"),
				},
			}, nil)
		})

		It("should return a random set of 3 available AZs", func() {
			zones, err := az.GetAvailabilityZones(context.Background(), p.MockEC2(), region, spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(zones).To(HaveLen(3))
			Expect(zonesAreUnique(zones)).To(BeTrue())
		})
	})

	When("instance types are defined", func() {
		BeforeEach(func() {
			spec = api.NewClusterConfig()
			spec.NodeGroups = []*api.NodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:         "test-az-1",
						InstanceType: "t2.small",
					},
				},
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name:         "test-az-2",
						InstanceType: "t2.medium",
					},
				},
			}
			p.MockEC2().On("DescribeAvailabilityZones", mock.Anything, &ec2.DescribeAvailabilityZonesInput{
				Filters: []ec2types.Filter{
					{
						Name:   aws.String("region-name"),
						Values: []string{region},
					},
					{
						Name:   aws.String("state"),
						Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
					},
					{
						Name:   aws.String("zone-type"),
						Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
					},
				},
			}).Return(&ec2.DescribeAvailabilityZonesOutput{
				AvailabilityZones: []ec2types.AvailabilityZone{
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone1"),
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone2"),
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone3"),
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone4"),
				},
			}, nil)
			// split DescribeInstanceTypeOfferings response in two pages so we unit test the use of the paginator at the same time
			p.MockEC2().On("DescribeInstanceTypeOfferings", mock.Anything, &ec2.DescribeInstanceTypeOfferingsInput{
				Filters: []ec2types.Filter{
					{
						Name:   aws.String("instance-type"),
						Values: []string{"t2.small", "t2.medium"},
					},
					{
						Name:   aws.String("location"),
						Values: []string{"zone1", "zone2", "zone3", "zone4"},
					},
				},
				LocationType: ec2types.LocationTypeAvailabilityZone,
				MaxResults:   aws.Int32(100),
			}, mock.Anything).Return(&ec2.DescribeInstanceTypeOfferingsOutput{
				NextToken: aws.String("token"),
				InstanceTypeOfferings: []ec2types.InstanceTypeOffering{
					{
						InstanceType: "t2.small",
						Location:     aws.String("zone1"),
						LocationType: "availability-zone",
					},
					{
						InstanceType: "t2.small",
						Location:     aws.String("zone2"),
						LocationType: "availability-zone",
					},
					{
						InstanceType: "t2.small",
						Location:     aws.String("zone4"),
						LocationType: "availability-zone",
					},
				},
			}, nil)
			p.MockEC2().On("DescribeInstanceTypeOfferings", mock.Anything, &ec2.DescribeInstanceTypeOfferingsInput{
				Filters: []ec2types.Filter{
					{
						Name:   aws.String("instance-type"),
						Values: []string{"t2.small", "t2.medium"},
					},
					{
						Name:   aws.String("location"),
						Values: []string{"zone1", "zone2", "zone3", "zone4"},
					},
				},
				LocationType: ec2types.LocationTypeAvailabilityZone,
				MaxResults:   aws.Int32(100),
				NextToken:    aws.String("token"),
			}, mock.Anything).Return(&ec2.DescribeInstanceTypeOfferingsOutput{
				InstanceTypeOfferings: []ec2types.InstanceTypeOffering{
					{
						InstanceType: "t2.medium",
						Location:     aws.String("zone1"),
						LocationType: "availability-zone",
					},
					{
						InstanceType: "t2.medium",
						Location:     aws.String("zone2"),
						LocationType: "availability-zone",
					},
					{
						InstanceType: "t2.medium",
						Location:     aws.String("zone4"),
						LocationType: "availability-zone",
					},
				},
			}, nil)
		})

		It("should return only the zones which support all selected instance types", func() {
			zones, err := az.GetAvailabilityZones(context.Background(), p.MockEC2(), region, spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(zones).To(ConsistOf([]string{"zone1", "zone2", "zone4"}))
			Expect(zonesAreUnique(zones)).To(BeTrue())
		})
	})

	When("fetching the AZs errors", func() {
		BeforeEach(func() {
			p.MockEC2().On("DescribeAvailabilityZones", mock.Anything, &ec2.DescribeAvailabilityZonesInput{
				Filters: []ec2types.Filter{
					{
						Name:   aws.String("region-name"),
						Values: []string{region},
					},
					{
						Name:   aws.String("state"),
						Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
					},
					{
						Name:   aws.String("zone-type"),
						Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
					},
				},
			}).Return(&ec2.DescribeAvailabilityZonesOutput{}, fmt.Errorf("foo"))
		})

		It("errors", func() {
			_, err := az.GetAvailabilityZones(context.Background(), p.MockEC2(), region, spec)
			Expect(err).To(MatchError(fmt.Sprintf("error getting availability zones for region %s: foo", region)))
		})
	})

	type unsupportedZoneEntry struct {
		region        string
		zoneNameToIDs map[string]string
		expectedZones []string
	}
	DescribeTable("region with unsupported zone IDs", func(e unsupportedZoneEntry) {
		var azs []ec2types.AvailabilityZone
		for zoneName, zoneID := range e.zoneNameToIDs {
			azs = append(azs, createAvailabilityZoneWithID(e.region, ec2types.AvailabilityZoneStateAvailable, zoneName, zoneID))
		}
		mockProvider := mockprovider.NewMockProvider()
		mockProvider.MockEC2().On("DescribeAvailabilityZones", mock.Anything, &ec2.DescribeAvailabilityZonesInput{
			Filters: []ec2types.Filter{
				{
					Name:   aws.String("region-name"),
					Values: []string{e.region},
				},
				{
					Name:   aws.String("state"),
					Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
				},
				{
					Name:   aws.String("zone-type"),
					Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
				},
			},
		}).Return(&ec2.DescribeAvailabilityZonesOutput{
			AvailabilityZones: azs,
		}, nil)
		mockProvider.MockEC2().On("DescribeInstanceTypeOfferings", mock.Anything, &ec2.DescribeInstanceTypeOfferingsInput{
			Filters: []ec2types.Filter{
				{
					Name:   aws.String("instance-type"),
					Values: []string{"t2.small", "t2.medium"},
				},
				{
					Name:   aws.String("location"),
					Values: []string{"zone1", "zone2", "zone3", "zone4"},
				},
			},
			LocationType: ec2types.LocationTypeAvailabilityZone,
			MaxResults:   aws.Int32(100),
		}, mock.Anything).Return(&ec2.DescribeInstanceTypeOfferingsOutput{
			NextToken: aws.String("token"),
			InstanceTypeOfferings: []ec2types.InstanceTypeOffering{
				{
					InstanceType: "t2.small",
					Location:     aws.String("zone1"),
					LocationType: "availability-zone",
				},
				{
					InstanceType: "t2.small",
					Location:     aws.String("zone2"),
					LocationType: "availability-zone",
				},
				{
					InstanceType: "t2.small",
					Location:     aws.String("zone4"),
					LocationType: "availability-zone",
				},
				{
					InstanceType: "t2.small",
					Location:     aws.String("zone3"),
					LocationType: "availability-zone",
				},
			},
		}, nil)
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Region = e.region
		clusterConfig.NodeGroups = []*api.NodeGroup{
			{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "test-az-1",
				},
			},
			{
				NodeGroupBase: &api.NodeGroupBase{
					Name: "test-az-2",
				},
			},
		}
		zones, err := az.GetAvailabilityZones(context.Background(), mockProvider.MockEC2(), e.region, clusterConfig)
		Expect(err).NotTo(HaveOccurred())
		Expect(zones).To(ConsistOf(e.expectedZones))
	},
		Entry(api.RegionCNNorth1, unsupportedZoneEntry{
			region: api.RegionCNNorth1,
			zoneNameToIDs: map[string]string{
				"zone1": "cnn1-az1",
				"zone2": "cnn1-az2",
				"zone4": "cnn1-az4",
			},
			expectedZones: []string{"zone1", "zone2"},
		}),
		Entry(api.RegionUSEast1, unsupportedZoneEntry{
			region: api.RegionUSEast1,
			zoneNameToIDs: map[string]string{
				"zone1": "use1-az1",
				"zone2": "use1-az3",
				"zone3": "use1-az2",
			},
			expectedZones: []string{"zone1", "zone3"},
		}),
		Entry(api.RegionUSWest1, unsupportedZoneEntry{
			region: api.RegionUSWest1,
			zoneNameToIDs: map[string]string{
				"zone1": "usw1-az2",
				"zone2": "usw1-az1",
				"zone3": "usw1-az3",
			},
			expectedZones: []string{"zone2", "zone3"},
		}),
		Entry(api.RegionCACentral1, unsupportedZoneEntry{
			region: api.RegionCACentral1,
			zoneNameToIDs: map[string]string{
				"zone1": "cac1-az1",
				"zone2": "cac1-az2",
				"zone3": "cac1-az3",
			},
			expectedZones: []string{"zone1", "zone2"},
		}),
	)
	When("the region contains zones that are denylisted", func() {
		BeforeEach(func() {
			region = api.RegionCNNorth1

			p.MockEC2().On("DescribeAvailabilityZones", mock.Anything, &ec2.DescribeAvailabilityZonesInput{
				Filters: []ec2types.Filter{
					{
						Name:   aws.String("region-name"),
						Values: []string{region},
					},
					{
						Name:   aws.String("state"),
						Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
					},
					{
						Name:   aws.String("zone-type"),
						Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
					},
				},
			}).Return(&ec2.DescribeAvailabilityZonesOutput{
				AvailabilityZones: []ec2types.AvailabilityZone{
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone1"),
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone2"),
					createAvailabilityZoneWithID(region, ec2types.AvailabilityZoneStateAvailable, "zone3", "cnn1-az4"),
				},
			}, nil)
		})

		It("should not use the denylisted zones", func() {
			zones, err := az.GetAvailabilityZones(context.Background(), p.MockEC2(), region, spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(zones).To(HaveLen(2))
			Expect(zones).To(ConsistOf("zone1", "zone2"))
		})
	})

	When("using us-east-1", func() {
		BeforeEach(func() {
			region = "us-east-1"

			p.MockEC2().On("DescribeAvailabilityZones", mock.Anything, &ec2.DescribeAvailabilityZonesInput{
				Filters: []ec2types.Filter{
					{
						Name:   aws.String("region-name"),
						Values: []string{region},
					},
					{
						Name:   aws.String("state"),
						Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
					},
					{
						Name:   aws.String("zone-type"),
						Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
					},
				},
			}).Return(&ec2.DescribeAvailabilityZonesOutput{
				AvailabilityZones: []ec2types.AvailabilityZone{
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone1"),
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone2"),
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone3"),
					createAvailabilityZone(region, ec2types.AvailabilityZoneStateAvailable, "zone4"),
				},
			}, nil)
		})

		It("should only use 2 AZs, rather than the default 3", func() {
			zones, err := az.GetAvailabilityZones(context.Background(), p.MockEC2(), region, spec)
			Expect(err).NotTo(HaveOccurred())
			Expect(zones).To(HaveLen(2))
			Expect(zonesAreUnique(zones)).To(BeTrue())
		})
	})
})

func zonesAreUnique(zones []string) bool {
	mapZones := make(map[string]interface{})
	for _, z := range zones {
		mapZones[z] = nil
	}
	return len(mapZones) == len(zones)
}

func createAvailabilityZone(region string, state ec2types.AvailabilityZoneState, zone string) ec2types.AvailabilityZone {
	return createAvailabilityZoneWithID(region, state, zone, "id-"+zone)
}

func createAvailabilityZoneWithID(region string, state ec2types.AvailabilityZoneState, zone, zoneID string) ec2types.AvailabilityZone {
	return ec2types.AvailabilityZone{
		RegionName: aws.String(region),
		State:      state,
		ZoneName:   aws.String(zone),
		ZoneId:     aws.String(zoneID),
	}
}
