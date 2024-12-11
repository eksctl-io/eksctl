package appmesh

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// GatewayRoute_GatewayRouteSpec AWS CloudFormation Resource (AWS::AppMesh::GatewayRoute.GatewayRouteSpec)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-gatewayroutespec.html
type GatewayRoute_GatewayRouteSpec struct {

	// GrpcRoute AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-gatewayroutespec.html#cfn-appmesh-gatewayroute-gatewayroutespec-grpcroute
	GrpcRoute *GatewayRoute_GrpcGatewayRoute `json:"GrpcRoute,omitempty"`

	// Http2Route AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-gatewayroutespec.html#cfn-appmesh-gatewayroute-gatewayroutespec-http2route
	Http2Route *GatewayRoute_HttpGatewayRoute `json:"Http2Route,omitempty"`

	// HttpRoute AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-gatewayroutespec.html#cfn-appmesh-gatewayroute-gatewayroutespec-httproute
	HttpRoute *GatewayRoute_HttpGatewayRoute `json:"HttpRoute,omitempty"`

	// Priority AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-gatewayroutespec.html#cfn-appmesh-gatewayroute-gatewayroutespec-priority
	Priority *types.Value `json:"Priority,omitempty"`

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
func (r *GatewayRoute_GatewayRouteSpec) AWSCloudFormationType() string {
	return "AWS::AppMesh::GatewayRoute.GatewayRouteSpec"
}
