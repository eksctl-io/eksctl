package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_VideoBlackFailoverSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.VideoBlackFailoverSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-videoblackfailoversettings.html
type Channel_VideoBlackFailoverSettings struct {

	// BlackDetectThreshold AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-videoblackfailoversettings.html#cfn-medialive-channel-videoblackfailoversettings-blackdetectthreshold
	BlackDetectThreshold *types.Value `json:"BlackDetectThreshold,omitempty"`

	// VideoBlackThresholdMsec AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-videoblackfailoversettings.html#cfn-medialive-channel-videoblackfailoversettings-videoblackthresholdmsec
	VideoBlackThresholdMsec *types.Value `json:"VideoBlackThresholdMsec,omitempty"`

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
func (r *Channel_VideoBlackFailoverSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.VideoBlackFailoverSettings"
}
