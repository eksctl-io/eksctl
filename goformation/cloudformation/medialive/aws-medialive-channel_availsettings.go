package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_AvailSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.AvailSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-availsettings.html
type Channel_AvailSettings struct {

	// Scte35SpliceInsert AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-availsettings.html#cfn-medialive-channel-availsettings-scte35spliceinsert
	Scte35SpliceInsert *Channel_Scte35SpliceInsert `json:"Scte35SpliceInsert,omitempty"`

	// Scte35TimeSignalApos AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-availsettings.html#cfn-medialive-channel-availsettings-scte35timesignalapos
	Scte35TimeSignalApos *Channel_Scte35TimeSignalApos `json:"Scte35TimeSignalApos,omitempty"`

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
func (r *Channel_AvailSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.AvailSettings"
}
