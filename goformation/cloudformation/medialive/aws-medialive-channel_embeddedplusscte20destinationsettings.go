package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_EmbeddedPlusScte20DestinationSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.EmbeddedPlusScte20DestinationSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-embeddedplusscte20destinationsettings.html
type Channel_EmbeddedPlusScte20DestinationSettings struct {

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
func (r *Channel_EmbeddedPlusScte20DestinationSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.EmbeddedPlusScte20DestinationSettings"
}
