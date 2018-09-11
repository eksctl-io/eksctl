package ami

import "github.com/pkg/errors"

var (
	// DefaultAMIResolvers contains a list of resolvers to try in order
	DefaultAMIResolvers = []Resolver{&StaticGPUResolver{}, &StaticDefaultResolver{}}
)

// ResolveAMI will resolve an AMI from the supplied region
// and instance type. It will invoke a specific resolver
// to do the actual detrminng of AMI.
func ResolveAMI(region string, instanceType string) (string, error) {
	for _, resolver := range DefaultAMIResolvers {
		ami, err := resolver.Resolve(region, instanceType)
		if err != nil {
			errors.Wrap(err, "error whilst resolving AMI")
		}
		if ami != "" {
			return ami, nil
		}
	}

	return "", NewErrFailedAMIResolution(region, instanceType)
}

// Resolver provides an interface to enable implementing multiple
// ways to determine which AMI to use from the region/instance type.
type Resolver interface {
	Resolve(region string, instanceType string) (string, error)
}
