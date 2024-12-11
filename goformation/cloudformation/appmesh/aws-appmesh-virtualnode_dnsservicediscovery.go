package appmesh

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// VirtualNode_DnsServiceDiscovery AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.DnsServiceDiscovery)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-dnsservicediscovery.html
type VirtualNode_DnsServiceDiscovery struct {

	// Hostname AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-dnsservicediscovery.html#cfn-appmesh-virtualnode-dnsservicediscovery-hostname
	Hostname *types.Value `json:"Hostname,omitempty"`

	// ResponseType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-dnsservicediscovery.html#cfn-appmesh-virtualnode-dnsservicediscovery-responsetype
	ResponseType *types.Value `json:"ResponseType,omitempty"`

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
func (r *VirtualNode_DnsServiceDiscovery) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.DnsServiceDiscovery"
}
