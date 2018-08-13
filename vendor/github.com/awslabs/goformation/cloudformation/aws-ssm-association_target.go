package cloudformation

import (
	"encoding/json"
)

// AWSSSMAssociation_Target AWS CloudFormation Resource (AWS::SSM::Association.Target)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-association-target.html
type AWSSSMAssociation_Target struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-association-target.html#cfn-ssm-association-target-key
	Key *Value `json:"Key,omitempty"`

	// Values AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-association-target.html#cfn-ssm-association-target-values
	Values []*Value `json:"Values,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSSMAssociation_Target) AWSCloudFormationType() string {
	return "AWS::SSM::Association.Target"
}

func (r *AWSSSMAssociation_Target) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
