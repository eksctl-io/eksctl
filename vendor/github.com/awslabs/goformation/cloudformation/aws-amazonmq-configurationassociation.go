package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSAmazonMQConfigurationAssociation AWS CloudFormation Resource (AWS::AmazonMQ::ConfigurationAssociation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-amazonmq-configurationassociation.html
type AWSAmazonMQConfigurationAssociation struct {

	// Broker AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-amazonmq-configurationassociation.html#cfn-amazonmq-configurationassociation-broker
	Broker *Value `json:"Broker,omitempty"`

	// Configuration AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-amazonmq-configurationassociation.html#cfn-amazonmq-configurationassociation-configuration
	Configuration *AWSAmazonMQConfigurationAssociation_ConfigurationId `json:"Configuration,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAmazonMQConfigurationAssociation) AWSCloudFormationType() string {
	return "AWS::AmazonMQ::ConfigurationAssociation"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSAmazonMQConfigurationAssociation) MarshalJSON() ([]byte, error) {
	type Properties AWSAmazonMQConfigurationAssociation
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
func (r *AWSAmazonMQConfigurationAssociation) UnmarshalJSON(b []byte) error {
	type Properties AWSAmazonMQConfigurationAssociation
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
		*r = AWSAmazonMQConfigurationAssociation(*res.Properties)
	}

	return nil
}

// GetAllAWSAmazonMQConfigurationAssociationResources retrieves all AWSAmazonMQConfigurationAssociation items from an AWS CloudFormation template
func (t *Template) GetAllAWSAmazonMQConfigurationAssociationResources() map[string]AWSAmazonMQConfigurationAssociation {
	results := map[string]AWSAmazonMQConfigurationAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSAmazonMQConfigurationAssociation:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AmazonMQ::ConfigurationAssociation" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSAmazonMQConfigurationAssociation{}
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

// GetAWSAmazonMQConfigurationAssociationWithName retrieves all AWSAmazonMQConfigurationAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSAmazonMQConfigurationAssociationWithName(name string) (AWSAmazonMQConfigurationAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSAmazonMQConfigurationAssociation:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AmazonMQ::ConfigurationAssociation" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSAmazonMQConfigurationAssociation{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSAmazonMQConfigurationAssociation{}, errors.New("resource not found")
}
