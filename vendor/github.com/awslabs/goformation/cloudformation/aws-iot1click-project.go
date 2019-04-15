package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSIoT1ClickProject AWS CloudFormation Resource (AWS::IoT1Click::Project)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iot1click-project.html
type AWSIoT1ClickProject struct {

	// Description AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iot1click-project.html#cfn-iot1click-project-description
	Description *Value `json:"Description,omitempty"`

	// PlacementTemplate AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iot1click-project.html#cfn-iot1click-project-placementtemplate
	PlacementTemplate *AWSIoT1ClickProject_PlacementTemplate `json:"PlacementTemplate,omitempty"`

	// ProjectName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iot1click-project.html#cfn-iot1click-project-projectname
	ProjectName *Value `json:"ProjectName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoT1ClickProject) AWSCloudFormationType() string {
	return "AWS::IoT1Click::Project"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSIoT1ClickProject) MarshalJSON() ([]byte, error) {
	type Properties AWSIoT1ClickProject
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
func (r *AWSIoT1ClickProject) UnmarshalJSON(b []byte) error {
	type Properties AWSIoT1ClickProject
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
		*r = AWSIoT1ClickProject(*res.Properties)
	}

	return nil
}

// GetAllAWSIoT1ClickProjectResources retrieves all AWSIoT1ClickProject items from an AWS CloudFormation template
func (t *Template) GetAllAWSIoT1ClickProjectResources() map[string]AWSIoT1ClickProject {
	results := map[string]AWSIoT1ClickProject{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSIoT1ClickProject:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::IoT1Click::Project" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSIoT1ClickProject{}
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

// GetAWSIoT1ClickProjectWithName retrieves all AWSIoT1ClickProject items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSIoT1ClickProjectWithName(name string) (AWSIoT1ClickProject, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSIoT1ClickProject:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::IoT1Click::Project" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSIoT1ClickProject{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSIoT1ClickProject{}, errors.New("resource not found")
}
