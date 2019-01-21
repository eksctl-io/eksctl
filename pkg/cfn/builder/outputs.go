package builder

import (
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"
	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

// makeImportValue imports output of another stack
func makeImportValue(stackName, output string) *gfn.Value {
	return gfn.MakeFnImportValueString(fmt.Sprintf("%s::%s", stackName, output))
}

func (r *resourceSet) defineOutput(name string, value interface{}, export bool, fn outputs.Collector) {
	r.outputs.Define(r.template, name, value, export, fn)
}

func (r *resourceSet) defineJoinedOutput(name string, values []*gfn.Value, export bool, fn outputs.Collector) {
	r.outputs.DefineJoined(r.template, name, values, export, fn)
}

func (r *resourceSet) defineOutputFromAtt(name, att string, export bool, fn outputs.Collector) {
	r.outputs.DefineFromAtt(r.template, name, att, export, fn)
}

func (r *resourceSet) defineOutputWithoutCollector(name string, value interface{}, export bool) {
	r.outputs.DefineWithoutCollector(r.template, name, value, export)
}

// GetAllOutputs collects all outputs from an instance of an active stack,
// the outputs are defined by the current resourceSet
func (r *resourceSet) GetAllOutputs(stack cfn.Stack) error {
	logger.Debug("processing stack outputs")
	return r.outputs.MustCollect(stack)
}
