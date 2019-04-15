package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSRoboMakerRobotApplication AWS CloudFormation Resource (AWS::RoboMaker::RobotApplication)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-robomaker-robotapplication.html
type AWSRoboMakerRobotApplication struct {

	// CurrentRevisionId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-robomaker-robotapplication.html#cfn-robomaker-robotapplication-currentrevisionid
	CurrentRevisionId *Value `json:"CurrentRevisionId,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-robomaker-robotapplication.html#cfn-robomaker-robotapplication-name
	Name *Value `json:"Name,omitempty"`

	// RobotSoftwareSuite AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-robomaker-robotapplication.html#cfn-robomaker-robotapplication-robotsoftwaresuite
	RobotSoftwareSuite *AWSRoboMakerRobotApplication_RobotSoftwareSuite `json:"RobotSoftwareSuite,omitempty"`

	// Sources AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-robomaker-robotapplication.html#cfn-robomaker-robotapplication-sources
	Sources []AWSRoboMakerRobotApplication_SourceConfig `json:"Sources,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-robomaker-robotapplication.html#cfn-robomaker-robotapplication-tags
	Tags interface{} `json:"Tags,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoboMakerRobotApplication) AWSCloudFormationType() string {
	return "AWS::RoboMaker::RobotApplication"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSRoboMakerRobotApplication) MarshalJSON() ([]byte, error) {
	type Properties AWSRoboMakerRobotApplication
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
func (r *AWSRoboMakerRobotApplication) UnmarshalJSON(b []byte) error {
	type Properties AWSRoboMakerRobotApplication
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
		*r = AWSRoboMakerRobotApplication(*res.Properties)
	}

	return nil
}

// GetAllAWSRoboMakerRobotApplicationResources retrieves all AWSRoboMakerRobotApplication items from an AWS CloudFormation template
func (t *Template) GetAllAWSRoboMakerRobotApplicationResources() map[string]AWSRoboMakerRobotApplication {
	results := map[string]AWSRoboMakerRobotApplication{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSRoboMakerRobotApplication:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::RoboMaker::RobotApplication" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSRoboMakerRobotApplication{}
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

// GetAWSRoboMakerRobotApplicationWithName retrieves all AWSRoboMakerRobotApplication items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSRoboMakerRobotApplicationWithName(name string) (AWSRoboMakerRobotApplication, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSRoboMakerRobotApplication:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::RoboMaker::RobotApplication" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSRoboMakerRobotApplication{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSRoboMakerRobotApplication{}, errors.New("resource not found")
}
