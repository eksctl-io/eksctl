package evidently

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Launch_LaunchGroupObject AWS CloudFormation Resource (AWS::Evidently::Launch.LaunchGroupObject)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-launch-launchgroupobject.html
type Launch_LaunchGroupObject struct {

	// Description AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-launch-launchgroupobject.html#cfn-evidently-launch-launchgroupobject-description
	Description *types.Value `json:"Description,omitempty"`

	// Feature AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-launch-launchgroupobject.html#cfn-evidently-launch-launchgroupobject-feature
	Feature *types.Value `json:"Feature,omitempty"`

	// GroupName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-launch-launchgroupobject.html#cfn-evidently-launch-launchgroupobject-groupname
	GroupName *types.Value `json:"GroupName,omitempty"`

	// Variation AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-launch-launchgroupobject.html#cfn-evidently-launch-launchgroupobject-variation
	Variation *types.Value `json:"Variation,omitempty"`

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
func (r *Launch_LaunchGroupObject) AWSCloudFormationType() string {
	return "AWS::Evidently::Launch.LaunchGroupObject"
}
