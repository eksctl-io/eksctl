package builder

import (
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"
	"github.com/kris-nova/logger"
)

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

// GetAllOutputs collects all outputs from an instance of an active stack,
// the outputs are defined by the current resourceSet, and are generally
// private to how builder chooses to define them
func (r *resourceSet) GetAllOutputs(stack cfn.Stack, outputs map[string]string) error {
	logger.Debug("processing stack outputs")
	for _, key := range r.outputs {
		var value *string
		for _, x := range stack.Outputs {
			if *x.OutputKey == key {
				value = x.OutputValue
				break
			}
		}
		if value == nil {
			return fmt.Errorf("%s is nil", key)
		}
		outputs[key] = *value
	}
	return nil
}
