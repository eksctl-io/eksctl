package cloudformation

import (
	"encoding/json"
)

// AWSEC2LaunchTemplate_LaunchTemplateElasticInferenceAccelerator AWS CloudFormation Resource (AWS::EC2::LaunchTemplate.LaunchTemplateElasticInferenceAccelerator)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplateelasticinferenceaccelerator.html
type AWSEC2LaunchTemplate_LaunchTemplateElasticInferenceAccelerator struct {

	// Type AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-launchtemplate-launchtemplateelasticinferenceaccelerator.html#cfn-ec2-launchtemplate-launchtemplateelasticinferenceaccelerator-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2LaunchTemplate_LaunchTemplateElasticInferenceAccelerator) AWSCloudFormationType() string {
	return "AWS::EC2::LaunchTemplate.LaunchTemplateElasticInferenceAccelerator"
}

func (r *AWSEC2LaunchTemplate_LaunchTemplateElasticInferenceAccelerator) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
