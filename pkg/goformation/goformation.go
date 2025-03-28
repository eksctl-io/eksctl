package goformation

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation"
	"github.com/weaveworks/eksctl/pkg/goformation/intrinsics"
)

//go:generate generate/generate.sh

// Open and parse a AWS CloudFormation template from file.
// Works with either JSON or YAML formatted templates.
func Open(filename string) (*cloudformation.Template, error) {
	return OpenWithOptions(filename, nil)
}

// OpenWithOptions opens and parse a AWS CloudFormation template from file.
// Works with either JSON or YAML formatted templates.
// Parsing can be tweaked via the specified options.
func OpenWithOptions(filename string, options *intrinsics.ProcessorOptions) (*cloudformation.Template, error) {

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(filename, ".json") {
		// This is definitely JSON
		return ParseJSONWithOptions(data, options)
	}

	return ParseYAMLWithOptions(data, options)
}

// ParseYAML an AWS CloudFormation template (expects a []byte of valid YAML)
func ParseYAML(data []byte) (*cloudformation.Template, error) {
	return ParseYAMLWithOptions(data, nil)
}

// ParseYAMLWithOptions an AWS CloudFormation template (expects a []byte of valid YAML)
// Parsing can be tweaked via the specified options.
func ParseYAMLWithOptions(data []byte, options *intrinsics.ProcessorOptions) (*cloudformation.Template, error) {
	// Process all AWS CloudFormation intrinsic functions (e.g. Fn::Join)
	intrinsified, err := intrinsics.ProcessYAML(data, options)
	if err != nil {
		return nil, err
	}

	return unmarshal(intrinsified)

}

// ParseJSON an AWS CloudFormation template (expects a []byte of valid JSON)
func ParseJSON(data []byte) (*cloudformation.Template, error) {
	return ParseJSONWithOptions(data, nil)
}

// ParseJSONWithOptions an AWS CloudFormation template (expects a []byte of valid JSON)
// Parsing can be tweaked via the specified options.
func ParseJSONWithOptions(data []byte, options *intrinsics.ProcessorOptions) (*cloudformation.Template, error) {
	return unmarshal(data)

}

func unmarshal(data []byte) (*cloudformation.Template, error) {

	template := &cloudformation.Template{}
	if err := json.Unmarshal(data, template); err != nil {
		return nil, err
	}

	return template, nil

}
