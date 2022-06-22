package az

import (
	"context"
	"fmt"
	"math/rand"
	gostrings "strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/utils/nodes"
	"github.com/weaveworks/eksctl/pkg/utils/strings"
)

var zoneIDsToAvoid = map[string][]string{
	api.RegionCNNorth1: {"cnn1-az4"}, // https://github.com/weaveworks/eksctl/issues/3916
}

func GetAvailabilityZones(ctx context.Context, ec2API awsapi.EC2, region string, spec *api.ClusterConfig) ([]string, error) {
	zones, err := getZones(ctx, ec2API, region, spec)
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

func getZones(ctx context.Context, ec2API awsapi.EC2, region string, spec *api.ClusterConfig) ([]string, error) {
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

	filteredZones := filterZones(region, output.AvailabilityZones)
	return filterBasedOnAvailability(ctx, filteredZones, spec, ec2API)
}

func filterBasedOnAvailability(ctx context.Context, zones []string, spec *api.ClusterConfig, ec2API awsapi.EC2) ([]string, error) {
	var (
		instanceList []string
		instances    = make(map[string]struct{}, 0)
	)

	nodePools := nodes.ToNodePools(spec)
	for _, ng := range nodePools {
		for _, i := range ng.InstanceTypeList() {
			if _, ok := instances[i]; !ok {
				instanceList = append(instanceList, i)
			}
			instances[i] = struct{}{}
		}
	}

	// Do an early exit if we don't have anything.
	if len(instances) == 0 {
		// nothing to do
		return zones, nil
	}

	// This list count will not exceed a 100, so it's not necessary to paginate it.
	// I doubt that there is a config out there with 20 different distinct instances defined in them.
	// If you find yourself reading this comment with that exact problem... You deserve a cookie.
	input := &ec2.DescribeInstanceTypeOfferingsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("instance-type"),
				Values: instanceList,
			},
			{
				Name:   aws.String("location"),
				Values: zones,
			},
		},
		LocationType: ec2types.LocationTypeAvailabilityZone,
		MaxResults:   aws.Int32(100),
	}
	output, err := ec2API.DescribeInstanceTypeOfferings(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("unable to list offerings for instance types: %w", err)
	}

	// zoneToInstanceMap['us-west-1b']['t2.small']=struct{}{}
	// zoneToInstanceMap['us-west-1b']['t2.large']=struct{}{}
	zoneToInstanceMap := make(map[string]map[string]struct{})
	for _, offer := range output.InstanceTypeOfferings {
		if _, ok := zoneToInstanceMap[aws.ToString(offer.Location)]; !ok {
			zoneToInstanceMap[aws.ToString(offer.Location)] = make(map[string]struct{})
		}
		zoneToInstanceMap[aws.ToString(offer.Location)][string(offer.InstanceType)] = struct{}{}
	}

	// check if a randomly selected zone supports all selected instances.
	// If we find an instance that is not supported by the selected zone,
	// we do not return that zone.
	var filteredList []string
	for _, zone := range zones {
		var noSupport []string
		for _, instance := range instanceList {
			if _, ok := zoneToInstanceMap[zone][instance]; !ok {
				noSupport = append(noSupport, instance)
			}
		}
		if len(noSupport) == 0 {
			filteredList = append(filteredList, zone)
		} else {
			logger.Info("skipping %s from selection because it doesn't support the following instance type(s): %s", zone, gostrings.Join(noSupport, ","))
		}
	}
	return filteredList, nil
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
