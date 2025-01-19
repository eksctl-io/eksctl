package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_FrameCaptureCdnSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.FrameCaptureCdnSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-framecapturecdnsettings.html
type Channel_FrameCaptureCdnSettings struct {

	// FrameCaptureS3Settings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-framecapturecdnsettings.html#cfn-medialive-channel-framecapturecdnsettings-framecaptures3settings
	FrameCaptureS3Settings *Channel_FrameCaptureS3Settings `json:"FrameCaptureS3Settings,omitempty"`

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
func (r *Channel_FrameCaptureCdnSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.FrameCaptureCdnSettings"
}
