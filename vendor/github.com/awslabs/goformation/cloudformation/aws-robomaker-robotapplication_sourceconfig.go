package cloudformation

import (
	"encoding/json"
)

// AWSRoboMakerRobotApplication_SourceConfig AWS CloudFormation Resource (AWS::RoboMaker::RobotApplication.SourceConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-robotapplication-sourceconfig.html
type AWSRoboMakerRobotApplication_SourceConfig struct {

	// Architecture AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-robotapplication-sourceconfig.html#cfn-robomaker-robotapplication-sourceconfig-architecture
	Architecture *Value `json:"Architecture,omitempty"`

	// S3Bucket AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-robotapplication-sourceconfig.html#cfn-robomaker-robotapplication-sourceconfig-s3bucket
	S3Bucket *Value `json:"S3Bucket,omitempty"`

	// S3Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-robotapplication-sourceconfig.html#cfn-robomaker-robotapplication-sourceconfig-s3key
	S3Key *Value `json:"S3Key,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoboMakerRobotApplication_SourceConfig) AWSCloudFormationType() string {
	return "AWS::RoboMaker::RobotApplication.SourceConfig"
}

func (r *AWSRoboMakerRobotApplication_SourceConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
