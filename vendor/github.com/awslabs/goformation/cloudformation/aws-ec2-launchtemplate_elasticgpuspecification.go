package cloudformation

import (
	"encoding/json"
)

// AWSEC2LaunchTemplate_ElasticGpuSpecification AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.ElasticGpuSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-elasticgpuspecification.html
type AWSEC2LaunchTemplate_ElasticGpuSpecification struct {

	// Type AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-elasticgpuspecification.html#cfn-ec2-launchtemplate-elasticgpuspecification-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_ElasticGpuSpecification) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.ElasticGpuSpecification"
}

func (r *AWSEC2LaunchTemplate_ElasticGpuSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
