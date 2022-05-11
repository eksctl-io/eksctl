package az

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
)

var zoneIDsToAvoid = map[string][]string{
	api.RegionCNNorth1: {"cnn1-az4"}, // https://github.com/weaveworks/eksctl/issues/3916
}

func GetAvailabilityZones(ctx context.Context, ec2API awsapi.EC2, region string) ([]string, error) {
	zones, err := getZones(ctx, ec2API, region)
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

func getZones(ctx context.Context, ec2API awsapi.EC2, region string) ([]string, error) {
	input := &ec2.DescribeAvailabilityZonesInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("region-name"),
				Values: []string{region},
			}, {
				Name:   aws.String("state"),
				Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
			}, {
				Name:   aws.String("zone-type"),
				Values: []string{string(ec2types.LocationTypeAvailabilityZone)},
			},
		},
	}

	output, err := ec2API.DescribeAvailabilityZones(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting availability zones for region %s: %w", region, err)
	}

	return filterZones(region, output.AvailabilityZones), nil
}

func filterZones(region string, zones []ec2types.AvailabilityZone) []string {
	var filteredZones []string
	azsToAvoid := zoneIDsToAvoid[region]
	for _, z := range zones {
		if !strings.Contains(azsToAvoid, *z.ZoneId) {
			filteredZones = append(filteredZones, *z.ZoneName)
		}
	}

	return filteredZones
}
