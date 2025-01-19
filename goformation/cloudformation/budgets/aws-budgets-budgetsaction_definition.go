package budgets

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// BudgetsAction_Definition AWS CloudFormation Resource (AWS::Budgets::BudgetsAction.Definition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budgetsaction-definition.html
type BudgetsAction_Definition struct {

	// IamActionDefinition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budgetsaction-definition.html#cfn-budgets-budgetsaction-definition-iamactiondefinition
	IamActionDefinition *BudgetsAction_IamActionDefinition `json:"IamActionDefinition,omitempty"`

	// ScpActionDefinition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budgetsaction-definition.html#cfn-budgets-budgetsaction-definition-scpactiondefinition
	ScpActionDefinition *BudgetsAction_ScpActionDefinition `json:"ScpActionDefinition,omitempty"`

	// SsmActionDefinition AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budgetsaction-definition.html#cfn-budgets-budgetsaction-definition-ssmactiondefinition
	SsmActionDefinition *BudgetsAction_SsmActionDefinition `json:"SsmActionDefinition,omitempty"`

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
func (r *BudgetsAction_Definition) AWSCloudFormationType() string {
	return "AWS::Budgets::BudgetsAction.Definition"
}
