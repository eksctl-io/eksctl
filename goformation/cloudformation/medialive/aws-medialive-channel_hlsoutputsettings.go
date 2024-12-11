package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_HlsOutputSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.HlsOutputSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsoutputsettings.html
type Channel_HlsOutputSettings struct {

	// H265PackagingType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsoutputsettings.html#cfn-medialive-channel-hlsoutputsettings-h265packagingtype
	H265PackagingType *types.Value `json:"H265PackagingType,omitempty"`

	// HlsSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsoutputsettings.html#cfn-medialive-channel-hlsoutputsettings-hlssettings
	HlsSettings *Channel_HlsSettings `json:"HlsSettings,omitempty"`

	// NameModifier AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsoutputsettings.html#cfn-medialive-channel-hlsoutputsettings-namemodifier
	NameModifier *types.Value `json:"NameModifier,omitempty"`

	// SegmentModifier AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlsoutputsettings.html#cfn-medialive-channel-hlsoutputsettings-segmentmodifier
	SegmentModifier *types.Value `json:"SegmentModifier,omitempty"`

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
func (r *Channel_HlsOutputSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.HlsOutputSettings"
}
