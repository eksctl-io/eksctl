package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSRoboMakerSimulationApplication AWS CloudFormation Resource (AWS::RoboMaker::SimulationApplication)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-robomaker-simulationapplication.html
type AWSRoboMakerSimulationApplication struct {

	// CurrentRevisionId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-robomaker-simulationapplication.html#cfn-robomaker-simulationapplication-currentrevisionid
	CurrentRevisionId *Value `json:"CurrentRevisionId,omitempty"`

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-robomaker-simulationapplication.html#cfn-robomaker-simulationapplication-name
	Name *Value `json:"Name,omitempty"`

	// RenderingEngine AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-robomaker-simulationapplication.html#cfn-robomaker-simulationapplication-renderingengine
	RenderingEngine *AWSRoboMakerSimulationApplication_RenderingEngine `json:"RenderingEngine,omitempty"`

	// RobotSoftwareSuite AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-robomaker-simulationapplication.html#cfn-robomaker-simulationapplication-robotsoftwaresuite
	RobotSoftwareSuite *AWSRoboMakerSimulationApplication_RobotSoftwareSuite `json:"RobotSoftwareSuite,omitempty"`

	// SimulationSoftwareSuite AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-robomaker-simulationapplication.html#cfn-robomaker-simulationapplication-simulationsoftwaresuite
	SimulationSoftwareSuite *AWSRoboMakerSimulationApplication_SimulationSoftwareSuite `json:"SimulationSoftwareSuite,omitempty"`

	// Sources AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-robomaker-simulationapplication.html#cfn-robomaker-simulationapplication-sources
	Sources []AWSRoboMakerSimulationApplication_SourceConfig `json:"Sources,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-robomaker-simulationapplication.html#cfn-robomaker-simulationapplication-tags
	Tags interface{} `json:"Tags,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoboMakerSimulationApplication) AWSCloudFormationType() string {
	return "AWS::RoboMaker::SimulationApplication"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSRoboMakerSimulationApplication) MarshalJSON() ([]byte, error) {
	type Properties AWSRoboMakerSimulationApplication
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
func (r *AWSRoboMakerSimulationApplication) UnmarshalJSON(b []byte) error {
	type Properties AWSRoboMakerSimulationApplication
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
		*r = AWSRoboMakerSimulationApplication(*res.Properties)
	}

	return nil
}

// GetAllAWSRoboMakerSimulationApplicationResources retrieves all AWSRoboMakerSimulationApplication items from an AWS CloudFormation template
func (t *Template) GetAllAWSRoboMakerSimulationApplicationResources() map[string]AWSRoboMakerSimulationApplication {
	results := map[string]AWSRoboMakerSimulationApplication{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSRoboMakerSimulationApplication:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::RoboMaker::SimulationApplication" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSRoboMakerSimulationApplication{}
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

// GetAWSRoboMakerSimulationApplicationWithName retrieves all AWSRoboMakerSimulationApplication items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSRoboMakerSimulationApplicationWithName(name string) (AWSRoboMakerSimulationApplication, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSRoboMakerSimulationApplication:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::RoboMaker::SimulationApplication" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSRoboMakerSimulationApplication{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSRoboMakerSimulationApplication{}, errors.New("resource not found")
}
