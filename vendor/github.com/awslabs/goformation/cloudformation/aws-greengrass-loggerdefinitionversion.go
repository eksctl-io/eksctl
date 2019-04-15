package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSGreengrassLoggerDefinitionVersion AWS CloudFormation Resource (AWS::Greengrass::LoggerDefinitionVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-loggerdefinitionversion.html
type AWSGreengrassLoggerDefinitionVersion struct {

	// LoggerDefinitionId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-loggerdefinitionversion.html#cfn-greengrass-loggerdefinitionversion-loggerdefinitionid
	LoggerDefinitionId *Value `json:"LoggerDefinitionId,omitempty"`

	// Loggers AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-loggerdefinitionversion.html#cfn-greengrass-loggerdefinitionversion-loggers
	Loggers []AWSGreengrassLoggerDefinitionVersion_Logger `json:"Loggers,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassLoggerDefinitionVersion) AWSCloudFormationType() string {
	return "AWS::Greengrass::LoggerDefinitionVersion"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSGreengrassLoggerDefinitionVersion) MarshalJSON() ([]byte, error) {
	type Properties AWSGreengrassLoggerDefinitionVersion
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
func (r *AWSGreengrassLoggerDefinitionVersion) UnmarshalJSON(b []byte) error {
	type Properties AWSGreengrassLoggerDefinitionVersion
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
		*r = AWSGreengrassLoggerDefinitionVersion(*res.Properties)
	}

	return nil
}

// GetAllAWSGreengrassLoggerDefinitionVersionResources retrieves all AWSGreengrassLoggerDefinitionVersion items from an AWS CloudFormation template
func (t *Template) GetAllAWSGreengrassLoggerDefinitionVersionResources() map[string]AWSGreengrassLoggerDefinitionVersion {
	results := map[string]AWSGreengrassLoggerDefinitionVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSGreengrassLoggerDefinitionVersion:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::LoggerDefinitionVersion" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassLoggerDefinitionVersion{}
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

// GetAWSGreengrassLoggerDefinitionVersionWithName retrieves all AWSGreengrassLoggerDefinitionVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSGreengrassLoggerDefinitionVersionWithName(name string) (AWSGreengrassLoggerDefinitionVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSGreengrassLoggerDefinitionVersion:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::LoggerDefinitionVersion" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassLoggerDefinitionVersion{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSGreengrassLoggerDefinitionVersion{}, errors.New("resource not found")
}
