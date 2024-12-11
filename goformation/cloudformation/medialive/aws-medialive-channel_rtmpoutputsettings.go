package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_RtmpOutputSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.RtmpOutputSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-rtmpoutputsettings.html
type Channel_RtmpOutputSettings struct {

	// CertificateMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-rtmpoutputsettings.html#cfn-medialive-channel-rtmpoutputsettings-certificatemode
	CertificateMode *types.Value `json:"CertificateMode,omitempty"`

	// ConnectionRetryInterval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-rtmpoutputsettings.html#cfn-medialive-channel-rtmpoutputsettings-connectionretryinterval
	ConnectionRetryInterval *types.Value `json:"ConnectionRetryInterval,omitempty"`

	// Destination AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-rtmpoutputsettings.html#cfn-medialive-channel-rtmpoutputsettings-destination
	Destination *Channel_OutputLocationRef `json:"Destination,omitempty"`

	// NumRetries AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-rtmpoutputsettings.html#cfn-medialive-channel-rtmpoutputsettings-numretries
	NumRetries *types.Value `json:"NumRetries,omitempty"`

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
func (r *Channel_RtmpOutputSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.RtmpOutputSettings"
}
