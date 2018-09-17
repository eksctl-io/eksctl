package cloudformation

import (
	"encoding/json"
)

// AWSGlueJob_ExecutionProperty AWS CloudFormation Resource (AWS::Glue::Job.ExecutionProperty)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-job-executionproperty.html
type AWSGlueJob_ExecutionProperty struct {

	// MaxConcurrentRuns AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-job-executionproperty.html#cfn-glue-job-executionproperty-maxconcurrentruns
	MaxConcurrentRuns *Value `json:"MaxConcurrentRuns,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGlueJob_ExecutionProperty) AWSCloudFormationType() string {
	return "AWS::Glue::Job.ExecutionProperty"
}

func (r *AWSGlueJob_ExecutionProperty) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
