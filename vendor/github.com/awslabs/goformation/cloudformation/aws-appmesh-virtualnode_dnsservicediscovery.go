package cloudformation

import (
	"encoding/json"
)

// AWSAppMeshVirtualNode_DnsServiceDiscovery AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.DnsServiceDiscovery)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-dnsservicediscovery.html
type AWSAppMeshVirtualNode_DnsServiceDiscovery struct {

	// Hostname AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-dnsservicediscovery.html#cfn-appmesh-virtualnode-dnsservicediscovery-hostname
	Hostname *Value `json:"Hostname,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSAppMeshVirtualNode_DnsServiceDiscovery) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.DnsServiceDiscovery"
}

func (r *AWSAppMeshVirtualNode_DnsServiceDiscovery) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
