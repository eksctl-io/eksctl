package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSLambdaLayerVersionPermission AWS CloudFormation Resource (AWS::Lambda::LayerVersionPermission)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-layerversionpermission.html
type AWSLambdaLayerVersionPermission struct {

	// Action AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-layerversionpermission.html#cfn-lambda-layerversionpermission-action
	Action *Value `json:"Action,omitempty"`

	// LayerVersionArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-layerversionpermission.html#cfn-lambda-layerversionpermission-layerversionarn
	LayerVersionArn *Value `json:"LayerVersionArn,omitempty"`

	// OrganizationId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-layerversionpermission.html#cfn-lambda-layerversionpermission-organizationid
	OrganizationId *Value `json:"OrganizationId,omitempty"`

	// Principal AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-layerversionpermission.html#cfn-lambda-layerversionpermission-principal
	Principal *Value `json:"Principal,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSLambdaLayerVersionPermission) AWSCloudFormationType() string {
	return "AWS::Lambda::LayerVersionPermission"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSLambdaLayerVersionPermission) MarshalJSON() ([]byte, error) {
	type Properties AWSLambdaLayerVersionPermission
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
func (r *AWSLambdaLayerVersionPermission) UnmarshalJSON(b []byte) error {
	type Properties AWSLambdaLayerVersionPermission
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
		*r = AWSLambdaLayerVersionPermission(*res.Properties)
	}

	return nil
}

// GetAllAWSLambdaLayerVersionPermissionResources retrieves all AWSLambdaLayerVersionPermission items from an AWS CloudFormation template
func (t *Template) GetAllAWSLambdaLayerVersionPermissionResources() map[string]AWSLambdaLayerVersionPermission {
	results := map[string]AWSLambdaLayerVersionPermission{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSLambdaLayerVersionPermission:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Lambda::LayerVersionPermission" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSLambdaLayerVersionPermission{}
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

// GetAWSLambdaLayerVersionPermissionWithName retrieves all AWSLambdaLayerVersionPermission items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSLambdaLayerVersionPermissionWithName(name string) (AWSLambdaLayerVersionPermission, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSLambdaLayerVersionPermission:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Lambda::LayerVersionPermission" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSLambdaLayerVersionPermission{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSLambdaLayerVersionPermission{}, errors.New("resource not found")
}
