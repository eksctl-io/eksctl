package cloudformation

import (
	"encoding/json"
)

// AWSWAFRegionalSizeConstraintSet_SizeConstraint AWS CloudFormation Resource (AWS::WAFRegional::SizeConstraintSet.SizeConstraint)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-sizeconstraintset-sizeconstraint.html
type AWSWAFRegionalSizeConstraintSet_SizeConstraint struct {

	// ComparisonOperator AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-sizeconstraintset-sizeconstraint.html#cfn-wafregional-sizeconstraintset-sizeconstraint-comparisonoperator
	ComparisonOperator *Value `json:"ComparisonOperator,omitempty"`

	// FieldToMatch AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-sizeconstraintset-sizeconstraint.html#cfn-wafregional-sizeconstraintset-sizeconstraint-fieldtomatch
	FieldToMatch *AWSWAFRegionalSizeConstraintSet_FieldToMatch `json:"FieldToMatch,omitempty"`

	// Size AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-sizeconstraintset-sizeconstraint.html#cfn-wafregional-sizeconstraintset-sizeconstraint-size
	Size *Value `json:"Size,omitempty"`

	// TextTransformation AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-sizeconstraintset-sizeconstraint.html#cfn-wafregional-sizeconstraintset-sizeconstraint-texttransformation
	TextTransformation *Value `json:"TextTransformation,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSWAFRegionalSizeConstraintSet_SizeConstraint) AWSCloudFormationType() string {
	return "AWS::WAFRegional::SizeConstraintSet.SizeConstraint"
}

func (r *AWSWAFRegionalSizeConstraintSet_SizeConstraint) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
