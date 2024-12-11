package lambda

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Alias_ProvisionedConcurrencyConfiguration AWS CloudFormation Resource (AWS::Lambda::Alias.ProvisionedConcurrencyConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-alias-provisionedconcurrencyconfiguration.html
type Alias_ProvisionedConcurrencyConfiguration struct {

	// ProvisionedConcurrentExecutions AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-lambda-alias-provisionedconcurrencyconfiguration.html#cfn-lambda-alias-provisionedconcurrencyconfiguration-provisionedconcurrentexecutions
	ProvisionedConcurrentExecutions *types.Value `json:"ProvisionedConcurrentExecutions"`

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
func (r *Alias_ProvisionedConcurrencyConfiguration) AWSCloudFormationType() string {
	return "AWS::Lambda::Alias.ProvisionedConcurrencyConfiguration"
}
