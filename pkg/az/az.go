package az

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var zoneIDsToAvoid = map[string][]string{
	api.RegionCNNorth1: {"cnn1-az4"}, // https://github.com/weaveworks/eksctl/issues/3916
}

func GetAvailabilityZones(ec2API ec2iface.EC2API, region string) ([]string, error) {
	zones, err := getZones(ec2API, region)
	if err != nil {
		return nil, err
	}

	numberOfZones := len(zones)
	if numberOfZones < api.MinRequiredAvailabilityZones {
		return nil, fmt.Errorf("only %d zones discovered %v, at least %d are required", numberOfZones, zones, api.MinRequiredAvailabilityZones)
	}

	if numberOfZones < api.RecommendedAvailabilityZones {
		return zones, nil
	}

	return randomSelectionOfZones(region, zones), nil
}

func randomSelectionOfZones(region string, availableZones []string) []string {
	var zones []string
	desiredNumberOfAZs := api.RecommendedAvailabilityZones
	if region == api.RegionUSEast1 {
		desiredNumberOfAZs = api.MinRequiredAvailabilityZones
	}

	for len(zones) < desiredNumberOfAZs {
		rand := rand.New(rand.NewSource(time.Now().UnixNano()))
		for _, rn := range rand.Perm(len(availableZones)) {
			zones = append(zones, availableZones[rn])
			if len(zones) == desiredNumberOfAZs {
				break
			}
		}
	}

	return zones
}

func getZones(ec2API ec2iface.EC2API, region string) ([]string, error) {
	regionFilter := &ec2.Filter{
		Name:   aws.String("region-name"),
		Values: []*string{aws.String(region)},
	}
	stateFilter := &ec2.Filter{
		Name:   aws.String("state"),
		Values: []*string{aws.String(ec2.AvailabilityZoneStateAvailable)},
	}

	input := &ec2.DescribeAvailabilityZonesInput{
		Filters: []*ec2.Filter{regionFilter, stateFilter},
	}

	output, err := ec2API.DescribeAvailabilityZones(input)
	if err != nil {
		return nil, fmt.Errorf("error getting availability zones for region %s: %w", region, err)
	}

	return filterZones(region, output.AvailabilityZones), nil
}

func filterZones(region string, zones []*ec2.AvailabilityZone) []string {
	var filteredZones []string
	azsToAvoid := zoneIDsToAvoid[region]
	for _, z := range zones {
		if !contains(azsToAvoid, *z.ZoneId) {
			filteredZones = append(filteredZones, *z.ZoneName)
		}
	}

	return filteredZones
}

func contains(list []string, value string) bool {
	for _, l := range list {
		if l == value {
			return true
		}
	}
	return false
}
