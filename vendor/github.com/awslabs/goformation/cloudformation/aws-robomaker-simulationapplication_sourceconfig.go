package cloudformation

import (
	"encoding/json"
)

// AWSRoboMakerSimulationApplication_SourceConfig AWS CloudFormation Resource (AWS::RoboMaker::SimulationApplication.SourceConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-simulationapplication-sourceconfig.html
type AWSRoboMakerSimulationApplication_SourceConfig struct {

	// Architecture AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-simulationapplication-sourceconfig.html#cfn-robomaker-simulationapplication-sourceconfig-architecture
	Architecture *Value `json:"Architecture,omitempty"`

	// S3Bucket AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-simulationapplication-sourceconfig.html#cfn-robomaker-simulationapplication-sourceconfig-s3bucket
	S3Bucket *Value `json:"S3Bucket,omitempty"`

	// S3Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-robomaker-simulationapplication-sourceconfig.html#cfn-robomaker-simulationapplication-sourceconfig-s3key
	S3Key *Value `json:"S3Key,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRoboMakerSimulationApplication_SourceConfig) AWSCloudFormationType() string {
	return "AWS::RoboMaker::SimulationApplication.SourceConfig"
}

func (r *AWSRoboMakerSimulationApplication_SourceConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
