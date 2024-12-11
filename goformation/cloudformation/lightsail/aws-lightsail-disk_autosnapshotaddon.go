package lightsail

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Disk_AutoSnapshotAddOn AWS CloudFormation Resource (AWS::Lightsail::Disk.AutoSnapshotAddOn)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-disk-autosnapshotaddon.html
type Disk_AutoSnapshotAddOn struct {

	// SnapshotTimeOfDay AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-disk-autosnapshotaddon.html#cfn-lightsail-disk-autosnapshotaddon-snapshottimeofday
	SnapshotTimeOfDay *types.Value `json:"SnapshotTimeOfDay,omitempty"`

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
func (r *Disk_AutoSnapshotAddOn) AWSCloudFormationType() string {
	return "AWS::Lightsail::Disk.AutoSnapshotAddOn"
}
