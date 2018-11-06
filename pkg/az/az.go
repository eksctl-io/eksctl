package az

import (
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/eks/api"
)

const (
	// RecommendedAvailabilityZones defines the default number of required availability zones
	RecommendedAvailabilityZones = api.RecommendedSubnets
	// MinRequiredAvailabilityZones defines the minimum number of required availability zones
	MinRequiredAvailabilityZones = api.MinRequiredSubnets
)

// SelectionStrategy provides an interface to allow changing the strategy used to
// select availability zones to use from a list available.
type SelectionStrategy interface {
	Select(availableZones []string) []string
}

// RequiredNumberRandomStrategy selects az zones randomly up to a required amount
// of zones.
type RequiredNumberRandomStrategy struct {
	RequiredAvailabilityZones int
}

// Select will randomly select az from the supplied list. The number of az's
// selected will be controlled by RequiredAvailabilityZones.
func (r *RequiredNumberRandomStrategy) Select(availableZones []string) []string {
	zones := []string{}
	for len(zones) < r.RequiredAvailabilityZones {
		rand := rand.New(rand.NewSource(time.Now().UnixNano()))
		for _, rn := range rand.Perm(len(availableZones)) {
			zones = append(zones, availableZones[rn])
			if len(zones) == r.RequiredAvailabilityZones {
				break
			}
		}
	}
	return zones
}

// NewRecommendedNumberRandomStrategy returns a RequiredNumberRandomStrategy that
// has the number of required zones set to the default (RecommendedAvailabilityZones)
func NewRecommendedNumberRandomStrategy() *RequiredNumberRandomStrategy {
	return &RequiredNumberRandomStrategy{RequiredAvailabilityZones: RecommendedAvailabilityZones}
}

// NewMinRequiredNumberRandomStrategy returns a RequiredNumberRandomStrategy that
// has the number of required zones set to the default (MinRequiredAvailabilityZones)
func NewMinRequiredNumberRandomStrategy() *RequiredNumberRandomStrategy {
	return &RequiredNumberRandomStrategy{RequiredAvailabilityZones: MinRequiredAvailabilityZones}
}

// ZoneUsageRule provides an interface to enable rules to determine if a
// zone should be used.
type ZoneUsageRule interface {
	CanUseZone(zone *ec2.AvailabilityZone) bool
}

// ZonesToAvoidRule can be used to ensure that certain az aren't used. This can be used
// to avoid zones that are know to be overpopulated.
type ZonesToAvoidRule struct {
	zonesToAvoid map[string]bool
}

// CanUseZone checks if the supplied zone is in the list of zones to be avoided.
func (za *ZonesToAvoidRule) CanUseZone(zone *ec2.AvailabilityZone) bool {
	_, avoidZone := za.zonesToAvoid[*zone.ZoneName]
	return !avoidZone
}

// NewZonesToAvoidRule returns a new ZonesToAvoidRule with the supplied
// zones set to avoid
func NewZonesToAvoidRule(zonesToAvoid map[string]bool) *ZonesToAvoidRule {
	return &ZonesToAvoidRule{zonesToAvoid: zonesToAvoid}
}

// AvailabilityZoneSelector used to select availability zones to use
type AvailabilityZoneSelector struct {
	ec2api   ec2iface.EC2API
	strategy SelectionStrategy
	rules    []ZoneUsageRule
}

// NewSelectorWithDefaults create a new AvailabilityZoneSelector with the
// default selection strategy and usage rules
func NewSelectorWithDefaults(ec2api ec2iface.EC2API) *AvailabilityZoneSelector {
	avoidZones := map[string]bool{}

	return &AvailabilityZoneSelector{
		ec2api:   ec2api,
		strategy: NewRecommendedNumberRandomStrategy(),
		rules:    []ZoneUsageRule{NewZonesToAvoidRule(avoidZones)},
	}
}

// NewSelectorWithMinRequired create a new AvailabilityZoneSelector with the
// minimum required selection strategy and usage rules
func NewSelectorWithMinRequired(ec2api ec2iface.EC2API) *AvailabilityZoneSelector {
	avoidZones := map[string]bool{}

	return &AvailabilityZoneSelector{
		ec2api:   ec2api,
		strategy: NewMinRequiredNumberRandomStrategy(),
		rules:    []ZoneUsageRule{NewZonesToAvoidRule(avoidZones)},
	}
}

// SelectZones returns a list fo az zones to use for the supplied region
func (a *AvailabilityZoneSelector) SelectZones(regionName string) ([]string, error) {
	availableZones, err := a.getZonesForRegion(regionName)
	if err != nil {
		return nil, err
	}

	usableZones := a.getUsableZones(availableZones)

	return a.strategy.Select(usableZones), nil
}

func (a *AvailabilityZoneSelector) getUsableZones(availableZones []*ec2.AvailabilityZone) []string {
	usableZones := []string{}
	for _, zone := range availableZones {
		zoneUsable := true
		for _, rule := range a.rules {
			if !rule.CanUseZone(zone) {
				zoneUsable = false
				break
			}
		}
		if zoneUsable {
			usableZones = append(usableZones, *zone.ZoneName)
		}
	}

	return usableZones
}

func (a *AvailabilityZoneSelector) getZonesForRegion(regionName string) ([]*ec2.AvailabilityZone, error) {
	regionFilter := &ec2.Filter{
		Name:   aws.String("region-name"),
		Values: []*string{aws.String(regionName)},
	}
	stateFilter := &ec2.Filter{
		Name:   aws.String("state"),
		Values: []*string{aws.String(ec2.AvailabilityZoneStateAvailable)},
	}

	input := &ec2.DescribeAvailabilityZonesInput{
		Filters: []*ec2.Filter{regionFilter, stateFilter},
	}

	output, err := a.ec2api.DescribeAvailabilityZones(input)
	if err != nil {
		return nil, errors.Wrapf(err, "getting availability zones for %s", regionName)
	}

	return output.AvailabilityZones, nil
}
