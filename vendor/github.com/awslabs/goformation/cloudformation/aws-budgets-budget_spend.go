package cloudformation

// AWSBudgetsBudget_Spend AWS CloudFormation Resource (AWS::Budgets::Budget.Spend)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-spend.html
type AWSBudgetsBudget_Spend struct {

	// Amount AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-spend.html#cfn-budgets-budget-spend-amount
	Amount float64 `json:"Amount,omitempty"`

	// Unit AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-spend.html#cfn-budgets-budget-spend-unit
	Unit *StringIntrinsic `json:"Unit,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSBudgetsBudget_Spend) AWSCloudFormationType() string {
	return "AWS::Budgets::Budget.Spend"
}
