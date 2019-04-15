package cloudformation

import (
	"encoding/json"
)

// AWSBatchJobDefinition_NodeRangeProperty AWS CloudFormation Resource (AWS::Batch::JobDefinition.NodeRangeProperty)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-noderangeproperty.html
type AWSBatchJobDefinition_NodeRangeProperty struct {

	// Container AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-noderangeproperty.html#cfn-batch-jobdefinition-noderangeproperty-container
	Container *AWSBatchJobDefinition_ContainerProperties `json:"Container,omitempty"`

	// TargetNodes AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-noderangeproperty.html#cfn-batch-jobdefinition-noderangeproperty-targetnodes
	TargetNodes *Value `json:"TargetNodes,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSBatchJobDefinition_NodeRangeProperty) AWSCloudFormationType() string {
	return "AWS::Batch::JobDefinition.NodeRangeProperty"
}

func (r *AWSBatchJobDefinition_NodeRangeProperty) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
