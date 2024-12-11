package elasticache

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// GlobalReplicationGroup_GlobalReplicationGroupMember AWS CloudFormation Resource (AWS::ElastiCache::GlobalReplicationGroup.GlobalReplicationGroupMember)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-globalreplicationgroup-globalreplicationgroupmember.html
type GlobalReplicationGroup_GlobalReplicationGroupMember struct {

	// ReplicationGroupId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-globalreplicationgroup-globalreplicationgroupmember.html#cfn-elasticache-globalreplicationgroup-globalreplicationgroupmember-replicationgroupid
	ReplicationGroupId *types.Value `json:"ReplicationGroupId,omitempty"`

	// ReplicationGroupRegion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-globalreplicationgroup-globalreplicationgroupmember.html#cfn-elasticache-globalreplicationgroup-globalreplicationgroupmember-replicationgroupregion
	ReplicationGroupRegion *types.Value `json:"ReplicationGroupRegion,omitempty"`

	// Role AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticache-globalreplicationgroup-globalreplicationgroupmember.html#cfn-elasticache-globalreplicationgroup-globalreplicationgroupmember-role
	Role *types.Value `json:"Role,omitempty"`

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
func (r *GlobalReplicationGroup_GlobalReplicationGroupMember) AWSCloudFormationType() string {
	return "AWS::ElastiCache::GlobalReplicationGroup.GlobalReplicationGroupMember"
}
