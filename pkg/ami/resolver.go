package ami

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/kris-nova/logger"
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
			if _, ok := err.(*UnsupportedQueryError); ok {
				logger.Debug(err.Error())
				continue
			}
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

// NewStaticResolver returns a static resolver that delegates on
// StaticGPUResolver, StaticBottlerocketResolver, and StaticDefaultResolver.
func NewStaticResolver() Resolver {
	return &MultiResolver{
		delegates: []Resolver{&StaticGPUResolver{}, &StaticBottlerocketResolver{}, &StaticDefaultResolver{}},
	}
}

// NewMultiResolver creates and returns a MultiResolver with the specified delegates
func NewMultiResolver(delegates ...Resolver) *MultiResolver {
	return &MultiResolver{
		delegates: delegates,
	}
}

// NewAutoResolver creates a new AutoResolver
func NewAutoResolver(api ec2iface.EC2API) Resolver {
	return &AutoResolver{api: api}
}

// NewSSMResolver creates a new AutoResolver.
func NewSSMResolver(api ssmiface.SSMAPI) Resolver {
	return &SSMResolver{ssmAPI: api}
}

// UnsupportedQueryError represents an unsupported AMI query error
type UnsupportedQueryError struct {
	msg string
}

// Error returns the error string
func (ue *UnsupportedQueryError) Error() string {
	return ue.msg
}
