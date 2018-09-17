package cloudformation

import (
	"encoding/json"
)

// AWSGuardDutyFilter_Condition AWS CloudFormation Resource (AWS::GuardDuty::Filter.Condition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-guardduty-filter-condition.html
type AWSGuardDutyFilter_Condition struct {

	// Eq AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-guardduty-filter-condition.html#cfn-guardduty-filter-condition-eq
	Eq []*Value `json:"Eq,omitempty"`

	// Gte AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-guardduty-filter-condition.html#cfn-guardduty-filter-condition-gte
	Gte *Value `json:"Gte,omitempty"`

	// Lt AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-guardduty-filter-condition.html#cfn-guardduty-filter-condition-lt
	Lt *Value `json:"Lt,omitempty"`

	// Lte AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-guardduty-filter-condition.html#cfn-guardduty-filter-condition-lte
	Lte *Value `json:"Lte,omitempty"`

	// Neq AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-guardduty-filter-condition.html#cfn-guardduty-filter-condition-neq
	Neq []*Value `json:"Neq,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGuardDutyFilter_Condition) AWSCloudFormationType() string {
	return "AWS::GuardDuty::Filter.Condition"
}

func (r *AWSGuardDutyFilter_Condition) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
