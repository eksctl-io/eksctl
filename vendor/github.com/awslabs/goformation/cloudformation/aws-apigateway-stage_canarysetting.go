package cloudformation

import (
	"encoding/json"
)

// AWSApiGatewayStage_CanarySetting AWS CloudFormation Resource (AWS::ApiGateway::Stage.CanarySetting)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-stage-canarysetting.html
type AWSApiGatewayStage_CanarySetting struct {

	// DeploymentId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-stage-canarysetting.html#cfn-apigateway-stage-canarysetting-deploymentid
	DeploymentId *Value `json:"DeploymentId,omitempty"`

	// PercentTraffic AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-stage-canarysetting.html#cfn-apigateway-stage-canarysetting-percenttraffic
	PercentTraffic *Value `json:"PercentTraffic,omitempty"`

	// StageVariableOverrides AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-stage-canarysetting.html#cfn-apigateway-stage-canarysetting-stagevariableoverrides
	StageVariableOverrides map[string]*Value `json:"StageVariableOverrides,omitempty"`

	// UseStageCache AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-stage-canarysetting.html#cfn-apigateway-stage-canarysetting-usestagecache
	UseStageCache *Value `json:"UseStageCache,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayStage_CanarySetting) AWSCloudFormationType() string {
	return "AWS::ApiGateway::Stage.CanarySetting"
}

func (r *AWSApiGatewayStage_CanarySetting) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
