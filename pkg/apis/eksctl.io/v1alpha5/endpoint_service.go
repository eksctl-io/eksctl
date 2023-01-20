package v1alpha5

import (
	"fmt"
)

// EndpointService represents a VPC endpoint service.
type EndpointService struct {
	// Name is the name of the endpoint service.
	Name string
	// Optional specifies whether the service is optional.
	Optional bool
	// OutpostsOnly specifies whether the endpoint is required only for Outposts clusters.
	OutpostsOnly bool
	// RequiresChinaPrefix is true if the endpoint service requires a prefix for China regions.
	RequiresChinaPrefix bool
}

var (
	// EndpointServiceS3 is an EndpointService for S3.
	EndpointServiceS3 = EndpointService{
		Name: "s3",
	}
	// EndpointServiceCloudWatch is an EndpointService for CloudWatch Logs.
	EndpointServiceCloudWatch = EndpointService{
		Name:     "logs",
		Optional: true,
	}
)

// EndpointServices is a list of supported endpoint services.
var EndpointServices = []EndpointService{
	{
		Name:                "ec2",
		RequiresChinaPrefix: true,
	},
	{
		Name:                "ecr.api",
		RequiresChinaPrefix: true,
	},
	{
		Name:                "ecr.dkr",
		RequiresChinaPrefix: true,
	},
	EndpointServiceS3,
	{
		Name:                "sts",
		RequiresChinaPrefix: true,
	},
	{
		Name:         "ssm",
		OutpostsOnly: true,
	},
	{
		Name:         "ssmmessages",
		OutpostsOnly: true,
	},
	{
		Name:         "ec2messages",
		OutpostsOnly: true,
	},
	{
		Name:         "secretsmanager",
		OutpostsOnly: true,
	},
	{
		Name:                "cloudformation",
		Optional:            true,
		RequiresChinaPrefix: true,
	},
	{
		Name:     "autoscaling",
		Optional: true,
	},
	EndpointServiceCloudWatch,
}

// RequiredEndpointServices returns a list of endpoint services that are required for a fully-private cluster.
func RequiredEndpointServices(controlPlaneOnOutposts bool) []EndpointService {
	var requiredServices []EndpointService
	for _, es := range EndpointServices {
		if !es.Optional && (controlPlaneOnOutposts || !es.OutpostsOnly) {
			requiredServices = append(requiredServices, es)
		}
	}
	return requiredServices
}

// MapOptionalEndpointServices maps a list of endpoint service names to []EndpointService.
func MapOptionalEndpointServices(endpointServiceNames []string, cloudWatchLoggingEnabled bool) ([]EndpointService, error) {
	optionalServices := getOptionalEndpointServices()
	var mapped []EndpointService
	hasCloudWatchLogs := false
	for _, es := range endpointServiceNames {
		endpointService, ok := optionalServices[es]
		if !ok {
			return nil, fmt.Errorf("invalid optional endpoint service: %q", es)
		}
		mapped = append(mapped, endpointService)
		if es == EndpointServiceCloudWatch.Name {
			hasCloudWatchLogs = true
		}
	}
	if cloudWatchLoggingEnabled && !hasCloudWatchLogs {
		mapped = append(mapped, optionalServices[EndpointServiceCloudWatch.Name])
	}
	return mapped, nil
}

func getOptionalEndpointServices() map[string]EndpointService {
	ret := map[string]EndpointService{}
	for _, es := range EndpointServices {
		if es.Optional {
			ret[es.Name] = es
		}
	}
	return ret
}

// ValidateAdditionalEndpointServices validates support for the specified additional endpoint services.
func ValidateAdditionalEndpointServices(serviceNames []string) error {
	seen := make(map[string]struct{})
	optionalServices := getOptionalEndpointServices()
	for _, service := range serviceNames {
		if _, ok := optionalServices[service]; !ok {
			return fmt.Errorf("unsupported endpoint service %q", service)
		}
		if _, ok := seen[service]; ok {
			return fmt.Errorf("found duplicate endpoint service: %q", service)
		}
		seen[service] = struct{}{}
	}
	return nil
}
