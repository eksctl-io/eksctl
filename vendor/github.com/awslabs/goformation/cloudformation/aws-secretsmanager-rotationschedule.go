package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSSecretsManagerRotationSchedule AWS CloudFormation Resource (AWS::SecretsManager::RotationSchedule)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-rotationschedule.html
type AWSSecretsManagerRotationSchedule struct {

	// RotationLambdaARN AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-rotationschedule.html#cfn-secretsmanager-rotationschedule-rotationlambdaarn
	RotationLambdaARN *Value `json:"RotationLambdaARN,omitempty"`

	// RotationRules AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-rotationschedule.html#cfn-secretsmanager-rotationschedule-rotationrules
	RotationRules *AWSSecretsManagerRotationSchedule_RotationRules `json:"RotationRules,omitempty"`

	// SecretId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-rotationschedule.html#cfn-secretsmanager-rotationschedule-secretid
	SecretId *Value `json:"SecretId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSecretsManagerRotationSchedule) AWSCloudFormationType() string {
	return "AWS::SecretsManager::RotationSchedule"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSSecretsManagerRotationSchedule) MarshalJSON() ([]byte, error) {
	type Properties AWSSecretsManagerRotationSchedule
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
func (r *AWSSecretsManagerRotationSchedule) UnmarshalJSON(b []byte) error {
	type Properties AWSSecretsManagerRotationSchedule
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
		*r = AWSSecretsManagerRotationSchedule(*res.Properties)
	}

	return nil
}

// GetAllAWSSecretsManagerRotationScheduleResources retrieves all AWSSecretsManagerRotationSchedule items from an AWS CloudFormation template
func (t *Template) GetAllAWSSecretsManagerRotationScheduleResources() map[string]AWSSecretsManagerRotationSchedule {
	results := map[string]AWSSecretsManagerRotationSchedule{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSSecretsManagerRotationSchedule:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::SecretsManager::RotationSchedule" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSSecretsManagerRotationSchedule{}
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

// GetAWSSecretsManagerRotationScheduleWithName retrieves all AWSSecretsManagerRotationSchedule items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSSecretsManagerRotationScheduleWithName(name string) (AWSSecretsManagerRotationSchedule, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSSecretsManagerRotationSchedule:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::SecretsManager::RotationSchedule" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSSecretsManagerRotationSchedule{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSSecretsManagerRotationSchedule{}, errors.New("resource not found")
}
