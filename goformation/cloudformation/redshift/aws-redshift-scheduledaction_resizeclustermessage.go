package redshift

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ScheduledAction_ResizeClusterMessage AWS CloudFormation Resource (AWS::Redshift::ScheduledAction.ResizeClusterMessage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-redshift-scheduledaction-resizeclustermessage.html
type ScheduledAction_ResizeClusterMessage struct {

	// Classic AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-redshift-scheduledaction-resizeclustermessage.html#cfn-redshift-scheduledaction-resizeclustermessage-classic
	Classic *types.Value `json:"Classic,omitempty"`

	// ClusterIdentifier AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-redshift-scheduledaction-resizeclustermessage.html#cfn-redshift-scheduledaction-resizeclustermessage-clusteridentifier
	ClusterIdentifier *types.Value `json:"ClusterIdentifier,omitempty"`

	// ClusterType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-redshift-scheduledaction-resizeclustermessage.html#cfn-redshift-scheduledaction-resizeclustermessage-clustertype
	ClusterType *types.Value `json:"ClusterType,omitempty"`

	// NodeType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-redshift-scheduledaction-resizeclustermessage.html#cfn-redshift-scheduledaction-resizeclustermessage-nodetype
	NodeType *types.Value `json:"NodeType,omitempty"`

	// NumberOfNodes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-redshift-scheduledaction-resizeclustermessage.html#cfn-redshift-scheduledaction-resizeclustermessage-numberofnodes
	NumberOfNodes *types.Value `json:"NumberOfNodes,omitempty"`

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
func (r *ScheduledAction_ResizeClusterMessage) AWSCloudFormationType() string {
	return "AWS::Redshift::ScheduledAction.ResizeClusterMessage"
}
