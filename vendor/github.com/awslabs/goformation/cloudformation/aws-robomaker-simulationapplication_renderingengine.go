package cloudformation

import (
	"encoding/json"
)

// AWSRoboMakerSimulationApplication_RenderingEngine AWS CloudFormation Resource (AWS::RoboMaker::SimulationApplication.RenderingEngine)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-simulationapplication-renderingengine.html
type AWSRoboMakerSimulationApplication_RenderingEngine struct {

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-simulationapplication-renderingengine.html#cfn-robomaker-simulationapplication-renderingengine-name
	Name *Value `json:"Name,omitempty"`

	// Version AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-simulationapplication-renderingengine.html#cfn-robomaker-simulationapplication-renderingengine-version
	Version *Value `json:"Version,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoboMakerSimulationApplication_RenderingEngine) AWSCloudFormationType() string {
	return "AWS::RoboMaker::SimulationApplication.RenderingEngine"
}

func (r *AWSRoboMakerSimulationApplication_RenderingEngine) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
