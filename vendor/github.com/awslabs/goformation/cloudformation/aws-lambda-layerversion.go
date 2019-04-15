package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSLambdaLayerVersion AWS CloudFormation Resource (AWS::Lambda::LayerVersion)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-layerversion.html
type AWSLambdaLayerVersion struct {

	// CompatibleRuntimes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-layerversion.html#cfn-lambda-layerversion-compatibleruntimes
	CompatibleRuntimes []*Value `json:"CompatibleRuntimes,omitempty"`

	// Content AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-layerversion.html#cfn-lambda-layerversion-content
	Content *AWSLambdaLayerVersion_Content `json:"Content,omitempty"`

	// Description AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-layerversion.html#cfn-lambda-layerversion-description
	Description *Value `json:"Description,omitempty"`

	// LayerName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-layerversion.html#cfn-lambda-layerversion-layername
	LayerName *Value `json:"LayerName,omitempty"`

	// LicenseInfo AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-layerversion.html#cfn-lambda-layerversion-licenseinfo
	LicenseInfo *Value `json:"LicenseInfo,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSLambdaLayerVersion) AWSCloudFormationType() string {
	return "AWS::Lambda::LayerVersion"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSLambdaLayerVersion) MarshalJSON() ([]byte, error) {
	type Properties AWSLambdaLayerVersion
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
func (r *AWSLambdaLayerVersion) UnmarshalJSON(b []byte) error {
	type Properties AWSLambdaLayerVersion
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
		*r = AWSLambdaLayerVersion(*res.Properties)
	}

	return nil
}

// GetAllAWSLambdaLayerVersionResources retrieves all AWSLambdaLayerVersion items from an AWS CloudFormation template
func (t *Template) GetAllAWSLambdaLayerVersionResources() map[string]AWSLambdaLayerVersion {
	results := map[string]AWSLambdaLayerVersion{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSLambdaLayerVersion:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Lambda::LayerVersion" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSLambdaLayerVersion{}
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

// GetAWSLambdaLayerVersionWithName retrieves all AWSLambdaLayerVersion items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSLambdaLayerVersionWithName(name string) (AWSLambdaLayerVersion, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSLambdaLayerVersion:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Lambda::LayerVersion" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSLambdaLayerVersion{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSLambdaLayerVersion{}, errors.New("resource not found")
}
