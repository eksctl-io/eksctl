package cloudformation

import (
	"encoding/json"
)

// AWSApiGatewayMethod_IntegrationResponse AWS CloudFormation Resource (AWS::ApiGateway::Method.IntegrationResponse)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apitgateway-method-integration-integrationresponse.html
type AWSApiGatewayMethod_IntegrationResponse struct {

	// ContentHandling AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apitgateway-method-integration-integrationresponse.html#cfn-apigateway-method-integrationresponse-contenthandling
	ContentHandling *Value `json:"ContentHandling,omitempty"`

	// ResponseParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apitgateway-method-integration-integrationresponse.html#cfn-apigateway-method-integration-integrationresponse-responseparameters
	ResponseParameters map[string]*Value `json:"ResponseParameters,omitempty"`

	// ResponseTemplates AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apitgateway-method-integration-integrationresponse.html#cfn-apigateway-method-integration-integrationresponse-responsetemplates
	ResponseTemplates map[string]*Value `json:"ResponseTemplates,omitempty"`

	// SelectionPattern AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apitgateway-method-integration-integrationresponse.html#cfn-apigateway-method-integration-integrationresponse-selectionpattern
	SelectionPattern *Value `json:"SelectionPattern,omitempty"`

	// StatusCode AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apitgateway-method-integration-integrationresponse.html#cfn-apigateway-method-integration-integrationresponse-statuscode
	StatusCode *Value `json:"StatusCode,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayMethod_IntegrationResponse) AWSCloudFormationType() string {
	return "AWS::ApiGateway::Method.IntegrationResponse"
}

func (r *AWSApiGatewayMethod_IntegrationResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
