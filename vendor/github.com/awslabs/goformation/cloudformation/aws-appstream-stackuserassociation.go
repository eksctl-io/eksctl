package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSAppStreamStackUserAssociation AWS CloudFormation Resource (AWS::AppStream::StackUserAssociation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appstream-stackuserassociation.html
type AWSAppStreamStackUserAssociation struct {

	// AuthenticationType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appstream-stackuserassociation.html#cfn-appstream-stackuserassociation-authenticationtype
	AuthenticationType *Value `json:"AuthenticationType,omitempty"`

	// SendEmailNotification AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appstream-stackuserassociation.html#cfn-appstream-stackuserassociation-sendemailnotification
	SendEmailNotification *Value `json:"SendEmailNotification,omitempty"`

	// StackName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appstream-stackuserassociation.html#cfn-appstream-stackuserassociation-stackname
	StackName *Value `json:"StackName,omitempty"`

	// UserName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-appstream-stackuserassociation.html#cfn-appstream-stackuserassociation-username
	UserName *Value `json:"UserName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppStreamStackUserAssociation) AWSCloudFormationType() string {
	return "AWS::AppStream::StackUserAssociation"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSAppStreamStackUserAssociation) MarshalJSON() ([]byte, error) {
	type Properties AWSAppStreamStackUserAssociation
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
func (r *AWSAppStreamStackUserAssociation) UnmarshalJSON(b []byte) error {
	type Properties AWSAppStreamStackUserAssociation
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
		*r = AWSAppStreamStackUserAssociation(*res.Properties)
	}

	return nil
}

// GetAllAWSAppStreamStackUserAssociationResources retrieves all AWSAppStreamStackUserAssociation items from an AWS CloudFormation template
func (t *Template) GetAllAWSAppStreamStackUserAssociationResources() map[string]AWSAppStreamStackUserAssociation {
	results := map[string]AWSAppStreamStackUserAssociation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSAppStreamStackUserAssociation:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AppStream::StackUserAssociation" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSAppStreamStackUserAssociation{}
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

// GetAWSAppStreamStackUserAssociationWithName retrieves all AWSAppStreamStackUserAssociation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSAppStreamStackUserAssociationWithName(name string) (AWSAppStreamStackUserAssociation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSAppStreamStackUserAssociation:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::AppStream::StackUserAssociation" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSAppStreamStackUserAssociation{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSAppStreamStackUserAssociation{}, errors.New("resource not found")
}
