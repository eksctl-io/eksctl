package builder

import (
	gfnt "goformation/v4/cloudformation/types"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

func (r *resourceSet) defineOutput(name string, value interface{}, export bool, fn outputs.Collector) {
	r.outputs.Define(r.template, name, value, export, fn)
}

func (r *resourceSet) defineJoinedOutput(name string, values []*gfnt.Value, export bool, fn outputs.Collector) {
	r.outputs.DefineJoined(r.template, name, values, export, fn)
}

func (r *resourceSet) defineOutputFromAtt(name, logicalName, att string, export bool, fn outputs.Collector) {
	r.outputs.DefineFromAtt(r.template, name, logicalName, att, export, fn)
}

func (r *resourceSet) defineOutputWithoutCollector(name string, value interface{}, export bool) {
	r.outputs.DefineWithoutCollector(r.template, name, value, export)
}

// GetAllOutputs collects all outputs from an instance of an active stack,
// the outputs are defined by the current resourceSet
func (r *resourceSet) GetAllOutputs(stack types.Stack) error {
	logger.Debug("processing stack outputs")
	return r.outputs.MustCollect(stack)
}
