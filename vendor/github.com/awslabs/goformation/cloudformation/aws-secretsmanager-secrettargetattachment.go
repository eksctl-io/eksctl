package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSSecretsManagerSecretTargetAttachment AWS CloudFormation Resource (AWS::SecretsManager::SecretTargetAttachment)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-secrettargetattachment.html
type AWSSecretsManagerSecretTargetAttachment struct {

	// SecretId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-secrettargetattachment.html#cfn-secretsmanager-secrettargetattachment-secretid
	SecretId *Value `json:"SecretId,omitempty"`

	// TargetId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-secrettargetattachment.html#cfn-secretsmanager-secrettargetattachment-targetid
	TargetId *Value `json:"TargetId,omitempty"`

	// TargetType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-secretsmanager-secrettargetattachment.html#cfn-secretsmanager-secrettargetattachment-targettype
	TargetType *Value `json:"TargetType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSecretsManagerSecretTargetAttachment) AWSCloudFormationType() string {
	return "AWS::SecretsManager::SecretTargetAttachment"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSSecretsManagerSecretTargetAttachment) MarshalJSON() ([]byte, error) {
	type Properties AWSSecretsManagerSecretTargetAttachment
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
func (r *AWSSecretsManagerSecretTargetAttachment) UnmarshalJSON(b []byte) error {
	type Properties AWSSecretsManagerSecretTargetAttachment
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
		*r = AWSSecretsManagerSecretTargetAttachment(*res.Properties)
	}

	return nil
}

// GetAllAWSSecretsManagerSecretTargetAttachmentResources retrieves all AWSSecretsManagerSecretTargetAttachment items from an AWS CloudFormation template
func (t *Template) GetAllAWSSecretsManagerSecretTargetAttachmentResources() map[string]AWSSecretsManagerSecretTargetAttachment {
	results := map[string]AWSSecretsManagerSecretTargetAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSSecretsManagerSecretTargetAttachment:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::SecretsManager::SecretTargetAttachment" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSSecretsManagerSecretTargetAttachment{}
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

// GetAWSSecretsManagerSecretTargetAttachmentWithName retrieves all AWSSecretsManagerSecretTargetAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSSecretsManagerSecretTargetAttachmentWithName(name string) (AWSSecretsManagerSecretTargetAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSSecretsManagerSecretTargetAttachment:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::SecretsManager::SecretTargetAttachment" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSSecretsManagerSecretTargetAttachment{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSSecretsManagerSecretTargetAttachment{}, errors.New("resource not found")
}
