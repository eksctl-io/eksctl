package cloudformation

import (
	"encoding/json"
)

// AWSBatchJobDefinition_Timeout AWS CloudFormation Resource (AWS::Batch::JobDefinition.Timeout)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-timeout.html
type AWSBatchJobDefinition_Timeout struct {

	// AttemptDurationSeconds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-timeout.html#cfn-batch-jobdefinition-timeout-attemptdurationseconds
	AttemptDurationSeconds *Value `json:"AttemptDurationSeconds,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSBatchJobDefinition_Timeout) AWSCloudFormationType() string {
	return "AWS::Batch::JobDefinition.Timeout"
}

func (r *AWSBatchJobDefinition_Timeout) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
