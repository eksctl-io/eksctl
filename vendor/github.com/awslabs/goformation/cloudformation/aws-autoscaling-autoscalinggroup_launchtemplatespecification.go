package cloudformation

import (
	"encoding/json"
)

// AWSAutoScalingAutoScalingGroup_LaunchTemplateSpecification AWS CloudFormation Resource (AWS::AutoScaling::AutoScalingGroup.LaunchTemplateSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-autoscalinggroup-launchtemplatespecification.html
type AWSAutoScalingAutoScalingGroup_LaunchTemplateSpecification struct {

	// LaunchTemplateId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-autoscalinggroup-launchtemplatespecification.html#cfn-autoscaling-autoscalinggroup-launchtemplatespecification-launchtemplateid
	LaunchTemplateId *Value `json:"LaunchTemplateId,omitempty"`

	// LaunchTemplateName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-autoscalinggroup-launchtemplatespecification.html#cfn-autoscaling-autoscalinggroup-launchtemplatespecification-launchtemplatename
	LaunchTemplateName *Value `json:"LaunchTemplateName,omitempty"`

	// Version AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-autoscaling-autoscalinggroup-launchtemplatespecification.html#cfn-autoscaling-autoscalinggroup-launchtemplatespecification-version
	Version *Value `json:"Version,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAutoScalingAutoScalingGroup_LaunchTemplateSpecification) AWSCloudFormationType() string {
	return "AWS::AutoScaling::AutoScalingGroup.LaunchTemplateSpecification"
}

func (r *AWSAutoScalingAutoScalingGroup_LaunchTemplateSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
