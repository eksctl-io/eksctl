package transfer

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// User_PosixProfile AWS CloudFormation Resource (AWS::Transfer::User.PosixProfile)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-transfer-user-posixprofile.html
type User_PosixProfile struct {

	// Gid AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-transfer-user-posixprofile.html#cfn-transfer-user-posixprofile-gid
	Gid *types.Value `json:"Gid"`

	// SecondaryGids AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-transfer-user-posixprofile.html#cfn-transfer-user-posixprofile-secondarygids
	SecondaryGids *types.Value `json:"SecondaryGids,omitempty"`

	// Uid AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-transfer-user-posixprofile.html#cfn-transfer-user-posixprofile-uid
	Uid *types.Value `json:"Uid"`

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
func (r *User_PosixProfile) AWSCloudFormationType() string {
	return "AWS::Transfer::User.PosixProfile"
}
