package cloudformation

import (
	"encoding/json"
)

// AWSEC2Instance_ElasticInferenceAccelerator AWS CloudFormation Resource (AWS::EC2::Instance.ElasticInferenceAccelerator)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-instance-elasticinferenceaccelerator.html
type AWSEC2Instance_ElasticInferenceAccelerator struct {

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-instance-elasticinferenceaccelerator.html#cfn-ec2-instance-elasticinferenceaccelerator-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2Instance_ElasticInferenceAccelerator) AWSCloudFormationType() string {
	return "AWS::EC2::Instance.ElasticInferenceAccelerator"
}

func (r *AWSEC2Instance_ElasticInferenceAccelerator) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
