package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_NielsenConfiguration AWS CloudFormation Resource (AWS::MediaLive::Channel.NielsenConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-nielsenconfiguration.html
type Channel_NielsenConfiguration struct {

	// DistributorId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-nielsenconfiguration.html#cfn-medialive-channel-nielsenconfiguration-distributorid
	DistributorId *types.Value `json:"DistributorId,omitempty"`

	// NielsenPcmToId3Tagging AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-nielsenconfiguration.html#cfn-medialive-channel-nielsenconfiguration-nielsenpcmtoid3tagging
	NielsenPcmToId3Tagging *types.Value `json:"NielsenPcmToId3Tagging,omitempty"`

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
func (r *Channel_NielsenConfiguration) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.NielsenConfiguration"
}
