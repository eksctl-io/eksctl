package outputs

import (
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
)

// Stack output names
const (
	// outputs from cluster stack
	ClusterVPC            = "VPC"
	ClusterSecurityGroup  = "SecurityGroup"
	ClusterSubnets        = "Subnets"
	ClusterSubnetsPrivate = string(ClusterSubnets + api.SubnetTopologyPrivate)
	ClusterSubnetsPublic  = string(ClusterSubnets + api.SubnetTopologyPublic)

	ClusterCertificateAuthorityData = "CertificateAuthorityData"
	ClusterEndpoint                 = "Endpoint"
	ClusterARN                      = "ARN"
	ClusterStackName                = "ClusterStackName"
	ClusterSharedNodeSecurityGroup  = "SharedNodeSecurityGroup"

	// outputs from nodegroup stack
	NodeGroupInstanceRoleARN = "InstanceRoleARN"
	// outputs to indicate configuration attributes that may have critical effect
	// on critical effect on forward-compatibility with respect to overal functionality
	// and integrity, e.g. networking
	NodeGroupFeaturePrivateNetworking   = "FeaturePrivateNetworking"
	NodeGroupFeatureSharedSecurityGroup = "FeatureSharedSecurityGroup"
	NodeGroupFeatureLocalSecurityGroup  = "FeatureLocalSecurityGroup"
)

// MustCollect will use each of the keys and attempt to find an output in the given
// stack, if any of the keys are not preset it will return an error
func MustCollect(stack cfn.Stack, keys []string, results map[string]string) error {
	for _, key := range keys {
		var value *string
		for _, x := range stack.Outputs {
			if *x.OutputKey == key {
				value = x.OutputValue
				break
			}
		}
		if value == nil {
			return fmt.Errorf("no ouput %q in stack %q", key, *stack.StackName)
		}
		results[key] = *value
	}
	return nil
}

type (
	// Collector is a callback function that takes an output value
	// and may return an error
	Collector func(string) error
	// Collectors are a map of ouput names to Collector callbacks
	collectors = map[string]Collector
	// CollectorSet is a wrapper to define methods for collectors
	CollectorSet struct{ set collectors }
)

// NewCollectorSet creates a new CollectorSet based on a map of
// ouput names to Collector callbacks
func NewCollectorSet(set map[string]Collector) *CollectorSet {
	if set == nil {
		return &CollectorSet{make(collectors)}
	}
	return &CollectorSet{set}
}

func (c *CollectorSet) doCollect(must bool, stack cfn.Stack) error {
	for key, collect := range c.set {
		var value *string
		for _, x := range stack.Outputs {
			if *x.OutputKey == key {
				value = x.OutputValue
				break
			}
		}
		if value == nil {
			if !must {
				return nil
			}
			err := fmt.Errorf("no ouput %q", key)
			if stack.StackName != nil {
				return fmt.Errorf("%s in stack %q", err.Error(), *stack.StackName)
			}
			return err
		}
		if err := collect(*value); err != nil {
			return err
		}
	}
	return nil
}

// Collect the outputs of a stack using required and optional CollectorSets
func Collect(stack cfn.Stack, required, optional map[string]Collector) error {
	if err := NewCollectorSet(optional).doCollect(false, stack); err != nil {
		return err
	}
	return NewCollectorSet(required).doCollect(true, stack)
}

// MustCollect will error if any of the outputs are missing
func (c *CollectorSet) MustCollect(stack cfn.Stack) error {
	return c.doCollect(true, stack)
}

// Define a new output, if template is given it will be declared in the it
// and optionally export too
func (c *CollectorSet) Define(template *gfn.Template, name string, value interface{}, export bool, fn Collector) {
	if template != nil {
		o := map[string]interface{}{"Value": value}
		if export {
			o["Export"] = map[string]*gfn.Value{
				"Name": gfn.MakeFnSubString(fmt.Sprintf("${%s}::%s", gfn.StackName, name)),
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
func (c *CollectorSet) DefineJoined(template *gfn.Template, name string, values []*gfn.Value, export bool, fn Collector) {
	c.Define(template, name, gfn.MakeFnJoin(",", values), export, fn)
}

// DefineFromAtt - a new output from an attributes
func (c *CollectorSet) DefineFromAtt(template *gfn.Template, name, att string, export bool, fn Collector) {
	c.Define(template, name, gfn.MakeFnGetAttString(att), export, fn)
}
