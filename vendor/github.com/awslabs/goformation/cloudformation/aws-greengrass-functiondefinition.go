package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSGreengrassFunctionDefinition AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-functiondefinition.html
type AWSGreengrassFunctionDefinition struct {

	// InitialVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-functiondefinition.html#cfn-greengrass-functiondefinition-initialversion
	InitialVersion *AWSGreengrassFunctionDefinition_FunctionDefinitionVersion `json:"InitialVersion,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-functiondefinition.html#cfn-greengrass-functiondefinition-name
	Name *Value `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinition) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinition"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSGreengrassFunctionDefinition) MarshalJSON() ([]byte, error) {
	type Properties AWSGreengrassFunctionDefinition
	return json.Marshal(&struct {
		Type       string
		Properties Properties
	}{
		Type:       r.AWSCloudFormationType(),
		Properties: (Properties)(*r),
	})
}

// UnmarshalJSON is a custom JSON unmarshalling hook that strips the outer
// AWS CloudFormation resource object, and just keeps the 'Properties' field.
func (r *AWSGreengrassFunctionDefinition) UnmarshalJSON(b []byte) error {
	type Properties AWSGreengrassFunctionDefinition
	res := &struct {
		Type       string
		Properties *Properties
	}{}
	if err := json.Unmarshal(b, &res); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return err
	}

	// If the resource has no Properties set, it could be nil
	if res.Properties != nil {
		*r = AWSGreengrassFunctionDefinition(*res.Properties)
	}

	return nil
}

// GetAllAWSGreengrassFunctionDefinitionResources retrieves all AWSGreengrassFunctionDefinition items from an AWS CloudFormation template
func (t *Template) GetAllAWSGreengrassFunctionDefinitionResources() map[string]AWSGreengrassFunctionDefinition {
	results := map[string]AWSGreengrassFunctionDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSGreengrassFunctionDefinition:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::FunctionDefinition" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassFunctionDefinition{}
						if err := result.UnmarshalJSON(b); err == nil {
							results[name] = *result
						}
					}
				}
			}
		}
	}
	return results
}

// GetAWSGreengrassFunctionDefinitionWithName retrieves all AWSGreengrassFunctionDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSGreengrassFunctionDefinitionWithName(name string) (AWSGreengrassFunctionDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSGreengrassFunctionDefinition:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::FunctionDefinition" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassFunctionDefinition{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSGreengrassFunctionDefinition{}, errors.New("resource not found")
}
