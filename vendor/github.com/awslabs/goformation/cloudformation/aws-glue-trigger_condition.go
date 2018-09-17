package cloudformation

import (
	"encoding/json"
)

// AWSGlueTrigger_Condition AWS CloudFormation Resource (AWS::Glue::Trigger.Condition)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-trigger-condition.html
type AWSGlueTrigger_Condition struct {

	// JobName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-trigger-condition.html#cfn-glue-trigger-condition-jobname
	JobName *Value `json:"JobName,omitempty"`

	// LogicalOperator AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-trigger-condition.html#cfn-glue-trigger-condition-logicaloperator
	LogicalOperator *Value `json:"LogicalOperator,omitempty"`

	// State AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-trigger-condition.html#cfn-glue-trigger-condition-state
	State *Value `json:"State,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGlueTrigger_Condition) AWSCloudFormationType() string {
	return "AWS::Glue::Trigger.Condition"
}

func (r *AWSGlueTrigger_Condition) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
