package cloudformation

import (
	"encoding/json"
)

// AWSApiGatewayStage_AccessLogSetting AWS CloudFormation Resource (AWS::ApiGateway::Stage.AccessLogSetting)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-stage-accesslogsetting.html
type AWSApiGatewayStage_AccessLogSetting struct {

	// DestinationArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-stage-accesslogsetting.html#cfn-apigateway-stage-accesslogsetting-destinationarn
	DestinationArn *Value `json:"DestinationArn,omitempty"`

	// Format AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-stage-accesslogsetting.html#cfn-apigateway-stage-accesslogsetting-format
	Format *Value `json:"Format,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayStage_AccessLogSetting) AWSCloudFormationType() string {
	return "AWS::ApiGateway::Stage.AccessLogSetting"
}

func (r *AWSApiGatewayStage_AccessLogSetting) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
