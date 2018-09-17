package cloudformation

import (
	"encoding/json"
)

// AWSApiGatewayApiKey_StageKey AWS CloudFormation Resource (AWS::ApiGateway::ApiKey.StageKey)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-apikey-stagekey.html
type AWSApiGatewayApiKey_StageKey struct {

	// RestApiId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-apikey-stagekey.html#cfn-apigateway-apikey-stagekey-restapiid
	RestApiId *Value `json:"RestApiId,omitempty"`

	// StageName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-apikey-stagekey.html#cfn-apigateway-apikey-stagekey-stagename
	StageName *Value `json:"StageName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayApiKey_StageKey) AWSCloudFormationType() string {
	return "AWS::ApiGateway::ApiKey.StageKey"
}

func (r *AWSApiGatewayApiKey_StageKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
