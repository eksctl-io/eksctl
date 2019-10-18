package ami

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
)

// MultiResolver is a Resolver that delegates to one or more Resolvers.
// It iterates over the delegate resolvers and returns the first AMI found
type MultiResolver struct {
	delegates []Resolver
}

// Resolve will resolve an AMI from the supplied region
// and instance type. It will invoke a specific resolver
// to do the actual determining of AMI.
func (r *MultiResolver) Resolve(region, version, instanceType, imageFamily string) (string, error) {
	for _, resolver := range r.delegates {
		ami, err := resolver.Resolve(region, version, instanceType, imageFamily)
		if err != nil {
			return "", err
		}
		if ami != "" {
			return ami, nil
		}
	}

	return "", NewErrFailedResolution(region, version, instanceType, imageFamily)
}

// Resolver provides an interface to enable implementing multiple
// ways to determine which AMI to use from the region/instance type/image family.
type Resolver interface {
	Resolve(region, version, instanceType, imageFamily string) (string, error)
}

// NewDefaultResolver returns a static resolver that delegates on StaticGPUResolver and StaticDefaultResolver
func NewDefaultResolver() Resolver {
	return &MultiResolver{
		delegates: []Resolver{&StaticGPUResolver{}, &StaticDefaultResolver{}},
	}
}

// NewAutoResolver creates a new AutoResolver
func NewAutoResolver(api ec2iface.EC2API) Resolver {
	return &AutoResolver{api: api}
}

// NewSSMResolver creates a new AutoResolver
func NewSSMResolver(api ssmiface.SSMAPI) Resolver {
	return &SSMResolver{ssmAPI: api}
}
