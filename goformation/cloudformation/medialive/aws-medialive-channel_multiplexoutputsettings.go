package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_MultiplexOutputSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.MultiplexOutputSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-multiplexoutputsettings.html
type Channel_MultiplexOutputSettings struct {

	// Destination AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-multiplexoutputsettings.html#cfn-medialive-channel-multiplexoutputsettings-destination
	Destination *Channel_OutputLocationRef `json:"Destination,omitempty"`

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
func (r *Channel_MultiplexOutputSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.MultiplexOutputSettings"
}
