package apigatewayv2

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Integration_ResponseParameterList AWS CloudFormation Resource (AWS::ApiGatewayV2::Integration.ResponseParameterList)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-integration-responseparameterlist.html
type Integration_ResponseParameterList struct {

	// ResponseParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-integration-responseparameterlist.html#cfn-apigatewayv2-integration-responseparameterlist-responseparameters
	ResponseParameters []Integration_ResponseParameter `json:"ResponseParameters,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Integration_ResponseParameterList) AWSCloudFormationType() string {
	return "AWS::ApiGatewayV2::Integration.ResponseParameterList"
}
