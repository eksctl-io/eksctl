package appmesh

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// GatewayRoute_HttpGatewayRouteHeader AWS CloudFormation Resource (AWS::AppMesh::GatewayRoute.HttpGatewayRouteHeader)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-httpgatewayrouteheader.html
type GatewayRoute_HttpGatewayRouteHeader struct {

	// Invert AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-httpgatewayrouteheader.html#cfn-appmesh-gatewayroute-httpgatewayrouteheader-invert
	Invert *types.Value `json:"Invert,omitempty"`

	// Match AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-httpgatewayrouteheader.html#cfn-appmesh-gatewayroute-httpgatewayrouteheader-match
	Match *GatewayRoute_HttpGatewayRouteHeaderMatch `json:"Match,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-httpgatewayrouteheader.html#cfn-appmesh-gatewayroute-httpgatewayrouteheader-name
	Name *types.Value `json:"Name,omitempty"`

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
func (r *GatewayRoute_HttpGatewayRouteHeader) AWSCloudFormationType() string {
	return "AWS::AppMesh::GatewayRoute.HttpGatewayRouteHeader"
}
