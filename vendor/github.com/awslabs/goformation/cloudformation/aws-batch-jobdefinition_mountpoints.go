package cloudformation

import (
	"encoding/json"
)

// AWSBatchJobDefinition_MountPoints AWS CloudFormation Resource (AWS::Batch::JobDefinition.MountPoints)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-mountpoints.html
type AWSBatchJobDefinition_MountPoints struct {

	// ContainerPath AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-mountpoints.html#cfn-batch-jobdefinition-mountpoints-containerpath
	ContainerPath *Value `json:"ContainerPath,omitempty"`

	// ReadOnly AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-mountpoints.html#cfn-batch-jobdefinition-mountpoints-readonly
	ReadOnly *Value `json:"ReadOnly,omitempty"`

	// SourceVolume AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-batch-jobdefinition-mountpoints.html#cfn-batch-jobdefinition-mountpoints-sourcevolume
	SourceVolume *Value `json:"SourceVolume,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSBatchJobDefinition_MountPoints) AWSCloudFormationType() string {
	return "AWS::Batch::JobDefinition.MountPoints"
}

func (r *AWSBatchJobDefinition_MountPoints) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
