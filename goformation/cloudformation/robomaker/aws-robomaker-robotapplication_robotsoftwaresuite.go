package robomaker

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// RobotApplication_RobotSoftwareSuite AWS CloudFormation Resource (AWS::RoboMaker::RobotApplication.RobotSoftwareSuite)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-robotapplication-robotsoftwaresuite.html
type RobotApplication_RobotSoftwareSuite struct {

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-robotapplication-robotsoftwaresuite.html#cfn-robomaker-robotapplication-robotsoftwaresuite-name
	Name *types.Value `json:"Name,omitempty"`

	// Version AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-robotapplication-robotsoftwaresuite.html#cfn-robomaker-robotapplication-robotsoftwaresuite-version
	Version *types.Value `json:"Version,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *RobotApplication_RobotSoftwareSuite) AWSCloudFormationType() string {
	return "AWS::RoboMaker::RobotApplication.RobotSoftwareSuite"
}
