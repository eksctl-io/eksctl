package cloudformation

import (
	"encoding/json"
)

// AWSBudgetsBudget_Spend AWS CloudFormation Resource (AWS::Budgets::Budget.Spend)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-spend.html
type AWSBudgetsBudget_Spend struct {

	// Amount AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-spend.html#cfn-budgets-budget-spend-amount
	Amount *Value `json:"Amount,omitempty"`

	// Unit AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-spend.html#cfn-budgets-budget-spend-unit
	Unit *Value `json:"Unit,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSBudgetsBudget_Spend) AWSCloudFormationType() string {
	return "AWS::Budgets::Budget.Spend"
}

func (r *AWSBudgetsBudget_Spend) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
