package template

import (
	"encoding/json"
)

// Commonly-used constants
const (
	AWSTemplateFormatVersion = "2010-09-09"
)

type (
	// MapOfInterfaces is an alias for map[string]interface{}
	MapOfInterfaces = map[string]interface{}
	// SliceOfInterfaces is an alias for []interface{}
	SliceOfInterfaces = []interface{}
)

// Template is a representation of a CloudFormation template
type Template struct {
	AWSTemplateFormatVersion string

	Description string
	Resources   map[string]AnyResource `json:",omitempty"`
	Outputs     map[string]Output      `json:",omitempty"`
}

// AnyResource represents a generic CloudFormation resource
type AnyResource struct {
	Type       string
	Properties interface{}
}

// Output represents a CloudFormation output definition
type Output struct {
	Description string          `json:",omitempty"`
	Value       MapOfInterfaces `json:",omitempty"`
}

// Resource defines the interface that every resource should implements
type Resource interface {
	Type() string
	Properties() interface{}
}

// NewTemplate constructs a new Template and returns the reference to it
func NewTemplate() *Template {
	return &Template{
		AWSTemplateFormatVersion: AWSTemplateFormatVersion,

		Resources: make(map[string]AnyResource),
		Outputs:   make(map[string]Output),
	}
}

// NewResource adds a resource to the template and returns a CloudFormation reference
func (t *Template) NewResource(name string, resource Resource) *Value {
	maybeSetNameTag(name, resource)
	t.Resources[name] = AnyResource{
		Type:       resource.Type(),
		Properties: resource.Properties(),
	}
	return MakeRef(name)
}

// RenderJSON will serialise the template to JSON
func (t *Template) RenderJSON() ([]byte, error) {
	return json.Marshal(t)
}

// LoadJSON will deserialise a JSON template
func (t *Template) LoadJSON(data []byte) error {
	return json.Unmarshal(data, t)
}
