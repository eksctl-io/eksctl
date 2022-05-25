package outputs

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	gfn "github.com/weaveworks/goformation/v4/cloudformation"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// Stack output names
const (
	// outputs from cluster stack
	ClusterVPC                  = "VPC"
	ClusterDefaultSecurityGroup = "ClusterSecurityGroupId"
	ClusterSecurityGroup        = "SecurityGroup"
	ClusterSubnetsPrivate       = string("Subnets" + api.SubnetTopologyPrivate)
	ClusterSubnetsPublic        = string("Subnets" + api.SubnetTopologyPublic)
	ClusterSubnetsPrivateLocal  = string("SubnetsLocalZone" + api.SubnetTopologyPrivate)
	ClusterSubnetsPublicLocal   = string("SubnetsLocalZone" + api.SubnetTopologyPublic)
	ClusterFullyPrivate         = "ClusterFullyPrivate"

	ClusterSubnetsPublicLegacy = "Subnets"

	ClusterCertificateAuthorityData = "CertificateAuthorityData"
	ClusterEndpoint                 = "Endpoint"
	ClusterARN                      = "ARN"
	ClusterStackName                = "ClusterStackName"
	ClusterSharedNodeSecurityGroup  = "SharedNodeSecurityGroup"
	ClusterServiceRoleARN           = "ServiceRoleARN"
	ClusterFeatureNATMode           = "FeatureNATMode"

	// outputs from nodegroup stack
	NodeGroupInstanceRoleARN    = "InstanceRoleARN"
	NodeGroupInstanceProfileARN = "InstanceProfileARN"

	// outputs to indicate configuration attributes that may have critical effect
	// on critical effect on forward-compatibility with respect to overall functionality
	// and integrity, e.g. networking
	NodeGroupFeaturePrivateNetworking     = "FeaturePrivateNetworking"
	NodeGroupFeatureSharedSecurityGroup   = "FeatureSharedSecurityGroup"
	NodeGroupFeatureLocalSecurityGroup    = "FeatureLocalSecurityGroup"
	NodeGroupFeaturePrivateEndpointAccess = "FeaturePrivateEndpointAccess"

	// outputs from Fargate stack:
	FargatePodExecutionRoleARN = "FargatePodExecutionRoleARN"

	// IAMServiceAccountRoleName is the name of iamserviceaccount role resource and output.
	IAMServiceAccountRoleName = "Role1"
)

type (
	// Collector is a callback function that takes an output value
	// and may return an error
	Collector func(string) error
	// Collectors are a map of output names to Collector callbacks
	collectors = map[string]Collector
	// CollectorSet is a wrapper to define methods for collectors
	CollectorSet struct{ set collectors }
)

// NewCollectorSet creates a new CollectorSet based on a map of
// output names to Collector callbacks
func NewCollectorSet(set map[string]Collector) *CollectorSet {
	if set == nil {
		return &CollectorSet{make(collectors)}
	}
	return &CollectorSet{set}
}

func get(stack types.Stack, key string) *string {
	for _, x := range stack.Outputs {
		if *x.OutputKey == key {
			return x.OutputValue
		}
	}
	return nil
}

func (c *CollectorSet) doCollect(must bool, stack types.Stack) error {
	for key, collector := range c.set {
		value := get(stack, key)
		if value == nil {
			if must {
				err := fmt.Errorf("no output %q", key)
				if stack.StackName != nil {
					return fmt.Errorf("%s in stack %q", err.Error(), *stack.StackName)
				}
				return err
			}
			continue
		}
		if err := collector(*value); err != nil {
			return err
		}
	}
	return nil
}

// Exists checks if the stack has given output key
func Exists(stack types.Stack, key string) bool {
	return get(stack, key) != nil
}

// Collect the outputs of a stack using required and optional CollectorSets
func Collect(stack types.Stack, required, optional map[string]Collector) error {
	if err := NewCollectorSet(optional).doCollect(false, stack); err != nil {
		return err
	}
	return NewCollectorSet(required).doCollect(true, stack)
}

// MustCollect will error if any of the outputs are missing
func (c *CollectorSet) MustCollect(stack types.Stack) error {
	return c.doCollect(true, stack)
}

// Define a new output, if template is given it will be declared in the it
// and optionally export too
func (c *CollectorSet) Define(template *gfn.Template, name string, value interface{}, export bool, fn Collector) {
	if template != nil {
		o := gfn.Output{
			Value: value,
		}
		if export {
			o.Export = &gfn.Export{
				Name: gfnt.MakeFnSubString(fmt.Sprintf("${%s}::%s", gfnt.StackName, name)),
			}
		}
		template.Outputs[name] = o
	}
	c.set[name] = fn
}

// DefineWithoutCollector allows bare outputs
func (c *CollectorSet) DefineWithoutCollector(template *gfn.Template, name string, value interface{}, export bool) {
	c.Define(template, name, value, export, func(_ string) error { return nil })
}

// DefineJoined - a new output as comma-separated list
func (c *CollectorSet) DefineJoined(template *gfn.Template, name string, values []*gfnt.Value, export bool, fn Collector) {
	c.Define(template, name, gfnt.MakeFnJoin(",", values), export, fn)
}

// DefineFromAtt - a new output from an attributes
func (c *CollectorSet) DefineFromAtt(template *gfn.Template, name, logicalName, att string, export bool, fn Collector) {
	c.Define(template, name, gfnt.MakeFnGetAttString(logicalName, att), export, fn)
}
