package events

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Rule_RedshiftDataParameters AWS CloudFormation Resource (AWS::Events::Rule.RedshiftDataParameters)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-redshiftdataparameters.html
type Rule_RedshiftDataParameters struct {

	// Database AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-redshiftdataparameters.html#cfn-events-rule-redshiftdataparameters-database
	Database *types.Value `json:"Database,omitempty"`

	// DbUser AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-redshiftdataparameters.html#cfn-events-rule-redshiftdataparameters-dbuser
	DbUser *types.Value `json:"DbUser,omitempty"`

	// SecretManagerArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-redshiftdataparameters.html#cfn-events-rule-redshiftdataparameters-secretmanagerarn
	SecretManagerArn *types.Value `json:"SecretManagerArn,omitempty"`

	// Sql AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-redshiftdataparameters.html#cfn-events-rule-redshiftdataparameters-sql
	Sql *types.Value `json:"Sql,omitempty"`

	// StatementName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-redshiftdataparameters.html#cfn-events-rule-redshiftdataparameters-statementname
	StatementName *types.Value `json:"StatementName,omitempty"`

	// WithEvent AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-events-rule-redshiftdataparameters.html#cfn-events-rule-redshiftdataparameters-withevent
	WithEvent *types.Value `json:"WithEvent,omitempty"`

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
func (r *Rule_RedshiftDataParameters) AWSCloudFormationType() string {
	return "AWS::Events::Rule.RedshiftDataParameters"
}
