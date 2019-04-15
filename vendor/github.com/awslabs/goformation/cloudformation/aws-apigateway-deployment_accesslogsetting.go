package cloudformation

import (
	"encoding/json"
)

// AWSApiGatewayDeployment_AccessLogSetting AWS CloudFormation Resource (AWS::ApiGateway::Deployment.AccessLogSetting)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-deployment-accesslogsetting.html
type AWSApiGatewayDeployment_AccessLogSetting struct {

	// DestinationArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-deployment-accesslogsetting.html#cfn-apigateway-deployment-accesslogsetting-destinationarn
	DestinationArn *Value `json:"DestinationArn,omitempty"`

	// Format AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-deployment-accesslogsetting.html#cfn-apigateway-deployment-accesslogsetting-format
	Format *Value `json:"Format,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayDeployment_AccessLogSetting) AWSCloudFormationType() string {
	return "AWS::ApiGateway::Deployment.AccessLogSetting"
}

func (r *AWSApiGatewayDeployment_AccessLogSetting) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
