package medialive

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Channel_HlsCdnSettings AWS CloudFormation Resource (AWS::MediaLive::Channel.HlsCdnSettings)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlscdnsettings.html
type Channel_HlsCdnSettings struct {

	// HlsAkamaiSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlscdnsettings.html#cfn-medialive-channel-hlscdnsettings-hlsakamaisettings
	HlsAkamaiSettings *Channel_HlsAkamaiSettings `json:"HlsAkamaiSettings,omitempty"`

	// HlsBasicPutSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlscdnsettings.html#cfn-medialive-channel-hlscdnsettings-hlsbasicputsettings
	HlsBasicPutSettings *Channel_HlsBasicPutSettings `json:"HlsBasicPutSettings,omitempty"`

	// HlsMediaStoreSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlscdnsettings.html#cfn-medialive-channel-hlscdnsettings-hlsmediastoresettings
	HlsMediaStoreSettings *Channel_HlsMediaStoreSettings `json:"HlsMediaStoreSettings,omitempty"`

	// HlsS3Settings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlscdnsettings.html#cfn-medialive-channel-hlscdnsettings-hlss3settings
	HlsS3Settings *Channel_HlsS3Settings `json:"HlsS3Settings,omitempty"`

	// HlsWebdavSettings AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-medialive-channel-hlscdnsettings.html#cfn-medialive-channel-hlscdnsettings-hlswebdavsettings
	HlsWebdavSettings *Channel_HlsWebdavSettings `json:"HlsWebdavSettings,omitempty"`

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
func (r *Channel_HlsCdnSettings) AWSCloudFormationType() string {
	return "AWS::MediaLive::Channel.HlsCdnSettings"
}
