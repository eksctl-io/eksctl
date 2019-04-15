package cloudformation

import (
	"encoding/json"
)

// AWSApiGatewayV2Route_ParameterConstraints AWS CloudFormation Resource (AWS::ApiGatewayV2::Route.ParameterConstraints)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-route-parameterconstraints.html
type AWSApiGatewayV2Route_ParameterConstraints struct {

	// Required AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-route-parameterconstraints.html#cfn-apigatewayv2-route-parameterconstraints-required
	Required *Value `json:"Required,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayV2Route_ParameterConstraints) AWSCloudFormationType() string {
	return "AWS::ApiGatewayV2::Route.ParameterConstraints"
}

func (r *AWSApiGatewayV2Route_ParameterConstraints) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
