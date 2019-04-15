package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualNode_Listener AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.Listener)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-listener.html
type AWSAppMeshVirtualNode_Listener struct {

	// HealthCheck AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-listener.html#cfn-appmesh-virtualnode-listener-healthcheck
	HealthCheck *AWSAppMeshVirtualNode_HealthCheck `json:"HealthCheck,omitempty"`

	// PortMapping AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-listener.html#cfn-appmesh-virtualnode-listener-portmapping
	PortMapping *AWSAppMeshVirtualNode_PortMapping `json:"PortMapping,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualNode_Listener) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.Listener"
}

func (r *AWSAppMeshVirtualNode_Listener) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
