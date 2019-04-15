package cloudformation

import (
	"encoding/json"
)

// AWSApiGatewayV2Stage_RouteSettings AWS CloudFormation Resource (AWS::ApiGatewayV2::Stage.RouteSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-routesettings.html
type AWSApiGatewayV2Stage_RouteSettings struct {

	// DataTraceEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-routesettings.html#cfn-apigatewayv2-stage-routesettings-datatraceenabled
	DataTraceEnabled *Value `json:"DataTraceEnabled,omitempty"`

	// DetailedMetricsEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-routesettings.html#cfn-apigatewayv2-stage-routesettings-detailedmetricsenabled
	DetailedMetricsEnabled *Value `json:"DetailedMetricsEnabled,omitempty"`

	// LoggingLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-routesettings.html#cfn-apigatewayv2-stage-routesettings-logginglevel
	LoggingLevel *Value `json:"LoggingLevel,omitempty"`

	// ThrottlingBurstLimit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-routesettings.html#cfn-apigatewayv2-stage-routesettings-throttlingburstlimit
	ThrottlingBurstLimit *Value `json:"ThrottlingBurstLimit,omitempty"`

	// ThrottlingRateLimit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-stage-routesettings.html#cfn-apigatewayv2-stage-routesettings-throttlingratelimit
	ThrottlingRateLimit *Value `json:"ThrottlingRateLimit,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayV2Stage_RouteSettings) AWSCloudFormationType() string {
	return "AWS::ApiGatewayV2::Stage.RouteSettings"
}

func (r *AWSApiGatewayV2Stage_RouteSettings) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
