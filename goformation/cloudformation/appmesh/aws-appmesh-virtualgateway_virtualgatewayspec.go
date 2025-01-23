package appmesh

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// VirtualGateway_VirtualGatewaySpec AWS CloudFormation Resource (AWS::AppMesh::VirtualGateway.VirtualGatewaySpec)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewayspec.html
type VirtualGateway_VirtualGatewaySpec struct {

	// BackendDefaults AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewayspec.html#cfn-appmesh-virtualgateway-virtualgatewayspec-backenddefaults
	BackendDefaults *VirtualGateway_VirtualGatewayBackendDefaults `json:"BackendDefaults,omitempty"`

	// Listeners AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewayspec.html#cfn-appmesh-virtualgateway-virtualgatewayspec-listeners
	Listeners []VirtualGateway_VirtualGatewayListener `json:"Listeners,omitempty"`

	// Logging AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualgateway-virtualgatewayspec.html#cfn-appmesh-virtualgateway-virtualgatewayspec-logging
	Logging *VirtualGateway_VirtualGatewayLogging `json:"Logging,omitempty"`

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
func (r *VirtualGateway_VirtualGatewaySpec) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualGateway.VirtualGatewaySpec"
}
