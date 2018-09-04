package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSSESConfigurationSet AWS CloudFormation Resource (AWS::SES::ConfigurationSet)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ses-configurationset.html
type AWSSESConfigurationSet struct {

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ses-configurationset.html#cfn-ses-configurationset-name
	Name *StringIntrinsic `json:"Name,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSESConfigurationSet) AWSCloudFormationType() string {
	return "AWS::SES::ConfigurationSet"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSSESConfigurationSet) MarshalJSON() ([]byte, error) {
	type Properties AWSSESConfigurationSet
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
func (r *AWSSESConfigurationSet) UnmarshalJSON(b []byte) error {
	type Properties AWSSESConfigurationSet
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
		*r = AWSSESConfigurationSet(*res.Properties)
	}

	return nil
}

// GetAllAWSSESConfigurationSetResources retrieves all AWSSESConfigurationSet items from an AWS CloudFormation template
func (t *Template) GetAllAWSSESConfigurationSetResources() map[string]AWSSESConfigurationSet {
	results := map[string]AWSSESConfigurationSet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSSESConfigurationSet:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::SES::ConfigurationSet" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						var result AWSSESConfigurationSet
						if err := json.Unmarshal(b, &result); err == nil {
							results[name] = result
						}
					}
				}
			}
		}
	}
	return results
}

// GetAWSSESConfigurationSetWithName retrieves all AWSSESConfigurationSet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSSESConfigurationSetWithName(name string) (AWSSESConfigurationSet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSSESConfigurationSet:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::SES::ConfigurationSet" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						var result AWSSESConfigurationSet
						if err := json.Unmarshal(b, &result); err == nil {
							return result, nil
						}
					}
				}
			}
		}
	}
	return AWSSESConfigurationSet{}, errors.New("resource not found")
}
