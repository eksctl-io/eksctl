package cloudformation

import (
	"encoding/json"
)

// AWSEC2Instance_ElasticGpuSpecification AWS CloudFormation Resource (AWS::EC2::Instance.ElasticGpuSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-instance-elasticgpuspecification.html
type AWSEC2Instance_ElasticGpuSpecification struct {

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-instance-elasticgpuspecification.html#cfn-ec2-instance-elasticgpuspecification-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2Instance_ElasticGpuSpecification) AWSCloudFormationType() string {
	return "AWS::EC2::Instance.ElasticGpuSpecification"
}

func (r *AWSEC2Instance_ElasticGpuSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
