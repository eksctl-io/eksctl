package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSGreengrassDeviceDefinition AWS CloudFormation Resource (AWS::Greengrass::DeviceDefinition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-devicedefinition.html
type AWSGreengrassDeviceDefinition struct {

	// InitialVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-devicedefinition.html#cfn-greengrass-devicedefinition-initialversion
	InitialVersion *AWSGreengrassDeviceDefinition_DeviceDefinitionVersion `json:"InitialVersion,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-devicedefinition.html#cfn-greengrass-devicedefinition-name
	Name *Value `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassDeviceDefinition) AWSCloudFormationType() string {
	return "AWS::Greengrass::DeviceDefinition"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSGreengrassDeviceDefinition) MarshalJSON() ([]byte, error) {
	type Properties AWSGreengrassDeviceDefinition
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
func (r *AWSGreengrassDeviceDefinition) UnmarshalJSON(b []byte) error {
	type Properties AWSGreengrassDeviceDefinition
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
		*r = AWSGreengrassDeviceDefinition(*res.Properties)
	}

	return nil
}

// GetAllAWSGreengrassDeviceDefinitionResources retrieves all AWSGreengrassDeviceDefinition items from an AWS CloudFormation template
func (t *Template) GetAllAWSGreengrassDeviceDefinitionResources() map[string]AWSGreengrassDeviceDefinition {
	results := map[string]AWSGreengrassDeviceDefinition{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSGreengrassDeviceDefinition:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::DeviceDefinition" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassDeviceDefinition{}
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

// GetAWSGreengrassDeviceDefinitionWithName retrieves all AWSGreengrassDeviceDefinition items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSGreengrassDeviceDefinitionWithName(name string) (AWSGreengrassDeviceDefinition, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSGreengrassDeviceDefinition:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::DeviceDefinition" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassDeviceDefinition{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSGreengrassDeviceDefinition{}, errors.New("resource not found")
}
