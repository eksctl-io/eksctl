package datasync

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// LocationHDFS_QopConfiguration AWS CloudFormation Resource (AWS::DataSync::LocationHDFS.QopConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-locationhdfs-qopconfiguration.html
type LocationHDFS_QopConfiguration struct {

	// DataTransferProtection AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-locationhdfs-qopconfiguration.html#cfn-datasync-locationhdfs-qopconfiguration-datatransferprotection
	DataTransferProtection *types.Value `json:"DataTransferProtection,omitempty"`

	// RpcProtection AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-datasync-locationhdfs-qopconfiguration.html#cfn-datasync-locationhdfs-qopconfiguration-rpcprotection
	RpcProtection *types.Value `json:"RpcProtection,omitempty"`

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
func (r *LocationHDFS_QopConfiguration) AWSCloudFormationType() string {
	return "AWS::DataSync::LocationHDFS.QopConfiguration"
}
