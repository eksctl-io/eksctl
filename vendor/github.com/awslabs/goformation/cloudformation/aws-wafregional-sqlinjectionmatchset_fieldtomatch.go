package cloudformation

import (
	"encoding/json"
)

// AWSWAFRegionalSqlInjectionMatchSet_FieldToMatch AWS CloudFormation Resource (AWS::WAFRegional::SqlInjectionMatchSet.FieldToMatch)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-sqlinjectionmatchset-fieldtomatch.html
type AWSWAFRegionalSqlInjectionMatchSet_FieldToMatch struct {

	// Data AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-sqlinjectionmatchset-fieldtomatch.html#cfn-wafregional-sqlinjectionmatchset-fieldtomatch-data
	Data *Value `json:"Data,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-wafregional-sqlinjectionmatchset-fieldtomatch.html#cfn-wafregional-sqlinjectionmatchset-fieldtomatch-type
	Type *Value `json:"Type,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSWAFRegionalSqlInjectionMatchSet_FieldToMatch) AWSCloudFormationType() string {
	return "AWS::WAFRegional::SqlInjectionMatchSet.FieldToMatch"
}

func (r *AWSWAFRegionalSqlInjectionMatchSet_FieldToMatch) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
