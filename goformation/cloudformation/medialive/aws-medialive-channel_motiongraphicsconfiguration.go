package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_MotionGraphicsConfiguration AWS CloudFormation Resource (AWS::MediaLive::Channel.MotionGraphicsConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-motiongraphicsconfiguration.html
type Channel_MotionGraphicsConfiguration struct {

	// MotionGraphicsInsertion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-motiongraphicsconfiguration.html#cfn-medialive-channel-motiongraphicsconfiguration-motiongraphicsinsertion
	MotionGraphicsInsertion *types.Value `json:"MotionGraphicsInsertion,omitempty"`

	// MotionGraphicsSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-motiongraphicsconfiguration.html#cfn-medialive-channel-motiongraphicsconfiguration-motiongraphicssettings
	MotionGraphicsSettings *Channel_MotionGraphicsSettings `json:"MotionGraphicsSettings,omitempty"`

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
func (r *Channel_MotionGraphicsConfiguration) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.MotionGraphicsConfiguration"
}
