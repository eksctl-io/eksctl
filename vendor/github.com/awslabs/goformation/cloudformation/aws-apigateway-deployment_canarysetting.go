package cloudformation

import (
	"encoding/json"
)

// AWSApiGatewayDeployment_CanarySetting AWS CloudFormation Resource (AWS::ApiGateway::Deployment.CanarySetting)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-deployment-canarysetting.html
type AWSApiGatewayDeployment_CanarySetting struct {

	// PercentTraffic AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-deployment-canarysetting.html#cfn-apigateway-deployment-canarysetting-percenttraffic
	PercentTraffic *Value `json:"PercentTraffic,omitempty"`

	// StageVariableOverrides AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-deployment-canarysetting.html#cfn-apigateway-deployment-canarysetting-stagevariableoverrides
	StageVariableOverrides map[string]*Value `json:"StageVariableOverrides,omitempty"`

	// UseStageCache AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-deployment-canarysetting.html#cfn-apigateway-deployment-canarysetting-usestagecache
	UseStageCache *Value `json:"UseStageCache,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayDeployment_CanarySetting) AWSCloudFormationType() string {
	return "AWS::ApiGateway::Deployment.CanarySetting"
}

func (r *AWSApiGatewayDeployment_CanarySetting) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
