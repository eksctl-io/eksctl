package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"
)

var typeToGo = map[string]string{
	"String":    "*types.Value",
	"Long":      "*types.Value",
	"Integer":   "*types.Value",
	"Double":    "*types.Value",
	"Boolean":   "*types.Value",
	"Timestamp": "string",
	"Json":      "interface{}",
	"Map":       "*types.Value",

	// Overrides to fix CF errors
	"ParameterValues": "interface{}", // fix for AWS::SSM::Association
}

var typeToPureGo = map[string]string{
	"String":    "string",
	"Long":      "int64",
	"Integer":   "int",
	"Double":    "float64",
	"Boolean":   "bool",
	"Timestamp": "string",
	"Json":      "interface{}",
	"Map":       "interface{}",

	// Overrides to fix CF errors
	"ParameterValues": "interface{}", // fix for AWS::SSM::Association
}

var typeToJSON = map[string]string{
	"String":    "string",
	"Long":      "number",
	"Integer":   "number",
	"Double":    "number",
	"Boolean":   "boolean",
	"Timestamp": "string",
	"Json":      "object",
	"Map":       "object",

	// Overrides to fix CF errors
	"ParameterValues": "object", // fix for AWS::SSM::Association
}

// Property represents an AWS CloudFormation resource property
type Property struct {

	// Documentation - A link to the AWS CloudFormation User Guide that provides information about the property.
	Documentation string `json:"Documentation"`

	// DuplicatesAllowed - If the value of the Type field is List, indicates whether AWS CloudFormation allows duplicate values.
	// If the value is true, AWS CloudFormation ignores duplicate values. If the value is false,
	// AWS CloudFormation returns an error if you submit duplicate values.
	DuplicatesAllowed bool `json:"DuplicatesAllowed"`

	// ItemType - If the value of the Type field is List or Map, indicates the type of list or map if they contain
	// non-primitive types. Otherwise, this field is omitted. For lists or maps that contain primitive
	// types, the PrimitiveItemType property indicates the valid value type.
	//
	// A subproperty name is a valid item type. For example, if the type value is List and the item type
	//  value is PortMapping, you can specify a list of port mapping properties.
	ItemType string `json:"ItemType"`

	// PrimitiveItemType - If the value of the Type field is List or Map, indicates the type of list or map
	// if they contain primitive types. Otherwise, this field is omitted. For lists or maps that contain
	// non-primitive types, the ItemType property indicates the valid value type.
	// The valid primitive types for lists and maps are String, Long, Integer, Double, Boolean, or Timestamp.
	// For example, if the type value is List and the item type value is String, you can specify a list of strings
	// for the property. If the type value is Map and the item type value is Boolean, you can specify a string
	// to Boolean mapping for the property.
	PrimitiveItemType string `json:"PrimitiveItemType"`

	// PrimitiveType - For primitive values, the valid primitive type for the property. A primitive type is a
	// basic data type for resource property values.
	// The valid primitive types are String, Long, Integer, Double, Boolean, Timestamp or Json.
	// If valid values are a non-primitive type, this field is omitted and the Type field indicates the valid value type.
	PrimitiveType string `json:"PrimitiveType"`

	// Required indicates whether the property is required.
	Required bool `json:"Required"`

	// Type - For non-primitive types, valid values for the property. The valid types are a subproperty name,
	// List or Map. If valid values are a primitive type, this field is omitted and the PrimitiveType field
	// indicates the valid value type. A list is a comma-separated list of values. A map is a set of key-value pairs,
	// where the keys are always strings. The value type for lists and maps are indicated by the ItemType
	// or PrimitiveItemType field.
	Type string `json:"Type"`

	// UpdateType - During a stack update, the update behavior when you add, remove, or modify the property.
	// AWS CloudFormation replaces the resource when you change Immutable properties. AWS CloudFormation doesn't
	// replace the resource when you change mutable properties. Conditional updates can be mutable or immutable,
	// depending on, for example, which other properties you updated. For more information, see the relevant
	// resource type documentation.
	UpdateType string `json:"UpdateType"`

	// Types - if a property can be different types, they will be listed here
	PrimitiveTypes     []string `json:"PrimitiveTypes"`
	PrimitiveItemTypes []string `json:"PrimitiveItemTypes"`
	ItemTypes          []string `json:"ItemTypes"`
	Types              []string `json:"Types"`
}

// Schema returns a JSON Schema for the resource (as a string)
func (p Property) Schema(name, parent string) string {

	// Open the schema template and setup a counter function that will
	// available in the template to be used to detect when trailing commas
	// are required in the JSON when looping through maps
	tmpl, err := template.New("schema-property.template").Funcs(template.FuncMap{
		"counter":           counter,
		"convertToJSONType": convertTypeToJSON,
	}).ParseFiles("generate/templates/schema-property.template")

	var buf bytes.Buffer
	parentpaths := strings.Split(parent, ".")

	templateData := struct {
		Name     string
		Parent   string
		Property Property
	}{
		Name:     name,
		Parent:   parentpaths[0],
		Property: p,
	}

	// Execute the template, writing it to the buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		fmt.Printf("Error: Failed to generate property %s\n%s\n", name, err)
		os.Exit(1)
	}

	return buf.String()

}

