package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_VideoSelectorColorSpaceSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.VideoSelectorColorSpaceSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-videoselectorcolorspacesettings.html
type Channel_VideoSelectorColorSpaceSettings struct {

	// Hdr10Settings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-videoselectorcolorspacesettings.html#cfn-medialive-channel-videoselectorcolorspacesettings-hdr10settings
	Hdr10Settings *Channel_Hdr10Settings `json:"Hdr10Settings,omitempty"`

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
func (r *Channel_VideoSelectorColorSpaceSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.VideoSelectorColorSpaceSettings"
}
