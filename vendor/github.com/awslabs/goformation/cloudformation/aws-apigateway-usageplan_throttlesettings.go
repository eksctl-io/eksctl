package cloudformation

import (
	"encoding/json"
)

// AWSApiGatewayUsagePlan_ThrottleSettings AWS CloudFormation Resource (AWS::ApiGateway::UsagePlan.ThrottleSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-usageplan-throttlesettings.html
type AWSApiGatewayUsagePlan_ThrottleSettings struct {

	// BurstLimit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-usageplan-throttlesettings.html#cfn-apigateway-usageplan-throttlesettings-burstlimit
	BurstLimit *Value `json:"BurstLimit,omitempty"`

	// RateLimit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-usageplan-throttlesettings.html#cfn-apigateway-usageplan-throttlesettings-ratelimit
	RateLimit *Value `json:"RateLimit,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayUsagePlan_ThrottleSettings) AWSCloudFormationType() string {
	return "AWS::ApiGateway::UsagePlan.ThrottleSettings"
}

func (r *AWSApiGatewayUsagePlan_ThrottleSettings) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
