package apigateway

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// RestApi_EndpointConfiguration AWS CloudFormation Resource (AWS::ApiGateway::RestApi.EndpointConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-restapi-endpointconfiguration.html
type RestApi_EndpointConfiguration struct {

	// Types AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-restapi-endpointconfiguration.html#cfn-apigateway-restapi-endpointconfiguration-types
	Types *types.Value `json:"Types,omitempty"`

	// VpcEndpointIds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-apigateway-restapi-endpointconfiguration.html#cfn-apigateway-restapi-endpointconfiguration-vpcendpointids
	VpcEndpointIds *types.Value `json:"VpcEndpointIds,omitempty"`

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
func (r *RestApi_EndpointConfiguration) AWSCloudFormationType() string {
	return "AWS::ApiGateway::RestApi.EndpointConfiguration"
}
