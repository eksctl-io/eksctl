package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSSecretsManagerResourcePolicy AWS CloudFormation Resource (AWS::SecretsManager::ResourcePolicy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-resourcepolicy.html
type AWSSecretsManagerResourcePolicy struct {

	// ResourcePolicy AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-resourcepolicy.html#cfn-secretsmanager-resourcepolicy-resourcepolicy
	ResourcePolicy interface{} `json:"ResourcePolicy,omitempty"`

	// SecretId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-resourcepolicy.html#cfn-secretsmanager-resourcepolicy-secretid
	SecretId *Value `json:"SecretId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSecretsManagerResourcePolicy) AWSCloudFormationType() string {
	return "AWS::SecretsManager::ResourcePolicy"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSSecretsManagerResourcePolicy) MarshalJSON() ([]byte, error) {
	type Properties AWSSecretsManagerResourcePolicy
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
func (r *AWSSecretsManagerResourcePolicy) UnmarshalJSON(b []byte) error {
	type Properties AWSSecretsManagerResourcePolicy
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
		*r = AWSSecretsManagerResourcePolicy(*res.Properties)
	}

	return nil
}

// GetAllAWSSecretsManagerResourcePolicyResources retrieves all AWSSecretsManagerResourcePolicy items from an AWS CloudFormation template
func (t *Template) GetAllAWSSecretsManagerResourcePolicyResources() map[string]AWSSecretsManagerResourcePolicy {
	results := map[string]AWSSecretsManagerResourcePolicy{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSSecretsManagerResourcePolicy:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::SecretsManager::ResourcePolicy" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSSecretsManagerResourcePolicy{}
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

// GetAWSSecretsManagerResourcePolicyWithName retrieves all AWSSecretsManagerResourcePolicy items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSSecretsManagerResourcePolicyWithName(name string) (AWSSecretsManagerResourcePolicy, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSSecretsManagerResourcePolicy:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::SecretsManager::ResourcePolicy" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSSecretsManagerResourcePolicy{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSSecretsManagerResourcePolicy{}, errors.New("resource not found")
}
