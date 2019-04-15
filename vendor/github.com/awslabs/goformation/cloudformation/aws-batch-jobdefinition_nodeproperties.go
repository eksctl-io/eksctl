package cloudformation

import (
	"encoding/json"
)

// AWSBatchJobDefinition_NodeProperties AWS CloudFormation Resource (AWS::Batch::JobDefinition.NodeProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-nodeproperties.html
type AWSBatchJobDefinition_NodeProperties struct {

	// MainNode AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-nodeproperties.html#cfn-batch-jobdefinition-nodeproperties-mainnode
	MainNode *Value `json:"MainNode,omitempty"`

	// NodeRangeProperties AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-nodeproperties.html#cfn-batch-jobdefinition-nodeproperties-noderangeproperties
	NodeRangeProperties []AWSBatchJobDefinition_NodeRangeProperty `json:"NodeRangeProperties,omitempty"`

	// NumNodes AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-nodeproperties.html#cfn-batch-jobdefinition-nodeproperties-numnodes
	NumNodes *Value `json:"NumNodes,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSBatchJobDefinition_NodeProperties) AWSCloudFormationType() string {
	return "AWS::Batch::JobDefinition.NodeProperties"
}

func (r *AWSBatchJobDefinition_NodeProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
