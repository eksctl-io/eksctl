package appmesh

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// GatewayRoute_GatewayRouteRangeMatch AWS CloudFormation Resource (AWS::AppMesh::GatewayRoute.GatewayRouteRangeMatch)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-gatewayrouterangematch.html
type GatewayRoute_GatewayRouteRangeMatch struct {

	// End AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-gatewayrouterangematch.html#cfn-appmesh-gatewayroute-gatewayrouterangematch-end
	End *types.Value `json:"End"`

	// Start AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-gatewayroute-gatewayrouterangematch.html#cfn-appmesh-gatewayroute-gatewayrouterangematch-start
	Start *types.Value `json:"Start"`

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
func (r *GatewayRoute_GatewayRouteRangeMatch) AWSCloudFormationType() string {
	return "AWS::AppMesh::GatewayRoute.GatewayRouteRangeMatch"
}
