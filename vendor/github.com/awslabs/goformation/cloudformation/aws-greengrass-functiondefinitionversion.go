package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSGreengrassFunctionDefinitionVersion AWS CloudFormation Resource (AWS::Greengrass::FunctionDefinitionVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-functiondefinitionversion.html
type AWSGreengrassFunctionDefinitionVersion struct {

	// DefaultConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-functiondefinitionversion.html#cfn-greengrass-functiondefinitionversion-defaultconfig
	DefaultConfig *AWSGreengrassFunctionDefinitionVersion_DefaultConfig `json:"DefaultConfig,omitempty"`

	// FunctionDefinitionId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-functiondefinitionversion.html#cfn-greengrass-functiondefinitionversion-functiondefinitionid
	FunctionDefinitionId *Value `json:"FunctionDefinitionId,omitempty"`

	// Functions AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-functiondefinitionversion.html#cfn-greengrass-functiondefinitionversion-functions
	Functions []AWSGreengrassFunctionDefinitionVersion_Function `json:"Functions,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassFunctionDefinitionVersion) AWSCloudFormationType() string {
	return "AWS::Greengrass::FunctionDefinitionVersion"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSGreengrassFunctionDefinitionVersion) MarshalJSON() ([]byte, error) {
	type Properties AWSGreengrassFunctionDefinitionVersion
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
func (r *AWSGreengrassFunctionDefinitionVersion) UnmarshalJSON(b []byte) error {
	type Properties AWSGreengrassFunctionDefinitionVersion
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
		*r = AWSGreengrassFunctionDefinitionVersion(*res.Properties)
	}

	return nil
}

// GetAllAWSGreengrassFunctionDefinitionVersionResources retrieves all AWSGreengrassFunctionDefinitionVersion items from an AWS CloudFormation template
func (t *Template) GetAllAWSGreengrassFunctionDefinitionVersionResources() map[string]AWSGreengrassFunctionDefinitionVersion {
	results := map[string]AWSGreengrassFunctionDefinitionVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSGreengrassFunctionDefinitionVersion:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::FunctionDefinitionVersion" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassFunctionDefinitionVersion{}
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

// GetAWSGreengrassFunctionDefinitionVersionWithName retrieves all AWSGreengrassFunctionDefinitionVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSGreengrassFunctionDefinitionVersionWithName(name string) (AWSGreengrassFunctionDefinitionVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSGreengrassFunctionDefinitionVersion:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::FunctionDefinitionVersion" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassFunctionDefinitionVersion{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSGreengrassFunctionDefinitionVersion{}, errors.New("resource not found")
}
