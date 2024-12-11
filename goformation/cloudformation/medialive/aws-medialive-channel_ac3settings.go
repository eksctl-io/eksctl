package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_Ac3Settings AWS CloudFormation Resource (AWS::MediaLive::Channel.Ac3Settings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-ac3settings.html
type Channel_Ac3Settings struct {

	// Bitrate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-ac3settings.html#cfn-medialive-channel-ac3settings-bitrate
	Bitrate *types.Value `json:"Bitrate,omitempty"`

	// BitstreamMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-ac3settings.html#cfn-medialive-channel-ac3settings-bitstreammode
	BitstreamMode *types.Value `json:"BitstreamMode,omitempty"`

	// CodingMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-ac3settings.html#cfn-medialive-channel-ac3settings-codingmode
	CodingMode *types.Value `json:"CodingMode,omitempty"`

	// Dialnorm AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-ac3settings.html#cfn-medialive-channel-ac3settings-dialnorm
	Dialnorm *types.Value `json:"Dialnorm,omitempty"`

	// DrcProfile AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-ac3settings.html#cfn-medialive-channel-ac3settings-drcprofile
	DrcProfile *types.Value `json:"DrcProfile,omitempty"`

	// LfeFilter AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-ac3settings.html#cfn-medialive-channel-ac3settings-lfefilter
	LfeFilter *types.Value `json:"LfeFilter,omitempty"`

	// MetadataControl AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-ac3settings.html#cfn-medialive-channel-ac3settings-metadatacontrol
	MetadataControl *types.Value `json:"MetadataControl,omitempty"`

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
func (r *Channel_Ac3Settings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.Ac3Settings"
}
