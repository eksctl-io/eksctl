package cloudformation

import (
	"encoding/json"
)

// AWSApiGatewayDeployment_DeploymentCanarySettings AWS CloudFormation Resource (AWS::ApiGateway::Deployment.DeploymentCanarySettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-deployment-deploymentcanarysettings.html
type AWSApiGatewayDeployment_DeploymentCanarySettings struct {

	// PercentTraffic AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-deployment-deploymentcanarysettings.html#cfn-apigateway-deployment-deploymentcanarysettings-percenttraffic
	PercentTraffic *Value `json:"PercentTraffic,omitempty"`

	// StageVariableOverrides AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-deployment-deploymentcanarysettings.html#cfn-apigateway-deployment-deploymentcanarysettings-stagevariableoverrides
	StageVariableOverrides map[string]*Value `json:"StageVariableOverrides,omitempty"`

	// UseStageCache AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-deployment-deploymentcanarysettings.html#cfn-apigateway-deployment-deploymentcanarysettings-usestagecache
	UseStageCache *Value `json:"UseStageCache,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayDeployment_DeploymentCanarySettings) AWSCloudFormationType() string {
	return "AWS::ApiGateway::Deployment.DeploymentCanarySettings"
}

func (r *AWSApiGatewayDeployment_DeploymentCanarySettings) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
