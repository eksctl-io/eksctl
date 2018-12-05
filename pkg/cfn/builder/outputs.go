package builder

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
)

// Outputs of CloudFormation stacks are collected into a struct with fields
// matching names of the outputs. Here is a set of reflect-based helpers that
// make this happen. Some data types get special treatment, e.g. string slices
// and byte slices.

const (
	// outputs that are destined for ClusterStackOutputs
	cfnOutputClusterVPC           = "VPC"
	cfnOutputClusterSecurityGroup = "SecurityGroup"
	cfnOutputClusterSubnets       = "Subnets"

	cfnOutputClusterCertificateAuthorityData = "CertificateAuthorityData"
	cfnOutputClusterEndpoint                 = "Endpoint"
	cfnOutputClusterARN                      = "ARN"
	cfnOutputClusterStackName                = "ClusterStackName"

	// this is set inside of NodeGroup
	cfnOutputInstanceRoleARN = "InstanceRoleARN"
)

// ClusterStackOutputs is a struct that hold all of cluster stack outputs,
// it's needed because some of the destination fields in ClusterConfig are
// deeply nested and we would need to do something complicated to handle
// those otherwise
type ClusterStackOutputs struct {
	VPC            string
	SecurityGroup  string
	SubnetsPrivate []string
	SubnetsPublic  []string

	ClusterStackName         string
	Endpoint                 string
	CertificateAuthorityData []byte
	ARN                      string
}

// newOutput defines a new output and optionally exports it
func (r *resourceSet) newOutput(name string, value interface{}, export bool) {
	o := map[string]interface{}{"Value": value}
	if export {
		o["Export"] = map[string]*gfn.Value{
			"Name": gfn.MakeFnSubString(fmt.Sprintf("${%s}::%s", gfn.StackName, name)),
		}
	}
	r.template.Outputs[name] = o
	r.outputs = append(r.outputs, name)
}

// newJoinedOutput defines a new output as comma-separated list
func (r *resourceSet) newJoinedOutput(name string, values []*gfn.Value, export bool) {
	r.newOutput(name, gfn.MakeFnJoin(",", values), export)
}

// newOutputFromAtt defines a new output from an attributes
func (r *resourceSet) newOutputFromAtt(name, att string, export bool) {
	r.newOutput(name, gfn.MakeFnGetAttString(att), export)
}

// makeImportValue imports output of another stack
func makeImportValue(stackName, output string) *gfn.Value {
	return gfn.MakeFnImportValueString(fmt.Sprintf("%s::%s", stackName, output))
}

// setOutput is the entrypoint that validates destination object
// and upon successful validation passes it to doSetOutput
func setOutput(obj interface{}, key, value string) error {
	e := reflect.ValueOf(obj).Elem()
	if e.Kind() != reflect.Struct {
		return fmt.Errorf("cannot use destination interface of type %q", e.Kind())
	}
	f := e.FieldByName(key)
	if !f.IsValid() && !f.CanSet() {
		return fmt.Errorf("cannot set destination field for %q", key)
	}
	return doSetOutput(f, key, value)
}

// doSetOutput handles string or slice output values
func doSetOutput(field reflect.Value, key, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
		return nil
	case reflect.Slice:
		return doSetOutputAsSlice(field, key, value)
	default:
		return fmt.Errorf("unexpected kind %q of destination field for %q", field.Kind(), key)
	}
}

// doSetOutputAsSlice sets output for fields of slice kind, it supports
// []string (for comma-separated lists defined with newJoinedOutput)
// and []byte (for BASE64-encoded values)
func doSetOutputAsSlice(field reflect.Value, key, value string) error {
	switch field.Type() {
	case reflect.ValueOf([]string{}).Type():
		// split comma-separated list and use the resulting slice
		field.Set(reflect.ValueOf(strings.Split(value, ",")))
		return nil
	case reflect.ValueOf([]byte{}).Type():
		// try to decode a string from BASE64, as certificates
		// are the only case where we expect []bytes
		data, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return errors.Wrapf(err, "decoding value of %q", key)
		}
		field.Set(reflect.ValueOf(data))
		return nil
	default:
		return fmt.Errorf("unexpected type %q of destination field for %q", field.Type(), key)
	}
}

// GetAllOutputs collects all outputs from an instance of an active stack,
// the outputs are defined by the current resourceSet, and are generally
// private to how builder chooses to define them. The destination obj is
// where outputs will be stored, it's fields are expected to match those
// that are known to the builder (namely, those are the cfnOutput* contants).
func (r *resourceSet) GetAllOutputs(stack cfn.Stack, obj interface{}) error {
	logger.Debug("processing stack outputs")
	for _, key := range r.outputs {
		value := doGetOutput(&stack, key)
		if value == nil {
			return fmt.Errorf("%s is nil", key)
		}
		if err := setOutput(obj, key, *value); err != nil {
			return errors.Wrap(err, "processing stack outputs")
		}
	}
	logger.Debug("obj = %#v", obj)
	return nil
}

// doGetOutput gets a value for a given output name, when output is not
// found in the given instance of an active stack, it will return nil
func doGetOutput(stack *cfn.Stack, key string) *string {
	for _, x := range stack.Outputs {
		if *x.OutputKey == key {
			return x.OutputValue
		}
	}
	return nil
}
