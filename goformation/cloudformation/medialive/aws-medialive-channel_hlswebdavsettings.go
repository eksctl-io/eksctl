package medialive

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Channel_HlsWebdavSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.HlsWebdavSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlswebdavsettings.html
type Channel_HlsWebdavSettings struct {

	// ConnectionRetryInterval AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlswebdavsettings.html#cfn-medialive-channel-hlswebdavsettings-connectionretryinterval
	ConnectionRetryInterval *types.Value `json:"ConnectionRetryInterval,omitempty"`

	// FilecacheDuration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlswebdavsettings.html#cfn-medialive-channel-hlswebdavsettings-filecacheduration
	FilecacheDuration *types.Value `json:"FilecacheDuration,omitempty"`

	// HttpTransferMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlswebdavsettings.html#cfn-medialive-channel-hlswebdavsettings-httptransfermode
	HttpTransferMode *types.Value `json:"HttpTransferMode,omitempty"`

	// NumRetries AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlswebdavsettings.html#cfn-medialive-channel-hlswebdavsettings-numretries
	NumRetries *types.Value `json:"NumRetries,omitempty"`

	// RestartDelay AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlswebdavsettings.html#cfn-medialive-channel-hlswebdavsettings-restartdelay
	RestartDelay *types.Value `json:"RestartDelay,omitempty"`

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
func (r *Channel_HlsWebdavSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.HlsWebdavSettings"
}
