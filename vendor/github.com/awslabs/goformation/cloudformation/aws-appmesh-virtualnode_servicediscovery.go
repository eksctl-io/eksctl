package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualNode_ServiceDiscovery AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.ServiceDiscovery)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-servicediscovery.html
type AWSAppMeshVirtualNode_ServiceDiscovery struct {

	// DNS AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-servicediscovery.html#cfn-appmesh-virtualnode-servicediscovery-dns
	DNS *AWSAppMeshVirtualNode_DnsServiceDiscovery `json:"DNS,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualNode_ServiceDiscovery) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.ServiceDiscovery"
}

func (r *AWSAppMeshVirtualNode_ServiceDiscovery) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
