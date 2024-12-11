package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_CaptionLanguageMapping AWS CloudFormation Resource (AWS::MediaLive::Channel.CaptionLanguageMapping)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captionlanguagemapping.html
type Channel_CaptionLanguageMapping struct {

	// CaptionChannel AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captionlanguagemapping.html#cfn-medialive-channel-captionlanguagemapping-captionchannel
	CaptionChannel *types.Value `json:"CaptionChannel,omitempty"`

	// LanguageCode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captionlanguagemapping.html#cfn-medialive-channel-captionlanguagemapping-languagecode
	LanguageCode *types.Value `json:"LanguageCode,omitempty"`

	// LanguageDescription AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-captionlanguagemapping.html#cfn-medialive-channel-captionlanguagemapping-languagedescription
	LanguageDescription *types.Value `json:"LanguageDescription,omitempty"`

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
func (r *Channel_CaptionLanguageMapping) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.CaptionLanguageMapping"
}
