package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSGreengrassConnectorDefinitionVersion AWS CloudFormation Resource (AWS::Greengrass::ConnectorDefinitionVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-connectordefinitionversion.html
type AWSGreengrassConnectorDefinitionVersion struct {

	// ConnectorDefinitionId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-connectordefinitionversion.html#cfn-greengrass-connectordefinitionversion-connectordefinitionid
	ConnectorDefinitionId *Value `json:"ConnectorDefinitionId,omitempty"`

	// Connectors AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-connectordefinitionversion.html#cfn-greengrass-connectordefinitionversion-connectors
	Connectors []AWSGreengrassConnectorDefinitionVersion_Connector `json:"Connectors,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassConnectorDefinitionVersion) AWSCloudFormationType() string {
	return "AWS::Greengrass::ConnectorDefinitionVersion"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSGreengrassConnectorDefinitionVersion) MarshalJSON() ([]byte, error) {
	type Properties AWSGreengrassConnectorDefinitionVersion
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
func (r *AWSGreengrassConnectorDefinitionVersion) UnmarshalJSON(b []byte) error {
	type Properties AWSGreengrassConnectorDefinitionVersion
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
		*r = AWSGreengrassConnectorDefinitionVersion(*res.Properties)
	}

	return nil
}

// GetAllAWSGreengrassConnectorDefinitionVersionResources retrieves all AWSGreengrassConnectorDefinitionVersion items from an AWS CloudFormation template
func (t *Template) GetAllAWSGreengrassConnectorDefinitionVersionResources() map[string]AWSGreengrassConnectorDefinitionVersion {
	results := map[string]AWSGreengrassConnectorDefinitionVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSGreengrassConnectorDefinitionVersion:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::ConnectorDefinitionVersion" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassConnectorDefinitionVersion{}
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

// GetAWSGreengrassConnectorDefinitionVersionWithName retrieves all AWSGreengrassConnectorDefinitionVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSGreengrassConnectorDefinitionVersionWithName(name string) (AWSGreengrassConnectorDefinitionVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSGreengrassConnectorDefinitionVersion:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::ConnectorDefinitionVersion" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassConnectorDefinitionVersion{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSGreengrassConnectorDefinitionVersion{}, errors.New("resource not found")
}
