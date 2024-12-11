package cassandra

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Table_EncryptionSpecification AWS CloudFormation Resource (AWS::Cassandra::Table.EncryptionSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cassandra-table-encryptionspecification.html
type Table_EncryptionSpecification struct {

	// EncryptionType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cassandra-table-encryptionspecification.html#cfn-cassandra-table-encryptionspecification-encryptiontype
	EncryptionType *types.Value `json:"EncryptionType,omitempty"`

	// KmsKeyIdentifier AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cassandra-table-encryptionspecification.html#cfn-cassandra-table-encryptionspecification-kmskeyidentifier
	KmsKeyIdentifier *types.Value `json:"KmsKeyIdentifier,omitempty"`

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
func (r *Table_EncryptionSpecification) AWSCloudFormationType() string {
	return "AWS::Cassandra::Table.EncryptionSpecification"
}
