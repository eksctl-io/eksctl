package cloudformation

import (
	"encoding/json"
)

// AWSWAFSqlInjectionMatchSet_SqlInjectionMatchTuple AWS CloudFormation Resource (AWS::WAF::SqlInjectionMatchSet.SqlInjectionMatchTuple)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-sqlinjectionmatchset-sqlinjectionmatchtuples.html
type AWSWAFSqlInjectionMatchSet_SqlInjectionMatchTuple struct {

	// FieldToMatch AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-sqlinjectionmatchset-sqlinjectionmatchtuples.html#cfn-waf-sqlinjectionmatchset-sqlinjectionmatchtuples-fieldtomatch
	FieldToMatch *AWSWAFSqlInjectionMatchSet_FieldToMatch `json:"FieldToMatch,omitempty"`

	// TextTransformation AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-waf-sqlinjectionmatchset-sqlinjectionmatchtuples.html#cfn-waf-sqlinjectionmatchset-sqlinjectionmatchtuples-texttransformation
	TextTransformation *Value `json:"TextTransformation,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSWAFSqlInjectionMatchSet_SqlInjectionMatchTuple) AWSCloudFormationType() string {
	return "AWS::WAF::SqlInjectionMatchSet.SqlInjectionMatchTuple"
}

func (r *AWSWAFSqlInjectionMatchSet_SqlInjectionMatchTuple) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
