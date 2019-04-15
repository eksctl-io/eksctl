package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshRoute_TcpRouteAction AWS CloudFormation Resource (AWS::AppMesh::Route.TcpRouteAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-tcprouteaction.html
type AWSAppMeshRoute_TcpRouteAction struct {

	// WeightedTargets AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-tcprouteaction.html#cfn-appmesh-route-tcprouteaction-weightedtargets
	WeightedTargets []AWSAppMeshRoute_WeightedTarget `json:"WeightedTargets,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshRoute_TcpRouteAction) AWSCloudFormationType() string {
	return "AWS::AppMesh::Route.TcpRouteAction"
}

func (r *AWSAppMeshRoute_TcpRouteAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
