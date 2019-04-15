package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshRoute_HttpRouteAction AWS CloudFormation Resource (AWS::AppMesh::Route.HttpRouteAction)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-httprouteaction.html
type AWSAppMeshRoute_HttpRouteAction struct {

	// WeightedTargets AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-route-httprouteaction.html#cfn-appmesh-route-httprouteaction-weightedtargets
	WeightedTargets []AWSAppMeshRoute_WeightedTarget `json:"WeightedTargets,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshRoute_HttpRouteAction) AWSCloudFormationType() string {
	return "AWS::AppMesh::Route.HttpRouteAction"
}

func (r *AWSAppMeshRoute_HttpRouteAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
