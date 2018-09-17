package cloudformation

import (
	"encoding/json"
)

// AWSCognitoUserPool_NumberAttributeConstraints AWS CloudFormation Resource (AWS::Cognito::UserPool.NumberAttributeConstraints)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-numberattributeconstraints.html
type AWSCognitoUserPool_NumberAttributeConstraints struct {

	// MaxValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-numberattributeconstraints.html#cfn-cognito-userpool-numberattributeconstraints-maxvalue
	MaxValue *Value `json:"MaxValue,omitempty"`

	// MinValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cognito-userpool-numberattributeconstraints.html#cfn-cognito-userpool-numberattributeconstraints-minvalue
	MinValue *Value `json:"MinValue,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCognitoUserPool_NumberAttributeConstraints) AWSCloudFormationType() string {
	return "AWS::Cognito::UserPool.NumberAttributeConstraints"
}

func (r *AWSCognitoUserPool_NumberAttributeConstraints) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
