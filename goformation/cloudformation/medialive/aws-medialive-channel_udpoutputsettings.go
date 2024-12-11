package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_UdpOutputSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.UdpOutputSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-udpoutputsettings.html
type Channel_UdpOutputSettings struct {

	// BufferMsec AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-udpoutputsettings.html#cfn-medialive-channel-udpoutputsettings-buffermsec
	BufferMsec *types.Value `json:"BufferMsec,omitempty"`

	// ContainerSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-udpoutputsettings.html#cfn-medialive-channel-udpoutputsettings-containersettings
	ContainerSettings *Channel_UdpContainerSettings `json:"ContainerSettings,omitempty"`

	// Destination AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-udpoutputsettings.html#cfn-medialive-channel-udpoutputsettings-destination
	Destination *Channel_OutputLocationRef `json:"Destination,omitempty"`

	// FecOutputSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-udpoutputsettings.html#cfn-medialive-channel-udpoutputsettings-fecoutputsettings
	FecOutputSettings *Channel_FecOutputSettings `json:"FecOutputSettings,omitempty"`

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
func (r *Channel_UdpOutputSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.UdpOutputSettings"
}
