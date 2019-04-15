package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSGreengrassGroupVersion AWS CloudFormation Resource (AWS::Greengrass::GroupVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-groupversion.html
type AWSGreengrassGroupVersion struct {

	// ConnectorDefinitionVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-groupversion.html#cfn-greengrass-groupversion-connectordefinitionversionarn
	ConnectorDefinitionVersionArn *Value `json:"ConnectorDefinitionVersionArn,omitempty"`

	// CoreDefinitionVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-groupversion.html#cfn-greengrass-groupversion-coredefinitionversionarn
	CoreDefinitionVersionArn *Value `json:"CoreDefinitionVersionArn,omitempty"`

	// DeviceDefinitionVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-groupversion.html#cfn-greengrass-groupversion-devicedefinitionversionarn
	DeviceDefinitionVersionArn *Value `json:"DeviceDefinitionVersionArn,omitempty"`

	// FunctionDefinitionVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-groupversion.html#cfn-greengrass-groupversion-functiondefinitionversionarn
	FunctionDefinitionVersionArn *Value `json:"FunctionDefinitionVersionArn,omitempty"`

	// GroupId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-groupversion.html#cfn-greengrass-groupversion-groupid
	GroupId *Value `json:"GroupId,omitempty"`

	// LoggerDefinitionVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-groupversion.html#cfn-greengrass-groupversion-loggerdefinitionversionarn
	LoggerDefinitionVersionArn *Value `json:"LoggerDefinitionVersionArn,omitempty"`

	// ResourceDefinitionVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-groupversion.html#cfn-greengrass-groupversion-resourcedefinitionversionarn
	ResourceDefinitionVersionArn *Value `json:"ResourceDefinitionVersionArn,omitempty"`

	// SubscriptionDefinitionVersionArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-groupversion.html#cfn-greengrass-groupversion-subscriptiondefinitionversionarn
	SubscriptionDefinitionVersionArn *Value `json:"SubscriptionDefinitionVersionArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassGroupVersion) AWSCloudFormationType() string {
	return "AWS::Greengrass::GroupVersion"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSGreengrassGroupVersion) MarshalJSON() ([]byte, error) {
	type Properties AWSGreengrassGroupVersion
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
func (r *AWSGreengrassGroupVersion) UnmarshalJSON(b []byte) error {
	type Properties AWSGreengrassGroupVersion
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
		*r = AWSGreengrassGroupVersion(*res.Properties)
	}

	return nil
}

// GetAllAWSGreengrassGroupVersionResources retrieves all AWSGreengrassGroupVersion items from an AWS CloudFormation template
func (t *Template) GetAllAWSGreengrassGroupVersionResources() map[string]AWSGreengrassGroupVersion {
	results := map[string]AWSGreengrassGroupVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSGreengrassGroupVersion:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::GroupVersion" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassGroupVersion{}
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

// GetAWSGreengrassGroupVersionWithName retrieves all AWSGreengrassGroupVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSGreengrassGroupVersionWithName(name string) (AWSGreengrassGroupVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSGreengrassGroupVersion:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::GroupVersion" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassGroupVersion{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSGreengrassGroupVersion{}, errors.New("resource not found")
}
