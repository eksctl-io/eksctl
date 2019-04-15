package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSGreengrassLoggerDefinition AWS CloudFormation Resource (AWS::Greengrass::LoggerDefinition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-loggerdefinition.html
type AWSGreengrassLoggerDefinition struct {

	// InitialVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-loggerdefinition.html#cfn-greengrass-loggerdefinition-initialversion
	InitialVersion *AWSGreengrassLoggerDefinition_LoggerDefinitionVersion `json:"InitialVersion,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-loggerdefinition.html#cfn-greengrass-loggerdefinition-name
	Name *Value `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassLoggerDefinition) AWSCloudFormationType() string {
	return "AWS::Greengrass::LoggerDefinition"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSGreengrassLoggerDefinition) MarshalJSON() ([]byte, error) {
	type Properties AWSGreengrassLoggerDefinition
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
func (r *AWSGreengrassLoggerDefinition) UnmarshalJSON(b []byte) error {
	type Properties AWSGreengrassLoggerDefinition
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
		*r = AWSGreengrassLoggerDefinition(*res.Properties)
	}

	return nil
}

// GetAllAWSGreengrassLoggerDefinitionResources retrieves all AWSGreengrassLoggerDefinition items from an AWS CloudFormation template
func (t *Template) GetAllAWSGreengrassLoggerDefinitionResources() map[string]AWSGreengrassLoggerDefinition {
	results := map[string]AWSGreengrassLoggerDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSGreengrassLoggerDefinition:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::LoggerDefinition" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassLoggerDefinition{}
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

// GetAWSGreengrassLoggerDefinitionWithName retrieves all AWSGreengrassLoggerDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSGreengrassLoggerDefinitionWithName(name string) (AWSGreengrassLoggerDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSGreengrassLoggerDefinition:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::LoggerDefinition" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassLoggerDefinition{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSGreengrassLoggerDefinition{}, errors.New("resource not found")
}
