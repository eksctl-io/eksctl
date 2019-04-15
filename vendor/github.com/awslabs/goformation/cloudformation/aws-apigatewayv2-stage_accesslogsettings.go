package cloudformation

import (
	"encoding/json"
)

// AWSApiGatewayV2Stage_AccessLogSettings AWS CloudFormation Resource (AWS::ApiGatewayV2::Stage.AccessLogSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-accesslogsettings.html
type AWSApiGatewayV2Stage_AccessLogSettings struct {

	// DestinationArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-accesslogsettings.html#cfn-apigatewayv2-stage-accesslogsettings-destinationarn
	DestinationArn *Value `json:"DestinationArn,omitempty"`

	// Format AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-accesslogsettings.html#cfn-apigatewayv2-stage-accesslogsettings-format
	Format *Value `json:"Format,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayV2Stage_AccessLogSettings) AWSCloudFormationType() string {
	return "AWS::ApiGatewayV2::Stage.AccessLogSettings"
}

func (r *AWSApiGatewayV2Stage_AccessLogSettings) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
