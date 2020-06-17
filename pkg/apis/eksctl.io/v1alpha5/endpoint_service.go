package v1alpha5

import "github.com/pkg/errors"

const (
	// Endpoint services that are required and cannot be disabled
	EndpointServiceEC2    = "ec2"
	EndpointServiceECRAPI = "ecr.api"
	EndpointServiceECRDKR = "ecr.dkr"
	EndpointServiceS3     = "s3"
	EndpointServiceSTS    = "sts"

	// Additional endpoint services
	EndpointServiceAutoscaling = "autoscaling"
	EndpointServiceCloudWatch  = "logs"
)

// DefaultEndpointServices returns a list of endpoint services that are enabled by default
func DefaultEndpointServices() []string {
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
	for _, service := range services {
		switch service {
		case EndpointServiceAutoscaling, EndpointServiceCloudWatch:
		default:
			return errors.Errorf("unsupported endpoint service %q", service)
		}
	}
	return nil
}