// UnmarshalJSON is a custom unmarshaller for CloudFormation Properties.
// It's only purpose, is to validate that every property has a valid type
// set in the CloudFormation Resource Specification. This is required as
// on occasions in the past, errors in the published spec have meant that some
// properties are missing types. See github.com/awslabs/goformation/issues/300.
// This method sets any properties that have missing types to 'Json' - which in
// turn defines them as interface{} in the generated Go structs.
func (p *Property) UnmarshalJSON(data []byte) error {

	// see: https://stackoverflow.com/questions/52433467/how-to-call-json-unmarshal-inside-unmarshaljson-without-causing-stack-overflow
	type TmpProperty Property

	var unmarshalled TmpProperty
	err := json.Unmarshal(data, &unmarshalled)
	if err != nil {
		return err
	}

	*p = Property(unmarshalled)

	if !p.HasValidType() {
		fmt.Printf("Warning: auto-fixing missing property type to 'Json' for %s\n", p.Documentation)
		p.PrimitiveType = "Json"
	}

	return nil

}

// HasValidType checks whether a property has a valid type defined
// It is possible that an invalid CloudFormation Resource Specification is published
// that does not have any type information for a property. If this happens, then
// generation should fail with an error message.
func (p Property) HasValidType() bool {
	invalid := p.ItemType == "" &&
		p.PrimitiveType == "" &&
		p.PrimitiveItemType == "" &&
		p.Type == "" &&
		len(p.ItemTypes) == 0 &&
		len(p.PrimitiveTypes) == 0 &&
		len(p.PrimitiveItemTypes) == 0 &&
		len(p.Types) == 0
	return !invalid
}

// IsPolymorphic checks whether a property can be multiple different types
func (p Property) IsPolymorphic() bool {
	return len(p.PrimitiveTypes) > 0 || len(p.PrimitiveItemTypes) > 0 || len(p.PrimitiveItemTypes) > 0 || len(p.ItemTypes) > 0 || len(p.Types) > 0
}

// IsPrimitive checks whether a property is a primitive type
func (p Property) IsPrimitive() bool {
	return p.PrimitiveType != ""
}

// IsNumeric checks whether a property is numeric
func (p Property) IsNumeric() bool {
	return p.IsPrimitive() &&
		(p.PrimitiveType == "Long" ||
			p.PrimitiveType == "Integer" ||
			p.PrimitiveType == "Double" ||
			p.PrimitiveType == "Boolean")
}

// IsMap checks whether a property should be a map (map[string]...)
func (p Property) IsMap() bool {
	return p.Type == "Map"
}

// IsList checks whether a property should be a list ([]...)
func (p Property) IsList() bool {
	return p.Type == "List"
}

// IsCustomType checks wither a property is a custom type
func (p Property) IsCustomType() bool {
	return p.PrimitiveType == "" && p.ItemType == "" && p.PrimitiveItemType == ""
}

// GoType returns the correct type for this property
// within a Go struct. For example, []string or map[string]AWSLambdaFunction_VpcConfig
func (p Property) GoType(typename, basename, name, packageName string) string {

	if p.ItemType == "Tag" {
		if packageName == "cloudformation" {
			return "[]Tag"
		}
		return "[]cloudformation.Tag"
	}

	if p.IsPolymorphic() {

		generatePolymorphicProperty(typename, basename+"_"+name, p)
		return basename + "_" + name

	}

	if p.IsMap() {

		if p.convertTypeToGo() != "" {
			return "map[string]" + p.convertTypeToGo()
		}

		return "map[string]" + basename + "_" + p.ItemType

	}

	if p.IsList() {

		if p.convertTypeToGo() != "" {
			if p.GoTypeIsValue() {
				return p.convertTypeToGo()
			}
			return "[]" + p.convertTypeToGo()
		}

		return "[]" + basename + "_" + p.ItemType

	}

	if p.IsCustomType() {
		return basename + "_" + p.Type
	}

	// Must be a primitive value
	return convertTypeToGo(p.PrimitiveType)

}

// GetJSONPrimitiveType returns the correct primitive property type for a JSON Schema.
// If the property is a list/map, then it will return the type of the items.
func (p Property) GetJSONPrimitiveType() string {
	return p.convertTypeToJSON()
}

// HasJSONPrimitiveType if GetJSONPrimitiveType is not ""
func (p Property) HasJSONPrimitiveType() bool {
	return p.convertTypeToJSON() != ""
}

func (p Property) convertTypeToGo() string {
	if p.PrimitiveType != "" {
		return convertTypeToGo(p.PrimitiveType)
	} else if p.PrimitiveItemType != "" {
		return convertTypeToGo(p.PrimitiveItemType)
	} else {
		return convertTypeToGo(p.ItemType)
	}
}

func (p Property) GoTypeIsValue() bool {
	var goType string
	if p.PrimitiveType != "" {
		goType = convertTypeToGo(p.PrimitiveType)
	} else if p.PrimitiveItemType != "" {
		goType = convertTypeToGo(p.PrimitiveItemType)
	} else {
		goType = convertTypeToGo(p.ItemType)
	}
	if strings.Contains(goType, "Value") {
		return true
	}
	for _, t := range append(append(p.PrimitiveTypes, p.PrimitiveItemTypes...), p.ItemTypes...) {
		if strings.Contains(convertTypeToGo(t), "Value") {
			return true
		}
	}
	return false
}

func (p Property) convertTypeToJSON() string {
	if p.PrimitiveType != "" {
		return convertTypeToJSON(p.PrimitiveType)
	} else if p.PrimitiveItemType != "" {
		return convertTypeToJSON(p.PrimitiveItemType)
	} else {
		return convertTypeToJSON(p.ItemType)
	}
}

func convertTypeToGo(name string) string {
	t, ok := typeToGo[name]
	if !ok {
		return ""
	}
	return t
}

func convertTypeToPureGo(name string) string {
	t, ok := typeToPureGo[name]
	if !ok {
		return ""
	}
	return t
}

func convertTypeToJSON(name string) string {
	t, ok := typeToJSON[name]
	if !ok {
		return ""
	}
	return t
}
