package intrinsics

import (
	"fmt"

	yamlwrapper "github.com/sanathkr/yaml"
)

// IntrinsicHandler is a function that applies an intrinsic function and returns
// the response that should be placed in it's place. An intrinsic handler function
// is passed the name of the intrinsic function (e.g. Fn::Join), and the object
// to apply it to (as an interface{}), and should return the resolved object (as an interface{}).
type IntrinsicHandler func(string, interface{}, interface{}) interface{}

// ProcessorOptions allows customisation of the intrinsic function processor behaviour.
// This allows disabling the processing of intrinsics,
// overriding of the handlers for each intrinsic function type,
// and overriding template parameters.
type ProcessorOptions struct {
	IntrinsicHandlerOverrides map[string]IntrinsicHandler
	ParameterOverrides        map[string]interface{}
	NoProcess                 bool
	ProcessOnlyGlobals        bool
	EvaluateConditions        bool
}

// ProcessYAML recursively searches through a byte array of JSON data for all
// AWS CloudFormation intrinsic functions, resolves them, and then returns
// the resulting  interface{} object.
func ProcessYAML(input []byte, options *ProcessorOptions) ([]byte, error) {

	// Convert short form intrinsic functions (e.g. !Sub) to long form
	registerTagMarshallers()

	data, err := yamlwrapper.YAMLToJSON(input)
	if err != nil {
		return nil, fmt.Errorf("invalid YAML template: %s", err)
	}
	return data, nil
}
