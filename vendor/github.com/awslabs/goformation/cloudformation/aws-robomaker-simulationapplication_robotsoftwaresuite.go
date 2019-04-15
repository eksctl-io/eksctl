package cloudformation

import (
	"encoding/json"
)

// AWSRoboMakerSimulationApplication_RobotSoftwareSuite AWS CloudFormation Resource (AWS::RoboMaker::SimulationApplication.RobotSoftwareSuite)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-simulationapplication-robotsoftwaresuite.html
type AWSRoboMakerSimulationApplication_RobotSoftwareSuite struct {

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-simulationapplication-robotsoftwaresuite.html#cfn-robomaker-simulationapplication-robotsoftwaresuite-name
	Name *Value `json:"Name,omitempty"`

	// Version AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-simulationapplication-robotsoftwaresuite.html#cfn-robomaker-simulationapplication-robotsoftwaresuite-version
	Version *Value `json:"Version,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoboMakerSimulationApplication_RobotSoftwareSuite) AWSCloudFormationType() string {
	return "AWS::RoboMaker::SimulationApplication.RobotSoftwareSuite"
}

func (r *AWSRoboMakerSimulationApplication_RobotSoftwareSuite) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
