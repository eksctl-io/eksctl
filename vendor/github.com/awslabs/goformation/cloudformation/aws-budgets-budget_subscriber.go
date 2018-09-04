package cloudformation

// AWSBudgetsBudget_Subscriber AWS CloudFormation Resource (AWS::Budgets::Budget.Subscriber)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-subscriber.html
type AWSBudgetsBudget_Subscriber struct {

	// Address AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-subscriber.html#cfn-budgets-budget-subscriber-address
	Address *StringIntrinsic `json:"Address,omitempty"`

	// SubscriptionType AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-subscriber.html#cfn-budgets-budget-subscriber-subscriptiontype
	SubscriptionType *StringIntrinsic `json:"SubscriptionType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSBudgetsBudget_Subscriber) AWSCloudFormationType() string {
	return "AWS::Budgets::Budget.Subscriber"
}
