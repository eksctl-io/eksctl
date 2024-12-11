package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_MultiplexProgramChannelDestinationSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.MultiplexProgramChannelDestinationSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-multiplexprogramchanneldestinationsettings.html
type Channel_MultiplexProgramChannelDestinationSettings struct {

	// MultiplexId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-multiplexprogramchanneldestinationsettings.html#cfn-medialive-channel-multiplexprogramchanneldestinationsettings-multiplexid
	MultiplexId *types.Value `json:"MultiplexId,omitempty"`

	// ProgramName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-multiplexprogramchanneldestinationsettings.html#cfn-medialive-channel-multiplexprogramchanneldestinationsettings-programname
	ProgramName *types.Value `json:"ProgramName,omitempty"`

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
func (r *Channel_MultiplexProgramChannelDestinationSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.MultiplexProgramChannelDestinationSettings"
}
