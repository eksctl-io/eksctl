package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_OutputDestination AWS CloudFormation Resource (AWS::MediaLive::Channel.OutputDestination)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputdestination.html
type Channel_OutputDestination struct {

	// Id AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputdestination.html#cfn-medialive-channel-outputdestination-id
	Id *types.Value `json:"Id,omitempty"`

	// MediaPackageSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputdestination.html#cfn-medialive-channel-outputdestination-mediapackagesettings
	MediaPackageSettings []Channel_MediaPackageOutputDestinationSettings `json:"MediaPackageSettings,omitempty"`

	// MultiplexSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputdestination.html#cfn-medialive-channel-outputdestination-multiplexsettings
	MultiplexSettings *Channel_MultiplexProgramChannelDestinationSettings `json:"MultiplexSettings,omitempty"`

	// Settings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-outputdestination.html#cfn-medialive-channel-outputdestination-settings
	Settings []Channel_OutputDestinationSettings `json:"Settings,omitempty"`

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
func (r *Channel_OutputDestination) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.OutputDestination"
}
