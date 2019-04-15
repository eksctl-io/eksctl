package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSGreengrassDeviceDefinitionVersion AWS CloudFormation Resource (AWS::Greengrass::DeviceDefinitionVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-devicedefinitionversion.html
type AWSGreengrassDeviceDefinitionVersion struct {

	// DeviceDefinitionId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-devicedefinitionversion.html#cfn-greengrass-devicedefinitionversion-devicedefinitionid
	DeviceDefinitionId *Value `json:"DeviceDefinitionId,omitempty"`

	// Devices AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-devicedefinitionversion.html#cfn-greengrass-devicedefinitionversion-devices
	Devices []AWSGreengrassDeviceDefinitionVersion_Device `json:"Devices,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassDeviceDefinitionVersion) AWSCloudFormationType() string {
	return "AWS::Greengrass::DeviceDefinitionVersion"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSGreengrassDeviceDefinitionVersion) MarshalJSON() ([]byte, error) {
	type Properties AWSGreengrassDeviceDefinitionVersion
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
func (r *AWSGreengrassDeviceDefinitionVersion) UnmarshalJSON(b []byte) error {
	type Properties AWSGreengrassDeviceDefinitionVersion
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
		*r = AWSGreengrassDeviceDefinitionVersion(*res.Properties)
	}

	return nil
}

// GetAllAWSGreengrassDeviceDefinitionVersionResources retrieves all AWSGreengrassDeviceDefinitionVersion items from an AWS CloudFormation template
func (t *Template) GetAllAWSGreengrassDeviceDefinitionVersionResources() map[string]AWSGreengrassDeviceDefinitionVersion {
	results := map[string]AWSGreengrassDeviceDefinitionVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSGreengrassDeviceDefinitionVersion:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::DeviceDefinitionVersion" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassDeviceDefinitionVersion{}
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

// GetAWSGreengrassDeviceDefinitionVersionWithName retrieves all AWSGreengrassDeviceDefinitionVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSGreengrassDeviceDefinitionVersionWithName(name string) (AWSGreengrassDeviceDefinitionVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSGreengrassDeviceDefinitionVersion:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::DeviceDefinitionVersion" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassDeviceDefinitionVersion{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSGreengrassDeviceDefinitionVersion{}, errors.New("resource not found")
}
