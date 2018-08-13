package cloudformation

import (
	"encoding/json"
)

// AWSEC2NetworkAclEntry_PortRange AWS CloudFormation Resource (AWS::EC2::NetworkAclEntry.PortRange)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkaclentry-portrange.html
type AWSEC2NetworkAclEntry_PortRange struct {

	// From AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkaclentry-portrange.html#cfn-ec2-networkaclentry-portrange-from
	From *Value `json:"From,omitempty"`

	// To AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-networkaclentry-portrange.html#cfn-ec2-networkaclentry-portrange-to
	To *Value `json:"To,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2NetworkAclEntry_PortRange) AWSCloudFormationType() string {
	return "AWS::EC2::NetworkAclEntry.PortRange"
}

func (r *AWSEC2NetworkAclEntry_PortRange) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
