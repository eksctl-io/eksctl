package lightsail

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Database_RelationalDatabaseParameter AWS CloudFormation Resource (AWS::Lightsail::Database.RelationalDatabaseParameter)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-database-relationaldatabaseparameter.html
type Database_RelationalDatabaseParameter struct {

	// AllowedValues AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-database-relationaldatabaseparameter.html#cfn-lightsail-database-relationaldatabaseparameter-allowedvalues
	AllowedValues *types.Value `json:"AllowedValues,omitempty"`

	// ApplyMethod AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-database-relationaldatabaseparameter.html#cfn-lightsail-database-relationaldatabaseparameter-applymethod
	ApplyMethod *types.Value `json:"ApplyMethod,omitempty"`

	// ApplyType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-database-relationaldatabaseparameter.html#cfn-lightsail-database-relationaldatabaseparameter-applytype
	ApplyType *types.Value `json:"ApplyType,omitempty"`

	// DataType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-database-relationaldatabaseparameter.html#cfn-lightsail-database-relationaldatabaseparameter-datatype
	DataType *types.Value `json:"DataType,omitempty"`

	// Description AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-database-relationaldatabaseparameter.html#cfn-lightsail-database-relationaldatabaseparameter-description
	Description *types.Value `json:"Description,omitempty"`

	// IsModifiable AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-database-relationaldatabaseparameter.html#cfn-lightsail-database-relationaldatabaseparameter-ismodifiable
	IsModifiable *types.Value `json:"IsModifiable,omitempty"`

	// ParameterName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-database-relationaldatabaseparameter.html#cfn-lightsail-database-relationaldatabaseparameter-parametername
	ParameterName *types.Value `json:"ParameterName,omitempty"`

	// ParameterValue AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lightsail-database-relationaldatabaseparameter.html#cfn-lightsail-database-relationaldatabaseparameter-parametervalue
	ParameterValue *types.Value `json:"ParameterValue,omitempty"`

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
func (r *Database_RelationalDatabaseParameter) AWSCloudFormationType() string {
	return "AWS::Lightsail::Database.RelationalDatabaseParameter"
}
