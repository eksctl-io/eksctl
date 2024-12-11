package connect

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// User_UserPhoneConfig AWS CloudFormation Resource (AWS::Connect::User.UserPhoneConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-user-userphoneconfig.html
type User_UserPhoneConfig struct {

	// AfterContactWorkTimeLimit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-user-userphoneconfig.html#cfn-connect-user-userphoneconfig-aftercontactworktimelimit
	AfterContactWorkTimeLimit *types.Value `json:"AfterContactWorkTimeLimit,omitempty"`

	// AutoAccept AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-user-userphoneconfig.html#cfn-connect-user-userphoneconfig-autoaccept
	AutoAccept *types.Value `json:"AutoAccept,omitempty"`

	// DeskPhoneNumber AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-user-userphoneconfig.html#cfn-connect-user-userphoneconfig-deskphonenumber
	DeskPhoneNumber *types.Value `json:"DeskPhoneNumber,omitempty"`

	// PhoneType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-connect-user-userphoneconfig.html#cfn-connect-user-userphoneconfig-phonetype
	PhoneType *types.Value `json:"PhoneType,omitempty"`

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
func (r *User_UserPhoneConfig) AWSCloudFormationType() string {
	return "AWS::Connect::User.UserPhoneConfig"
}
