package connect

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// QuickConnect_QuickConnectConfig AWS CloudFormation Resource (AWS::Connect::QuickConnect.QuickConnectConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-quickconnect-quickconnectconfig.html
type QuickConnect_QuickConnectConfig struct {

	// PhoneConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-quickconnect-quickconnectconfig.html#cfn-connect-quickconnect-quickconnectconfig-phoneconfig
	PhoneConfig *QuickConnect_PhoneNumberQuickConnectConfig `json:"PhoneConfig,omitempty"`

	// QueueConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-quickconnect-quickconnectconfig.html#cfn-connect-quickconnect-quickconnectconfig-queueconfig
	QueueConfig *QuickConnect_QueueQuickConnectConfig `json:"QueueConfig,omitempty"`

	// QuickConnectType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-quickconnect-quickconnectconfig.html#cfn-connect-quickconnect-quickconnectconfig-quickconnecttype
	QuickConnectType *types.Value `json:"QuickConnectType,omitempty"`

	// UserConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-quickconnect-quickconnectconfig.html#cfn-connect-quickconnect-quickconnectconfig-userconfig
	UserConfig *QuickConnect_UserQuickConnectConfig `json:"UserConfig,omitempty"`

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
func (r *QuickConnect_QuickConnectConfig) AWSCloudFormationType() string {
	return "AWS::Connect::QuickConnect.QuickConnectConfig"
}
