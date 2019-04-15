package cloudformation

import (
	"encoding/json"
)

// AWSRoboMakerRobotApplication_RobotSoftwareSuite AWS CloudFormation Resource (AWS::RoboMaker::RobotApplication.RobotSoftwareSuite)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-robotapplication-robotsoftwaresuite.html
type AWSRoboMakerRobotApplication_RobotSoftwareSuite struct {

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-robotapplication-robotsoftwaresuite.html#cfn-robomaker-robotapplication-robotsoftwaresuite-name
	Name *Value `json:"Name,omitempty"`

	// Version AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-robotapplication-robotsoftwaresuite.html#cfn-robomaker-robotapplication-robotsoftwaresuite-version
	Version *Value `json:"Version,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoboMakerRobotApplication_RobotSoftwareSuite) AWSCloudFormationType() string {
	return "AWS::RoboMaker::RobotApplication.RobotSoftwareSuite"
}

func (r *AWSRoboMakerRobotApplication_RobotSoftwareSuite) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
