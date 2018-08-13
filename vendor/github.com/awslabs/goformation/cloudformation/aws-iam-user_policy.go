package cloudformation

import (
	"encoding/json"
)

// AWSIAMUser_Policy AWS CloudFormation Resource (AWS::IAM::User.Policy)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iam-policy.html
type AWSIAMUser_Policy struct {

	// PolicyDocument AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iam-policy.html#cfn-iam-policies-policydocument
	PolicyDocument interface{} `json:"PolicyDocument,omitempty"`

	// PolicyName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iam-policy.html#cfn-iam-policies-policyname
	PolicyName *Value `json:"PolicyName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIAMUser_Policy) AWSCloudFormationType() string {
	return "AWS::IAM::User.Policy"
}

func (r *AWSIAMUser_Policy) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
