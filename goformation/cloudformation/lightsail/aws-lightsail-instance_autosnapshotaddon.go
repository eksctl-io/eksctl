package lightsail

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Instance_AutoSnapshotAddOn AWS CloudFormation Resource (AWS::Lightsail::Instance.AutoSnapshotAddOn)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-autosnapshotaddon.html
type Instance_AutoSnapshotAddOn struct {

	// SnapshotTimeOfDay AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-instance-autosnapshotaddon.html#cfn-lightsail-instance-autosnapshotaddon-snapshottimeofday
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
func (r *Instance_AutoSnapshotAddOn) AWSCloudFormationType() string {
	return "AWS::Lightsail::Instance.AutoSnapshotAddOn"
}
