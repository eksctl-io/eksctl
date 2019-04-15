package cloudformation

import (
	"encoding/json"
)

// AWSEC2EC2Fleet_TagSpecification AWS CloudFormation Resource (AWS::EC2::EC2Fleet.TagSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-tagspecification.html
type AWSEC2EC2Fleet_TagSpecification struct {

	// ResourceType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-tagspecification.html#cfn-ec2-ec2fleet-tagspecification-resourcetype
	ResourceType *Value `json:"ResourceType,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ec2-ec2fleet-tagspecification.html#cfn-ec2-ec2fleet-tagspecification-tags
	Tags []AWSEC2EC2Fleet_TagRequest `json:"Tags,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2EC2Fleet_TagSpecification) AWSCloudFormationType() string {
	return "AWS::EC2::EC2Fleet.TagSpecification"
}

func (r *AWSEC2EC2Fleet_TagSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
