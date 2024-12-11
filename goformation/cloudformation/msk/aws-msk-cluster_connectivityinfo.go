package msk

import (
	"goformation/v4/cloudformation/policies"
)

// Cluster_ConnectivityInfo AWS CloudFormation Resource (AWS::MSK::Cluster.ConnectivityInfo)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-msk-cluster-connectivityinfo.html
type Cluster_ConnectivityInfo struct {

	// PublicAccess AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-msk-cluster-connectivityinfo.html#cfn-msk-cluster-connectivityinfo-publicaccess
	PublicAccess *Cluster_PublicAccess `json:"PublicAccess,omitempty"`

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
func (r *Cluster_ConnectivityInfo) AWSCloudFormationType() string {
	return "AWS::MSK::Cluster.ConnectivityInfo"
}
