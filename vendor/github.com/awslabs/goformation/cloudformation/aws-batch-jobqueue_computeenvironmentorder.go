package cloudformation

import (
	"encoding/json"
)

// AWSBatchJobQueue_ComputeEnvironmentOrder AWS CloudFormation Resource (AWS::Batch::JobQueue.ComputeEnvironmentOrder)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobqueue-computeenvironmentorder.html
type AWSBatchJobQueue_ComputeEnvironmentOrder struct {

	// ComputeEnvironment AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobqueue-computeenvironmentorder.html#cfn-batch-jobqueue-computeenvironmentorder-computeenvironment
	ComputeEnvironment *Value `json:"ComputeEnvironment,omitempty"`

	// Order AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobqueue-computeenvironmentorder.html#cfn-batch-jobqueue-computeenvironmentorder-order
	Order *Value `json:"Order,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSBatchJobQueue_ComputeEnvironmentOrder) AWSCloudFormationType() string {
	return "AWS::Batch::JobQueue.ComputeEnvironmentOrder"
}

func (r *AWSBatchJobQueue_ComputeEnvironmentOrder) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
