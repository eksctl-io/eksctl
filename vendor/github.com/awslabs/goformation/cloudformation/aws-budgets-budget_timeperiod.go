package cloudformation

// AWSBudgetsBudget_TimePeriod AWS CloudFormation Resource (AWS::Budgets::Budget.TimePeriod)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-timeperiod.html
type AWSBudgetsBudget_TimePeriod struct {

	// End AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-timeperiod.html#cfn-budgets-budget-timeperiod-end
	End *StringIntrinsic `json:"End,omitempty"`

	// Start AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-timeperiod.html#cfn-budgets-budget-timeperiod-start
	Start *StringIntrinsic `json:"Start,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSBudgetsBudget_TimePeriod) AWSCloudFormationType() string {
	return "AWS::Budgets::Budget.TimePeriod"
}
