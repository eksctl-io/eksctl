package ami

import (
	"github.com/weaveworks/eksctl/pkg/utils"
)

var (
	// DefaultAMIResolvers contains a list of resolvers to try in order
	DefaultAMIResolvers = []Resolver{&GpuResolver{}, &DefaultResolver{}}
)

// ResolveAMI will resolve an AMI from the supplied region
// and instance type. It will invoke a specific resolver
// to do the actual detrminng of AMI.
func ResolveAMI(region string, instanceType string) (string, error) {
	for _, resolver := range DefaultAMIResolvers {
		ami := resolver.Resolve(region, instanceType)
		if ami != "" {
			return ami, nil
		}
	}

	return "", NewErrFailedAMIResolution(region, instanceType)
}

// Resolver provides an interface to enable implementing multiple
// ways to determine which AMI to use from the region/instance type.
type Resolver interface {
	Resolve(region string, instanceType string) string
}

// DefaultResolver resolves the AMi to the defaults for the region
type DefaultResolver struct {
}

// Resolve will return an AMI to use based on the default AMI for each region
// TODO: https://github.com/weaveworks/eksctl/issues/49
// currently source of truth for these is here:
// https://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html
func (r *DefaultResolver) Resolve(region string, instanceType string) string {
	switch region {
	case "us-west-2":
		return "ami-08cab282f9979fc7a"
	case "us-east-1":
		return "ami-0b2ae3c6bda8b5c06"
	case "eu-west-1":
		return "ami-066110c1a7466949e"
	default:
		return ""
	}
}

// GpuResolver resolves the AMI for GPU instances types.
type GpuResolver struct {
}

// Resolve will return an AMI based on the region for GPU instance types
func (r *GpuResolver) Resolve(region string, instanceType string) string {
	if !utils.IsGPUInstanceType(instanceType) {
		return ""
	}

	switch region {
	case "us-west-2":
		return "ami-0d20f2404b9a1c4d1"
	case "us-east-1":
		return "ami-09fe6fc9106bda972"
	case "eu-west-1":
		return "ami-09e0c6b3d3cf906f1"
	default:
		return ""
	}
}
