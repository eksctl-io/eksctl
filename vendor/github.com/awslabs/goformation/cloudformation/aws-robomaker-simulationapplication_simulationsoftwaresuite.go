package cloudformation

import (
	"encoding/json"
)

// AWSRoboMakerSimulationApplication_SimulationSoftwareSuite AWS CloudFormation Resource (AWS::RoboMaker::SimulationApplication.SimulationSoftwareSuite)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-simulationapplication-simulationsoftwaresuite.html
type AWSRoboMakerSimulationApplication_SimulationSoftwareSuite struct {

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-simulationapplication-simulationsoftwaresuite.html#cfn-robomaker-simulationapplication-simulationsoftwaresuite-name
	Name *Value `json:"Name,omitempty"`

	// Version AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-simulationapplication-simulationsoftwaresuite.html#cfn-robomaker-simulationapplication-simulationsoftwaresuite-version
	Version *Value `json:"Version,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoboMakerSimulationApplication_SimulationSoftwareSuite) AWSCloudFormationType() string {
	return "AWS::RoboMaker::SimulationApplication.SimulationSoftwareSuite"
}

func (r *AWSRoboMakerSimulationApplication_SimulationSoftwareSuite) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
