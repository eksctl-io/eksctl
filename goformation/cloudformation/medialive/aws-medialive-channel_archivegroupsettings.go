package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_ArchiveGroupSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.ArchiveGroupSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-archivegroupsettings.html
type Channel_ArchiveGroupSettings struct {

	// ArchiveCdnSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-archivegroupsettings.html#cfn-medialive-channel-archivegroupsettings-archivecdnsettings
	ArchiveCdnSettings *Channel_ArchiveCdnSettings `json:"ArchiveCdnSettings,omitempty"`

	// Destination AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-archivegroupsettings.html#cfn-medialive-channel-archivegroupsettings-destination
	Destination *Channel_OutputLocationRef `json:"Destination,omitempty"`

	// RolloverInterval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-archivegroupsettings.html#cfn-medialive-channel-archivegroupsettings-rolloverinterval
	RolloverInterval *types.Value `json:"RolloverInterval,omitempty"`

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
func (r *Channel_ArchiveGroupSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.ArchiveGroupSettings"
}
