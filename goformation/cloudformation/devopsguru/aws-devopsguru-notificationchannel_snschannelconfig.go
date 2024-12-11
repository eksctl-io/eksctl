package devopsguru

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// NotificationChannel_SnsChannelConfig AWS CloudFormation Resource (AWS::DevOpsGuru::NotificationChannel.SnsChannelConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-devopsguru-notificationchannel-snschannelconfig.html
type NotificationChannel_SnsChannelConfig struct {

	// TopicArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-devopsguru-notificationchannel-snschannelconfig.html#cfn-devopsguru-notificationchannel-snschannelconfig-topicarn
	TopicArn *types.Value `json:"TopicArn,omitempty"`

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
func (r *NotificationChannel_SnsChannelConfig) AWSCloudFormationType() string {
	return "AWS::DevOpsGuru::NotificationChannel.SnsChannelConfig"
}
