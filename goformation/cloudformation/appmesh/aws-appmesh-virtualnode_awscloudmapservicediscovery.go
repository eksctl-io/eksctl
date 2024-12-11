package appmesh

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// VirtualNode_AwsCloudMapServiceDiscovery AWS CloudFormation Resource (AWS::AppMesh::VirtualNode.AwsCloudMapServiceDiscovery)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-awscloudmapservicediscovery.html
type VirtualNode_AwsCloudMapServiceDiscovery struct {

	// Attributes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-awscloudmapservicediscovery.html#cfn-appmesh-virtualnode-awscloudmapservicediscovery-attributes
	Attributes []VirtualNode_AwsCloudMapInstanceAttribute `json:"Attributes,omitempty"`

	// NamespaceName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-awscloudmapservicediscovery.html#cfn-appmesh-virtualnode-awscloudmapservicediscovery-namespacename
	NamespaceName *types.Value `json:"NamespaceName,omitempty"`

	// ServiceName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-appmesh-virtualnode-awscloudmapservicediscovery.html#cfn-appmesh-virtualnode-awscloudmapservicediscovery-servicename
	ServiceName *types.Value `json:"ServiceName,omitempty"`

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
func (r *VirtualNode_AwsCloudMapServiceDiscovery) AWSCloudFormationType() string {
	return "AWS::AppMesh::VirtualNode.AwsCloudMapServiceDiscovery"
}
