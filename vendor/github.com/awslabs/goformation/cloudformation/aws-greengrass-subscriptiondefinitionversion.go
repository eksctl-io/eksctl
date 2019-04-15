package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSGreengrassSubscriptionDefinitionVersion AWS CloudFormation Resource (AWS::Greengrass::SubscriptionDefinitionVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-subscriptiondefinitionversion.html
type AWSGreengrassSubscriptionDefinitionVersion struct {

	// SubscriptionDefinitionId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-subscriptiondefinitionversion.html#cfn-greengrass-subscriptiondefinitionversion-subscriptiondefinitionid
	SubscriptionDefinitionId *Value `json:"SubscriptionDefinitionId,omitempty"`

	// Subscriptions AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-greengrass-subscriptiondefinitionversion.html#cfn-greengrass-subscriptiondefinitionversion-subscriptions
	Subscriptions []AWSGreengrassSubscriptionDefinitionVersion_Subscription `json:"Subscriptions,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGreengrassSubscriptionDefinitionVersion) AWSCloudFormationType() string {
	return "AWS::Greengrass::SubscriptionDefinitionVersion"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSGreengrassSubscriptionDefinitionVersion) MarshalJSON() ([]byte, error) {
	type Properties AWSGreengrassSubscriptionDefinitionVersion
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
func (r *AWSGreengrassSubscriptionDefinitionVersion) UnmarshalJSON(b []byte) error {
	type Properties AWSGreengrassSubscriptionDefinitionVersion
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
		*r = AWSGreengrassSubscriptionDefinitionVersion(*res.Properties)
	}

	return nil
}

// GetAllAWSGreengrassSubscriptionDefinitionVersionResources retrieves all AWSGreengrassSubscriptionDefinitionVersion items from an AWS CloudFormation template
func (t *Template) GetAllAWSGreengrassSubscriptionDefinitionVersionResources() map[string]AWSGreengrassSubscriptionDefinitionVersion {
	results := map[string]AWSGreengrassSubscriptionDefinitionVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSGreengrassSubscriptionDefinitionVersion:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::SubscriptionDefinitionVersion" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassSubscriptionDefinitionVersion{}
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

// GetAWSGreengrassSubscriptionDefinitionVersionWithName retrieves all AWSGreengrassSubscriptionDefinitionVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSGreengrassSubscriptionDefinitionVersionWithName(name string) (AWSGreengrassSubscriptionDefinitionVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSGreengrassSubscriptionDefinitionVersion:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Greengrass::SubscriptionDefinitionVersion" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSGreengrassSubscriptionDefinitionVersion{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSGreengrassSubscriptionDefinitionVersion{}, errors.New("resource not found")
}
