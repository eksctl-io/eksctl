package cloudformation

import (
	"encoding/json"
)

// AWSBudgetsBudget_CostTypes AWS CloudFormation Resource (AWS::Budgets::Budget.CostTypes)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-costtypes.html
type AWSBudgetsBudget_CostTypes struct {

	// IncludeCredit AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-costtypes.html#cfn-budgets-budget-costtypes-includecredit
	IncludeCredit *Value `json:"IncludeCredit,omitempty"`

	// IncludeDiscount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-costtypes.html#cfn-budgets-budget-costtypes-includediscount
	IncludeDiscount *Value `json:"IncludeDiscount,omitempty"`

	// IncludeOtherSubscription AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-costtypes.html#cfn-budgets-budget-costtypes-includeothersubscription
	IncludeOtherSubscription *Value `json:"IncludeOtherSubscription,omitempty"`

	// IncludeRecurring AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-costtypes.html#cfn-budgets-budget-costtypes-includerecurring
	IncludeRecurring *Value `json:"IncludeRecurring,omitempty"`

	// IncludeRefund AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-costtypes.html#cfn-budgets-budget-costtypes-includerefund
	IncludeRefund *Value `json:"IncludeRefund,omitempty"`

	// IncludeSubscription AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-costtypes.html#cfn-budgets-budget-costtypes-includesubscription
	IncludeSubscription *Value `json:"IncludeSubscription,omitempty"`

	// IncludeSupport AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-costtypes.html#cfn-budgets-budget-costtypes-includesupport
	IncludeSupport *Value `json:"IncludeSupport,omitempty"`

	// IncludeTax AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-costtypes.html#cfn-budgets-budget-costtypes-includetax
	IncludeTax *Value `json:"IncludeTax,omitempty"`

	// IncludeUpfront AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-costtypes.html#cfn-budgets-budget-costtypes-includeupfront
	IncludeUpfront *Value `json:"IncludeUpfront,omitempty"`

	// UseAmortized AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-costtypes.html#cfn-budgets-budget-costtypes-useamortized
	UseAmortized *Value `json:"UseAmortized,omitempty"`

	// UseBlended AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-budgets-budget-costtypes.html#cfn-budgets-budget-costtypes-useblended
	UseBlended *Value `json:"UseBlended,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSBudgetsBudget_CostTypes) AWSCloudFormationType() string {
	return "AWS::Budgets::Budget.CostTypes"
}

func (r *AWSBudgetsBudget_CostTypes) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
