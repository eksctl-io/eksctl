package nimblestudio

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// LaunchProfile_StreamConfiguration AWS CloudFormation Resource (AWS::NimbleStudio::LaunchProfile.StreamConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-launchprofile-streamconfiguration.html
type LaunchProfile_StreamConfiguration struct {

	// ClipboardMode AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-launchprofile-streamconfiguration.html#cfn-nimblestudio-launchprofile-streamconfiguration-clipboardmode
	ClipboardMode *types.Value `json:"ClipboardMode,omitempty"`

	// Ec2InstanceTypes AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-launchprofile-streamconfiguration.html#cfn-nimblestudio-launchprofile-streamconfiguration-ec2instancetypes
	Ec2InstanceTypes *types.Value `json:"Ec2InstanceTypes,omitempty"`

	// MaxSessionLengthInMinutes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-launchprofile-streamconfiguration.html#cfn-nimblestudio-launchprofile-streamconfiguration-maxsessionlengthinminutes
	MaxSessionLengthInMinutes *types.Value `json:"MaxSessionLengthInMinutes,omitempty"`

	// StreamingImageIds AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-nimblestudio-launchprofile-streamconfiguration.html#cfn-nimblestudio-launchprofile-streamconfiguration-streamingimageids
	StreamingImageIds *types.Value `json:"StreamingImageIds,omitempty"`

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
func (r *LaunchProfile_StreamConfiguration) AWSCloudFormationType() string {
	return "AWS::NimbleStudio::LaunchProfile.StreamConfiguration"
}
