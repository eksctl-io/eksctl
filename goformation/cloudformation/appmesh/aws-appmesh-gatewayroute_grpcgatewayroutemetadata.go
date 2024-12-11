package appmesh

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// GatewayRoute_GrpcGatewayRouteMetadata AWS CloudFormation Resource (AWS::AppMesh::GatewayRoute.GrpcGatewayRouteMetadata)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-grpcgatewayroutemetadata.html
type GatewayRoute_GrpcGatewayRouteMetadata struct {

	// Invert AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-grpcgatewayroutemetadata.html#cfn-appmesh-gatewayroute-grpcgatewayroutemetadata-invert
	Invert *types.Value `json:"Invert,omitempty"`

	// Match AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-grpcgatewayroutemetadata.html#cfn-appmesh-gatewayroute-grpcgatewayroutemetadata-match
	Match *GatewayRoute_GatewayRouteMetadataMatch `json:"Match,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-grpcgatewayroutemetadata.html#cfn-appmesh-gatewayroute-grpcgatewayroutemetadata-name
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
func (r *GatewayRoute_GrpcGatewayRouteMetadata) AWSCloudFormationType() string {
	return "AWS::AppMesh::GatewayRoute.GrpcGatewayRouteMetadata"
}
