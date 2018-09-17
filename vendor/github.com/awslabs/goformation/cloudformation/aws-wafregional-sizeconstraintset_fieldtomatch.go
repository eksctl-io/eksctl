package cloudformation

import (
	"encoding/json"
)

// AWSWAFRegionalSizeConstraintSet_FieldToMatch AWS CloudFormation Resource (AWS::WAFRegional::SizeConstraintSet.FieldToMatch)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-sizeconstraintset-fieldtomatch.html
type AWSWAFRegionalSizeConstraintSet_FieldToMatch struct {

	// Data AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-sizeconstraintset-fieldtomatch.html#cfn-wafregional-sizeconstraintset-fieldtomatch-data
	Data *Value `json:"Data,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-sizeconstraintset-fieldtomatch.html#cfn-wafregional-sizeconstraintset-fieldtomatch-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSWAFRegionalSizeConstraintSet_FieldToMatch) AWSCloudFormationType() string {
	return "AWS::WAFRegional::SizeConstraintSet.FieldToMatch"
}

func (r *AWSWAFRegionalSizeConstraintSet_FieldToMatch) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
