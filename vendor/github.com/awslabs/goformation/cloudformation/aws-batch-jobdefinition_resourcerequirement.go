package cloudformation

import (
	"encoding/json"
)

// AWSBatchJobDefinition_ResourceRequirement AWS CloudFormation Resource (AWS::Batch::JobDefinition.ResourceRequirement)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-resourcerequirement.html
type AWSBatchJobDefinition_ResourceRequirement struct {

	// Type AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-resourcerequirement.html#cfn-batch-jobdefinition-resourcerequirement-type
	Type *Value `json:"Type,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-resourcerequirement.html#cfn-batch-jobdefinition-resourcerequirement-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSBatchJobDefinition_ResourceRequirement) AWSCloudFormationType() string {
	return "AWS::Batch::JobDefinition.ResourceRequirement"
}

func (r *AWSBatchJobDefinition_ResourceRequirement) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
