package ami

var (
	// DefaultResolvers contains a list of resolvers to try in order
	DefaultResolvers = []Resolver{&StaticGPUResolver{}, &StaticDefaultResolver{}}
)

// Resolve will resolve an AMI from the supplied region
// and instance type. It will invoke a specific resolver
// to do the actual determining of AMI.
func Resolve(region string, instanceType string, imageFamily string) (string, error) {
	for _, resolver := range DefaultResolvers {
		ami, err := resolver.Resolve(region, instanceType, imageFamily)
		if err != nil {
			return "", err
		}
		if ami != "" {
			return ami, nil
		}
	}

	return "", NewErrFailedResolution(region, instanceType, imageFamily)
}

// Resolver provides an interface to enable implementing multiple
// ways to determine which AMI to use from the region/instance type/image family.
type Resolver interface {
	Resolve(region string, instanceType string, imageFamily string) (string, error)
}
