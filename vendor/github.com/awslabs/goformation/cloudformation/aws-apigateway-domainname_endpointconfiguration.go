package cloudformation

// AWSApiGatewayDomainName_EndpointConfiguration AWS CloudFormation Resource (AWS::ApiGateway::DomainName.EndpointConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-domainname-endpointconfiguration.html
type AWSApiGatewayDomainName_EndpointConfiguration struct {

	// Types AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-domainname-endpointconfiguration.html#cfn-apigateway-domainname-endpointconfiguration-types
	Types []*StringIntrinsic `json:"Types,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSApiGatewayDomainName_EndpointConfiguration) AWSCloudFormationType() string {
	return "AWS::ApiGateway::DomainName.EndpointConfiguration"
}
