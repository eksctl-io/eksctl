package cloudformation

import (
	"encoding/json"
)

// AWSApiGatewayV2RouteResponse_ParameterConstraints AWS CloudFormation Resource (AWS::ApiGatewayV2::RouteResponse.ParameterConstraints)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-routeresponse-parameterconstraints.html
type AWSApiGatewayV2RouteResponse_ParameterConstraints struct {

	// Required AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-routeresponse-parameterconstraints.html#cfn-apigatewayv2-routeresponse-parameterconstraints-required
	Required *Value `json:"Required,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayV2RouteResponse_ParameterConstraints) AWSCloudFormationType() string {
	return "AWS::ApiGatewayV2::RouteResponse.ParameterConstraints"
}

func (r *AWSApiGatewayV2RouteResponse_ParameterConstraints) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
