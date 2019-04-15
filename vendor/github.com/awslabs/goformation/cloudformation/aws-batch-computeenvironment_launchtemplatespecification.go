package cloudformation

import (
	"encoding/json"
)

// AWSBatchComputeEnvironment_LaunchTemplateSpecification AWS CloudFormation Resource (AWS::Batch::ComputeEnvironment.LaunchTemplateSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-computeenvironment-launchtemplatespecification.html
type AWSBatchComputeEnvironment_LaunchTemplateSpecification struct {

	// LaunchTemplateId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-computeenvironment-launchtemplatespecification.html#cfn-batch-computeenvironment-launchtemplatespecification-launchtemplateid
	LaunchTemplateId *Value `json:"LaunchTemplateId,omitempty"`

	// LaunchTemplateName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-computeenvironment-launchtemplatespecification.html#cfn-batch-computeenvironment-launchtemplatespecification-launchtemplatename
	LaunchTemplateName *Value `json:"LaunchTemplateName,omitempty"`

	// Version AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-computeenvironment-launchtemplatespecification.html#cfn-batch-computeenvironment-launchtemplatespecification-version
	Version *Value `json:"Version,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSBatchComputeEnvironment_LaunchTemplateSpecification) AWSCloudFormationType() string {
	return "AWS::Batch::ComputeEnvironment.LaunchTemplateSpecification"
}

func (r *AWSBatchComputeEnvironment_LaunchTemplateSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
