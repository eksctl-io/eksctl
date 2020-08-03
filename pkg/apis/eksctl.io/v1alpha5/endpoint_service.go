package v1alpha5

import (
	"github.com/pkg/errors"
)

const (
	// Endpoint services that are required and cannot be disabled
	EndpointServiceEC2    = "ec2"
	EndpointServiceECRAPI = "ecr.api"
	EndpointServiceECRDKR = "ecr.dkr"
	EndpointServiceS3     = "s3"
	EndpointServiceSTS    = "sts"
)

// Values for `AdditionalEndpointServices`
// Additional endpoint services
const (
	EndpointServiceCloudFormation = "cloudformation"
	EndpointServiceAutoscaling    = "autoscaling"
	EndpointServiceCloudWatch     = "logs"
)

// RequiredEndpointServices returns a list of endpoint services that are required for a fully-private cluster
func RequiredEndpointServices() []string {
	return []string{
		EndpointServiceEC2,
		EndpointServiceECRAPI,
		EndpointServiceECRDKR,
		EndpointServiceS3,
		EndpointServiceSTS,
	}
}

// ValidateAdditionalEndpointServices validates support for the specified additional endpoint services
func ValidateAdditionalEndpointServices(services []string) error {
	seen := make(map[string]struct{})
	for _, service := range services {
		switch service {
		case EndpointServiceCloudFormation, EndpointServiceAutoscaling, EndpointServiceCloudWatch:
			if _, ok := seen[service]; ok {
				return errors.Errorf("found duplicate endpoint service: %q", service)
			}
			seen[service] = struct{}{}
		default:
			return errors.Errorf("unsupported endpoint service %q", service)
		}
	}
	return nil
}
