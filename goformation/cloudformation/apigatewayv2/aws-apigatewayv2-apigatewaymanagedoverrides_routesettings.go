package apigatewayv2

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ApiGatewayManagedOverrides_RouteSettings AWS CloudFormation Resource (AWS::ApiGatewayV2::ApiGatewayManagedOverrides.RouteSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-apigatewaymanagedoverrides-routesettings.html
type ApiGatewayManagedOverrides_RouteSettings struct {

	// DataTraceEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-apigatewaymanagedoverrides-routesettings.html#cfn-apigatewayv2-apigatewaymanagedoverrides-routesettings-datatraceenabled
	DataTraceEnabled *types.Value `json:"DataTraceEnabled,omitempty"`

	// DetailedMetricsEnabled AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-apigatewaymanagedoverrides-routesettings.html#cfn-apigatewayv2-apigatewaymanagedoverrides-routesettings-detailedmetricsenabled
	DetailedMetricsEnabled *types.Value `json:"DetailedMetricsEnabled,omitempty"`

	// LoggingLevel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-apigatewaymanagedoverrides-routesettings.html#cfn-apigatewayv2-apigatewaymanagedoverrides-routesettings-logginglevel
	LoggingLevel *types.Value `json:"LoggingLevel,omitempty"`

	// ThrottlingBurstLimit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-apigatewaymanagedoverrides-routesettings.html#cfn-apigatewayv2-apigatewaymanagedoverrides-routesettings-throttlingburstlimit
	ThrottlingBurstLimit *types.Value `json:"ThrottlingBurstLimit,omitempty"`

	// ThrottlingRateLimit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigatewayv2-apigatewaymanagedoverrides-routesettings.html#cfn-apigatewayv2-apigatewaymanagedoverrides-routesettings-throttlingratelimit
	ThrottlingRateLimit *types.Value `json:"ThrottlingRateLimit,omitempty"`

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
func (r *ApiGatewayManagedOverrides_RouteSettings) AWSCloudFormationType() string {
	return "AWS::ApiGatewayV2::ApiGatewayManagedOverrides.RouteSettings"
}
